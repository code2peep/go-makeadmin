#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

require_text() {
    local file="$1"
    local text="$2"
    if ! grep -Fq "$text" "$file"; then
        echo "FAIL: expected $file to contain: $text"
        exit 1
    fi
}

require_text scripts/init-p1-db.sh 'ADMIN_PASSWORD="${ADMIN_PASSWORD:-123456}"'
require_text scripts/init-p1-db.sh "length must be 6-72 bytes"
require_text server/makeadmin/security/password.go "MinPasswordLength = 6"
require_text server/makeadmin/security/password_test.go 'ValidatePassword("123456")'
require_text docs/LOCAL_DEV.md "username: admin"
require_text docs/LOCAL_DEV.md "password: 123456"
require_text docs/P1_MINIMAL_SEED.md "admin / 123456"
require_text docs/P6_STATUS.md "admin / 123456"

echo "OK: local dev login contract passed"
