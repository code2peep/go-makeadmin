#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SOURCE_FILE="${SOURCE_FILE:-$ROOT_DIR/sql/install.sql}"
OUT_FILE="${OUT_FILE:-$ROOT_DIR/sql/install.core.sql}"

core_tables=(
    la_album
    la_album_cate
    la_dict_data
    la_dict_type
    la_gen_table
    la_gen_table_column
    la_system_auth_admin
    la_system_auth_dept
    la_system_auth_menu
    la_system_auth_perm
    la_system_auth_post
    la_system_auth_role
    la_system_config
    la_system_log_login
    la_system_log_operate
)

if [ ! -f "$SOURCE_FILE" ]; then
    echo "FAIL: source SQL not found: $SOURCE_FILE"
    exit 1
fi

tables_csv="$(IFS=,; echo "${core_tables[*]}")"
tmp_file="$(mktemp)"

{
    echo "SET NAMES utf8mb4;"
    echo "SET FOREIGN_KEY_CHECKS = 0;"
    echo
    echo "-- Generated from sql/install.sql by scripts/build-core-sql.sh."
    echo "-- Keeps only go-makeadmin P0 core backend tables and seed data."

    awk -v tables="$tables_csv" '
        BEGIN {
            split(tables, names, ",")
            for (i in names) {
                keep[names[i]] = 1
            }
        }

        function table_name(line) {
            if (match(line, /`[^`]+`/)) {
                return substr(line, RSTART + 1, RLENGTH - 2)
            }
            return ""
        }

        /^DROP TABLE IF EXISTS `/ {
            table = table_name($0)
            in_table = keep[table]
            if (in_table) {
                print ""
                print "-- ----------------------------"
                print "-- Table structure for " table
                print "-- ----------------------------"
                print $0
            }
            next
        }

        in_table {
            print $0
            if ($0 ~ /^\).*;$/) {
                in_table = 0
            }
            next
        }

        /^INSERT INTO `/ {
            table = table_name($0)
            if (keep[table]) {
                if (table == "la_system_config" && $0 !~ /VALUES \([0-9]+, '\''(storage|website|protocol)'\'',/) {
                    next
                }
                if (table == "la_system_config" && $0 ~ /VALUES \([0-9]+, '\''website'\'', '\''(shopName|shopLogo)'\'',/) {
                    next
                }
                print $0
            }
            next
        }
    ' "$SOURCE_FILE"

    echo
    echo "SET FOREIGN_KEY_CHECKS = 1;"
} > "$tmp_file"

mv "$tmp_file" "$OUT_FILE"
echo "Generated $OUT_FILE"
