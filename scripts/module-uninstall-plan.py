#!/usr/bin/env python3
"""Generate a dry-run uninstall SQL preview for module manifests."""

from __future__ import annotations

import argparse
import importlib.util
import os
import subprocess
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
REGISTRY_PLAN = ROOT / "scripts" / "module-registry-plan.py"
WRITE_ENV = "MAKEADMIN_ALLOW_MODULE_UNINSTALL_WRITE"


class PlanError(RuntimeError):
    pass


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Preview or apply module registry cleanup SQL.")
    parser.add_argument(
        "--manifest",
        default="examples/demo/manifest.json",
        help="module manifest path, default: examples/demo/manifest.json",
    )
    parser.add_argument("--apply", action="store_true", help=f"delete module registry rows; requires {WRITE_ENV}=1")
    parser.add_argument("--confirm-module", help="required with --apply; must equal manifest module")
    parser.add_argument("--confirm-delete", action="store_true", help="required with --apply")
    parser.add_argument("--mysql-host", default=os.environ.get("MYSQL_HOST", "127.0.0.1"))
    parser.add_argument("--mysql-port", default=os.environ.get("MYSQL_PORT", "3306"))
    parser.add_argument("--mysql-user", default=os.environ.get("MYSQL_USER", "root"))
    parser.add_argument("--mysql-database", default=os.environ.get("MYSQL_DATABASE", "go_makeadmin"))
    return parser.parse_args()


def load_script(name: str, path: Path) -> Any:
    spec = importlib.util.spec_from_file_location(name, path)
    if spec is None or spec.loader is None:
        raise PlanError(f"cannot load script: {path}")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def build_sql(registry: Any, manifest: dict[str, Any]) -> str:
    codes = [permission["code"] for permission in manifest["permissions"]]
    code_list = ", ".join(registry.sql_quote(code) for code in codes)
    route_name = registry.sql_quote(manifest["menu"]["routeName"])
    statements = [
        "SET @module_route_name = " + route_name + ";",
        "\n".join(
            [
                "DELETE rp FROM `ma_role_permission` AS rp",
                "INNER JOIN `ma_permission` AS p ON p.id = rp.permission_id",
                f"WHERE p.code IN ({code_list});",
            ]
        ),
        "\n".join(
            [
                "DELETE mp FROM `ma_menu_permission` AS mp",
                "LEFT JOIN `ma_menu` AS m ON m.id = mp.menu_id",
                "LEFT JOIN `ma_permission` AS p ON p.id = mp.permission_id",
                f"WHERE m.route_name = @module_route_name OR p.code IN ({code_list});",
            ]
        ),
        "\n".join(
            [
                "DELETE FROM `ma_menu`",
                "WHERE route_name = @module_route_name;",
            ]
        ),
        "\n".join(
            [
                "DELETE FROM `ma_permission`",
                f"WHERE code IN ({code_list});",
            ]
        ),
    ]
    return "\n\n".join(statements)


def validate_apply_gate(args: argparse.Namespace, manifest: dict[str, Any]) -> None:
    if os.environ.get(WRITE_ENV) != "1":
        raise PlanError(f"--apply requires {WRITE_ENV}=1; no database access was attempted")
    if not args.confirm_module:
        raise PlanError("--apply requires --confirm-module; no database access was attempted")
    if args.confirm_module != manifest["module"]:
        raise PlanError("--confirm-module must equal manifest module; no database access was attempted")
    if not args.confirm_delete:
        raise PlanError("--apply requires --confirm-delete; no database access was attempted")


def mysql_query(args: argparse.Namespace, query: str) -> str:
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
    return result.stdout.strip()


def snapshot_sql(registry: Any, manifest: dict[str, Any]) -> str:
    codes = [permission["code"] for permission in manifest["permissions"]]
    code_list = ", ".join(registry.sql_quote(code) for code in codes)
    route_name = registry.sql_quote(manifest["menu"]["routeName"])
    return "\n".join(
        [
            "SELECT",
            f"(SELECT COUNT(*) FROM `ma_permission` WHERE code IN ({code_list})),",
            f"(SELECT COUNT(*) FROM `ma_menu` WHERE route_name = {route_name}),",
            "(",
            "    SELECT COUNT(*) FROM `ma_menu_permission` AS mp",
            "    LEFT JOIN `ma_menu` AS m ON m.id = mp.menu_id",
            "    LEFT JOIN `ma_permission` AS p ON p.id = mp.permission_id",
            f"    WHERE m.route_name = {route_name} OR p.code IN ({code_list})",
            "),",
            "(",
            "    SELECT COUNT(*) FROM `ma_role_permission` AS rp",
            "    INNER JOIN `ma_permission` AS p ON p.id = rp.permission_id",
            f"    WHERE p.code IN ({code_list})",
            ");",
        ]
    )


def parse_snapshot(raw: str) -> tuple[int, int, int, int]:
    parts = raw.split("\t")
    if len(parts) != 4:
        raise PlanError(f"unexpected mysql snapshot output: {raw}")
    return tuple(int(part) for part in parts)  # type: ignore[return-value]


def apply_uninstall_sql(args: argparse.Namespace, registry: Any, manifest: dict[str, Any]) -> None:
    before = parse_snapshot(mysql_query(args, snapshot_sql(registry, manifest)))
    print(f"Module uninstall apply: module={manifest['module']}")
    print(f"Before: permissions={before[0]} menus={before[1]} menuPermissions={before[2]} rolePermissions={before[3]}")
    if sum(before) == 0:
        print("Result: no-op")
        return
    mysql_query(args, f"START TRANSACTION;\n\n{build_sql(registry, manifest)}\n\nCOMMIT;")
    after = parse_snapshot(mysql_query(args, snapshot_sql(registry, manifest)))
    print(f"After: permissions={after[0]} menus={after[1]} menuPermissions={after[2]} rolePermissions={after[3]}")


def main() -> int:
    args = parse_args()
    manifest_path = (ROOT / args.manifest).resolve()
    if not manifest_path.is_file():
        raise PlanError(f"manifest not found: {args.manifest}")
    if not manifest_path.is_relative_to(ROOT):
        raise PlanError("manifest must be inside repository")

    registry = load_script("module_registry_plan", REGISTRY_PLAN)
    manifest = registry.load_manifest(manifest_path)
    if args.apply:
        validate_apply_gate(args, manifest)
        apply_uninstall_sql(args, registry, manifest)
        return 0
    print("-- Dry-run uninstall SQL preview. Review manually before applying; this script did not execute it.")
    print(f"-- Module: {manifest['module']}")
    print(build_sql(registry, manifest))
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except PlanError as exc:
        print(f"FAIL: {exc}")
        raise SystemExit(1)
