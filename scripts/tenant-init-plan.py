#!/usr/bin/env python3
"""Build or apply a tenant initialization plan.

Dry-run mode only reads source and target tenant rows, then prints a SQL
preview. Apply mode requires explicit local write gates before any database
write is attempted.
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
WRITE_ENV = "MAKEADMIN_ALLOW_TENANT_INIT_WRITE"


class PlanError(RuntimeError):
    pass


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Preview or apply tenant initialization SQL.")
    parser.add_argument("--from-tenant", type=int, default=0, help="source tenant id, default: 0")
    parser.add_argument("--to-tenant", type=int, required=True, help="target tenant id")
    parser.add_argument("--confirm-to-tenant", type=int, help="required with --apply; must equal --to-tenant")
    parser.add_argument("--mysql-host", default=os.environ.get("MYSQL_HOST", "127.0.0.1"))
    parser.add_argument("--mysql-port", default=os.environ.get("MYSQL_PORT", "3306"))
    parser.add_argument("--mysql-user", default=os.environ.get("MYSQL_USER", "root"))
    parser.add_argument("--mysql-database", default=os.environ.get("MYSQL_DATABASE", "go_makeadmin"))
    parser.add_argument("--copy-secret", action="store_true", help="keep cloud storage accessKey/secretKey in SQL preview")
    parser.add_argument("--sql-only", action="store_true", help="print only SQL preview")
    parser.add_argument("--apply", action="store_true", help=f"write missing rows; requires {WRITE_ENV}=1")
    return parser.parse_args()


def mysql_run(args: argparse.Namespace, query: str) -> subprocess.CompletedProcess[str]:
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
        return subprocess.run(command, env=env, check=True, capture_output=True, text=True)
    except FileNotFoundError as exc:
        raise PlanError("mysql client is required") from exc
    except subprocess.CalledProcessError as exc:
        stderr = exc.stderr.strip()
        raise PlanError(f"mysql query failed: {stderr}") from exc


def mysql_json(args: argparse.Namespace, query: str) -> list[dict[str, Any]]:
    result = mysql_run(args, query)
    raw = result.stdout.strip()
    if not raw or raw == "null":
        return []
    payload = json.loads(raw)
    if not isinstance(payload, list):
        raise PlanError(f"mysql JSON result is not an array: {payload!r}")
    return payload


def mysql_exec(args: argparse.Namespace, query: str) -> None:
    mysql_run(args, query)


def validate_base_args(args: argparse.Namespace) -> None:
    if args.from_tenant < 0 or args.to_tenant < 0:
        raise PlanError("tenant ids must be non-negative")
    if args.from_tenant == args.to_tenant:
        raise PlanError("--from-tenant and --to-tenant must be different")
    if args.apply and args.sql_only:
        raise PlanError("--sql-only cannot be combined with --apply")


def validate_apply_gate(args: argparse.Namespace) -> None:
    if os.environ.get(WRITE_ENV) != "1":
        raise PlanError(f"--apply requires {WRITE_ENV}=1; no database access was attempted")
    if args.confirm_to_tenant is None:
        raise PlanError("--apply requires --confirm-to-tenant; no database access was attempted")
    if args.confirm_to_tenant != args.to_tenant:
        raise PlanError("--confirm-to-tenant must equal --to-tenant; no database access was attempted")


def ensure_target_tenant_enabled(args: argparse.Namespace) -> None:
    rows = mysql_json(
        args,
        f"""
        SELECT COALESCE(JSON_ARRAYAGG(JSON_OBJECT(
            'id', id,
            'status', status,
            'delete_time', delete_time
        )), JSON_ARRAY())
        FROM ma_tenant
        WHERE id = {args.to_tenant};
        """,
    )
    if not rows:
        raise PlanError(f"target tenant {args.to_tenant} does not exist")
    row = rows[0]
    if int(row.get("status") or 0) != 1 or int(row.get("delete_time") or 0) != 0:
        raise PlanError(f"target tenant {args.to_tenant} is not enabled")


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
    print(build_sql(args, plan))


def build_sql(args: argparse.Namespace, plan: dict[str, list[dict[str, Any]]]) -> str:
    statements = [
        f"SET @tenant_id = {args.to_tenant};",
        "SET @now = UNIX_TIMESTAMP();",
        "START TRANSACTION;",
    ]
    if plan["settings_insert"]:
        statements.append(setting_insert_sql(plan["settings_insert"]))
    if plan["categories_insert"]:
        for row in plan["categories_insert"]:
            statements.extend(category_insert_sql(row))
    if not plan["settings_insert"] and not plan["categories_insert"]:
        statements.append("-- No missing tenant initialization rows.")
    statements.append("COMMIT;")
    return "\n\n".join(statements)


def setting_insert_sql(rows: list[dict[str, Any]]) -> str:
    values = []
    for row in rows:
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
    return "\n".join(
        [
            "INSERT INTO `ma_setting`",
            "(`tenant_id`, `setting_group`, `setting_key`, `setting_value`, `value_type`, `is_public`, `remark`, `create_time`, `update_time`)",
            "VALUES",
            ",\n".join(values) + ";",
        ]
    )


def category_insert_sql(row: dict[str, Any]) -> list[str]:
    statements = [f"SET @parent_id = {parent_id_expression(row)};"]
    statements.append(
        "\n".join(
            [
                "INSERT INTO `ma_file_category`",
                "(`tenant_id`, `parent_id`, `code`, `name`, `file_type`, `status`, `sort`, `create_time`, `update_time`, `delete_time`)",
                "VALUES",
                "("
                "@tenant_id, "
                "@parent_id, "
                f"{sql_quote(row['code'])}, "
                f"{sql_quote(row['name'])}, "
                f"{sql_quote(row['file_type'])}, "
                f"{int(row.get('status') or 1)}, "
                f"{int(row.get('sort') or 0)}, "
                "@now, @now, 0"
                ");",
            ]
        )
    )
    return statements


def apply_plan(args: argparse.Namespace, plan: dict[str, list[dict[str, Any]]]) -> None:
    mysql_exec(args, build_sql(args, plan))
    print(f"Tenant init apply: from={args.from_tenant} to={args.to_tenant}")
    print(f"Settings: inserted={len(plan['settings_insert'])} skip_existing={len(plan['settings_skip'])}")
    print(f"File categories: inserted={len(plan['categories_insert'])} skip_existing={len(plan['categories_skip'])}")
    print("Transaction: committed")


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
    validate_base_args(args)
    if args.apply:
        validate_apply_gate(args)
        ensure_target_tenant_enabled(args)
    plan = build_plan(args)
    if args.apply:
        apply_plan(args, plan)
    else:
        print_plan(args, plan)
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except PlanError as exc:
        print(f"FAIL: {exc}", file=sys.stderr)
        raise SystemExit(1)
