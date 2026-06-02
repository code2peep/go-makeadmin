#!/usr/bin/env python3
"""Generate a dry-run uninstall SQL preview for module manifests."""

from __future__ import annotations

import argparse
import importlib.util
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
REGISTRY_PLAN = ROOT / "scripts" / "module-registry-plan.py"


class PlanError(RuntimeError):
    pass


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Preview module registry cleanup SQL.")
    parser.add_argument(
        "--manifest",
        default="examples/demo/manifest.json",
        help="module manifest path, default: examples/demo/manifest.json",
    )
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


def main() -> int:
    args = parse_args()
    manifest_path = (ROOT / args.manifest).resolve()
    if not manifest_path.is_file():
        raise PlanError(f"manifest not found: {args.manifest}")
    if not manifest_path.is_relative_to(ROOT):
        raise PlanError("manifest must be inside repository")

    registry = load_script("module_registry_plan", REGISTRY_PLAN)
    manifest = registry.load_manifest(manifest_path)
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
