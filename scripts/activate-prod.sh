#!/bin/bash
REV="${1:?Revision required}"
set -euo pipefail
[[ "$REV" =~ ^[0-9a-f]{40}$ ]] || exit 1
BASE=/usr/local/inventory
RELEASE="$BASE/releases/$REV"
LIVE="$BASE/bin/inventory"
WEB="$BASE/web/public"
PLIST=/Library/LaunchDaemons/com.rob.inventory.plist
BACKUP="$BASE/releases/.rollback-$(date +%Y%m%d%H%M%S)-$$"
[[ -f "$RELEASE/inventory" && -f "$RELEASE/SHA256SUMS" && -f "$RELEASE/REVISION" ]] || exit 1
[[ "$(cat "$RELEASE/REVISION")" == "$REV" ]] || exit 1
[[ -x "$LIVE" && -d "$WEB" && -f "$PLIST" ]] || exit 1
[[ -f /etc/inventory/config.json && -d /var/lib/inventory/photos ]] || exit 1
mkdir -p "$BACKUP/web"
cp -p "$LIVE" "$BACKUP/inventory"
cp -p "$PLIST" "$BACKUP/service.plist"
rsync -a "$WEB/" "$BACKUP/web/"
# Rollback always uses the actual immediately preceding live files, even if
# an older release directory has been removed.
activated=0
rollback() {
  status=$?
  trap - ERR EXIT
  if (( activated == 1 )); then
    echo 'ERROR: Activation failed; restoring previous executable, web assets and service.' >&2
    launchctl bootout system/com.rob.inventory || true
    cp -p "$BACKUP/inventory" "$LIVE"
    rsync -a --delete "$BACKUP/web/" "$WEB/"
    cp -p "$BACKUP/service.plist" "$PLIST"
    chown root:wheel "$PLIST"
    chmod 644 "$PLIST"
    launchctl bootstrap system "$PLIST" || echo 'CRITICAL: Automatic rollback restart failed' >&2
  fi
  exit "$status"
}
trap rollback ERR EXIT
# Stop before replacing executable or web assets.
activated=1
launchctl bootout system/com.rob.inventory
install -m 755 "$RELEASE/inventory" "$LIVE"
rsync -a --delete "$RELEASE/web/public/" "$WEB/"
# No service plist change is necessary for subsequent deployments.
# Verify it still points to the stable installed paths.
/usr/libexec/PlistBuddy -c 'Print :ProgramArguments:0' "$PLIST" | grep -Fx "$LIVE"
/usr/libexec/PlistBuddy -c 'Print :ProgramArguments:2' "$PLIST" | grep -Fx '/etc/inventory/config.json'
launchctl bootstrap system "$PLIST"
healthy=0
for i in $(seq 1 20); do
  if curl -fsS --max-time 3 http://localhost:8080/ >/dev/null &&
     pgrep -fl 'inventory.*server' | grep -Fq "$LIVE --config /etc/inventory/config.json server --port 8080"; then
    healthy=1
    break
  fi
  sleep 2
done
[[ "$healthy" == 1 ]] || { echo 'ERROR: HTTP/process health check failed' >&2; false; }
[[ "$(shasum -a 256 "$LIVE" | awk '{print $1}')" == "$(shasum -a 256 "$RELEASE/inventory" | awk '{print $1}')" ]]
diff -qr "$RELEASE/web/public" "$WEB" >/dev/null
launchctl print system/com.rob.inventory | grep -q 'state = running'
# Preserve rollback snapshot; never touch config, photos, database, or backups.
trap - ERR EXIT
printf '%s\n' "$REV" > "$BASE/ACTIVE_REVISION"
echo "DEPLOYMENT SUCCEEDED: $REV"
echo "Previous live files saved at: $BACKUP"
