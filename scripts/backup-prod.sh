#!/bin/bash
# Encrypted backup for the inventory production server.
set -euo pipefail
umask 077
export PATH="/usr/local/bin:/opt/homebrew/bin:/usr/local/mysql/bin:/usr/bin:/bin:/usr/sbin:/sbin"

LOG_DIR="/var/log/inventory"
mkdir -p "$LOG_DIR"
chmod 700 "$LOG_DIR"

exec > >(tee -a "$LOG_DIR/backup.log") 2>&1
echo "=== Backup started: $(date) ==="

CONFIG_FILE="/etc/inventory/config.json"
PHOTOS_DIR="/var/lib/inventory/photos"
LOCAL_DIR="$HOME/backups/inventory"
CLOUD_DIR="$HOME/Google Drive/My Drive/Inventory Backups"
RECIPIENT_FILE="$HOME/.config/age/inventory-backup-recipient.txt"
MYSQLDUMP="/usr/local/mysql/bin/mysqldump"

for cmd in age tar mktemp; do
  command -v "$cmd" >/dev/null || { echo "Missing command: $cmd" >&2; exit 1; }
done
[[ -x "$MYSQLDUMP" ]] || { echo "mysqldump not executable" >&2; exit 1; }
[[ -s "$RECIPIENT_FILE" ]] || { echo "Missing age recipient" >&2; exit 1; }
[[ -r "$CONFIG_FILE" ]] || { echo "Missing or unreadable config.json" >&2; exit 1; }
[[ -d "$PHOTOS_DIR" && -r "$PHOTOS_DIR" ]] || {
  echo "Missing or unreadable photos directory" >&2
  exit 1
}
mkdir -p "$LOCAL_DIR"
chmod 700 "$LOCAL_DIR"

stamp="$(date '+%Y-%m-%d_%H-%M-%S')"

if [[ "$(date +%u)" == "7" ]]; then
  backup_type="weekly"
else
  backup_type="daily"
fi
name="inventory-${backup_type}-${stamp}.tar.gz.age"

final="$LOCAL_DIR/$name"
partial="$LOCAL_DIR/.$name.partial"
work="$(mktemp -d "${TMPDIR:-/tmp/}inventory-backup.XXXXXXXX")"
cleanup() { rm -rf "$work"; rm -f "$partial"; }

finish() {
  status=$?

  cleanup

  if (( status == 0 )); then
    echo "=== Backup succeeded: $(date) ==="
  else
    echo "=== BACKUP FAILED (exit $status): $(date) ===" >&2
  fi
}

trap finish EXIT

# Temporary SQL stays in a private directory and is deleted on exit.
# This is not a plaintext backup in Google Drive.
"$MYSQLDUMP" --login-path=inventory-backup --single-transaction --quick \
  --no-tablespaces --routines --triggers --events inventory > "$work/database.sql"

# Stream the package through encryption, without creating an unencrypted tarball.
mkdir -p "$work/config" "$work/photos"

cp "$CONFIG_FILE" "$work/config/config.json"
cp -R "$PHOTOS_DIR/." "$work/photos/"

tar -czf - -C "$work" database.sql config photos | \
  age -R "$RECIPIENT_FILE" -o "$partial"

[[ -s "$partial" ]] || { echo "Empty encrypted backup" >&2; exit 1; }
mv "$partial" "$final"
echo "Local encrypted backup: $final"

prune_backups() {
  local dir="$1"
  local type="$2"
  local keep="$3"
  local files=()
  local i

  shopt -s nullglob
  files=("$dir"/inventory-"$type"-*.tar.gz.age)
  shopt -u nullglob

  # Filenames sort chronologically because timestamps use YYYY-MM-DD.
  if (( ${#files[@]} > keep )); then
    for (( i=0; i<${#files[@]}-keep; i++ )); do
      rm -- "${files[$i]}"
      echo "Removed old $type backup: ${files[$i]}"
    done
  fi
}

cloud_copy_ok=0

if [[ -d "$CLOUD_DIR" ]]; then
  cloud_partial="$CLOUD_DIR/.$name.partial"

  if cp "$final" "$cloud_partial" &&
     cmp -s "$final" "$cloud_partial" &&
     mv "$cloud_partial" "$CLOUD_DIR/$name"; then

    cloud_copy_ok=1
    echo "OFFSITE COPY: Saved to Google Drive folder."
    echo "Note: Google Drive upload/sync is not yet verified."
  else
    echo "WARNING: Google Drive copy failed; local backup retained." >&2
    rm -f "$cloud_partial" || true
  fi
else
  echo "WARNING: Google Drive folder unavailable; local backup retained." >&2
fi

# Retain 7 daily and 4 weekly backups.
# Preserve local history if the offsite copy fails.
if (( cloud_copy_ok == 1 )); then
  prune_backups "$LOCAL_DIR" daily 7
  prune_backups "$LOCAL_DIR" weekly 4

  prune_backups "$CLOUD_DIR" daily 7
  prune_backups "$CLOUD_DIR" weekly 4
else
  echo "WARNING: Local retention skipped because offsite copy failed."
fi

