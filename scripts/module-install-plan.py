#!/usr/bin/env python3
"""Generate a complete dry-run install plan for a module manifest."""

from __future__ import annotations

import argparse
import importlib.util
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
REGISTRY_PLAN = ROOT / "scripts" / "module-registry-plan.py"
ROLE_GRANT_PLAN = ROOT / "scripts" / "module-role-grant-plan.py"
DEMO_RUNTIME_ENV = "MAKEADMIN_ENABLE_DEMO_MODULE"


class PlanError(RuntimeError):
    pass


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Preview a full module install plan without database writes.")
    parser.add_argument(
        "--manifest",
        default="examples/demo/manifest.json",
        help="module manifest path, default: examples/demo/manifest.json",
    )
    parser.add_argument("--tenant-id", type=non_negative_int, default=0)
    parser.add_argument("--role-id", type=positive_int, help="include role grant SQL for this role id")
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
