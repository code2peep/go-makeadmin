#!/usr/bin/env python3
"""Generate a dry-run SQL preview for module registry manifests."""

from __future__ import annotations

import argparse
import importlib.util
import json
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
VALIDATOR = ROOT / "scripts" / "check-module-manifests.py"


class PlanError(RuntimeError):
    pass


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Preview module menu and permission registry SQL.")
    parser.add_argument(
        "--manifest",
        default="examples/demo/manifest.json",
        help="module manifest path, default: examples/demo/manifest.json",
    )
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
    menu = manifest["menu"]
    permissions = manifest["permissions"]
    menu_permission = menu["permission"]
    print("-- Dry-run SQL preview. Review manually before applying; this script did not execute it.")
    print("SET @now = UNIX_TIMESTAMP();")
    print(f"SET @parent_route_name = {sql_quote(menu['parent'])};")
    print("SET @parent_menu_id = COALESCE((")
    print("    SELECT id FROM `ma_menu`")
    print("    WHERE route_name = @parent_route_name AND delete_time = 0")
    print("    LIMIT 1")
    print("), 0);")
    print()
    for index, permission in enumerate(permissions):
        print(permission_insert_sql(permission, sort=1000 - index * 10))
        print()
    print(menu_insert_sql(menu))
    print()
    print("SET @module_menu_id = COALESCE((")
    print("    SELECT id FROM `ma_menu`")
    print(f"    WHERE route_name = {sql_quote(menu['routeName'])} AND delete_time = 0")
    print("    LIMIT 1")
    print("), 0);")
    print("SET @module_permission_id = COALESCE((")
    print("    SELECT id FROM `ma_permission`")
    print(f"    WHERE code = {sql_quote(menu_permission)}")
    print("    LIMIT 1")
    print("), 0);")
    print(menu_permission_insert_sql())


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
    print_sql(manifest)
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except PlanError as exc:
        print(f"FAIL: {exc}")
        raise SystemExit(1)
