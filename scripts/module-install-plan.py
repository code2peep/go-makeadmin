#!/usr/bin/env python3
"""Generate a complete dry-run install plan for a module manifest."""

from __future__ import annotations

import argparse
import importlib.util
import os
import subprocess
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
REGISTRY_PLAN = ROOT / "scripts" / "module-registry-plan.py"
ROLE_GRANT_PLAN = ROOT / "scripts" / "module-role-grant-plan.py"
DEMO_RUNTIME_ENV = "MAKEADMIN_ENABLE_DEMO_MODULE"
WRITE_ENV = "MAKEADMIN_ALLOW_MODULE_INSTALL_WRITE"


class PlanError(RuntimeError):
    pass


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Preview or apply a full module install plan.")
    parser.add_argument(
        "--manifest",
        default="examples/demo/manifest.json",
        help="module manifest path, default: examples/demo/manifest.json",
    )
    parser.add_argument("--tenant-id", type=non_negative_int, default=0)
    parser.add_argument("--role-id", type=positive_int, help="include role grant SQL for this role id")
    parser.add_argument("--apply", action="store_true", help=f"write module registry rows; requires {WRITE_ENV}=1")
    parser.add_argument("--confirm-module", help="required with --apply; must equal manifest module")
    parser.add_argument("--confirm-role-id", type=positive_int, help="required with --apply when --role-id is set")
    parser.add_argument("--mysql-host", default=os.environ.get("MYSQL_HOST", "127.0.0.1"))
    parser.add_argument("--mysql-port", default=os.environ.get("MYSQL_PORT", "3306"))
    parser.add_argument("--mysql-user", default=os.environ.get("MYSQL_USER", "root"))
    parser.add_argument("--mysql-database", default=os.environ.get("MYSQL_DATABASE", "go_makeadmin"))
    return parser.parse_args()


def non_negative_int(value: str) -> int:
    parsed = int(value)
    if parsed < 0:
        raise argparse.ArgumentTypeError("must be non-negative")
    return parsed


def positive_int(value: str) -> int:
    parsed = int(value)
    if parsed <= 0:
        raise argparse.ArgumentTypeError("must be positive")
    return parsed


def load_script(name: str, path: Path) -> Any:
    spec = importlib.util.spec_from_file_location(name, path)
    if spec is None or spec.loader is None:
        raise PlanError(f"cannot load script: {path}")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def runtime_hint(manifest: dict[str, Any]) -> str:
    if manifest["module"] == "article":
        return f"{DEMO_RUNTIME_ENV}=1"
    return "No runtime env gate is defined for this module yet."


def build_install_sql(registry: Any, role_grant: Any, manifest: dict[str, Any], tenant_id: int, role_id: int | None) -> str:
    statements = [registry.build_sql(manifest)]
    if role_id is not None:
        statements.append(role_grant.build_sql(manifest, tenant_id, role_id))
    return "\n\n".join(statements)


def validate_apply_gate(args: argparse.Namespace, manifest: dict[str, Any]) -> None:
    if os.environ.get(WRITE_ENV) != "1":
        raise PlanError(f"--apply requires {WRITE_ENV}=1; no database access was attempted")
    if not args.confirm_module:
        raise PlanError("--apply requires --confirm-module; no database access was attempted")
    if args.confirm_module != manifest["module"]:
        raise PlanError("--confirm-module must equal manifest module; no database access was attempted")
    if manifest["requiresSchema"]:
        raise PlanError("--apply does not create module schema; no database access was attempted")
    if args.role_id is not None and args.confirm_role_id != args.role_id:
        raise PlanError("--apply with --role-id requires matching --confirm-role-id; no database access was attempted")


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


def apply_install_sql(args: argparse.Namespace, install_sql: str, manifest: dict[str, Any]) -> None:
    mysql_exec(args, f"START TRANSACTION;\n\n{install_sql}\n\nCOMMIT;")
    print(f"Module install apply: module={manifest['module']}")
    print(f"Permissions: planned={len(manifest['permissions'])}")
    print("Menu: planned=1")
    print("Menu permissions: planned=1")
    if args.role_id is None:
        print("Role permissions: skipped")
    else:
        print(f"Role permissions: planned={len(manifest['permissions'])}")


def print_manifest_summary(manifest_path: Path, manifest: dict[str, Any]) -> None:
    print("## Manifest")
    print(f"Path: {manifest_path.relative_to(ROOT)}")
    print(f"Module: {manifest['module']}")
    print(f"Entity: {manifest['entity']}")
    print(f"Table: {manifest['table']}")
    print(f"Requires schema: {str(manifest['requiresSchema']).lower()}")
    print(f"Runtime registered by default: {str(manifest['runtimeRegistered']).lower()}")
    print()
    print("## Backend Routes")
    for route in manifest["backend"]["routes"]:
        print(f"- {route['method']} {route['path']} permission={route['permission']}")
    print()
    print("## Frontend")
    print(f"API: {manifest['frontend']['api']}")
    for view in manifest["frontend"]["views"]:
        print(f"View: {view}")
    print()


def main() -> int:
    args = parse_args()
    manifest_path = (ROOT / args.manifest).resolve()
    if not manifest_path.is_file():
        raise PlanError(f"manifest not found: {args.manifest}")
    if not manifest_path.is_relative_to(ROOT):
        raise PlanError("manifest must be inside repository")

    registry = load_script("module_registry_plan", REGISTRY_PLAN)
    role_grant = load_script("module_role_grant_plan", ROLE_GRANT_PLAN)
    manifest = registry.load_manifest(manifest_path)
    install_sql = build_install_sql(registry, role_grant, manifest, args.tenant_id, args.role_id)

    if args.apply:
        validate_apply_gate(args, manifest)
        apply_install_sql(args, install_sql, manifest)
        return 0

    print("# Module install dry-run plan")
    print("This command did not connect to a database and did not execute SQL.")
    print()
    print_manifest_summary(manifest_path, manifest)

    print("## Registry SQL")
    print(registry.build_sql(manifest))
    print()

    print("## Role Grant SQL")
    if args.role_id is None:
        print("No role grant SQL generated. Pass --role-id to include ma_role_permission grants.")
    else:
        print(role_grant.build_sql(manifest, args.tenant_id, args.role_id))
    print()

    print("## Runtime")
    print(runtime_hint(manifest))
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except PlanError as exc:
        print(f"FAIL: {exc}")
        raise SystemExit(1)
