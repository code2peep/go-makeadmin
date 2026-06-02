#!/usr/bin/env python3
"""Build a dry-run tenant initialization plan.

This script only reads the source and target tenant rows, then prints a SQL
preview. It never executes the generated SQL.
"""

from __future__ import annotations

import argparse
import json
import os
import subprocess
import sys
from typing import Any


SETTING_GROUPS = ("website", "protocol", "storage")
SECRET_FIELDS = {"secretKey", "accessKey"}


class PlanError(RuntimeError):
    pass


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Preview tenant initialization SQL without writing database rows.")
    parser.add_argument("--from-tenant", type=int, default=0, help="source tenant id, default: 0")
    parser.add_argument("--to-tenant", type=int, required=True, help="target tenant id")
    parser.add_argument("--mysql-host", default=os.environ.get("MYSQL_HOST", "127.0.0.1"))
    parser.add_argument("--mysql-port", default=os.environ.get("MYSQL_PORT", "3306"))
    parser.add_argument("--mysql-user", default=os.environ.get("MYSQL_USER", "root"))
    parser.add_argument("--mysql-database", default=os.environ.get("MYSQL_DATABASE", "go_makeadmin"))
    parser.add_argument("--copy-secret", action="store_true", help="keep cloud storage accessKey/secretKey in SQL preview")
    parser.add_argument("--sql-only", action="store_true", help="print only SQL preview")
    parser.add_argument("--apply", action="store_true", help="reserved write mode; currently fails before database access")
    return parser.parse_args()


def mysql_json(args: argparse.Namespace, query: str) -> list[dict[str, Any]]:
    command = [
        "mysql",
        "--host",
        args.mysql_host,
        "--port",
        str(args.mysql_port),
        "--user",
        args.mysql_user,
        "--database",
        args.mysql_database,
        "--batch",
        "--raw",
        "--skip-column-names",
        "--execute",
        query,
    ]
    env = os.environ.copy()
    if "MYSQL_PASSWORD" in env:
        env["MYSQL_PWD"] = env["MYSQL_PASSWORD"]
    try:
        result = subprocess.run(command, env=env, check=True, capture_output=True, text=True)
    except FileNotFoundError as exc:
        raise PlanError("mysql client is required") from exc
    except subprocess.CalledProcessError as exc:
        stderr = exc.stderr.strip()
        raise PlanError(f"mysql query failed: {stderr}") from exc
    raw = result.stdout.strip()
    if not raw or raw == "null":
        return []
    payload = json.loads(raw)
    if not isinstance(payload, list):
        raise PlanError(f"mysql JSON result is not an array: {payload!r}")
    return payload


def source_settings(args: argparse.Namespace) -> list[dict[str, Any]]:
    groups = ",".join(sql_quote(group) for group in SETTING_GROUPS)
    return mysql_json(
        args,
        f"""
        SELECT COALESCE(JSON_ARRAYAGG(JSON_OBJECT(
            'setting_group', setting_group,
            'setting_key', setting_key,
            'setting_value', setting_value,
            'value_type', value_type,
            'is_public', is_public,
            'remark', remark
        )), JSON_ARRAY())
        FROM (
            SELECT setting_group, setting_key, setting_value, value_type, is_public, remark
            FROM ma_setting
            WHERE tenant_id = {args.from_tenant}
              AND setting_group IN ({groups})
            ORDER BY setting_group ASC, setting_key ASC
        ) AS source_settings;
        """,
    )


def target_setting_keys(args: argparse.Namespace) -> set[tuple[str, str]]:
    rows = mysql_json(
        args,
        f"""
        SELECT COALESCE(JSON_ARRAYAGG(JSON_OBJECT(
            'setting_group', setting_group,
            'setting_key', setting_key
        )), JSON_ARRAY())
        FROM ma_setting
        WHERE tenant_id = {args.to_tenant};
        """,
    )
    return {(str(row["setting_group"]), str(row["setting_key"])) for row in rows}


def source_file_categories(args: argparse.Namespace) -> list[dict[str, Any]]:
    return mysql_json(
        args,
        f"""
        SELECT COALESCE(JSON_ARRAYAGG(JSON_OBJECT(
            'id', id,
            'parent_id', parent_id,
            'parent_code', parent_code,
            'code', code,
            'name', name,
            'file_type', file_type,
            'status', status,
            'sort', sort
        )), JSON_ARRAY())
        FROM (
            SELECT category.id,
                   category.parent_id,
                   COALESCE(parent.code, '') AS parent_code,
                   category.code,
                   category.name,
                   category.file_type,
                   category.status,
                   category.sort
            FROM ma_file_category AS category
            LEFT JOIN ma_file_category AS parent ON parent.id = category.parent_id
            WHERE category.tenant_id = {args.from_tenant}
              AND category.delete_time = 0
            ORDER BY category.parent_id ASC, category.sort DESC, category.id ASC
        ) AS source_categories;
        """,
    )


def target_file_category_codes(args: argparse.Namespace) -> set[str]:
    rows = mysql_json(
        args,
        f"""
        SELECT COALESCE(JSON_ARRAYAGG(JSON_OBJECT('code', code)), JSON_ARRAY())
        FROM ma_file_category
        WHERE tenant_id = {args.to_tenant}
          AND delete_time = 0;
        """,
    )
    return {str(row["code"]) for row in rows}


def scrub_setting(row: dict[str, Any], copy_secret: bool) -> dict[str, Any]:
    row = dict(row)
    if copy_secret:
        return row
    if row.get("setting_group") == "storage" and row.get("setting_key") != "default":
        value = row.get("setting_value") or "{}"
        try:
            payload = json.loads(str(value))
        except json.JSONDecodeError:
            return row
        if isinstance(payload, dict):
            for field in SECRET_FIELDS:
                if field in payload:
                    payload[field] = ""
            row["setting_value"] = json.dumps(payload, ensure_ascii=False, separators=(",", ":"))
    return row


def build_plan(args: argparse.Namespace) -> dict[str, list[dict[str, Any]]]:
    settings = [scrub_setting(row, args.copy_secret) for row in source_settings(args)]
    existing_settings = target_setting_keys(args)
    categories = source_file_categories(args)
    existing_category_codes = target_file_category_codes(args)
    return {
        "settings_insert": [
            row for row in settings if (str(row["setting_group"]), str(row["setting_key"])) not in existing_settings
        ],
        "settings_skip": [
            row for row in settings if (str(row["setting_group"]), str(row["setting_key"])) in existing_settings
        ],
        "categories_insert": [row for row in categories if str(row["code"]) not in existing_category_codes],
        "categories_skip": [row for row in categories if str(row["code"]) in existing_category_codes],
    }


def print_plan(args: argparse.Namespace, plan: dict[str, list[dict[str, Any]]]) -> None:
    if args.sql_only:
        print_sql(args, plan)
        return
    print(f"Tenant init dry-run: from={args.from_tenant} to={args.to_tenant}")
    print(f"Settings: insert={len(plan['settings_insert'])} skip_existing={len(plan['settings_skip'])}")
    for row in plan["settings_insert"]:
        print(f"  + setting {row['setting_group']}.{row['setting_key']}")
    for row in plan["settings_skip"]:
        print(f"  = setting {row['setting_group']}.{row['setting_key']}")
    print(f"File categories: insert={len(plan['categories_insert'])} skip_existing={len(plan['categories_skip'])}")
    for row in plan["categories_insert"]:
        print(f"  + category {row['code']} ({row['name']})")
    for row in plan["categories_skip"]:
        print(f"  = category {row['code']} ({row['name']})")
    print()
    print_sql(args, plan)


def print_sql(args: argparse.Namespace, plan: dict[str, list[dict[str, Any]]]) -> None:
    print("-- Dry-run SQL preview. Review manually before applying; this script did not execute it.")
    print(f"SET @tenant_id = {args.to_tenant};")
    print("SET @now = UNIX_TIMESTAMP();")
    print()
    if plan["settings_insert"]:
        print("INSERT INTO `ma_setting`")
        print("(`tenant_id`, `setting_group`, `setting_key`, `setting_value`, `value_type`, `is_public`, `remark`, `create_time`, `update_time`)")
        print("VALUES")
        values = []
        for row in plan["settings_insert"]:
            values.append(
                "("
                "@tenant_id, "
                f"{sql_quote(row['setting_group'])}, "
                f"{sql_quote(row['setting_key'])}, "
                f"{sql_quote(row.get('setting_value', ''))}, "
                f"{sql_quote(row.get('value_type', 'string'))}, "
                f"{int(row.get('is_public') or 0)}, "
                f"{sql_quote(row.get('remark', ''))}, "
                "@now, @now"
                ")"
            )
        print(",\n".join(values) + ";")
        print()
    if plan["categories_insert"]:
        print("INSERT INTO `ma_file_category`")
        print("(`tenant_id`, `parent_id`, `code`, `name`, `file_type`, `status`, `sort`, `create_time`, `update_time`, `delete_time`)")
        print("VALUES")
        values = []
        for row in plan["categories_insert"]:
            parent_id_sql = parent_id_expression(row)
            values.append(
                "("
                "@tenant_id, "
                f"{parent_id_sql}, "
                f"{sql_quote(row['code'])}, "
                f"{sql_quote(row['name'])}, "
                f"{sql_quote(row['file_type'])}, "
                f"{int(row.get('status') or 1)}, "
                f"{int(row.get('sort') or 0)}, "
                "@now, @now, 0"
                ")"
            )
        print(",\n".join(values) + ";")


def parent_id_expression(row: dict[str, Any]) -> str:
    if int(row.get("parent_id") or 0) == 0:
        return "0"
    parent_code = str(row.get("parent_code") or "")
    if not parent_code:
        return "0"
    return (
        "(SELECT COALESCE(MAX(parent.id), 0) "
        "FROM `ma_file_category` AS parent "
        f"WHERE parent.tenant_id = @tenant_id AND parent.code = {sql_quote(parent_code)} AND parent.delete_time = 0)"
    )


def sql_quote(value: Any) -> str:
    text = str(value)
    return "'" + text.replace("\\", "\\\\").replace("'", "''") + "'"


def main() -> int:
    args = parse_args()
    if args.from_tenant < 0 or args.to_tenant < 0:
        raise PlanError("tenant ids must be non-negative")
    if args.from_tenant == args.to_tenant:
        raise PlanError("--from-tenant and --to-tenant must be different")
    if args.apply:
        raise PlanError("--apply is intentionally disabled until DB write approval is granted; no database access was attempted")
    plan = build_plan(args)
    print_plan(args, plan)
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except PlanError as exc:
        print(f"FAIL: {exc}", file=sys.stderr)
        raise SystemExit(1)
