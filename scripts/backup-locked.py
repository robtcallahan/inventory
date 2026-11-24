#!/usr/local/bin/python3

import fcntl
import os
import subprocess
import sys

lock_path = os.path.expanduser("~/backups/inventory/.backup.lock")

with open(lock_path, "a") as lock_file:
    try:
        fcntl.flock(lock_file, fcntl.LOCK_EX | fcntl.LOCK_NB)
    except BlockingIOError:
        print("ERROR: Another inventory backup is running.", file=sys.stderr)
        sys.exit(1)

    result = subprocess.run(
        ["/bin/bash", "/usr/local/libexec/inventory/backup-prod.sh"]
    )
    sys.exit(result.returncode)
