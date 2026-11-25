#!/bin/bash
# Activate a staged inventory release on the production Mac.
# Installed root:wheel at /usr/local/libexec/inventory/activate-prod.sh.
set -euo pipefail

REV="${1:?Usage: activate-prod.sh REVISION}"
[[ "$REV" =~ ^[0-9a-f]{40}$ ]] || { echo 'ERROR: Invalid revision' >&2; exit 1; }
[[ "$(id -u)" == 0 ]] || { echo 'ERROR: Must run as root' >&2; exit 1; }

BASE=/usr/local/inventory
RELEASE="$BASE/releases/$REV"
LIVE="$BASE/bin/inventory"
WEB="$BASE/web/public"
PLIST=/Library/LaunchDaemons/com.rob.inventory.plist
SERVICE=system/com.rob.inventory
SAFE_STAGE=''
BACKUP=''
activated=0

# All release paths are fixed. Do not execute code from the release directory.
FILES=(inventory REVISION activate-prod.sh web/public/index.html web/public/css/style.css web/public/js/main.js)

verify_package() {
  local dir="$1" file expected
  [[ -f "$dir/SHA256SUMS" && ! -L "$dir/SHA256SUMS" ]] || return 1
  [[ -f "$dir/REVISION" && "$(cat "$dir/REVISION")" == "$REV" ]] || return 1
  [[ -d "$dir/web/public" && ! -L "$dir/web/public" ]] || return 1
  for file in "${FILES[@]}"; do
    [[ -f "$dir/$file" && ! -L "$dir/$file" ]] || return 1
  done
  # Reject extra files, symlinks and unexpected manifest entries.
  local listing expected
  listing="$(cd "$dir" && find . -type f | sed 's#^./##' | LC_ALL=C sort)"
  expected="$(printf '%s\n' "${FILES[@]}" SHA256SUMS | LC_ALL=C sort)"
  [[ "$listing" == "$expected" ]] || return 1
  [[ -z "$(find "$dir" -type l -print -quit)" ]] || return 1
  local -a names=()
  local hash name
  while read -r hash name; do
    [[ "$hash" =~ ^[0-9a-f]{64}$ && -n "$name" ]] || return 1
    names+=("$name")
  done < "$dir/SHA256SUMS"
  [[ "$(printf '%s\n' "${names[@]}" | LC_ALL=C sort)" == "$(printf '%s\n' "${FILES[@]}" | LC_ALL=C sort)" ]] || return 1
  (cd "$dir" && shasum -a 256 -c SHA256SUMS >/dev/null) || return 1
}

service_healthy() {
  launchctl print "$SERVICE" 2>/dev/null | grep -q 'state = running' &&
  curl -fsS --max-time 3 http://localhost:8080/ >/dev/null 2>&1
}

wait_healthy() {
  local i
  for ((i=0; i<20; i++)); do
    if service_healthy; then return 0; fi
    sleep 2
  done
  return 1
}

on_exit() {
  local status=$?
  trap - EXIT
  set +e
  if (( status != 0 && activated == 1 )); then
    echo 'ERROR: Activation failed; attempting rollback.' >&2
    launchctl bootout "$SERVICE" >/dev/null 2>&1 || true
    local restore_failed=0
    cp -p "$BACKUP/inventory" "$LIVE" || restore_failed=1
    rsync -a --delete "$BACKUP/web/" "$WEB/" || restore_failed=1
    cp -p "$BACKUP/service.plist" "$PLIST" || restore_failed=1
    chown root:wheel "$PLIST" || restore_failed=1
    chmod 644 "$PLIST" || restore_failed=1
    if (( restore_failed )); then
      echo 'CRITICAL: Some rollback file restorations failed.' >&2
    fi
    if launchctl bootstrap system "$PLIST" && wait_healthy; then
      if (( restore_failed == 0 )); then
        echo 'ROLLBACK VERIFIED: Previous files restored; service responds.' >&2
      else
        echo 'CRITICAL: Service responds but rollback restoration was incomplete.' >&2
      fi
    else
      echo 'CRITICAL: Rolled-back service did not pass health checks.' >&2
    fi
  fi
  if [[ -n "$SAFE_STAGE" && -d "$SAFE_STAGE" ]]; then
    rm -rf -- "$SAFE_STAGE" || echo 'WARNING: Could not remove protected stage.' >&2
  fi
  exit "$status"
}
trap on_exit EXIT

[[ -d "$RELEASE" && ! -L "$RELEASE" ]] || { echo 'ERROR: Release directory missing' >&2; exit 1; }
[[ -x "$LIVE" && -d "$WEB" && -f "$PLIST" ]] || { echo 'ERROR: Live installation incomplete' >&2; exit 1; }
[[ -f /etc/inventory/config.json && -d /var/lib/inventory/photos ]] || { echo 'ERROR: Config/photos missing' >&2; exit 1; }

# Root-controlled private copy. Only this copy is used for installation.
SAFE_STAGE="$(mktemp -d "$BASE/.activate.XXXXXXXX")"
chmod 700 "$SAFE_STAGE"
rsync -a "$RELEASE/" "$SAFE_STAGE/"
verify_package "$SAFE_STAGE" || { echo 'ERROR: Protected release validation failed' >&2; exit 1; }

# Check plist before stopping the service.
[[ "$(/usr/libexec/PlistBuddy -c 'Print :ProgramArguments:0' "$PLIST")" == "$LIVE" ]] || exit 1
[[ "$(/usr/libexec/PlistBuddy -c 'Print :ProgramArguments:2' "$PLIST")" == '/etc/inventory/config.json' ]] || exit 1

BACKUP="$(mktemp -d "$BASE/releases/.rollback.XXXXXXXX")"
mkdir -p "$BACKUP/web"
cp -p "$LIVE" "$BACKUP/inventory"
cp -p "$PLIST" "$BACKUP/service.plist"
rsync -a "$WEB/" "$BACKUP/web/"

echo "Activating $REV"
activated=1
launchctl bootout "$SERVICE"
install -m 755 "$SAFE_STAGE/inventory" "$LIVE"
rsync -a --delete "$SAFE_STAGE/web/public/" "$WEB/"
launchctl bootstrap system "$PLIST"
wait_healthy || { echo 'ERROR: New service failed health check' >&2; exit 1; }
[[ "$(shasum -a 256 "$LIVE" | awk '{print $1}')" == "$(shasum -a 256 "$SAFE_STAGE/inventory" | awk '{print $1}')" ]] || exit 1
diff -qr "$SAFE_STAGE/web/public" "$WEB" >/dev/null || exit 1
printf '%s\n' "$REV" > "$BASE/ACTIVE_REVISION"
echo "DEPLOYMENT SUCCEEDED: $REV"
echo "Previous live files saved at: $BACKUP"
