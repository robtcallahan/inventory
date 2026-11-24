#!/bin/bash
# Inventory deployment, executed on the development Mac.
# Usage: scripts/deploy-prod.sh [--dry-run|--deploy]
set -euo pipefail

MODE="${1:---dry-run}"
[[ $# -le 1 && ( "$MODE" == --dry-run || "$MODE" == --deploy ) ]] || {
  echo 'Usage: deploy-prod.sh [--dry-run|--deploy]' >&2; exit 2;
}

HOST='rob@ebay-mac.local'
BRANCH='production-layout'
ROOT="$(git rev-parse --show-toplevel)"
cd "$ROOT"
[[ "$(git branch --show-current)" == "$BRANCH" ]] || { echo "ERROR: Must be on $BRANCH" >&2; exit 1; }
[[ -z "$(git status --porcelain --untracked-files=all)" ]] || {
  echo 'ERROR: Working tree is not clean. Commit or remove changes first.' >&2
  git status --short >&2
  exit 1
}
REV="$(git rev-parse HEAD)"
[[ "$REV" =~ ^[0-9a-f]{40}$ ]] || exit 1
command -v go >/dev/null
command -v rsync >/dev/null
command -v ssh >/dev/null
command -v shasum >/dev/null

echo "Testing revision $REV"
go test ./...

STAGE="$(mktemp -d "${TMPDIR:-/tmp/}inventory-release.XXXXXXXX")"
trap 'rm -rf "$STAGE"' EXIT
mkdir -p "$STAGE/web/public"
echo "Building darwin/amd64 executable"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -buildvcs=true -o "$STAGE/inventory" .
META="$(go version -m "$STAGE/inventory")"
grep -Fq "vcs.revision=$REV" <<< "$META" || { echo 'ERROR: Binary revision mismatch' >&2; exit 1; }
grep -Fq 'vcs.modified=false' <<< "$META" || { echo 'ERROR: Build has local changes' >&2; exit 1; }
grep -Fq 'GOARCH=amd64' <<< "$META" || exit 1
grep -Fq 'GOOS=darwin' <<< "$META" || exit 1
rsync -a --delete web/public/ "$STAGE/web/public/"
printf '%s\n' "$REV" > "$STAGE/REVISION"
(
  cd "$STAGE"
  find inventory REVISION web/public -type f -print | LC_ALL=C sort | while IFS= read -r f; do
    shasum -a 256 "$f"
  done > SHA256SUMS
)
echo "Package: $REV ($(find "$STAGE/web/public" -type f | wc -l | tr -d ' ') web files)"
if [[ "$MODE" == --dry-run ]]; then
  echo 'DRY RUN PASSED: Tests, Intel build, revision and local checksums verified; production untouched.'
  exit 0
fi

# Preflight: confirm host, directories, running service and config before copying anything.
echo "Checking production host $HOST"
ssh -o BatchMode=yes "$HOST" 'set -eu; test "$(hostname -s)" = "ebay-mac"; test -f /etc/inventory/config.json; test -d /var/lib/inventory/photos; test -x /usr/local/inventory/bin/inventory; test -f /Library/LaunchDaemons/com.rob.inventory.plist; curl -fsS --max-time 5 http://localhost:8080/ >/dev/null'
REMOTE="/usr/local/inventory/releases/$REV"
echo "Staging $REMOTE"
ssh "$HOST" "mkdir -p '$REMOTE'"
cp "$ROOT/scripts/activate-prod.sh" "$STAGE/activate-prod.sh"
( cd "$STAGE" && shasum -a 256 activate-prod.sh >> SHA256SUMS )
rsync -a --delete "$STAGE/" "$HOST:$REMOTE/"
# Verify all files, and ensure no unexpected extra files exist.
ssh "$HOST" "REV='$REV' /bin/bash -s" <<'VERIFY'
set -euo pipefail
cd "/usr/local/inventory/releases/$REV"
while read -r hash file; do
  [[ -f "$file" ]] || exit 1
  [[ "$(shasum -a 256 "$file" | awk '{print $1}')" == "$hash" ]] || exit 1
done < SHA256SUMS
[[ "$(cat REVISION)" == "$REV" ]]
[[ "$(find inventory REVISION activate-prod.sh web/public -type f | wc -l | tr -d ' ')" == "$(wc -l < SHA256SUMS | tr -d ' ')" ]]
echo 'Remote checksums verified'
VERIFY

# Remote activation runs in one authenticated SSH session. sudo runs a fixed
# root-owned-location deployment transaction, with rollback on any error.
echo 'Activating release (remote sudo password may be requested)'
ssh -tt "$HOST" "sudo /bin/bash '$REMOTE/activate-prod.sh' '$REV'"
