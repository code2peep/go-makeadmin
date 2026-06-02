#!/usr/bin/env python3
"""Generate a dry-run role grant SQL preview for module manifests."""

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
    parser = argparse.ArgumentParser(description="Preview module role permission grant SQL.")
    parser.add_argument(
        "--manifest",
        default="examples/demo/manifest.json",
        help="module manifest path, default: examples/demo/manifest.json",
    )
    parser.add_argument("--tenant-id", type=non_negative_int, default=0)
    parser.add_argument("--role-id", type=positive_int, required=True)
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


def build_sql(manifest: dict[str, Any], tenant_id: int, role_id: int) -> str:
    codes = [permission["code"] for permission in manifest["permissions"]]
    code_list = ", ".join(sql_quote(code) for code in codes)
    statements = [
        "SET @now = UNIX_TIMESTAMP();",
        f"SET @tenant_id = {tenant_id};",
        f"SET @role_id = {role_id};",
        "\n".join(
            [
                "INSERT INTO `ma_role_permission`",
                "(`tenant_id`, `role_id`, `permission_id`, `create_time`)",
                "SELECT @tenant_id, @role_id, p.id, @now",
                "FROM `ma_permission` AS p",
                f"WHERE p.code IN ({code_list})",
                "  AND p.status = 1",
                "  AND EXISTS (",
                "      SELECT 1 FROM `ma_role` AS r",
                "      WHERE r.tenant_id = @tenant_id",
                "        AND r.id = @role_id",
                "        AND r.status = 1",
                "        AND r.delete_time = 0",
                "  )",
                "  AND NOT EXISTS (",
                "      SELECT 1 FROM `ma_role_permission` AS rp",
                "      WHERE rp.tenant_id = @tenant_id",
                "        AND rp.role_id = @role_id",
                "        AND rp.permission_id = p.id",
                "  );",
            ]
        ),
    ]
    return "\n\n".join(statements)


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
    print("-- Dry-run role grant SQL preview. Review manually before applying; this script did not execute it.")
    print(f"-- Module: {manifest['module']}")
    print(f"-- Tenant ID: {args.tenant_id}")
    print(f"-- Role ID: {args.role_id}")
    print(build_sql(manifest, args.tenant_id, args.role_id))
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except PlanError as exc:
        print(f"FAIL: {exc}")
        raise SystemExit(1)
