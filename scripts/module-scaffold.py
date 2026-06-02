#!/usr/bin/env python3
"""Scaffold a standard business module manifest and README."""

from __future__ import annotations

import argparse
import importlib.util
import json
import re
import tempfile
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
EXAMPLES = ROOT / "examples"
VALIDATOR = ROOT / "scripts" / "check-module-manifests.py"
REGISTRY_PLAN = ROOT / "scripts" / "module-registry-plan.py"
ROLE_GRANT_PLAN = ROOT / "scripts" / "module-role-grant-plan.py"
UNINSTALL_PLAN = ROOT / "scripts" / "module-uninstall-plan.py"
MODULE_RE = re.compile(r"^[a-z][a-z0-9_]*$")
ENTITY_RE = re.compile(r"^[A-Z][A-Za-z0-9]*$")
TABLE_RE = re.compile(r"^[a-z][a-z0-9_]*$")
ACTIONS: tuple[tuple[str, str, str], ...] = (
    ("list", "GET", "list"),
    ("detail", "GET", "detail"),
    ("add", "POST", "add"),
    ("edit", "POST", "edit"),
    ("del", "POST", "delete"),
)


class ScaffoldError(RuntimeError):
    pass


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Scaffold a business module manifest and README under examples/.")
    parser.add_argument("--module", required=True, help="module slug, for example: billing_invoice")
    parser.add_argument("--entity", help="Go entity name, default: PascalCase module")
    parser.add_argument("--table", help="database table name, default: ma_<module>")
    parser.add_argument("--backend-package", default="gencode")
    parser.add_argument("--menu-name", help="menu display name, default: title-cased module")
    parser.add_argument("--menu-parent", default="dev_tools")
    parser.add_argument("--menu-route-path", help="frontend route path, default: /<module>")
    parser.add_argument("--menu-route-name", help="frontend route name, default: <module>.index")
    parser.add_argument("--component", help="frontend component path, default: <module>/index")
    parser.add_argument("--requires-schema", action="store_true", help="mark manifest as requiring a business table")
    parser.add_argument("--runtime-registered", action="store_true", help="mark runtime routes as already registered")
    output = parser.add_mutually_exclusive_group()
    output.add_argument("--dry-run", action="store_true", help="print generated files instead of writing them")
    output.add_argument("--print-manifest", action="store_true", help="print generated manifest JSON only")
    return parser.parse_args()


def pascal_case(value: str) -> str:
    return "".join(part.capitalize() for part in value.split("_") if part)


def title_name(value: str) -> str:
    return " ".join(part.capitalize() for part in value.split("_") if part)


def validate_slug(value: str, field: str) -> None:
    if not MODULE_RE.match(value):
        raise ScaffoldError(f"{field} must match {MODULE_RE.pattern}")


def normalize_args(args: argparse.Namespace) -> dict[str, Any]:
    validate_slug(args.module, "--module")

    entity = args.entity or pascal_case(args.module)
    if not ENTITY_RE.match(entity):
        raise ScaffoldError("--entity must be PascalCase, for example: BillingInvoice")

    table = args.table or f"ma_{args.module}"
    if not TABLE_RE.match(table):
        raise ScaffoldError("--table must be a lower snake_case table name")

    menu_name = args.menu_name or title_name(args.module)
    route_path = args.menu_route_path or f"/{args.module}"
    if not route_path.startswith("/"):
        raise ScaffoldError("--menu-route-path must start with /")

    route_name = args.menu_route_name or f"{args.module}.index"
    component = args.component or f"{args.module}/index"
    if args.module not in route_name and args.module not in component:
        raise ScaffoldError("--menu-route-name or --component must include the module slug")

    return {
        "module": args.module,
        "entity": entity,
        "table": table,
        "backendPackage": args.backend_package,
        "menuName": menu_name,
        "menuParent": args.menu_parent,
        "menuRoutePath": route_path,
        "menuRouteName": route_name,
        "component": component,
        "requiresSchema": bool(args.requires_schema),
        "runtimeRegistered": bool(args.runtime_registered),
    }


def build_manifest(values: dict[str, Any]) -> dict[str, Any]:
    module = values["module"]
    permissions = [
        {
            "code": f"{module}:{action}",
            "name": f"{values['menuName']} {label}",
            "module": module,
            "resource": module,
            "action": action,
        }
        for action, _, label in ACTIONS
    ]
    return {
        "version": 1,
        "module": module,
        "entity": values["entity"],
        "table": values["table"],
        "backendPackage": values["backendPackage"],
        "backend": {
            "routes": [
                {
                    "method": method,
                    "path": f"/{module}/{action}",
                    "permission": f"{module}:{action}",
                }
                for action, method, _ in ACTIONS
            ]
        },
        "frontend": {
            "api": f"admin/src/api/{module}.ts",
            "views": [
                f"admin/src/views/{module}/index.vue",
                f"admin/src/views/{module}/edit.vue",
            ],
        },
        "menu": {
            "code": values["menuRouteName"],
            "parent": values["menuParent"],
            "type": "page",
            "name": values["menuName"],
            "routePath": values["menuRoutePath"],
            "routeName": values["menuRouteName"],
            "component": values["component"],
            "permission": f"{module}:list",
            "visible": False,
            "sort": 100,
        },
        "permissions": permissions,
        "runtimeRegistered": values["runtimeRegistered"],
        "requiresSchema": values["requiresSchema"],
    }


def build_readme(values: dict[str, Any]) -> str:
    module = values["module"]
    manifest_path = f"examples/{module}/manifest.json"
    lines = [
        f"# {values['menuName']} Module",
        "",
        "本目录是 go-makeadmin 标准业务模块脚手架输出，用来承载模块接入约定和后续生成结果。",
        "",
        "## 模块约定",
        "",
        f"- 模块名：`{module}`",
        f"- 实体名：`{values['entity']}`",
        f"- 表名：`{values['table']}`",
        f"- 后端包名：`{values['backendPackage']}`",
        f"- 菜单父级：`{values['menuParent']}`",
        f"- 菜单路由：`{values['menuRouteName']}`",
        f"- 前端 API：`admin/src/api/{module}.ts`",
        f"- 前端页面：`admin/src/views/{module}/`",
        f"- 需要 schema：`{str(values['requiresSchema']).lower()}`",
        "",
        "## 标准命令",
        "",
        "```bash",
        "python3 scripts/check-module-manifests.py",
        f"python3 scripts/module-install-plan.py --manifest {manifest_path} --tenant-id 0 --role-id 1",
        f"python3 scripts/module-uninstall-plan.py --manifest {manifest_path}",
        "```",
        "",
        "上述命令默认不连接数据库、不执行写入或删除。",
        "",
        "## 后续接入",
        "",
        "- 使用代码生成器生成 Go model、schema、service、route。",
        "- 使用前端模板生成 API、列表页和编辑页。",
        "- schema 由模块自己的 SQL 或迁移文档负责；模块安装脚本不会自动建表。",
        "- 安装和卸载写入必须继续使用 P2 已定义的显式环境变量门禁。",
        "",
    ]
    return "\n".join(lines)


def load_script(name: str, path: Path) -> Any:
    spec = importlib.util.spec_from_file_location(name, path)
    if spec is None or spec.loader is None:
        raise ScaffoldError(f"cannot load script: {path}")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def validate_manifest(manifest: dict[str, Any]) -> None:
    validator = load_script("module_manifest_validator", VALIDATOR)
    with tempfile.TemporaryDirectory() as temp_dir:
        manifest_path = Path(temp_dir) / "manifest.json"
        manifest_path.write_text(json.dumps(manifest, ensure_ascii=False, indent=2) + "\n")
        validator.validate_manifest(manifest_path)


def validate_lifecycle_builders(manifest: dict[str, Any]) -> None:
    registry = load_script("module_registry_plan", REGISTRY_PLAN)
    role_grant = load_script("module_role_grant_plan", ROLE_GRANT_PLAN)
    uninstall = load_script("module_uninstall_plan", UNINSTALL_PLAN)
    registry.build_sql(manifest)
    role_grant.build_sql(manifest, tenant_id=0, role_id=1)
    uninstall.build_sql(registry, manifest)


def write_files(module: str, manifest_text: str, readme: str) -> None:
    output_dir = EXAMPLES / module
    if output_dir.exists():
        raise ScaffoldError(f"module directory already exists: {output_dir.relative_to(ROOT)}")
    output_dir.mkdir(parents=False)
    (output_dir / "manifest.json").write_text(manifest_text)
    (output_dir / "README.md").write_text(readme)
    print(f"Created {output_dir.relative_to(ROOT)}/manifest.json")
    print(f"Created {output_dir.relative_to(ROOT)}/README.md")


def main() -> int:
    args = parse_args()
    values = normalize_args(args)
    manifest = build_manifest(values)
    validate_manifest(manifest)
    validate_lifecycle_builders(manifest)
    manifest_text = json.dumps(manifest, ensure_ascii=False, indent=2) + "\n"
    readme = build_readme(values)

    if args.print_manifest:
        print(manifest_text, end="")
        return 0

    if args.dry_run:
        module = values["module"]
        print(f"# Would write examples/{module}/manifest.json")
        print(manifest_text)
        print(f"# Would write examples/{module}/README.md")
        print(readme)
        return 0

    write_files(values["module"], manifest_text, readme)
    print()
    print("Next:")
    print("python3 scripts/check-module-manifests.py")
    print(f"python3 scripts/module-install-plan.py --manifest examples/{values['module']}/manifest.json --tenant-id 0 --role-id 1")
    print(f"python3 scripts/module-uninstall-plan.py --manifest examples/{values['module']}/manifest.json")
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except ScaffoldError as exc:
        print(f"FAIL: {exc}")
        raise SystemExit(1)
