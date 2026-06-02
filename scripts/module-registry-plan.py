#!/usr/bin/env python3
"""Generate or apply registry SQL for module manifests."""

from __future__ import annotations

import argparse
import importlib.util
import json
import os
import subprocess
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
VALIDATOR = ROOT / "scripts" / "check-module-manifests.py"
WRITE_ENV = "MAKEADMIN_ALLOW_MODULE_REGISTRY_WRITE"


class PlanError(RuntimeError):
    pass


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Preview or apply module menu and permission registry SQL.")
    parser.add_argument(
        "--manifest",
        default="examples/demo/manifest.json",
        help="module manifest path, default: examples/demo/manifest.json",
    )
    parser.add_argument("--apply", action="store_true", help=f"write registry rows; requires {WRITE_ENV}=1")
    parser.add_argument("--confirm-module", help="required with --apply; must equal manifest module")
    parser.add_argument("--mysql-host", default=os.environ.get("MYSQL_HOST", "127.0.0.1"))
    parser.add_argument("--mysql-port", default=os.environ.get("MYSQL_PORT", "3306"))
    parser.add_argument("--mysql-user", default=os.environ.get("MYSQL_USER", "root"))
    parser.add_argument("--mysql-database", default=os.environ.get("MYSQL_DATABASE", "go_makeadmin"))
    return parser.parse_args()


def load_validator() -> Any:
    spec = importlib.util.spec_from_file_location("module_manifest_validator", VALIDATOR)
    if spec is None or spec.loader is None:
        raise PlanError(f"cannot load validator: {VALIDATOR}")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def load_manifest(path: Path) -> dict[str, Any]:
    validator = load_validator()
    validator.validate_manifest(path)
    return json.loads(path.read_text())


def print_sql(manifest: dict[str, Any]) -> None:
    print("-- Dry-run SQL preview. Review manually before applying; this script did not execute it.")
    print(build_sql(manifest))


def build_sql(manifest: dict[str, Any]) -> str:
    menu = manifest["menu"]
    permissions = manifest["permissions"]
    menu_permission = menu["permission"]
    statements = [
        "SET @now = UNIX_TIMESTAMP();",
        f"SET @parent_route_name = {sql_quote(menu['parent'])};",
        "\n".join(
            [
                "SET @parent_menu_id = COALESCE((",
                "    SELECT id FROM `ma_menu`",
                "    WHERE route_name = @parent_route_name AND delete_time = 0",
                "    LIMIT 1",
                "), 0);",
            ]
        ),
    ]
    for index, permission in enumerate(permissions):
        statements.append(permission_insert_sql(permission, sort=1000 - index * 10))
    statements.append(menu_insert_sql(menu))
    statements.append(
        "\n".join(
            [
                "SET @module_menu_id = COALESCE((",
                "    SELECT id FROM `ma_menu`",
                f"    WHERE route_name = {sql_quote(menu['routeName'])} AND delete_time = 0",
                "    LIMIT 1",
                "), 0);",
            ]
        )
    )
    statements.append(
        "\n".join(
            [
                "SET @module_permission_id = COALESCE((",
                "    SELECT id FROM `ma_permission`",
                f"    WHERE code = {sql_quote(menu_permission)}",
                "    LIMIT 1",
                "), 0);",
            ]
        )
    )
    statements.append(menu_permission_insert_sql())
    return "\n\n".join(statements)


def validate_apply_gate(args: argparse.Namespace, manifest: dict[str, Any]) -> None:
    if os.environ.get(WRITE_ENV) != "1":
        raise PlanError(f"--apply requires {WRITE_ENV}=1; no database access was attempted")
    if not args.confirm_module:
        raise PlanError("--apply requires --confirm-module; no database access was attempted")
    if args.confirm_module != manifest["module"]:
        raise PlanError("--confirm-module must equal manifest module; no database access was attempted")


def mysql_exec(args: argparse.Namespace, query: str) -> None:
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
        "--execute",
        query,
    ]
    env = os.environ.copy()
    if "MYSQL_PASSWORD" in env:
        env["MYSQL_PWD"] = env["MYSQL_PASSWORD"]
    try:
        subprocess.run(command, env=env, check=True, capture_output=True, text=True)
    except FileNotFoundError as exc:
        raise PlanError("mysql client is required") from exc
    except subprocess.CalledProcessError as exc:
        stderr = exc.stderr.strip()
        raise PlanError(f"mysql query failed: {stderr}") from exc


def apply_sql(args: argparse.Namespace, manifest: dict[str, Any]) -> None:
    mysql_exec(args, build_sql(manifest))
    print(f"Module registry apply: module={manifest['module']}")
    print(f"Permissions: planned={len(manifest['permissions'])}")
    print("Menu: planned=1")
    print("Menu permissions: planned=1")


def permission_insert_sql(permission: dict[str, Any], sort: int) -> str:
    return "\n".join(
        [
            "INSERT INTO `ma_permission`",
            "(`code`, `name`, `module`, `resource`, `action`, `status`, `sort`, `create_time`, `update_time`)",
            "SELECT",
            f"{sql_quote(permission['code'])}, {sql_quote(permission['name'])}, {sql_quote(permission['module'])}, "
            f"{sql_quote(permission['resource'])}, {sql_quote(permission['action'])}, 1, {sort}, @now, @now",
            "WHERE NOT EXISTS (",
            f"    SELECT 1 FROM `ma_permission` WHERE `code` = {sql_quote(permission['code'])}",
            ");",
        ]
    )


def menu_insert_sql(menu: dict[str, Any]) -> str:
    visible = 1 if menu["visible"] else 0
    return "\n".join(
        [
            "INSERT INTO `ma_menu`",
            "(`parent_id`, `menu_type`, `name`, `icon`, `route_path`, `route_name`, `component`, `redirect`, "
            "`active_path`, `meta`, `is_visible`, `is_cache`, `status`, `sort`, `create_time`, `update_time`, `delete_time`)",
            "SELECT",
            f"@parent_menu_id, {sql_quote(menu['type'])}, {sql_quote(menu['name'])}, '', {sql_quote(menu['routePath'])}, "
            f"{sql_quote(menu['routeName'])}, {sql_quote(menu['component'])}, '', '', '{{}}', {visible}, 1, 1, "
            f"{int(menu['sort'])}, @now, @now, 0",
            "WHERE NOT EXISTS (",
            f"    SELECT 1 FROM `ma_menu` WHERE `route_name` = {sql_quote(menu['routeName'])} AND `delete_time` = 0",
            ");",
        ]
    )


def menu_permission_insert_sql() -> str:
    return "\n".join(
        [
            "INSERT INTO `ma_menu_permission`",
            "(`menu_id`, `permission_id`, `create_time`)",
            "SELECT @module_menu_id, @module_permission_id, @now",
            "WHERE @module_menu_id > 0",
            "  AND @module_permission_id > 0",
            "  AND NOT EXISTS (",
            "      SELECT 1 FROM `ma_menu_permission`",
            "      WHERE `menu_id` = @module_menu_id AND `permission_id` = @module_permission_id",
            "  );",
        ]
    )


def sql_quote(value: Any) -> str:
    text = str(value)
    return "'" + text.replace("\\", "\\\\").replace("'", "''") + "'"


def main() -> int:
    args = parse_args()
    manifest_path = (ROOT / args.manifest).resolve()
    if not manifest_path.is_file():
        raise PlanError(f"manifest not found: {args.manifest}")
    if not manifest_path.is_relative_to(ROOT):
        raise PlanError("manifest must be inside repository")
    manifest = load_manifest(manifest_path)
    if args.apply:
        validate_apply_gate(args, manifest)
        apply_sql(args, manifest)
    else:
        print_sql(manifest)
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except PlanError as exc:
        print(f"FAIL: {exc}")
        raise SystemExit(1)
