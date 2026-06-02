#!/usr/bin/env python3
"""Preview backend codegen table and column configuration from a module manifest."""

from __future__ import annotations

import argparse
import importlib.util
import json
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
VALIDATOR = ROOT / "scripts" / "check-module-manifests.py"
DEFAULT_OPTIONS = {
    "treePrimary": "",
    "treeParent": "",
    "treeName": "",
    "subTableName": "",
    "subTableFk": "",
}
HTML_DEFAULTS = {
    "input": {"columnType": "varchar", "columnLength": 255, "goType": "string", "queryType": "LIKE"},
    "number": {"columnType": "int", "columnLength": 0, "goType": "int", "queryType": "="},
    "textarea": {"columnType": "text", "columnLength": 0, "goType": "string", "queryType": "LIKE"},
    "select": {"columnType": "varchar", "columnLength": 100, "goType": "string", "queryType": "="},
    "radio": {"columnType": "varchar", "columnLength": 100, "goType": "string", "queryType": "="},
    "checkbox": {"columnType": "varchar", "columnLength": 255, "goType": "string", "queryType": "="},
    "datetime": {"columnType": "datetime", "columnLength": 0, "goType": "time.Time", "queryType": "="},
}


class PlanError(RuntimeError):
    pass


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Preview ma_codegen_* and legacy /gen config from a module manifest.")
    parser.add_argument(
        "--manifest",
        default="examples/demo/manifest.json",
        help="module manifest path, default: examples/demo/manifest.json",
    )
    parser.add_argument("--tenant-id", type=non_negative_int, default=0)
    parser.add_argument("--author", default="codepeep")
    parser.add_argument("--package-name", default="gencode")
    parser.add_argument("--format", choices=("summary", "json"), default="summary")
    return parser.parse_args()


def non_negative_int(value: str) -> int:
    parsed = int(value)
    if parsed < 0:
        raise argparse.ArgumentTypeError("must be non-negative")
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


def display_path(path: Path) -> str:
    try:
        return str(path.relative_to(ROOT))
    except ValueError:
        return str(path)


def build_plan(manifest_path: Path, manifest: dict[str, Any], args: argparse.Namespace) -> dict[str, Any]:
    module = manifest["module"]
    function_name = manifest["menu"]["name"]
    options_json = json.dumps(DEFAULT_OPTIONS, separators=(",", ":"), ensure_ascii=False)
    source = display_path(manifest_path)
    legacy_table = {
        "id": 0,
        "tableName": manifest["table"],
        "tableComment": function_name,
        "subTableName": "",
        "subTableFk": "",
        "authorName": args.author,
        "entityName": manifest["entity"],
        "moduleName": module,
        "functionName": function_name,
        "treePrimary": "",
        "treeParent": "",
        "treeName": "",
        "genTpl": "crud",
        "genType": 0,
        "genPath": "/",
        "remarks": f"generated from {source}",
    }
    legacy_columns = build_legacy_columns(manifest)
    return {
        "source": source,
        "module": module,
        "warning": "Dry-run preview only; this command did not connect to a database and did not write ma_codegen_* rows.",
        "makeadmin": {
            "table": {
                "id": 0,
                "tenantId": args.tenant_id,
                "sourceTable": manifest["table"],
                "tableComment": function_name,
                "moduleName": module,
                "packageName": args.package_name,
                "businessName": module,
                "entityName": manifest["entity"],
                "functionName": function_name,
                "authorName": args.author,
                "templateType": "crud",
                "generateType": "zip",
                "generatePath": "/",
                "options": options_json,
                "remark": f"generated from {source}",
                "deleteTime": 0,
            },
            "columns": [legacy_to_makeadmin_column(column) for column in legacy_columns],
        },
        "legacy": {
            "genTable": legacy_table,
            "genTableColumns": legacy_columns,
        },
        "next": {
            "validateManifest": "python3 scripts/check-module-manifests.py",
            "installPlan": f"python3 scripts/module-install-plan.py --manifest {source} --tenant-id {args.tenant_id} --role-id 1",
            "codegenLink": f"scripts/check-module-codegen.sh --manifest {source}",
        },
    }


def build_legacy_columns(manifest: dict[str, Any]) -> list[dict[str, Any]]:
    codegen = manifest.get("codegen") or {}
    if not isinstance(codegen, dict):
        raise PlanError("codegen must be an object")
    configured_columns = codegen.get("columns", [])
    if not configured_columns:
        return [
            primary_key_column(),
            {
                "id": 0,
                "tableId": 0,
                "columnName": "title",
                "columnComment": "Title",
                "columnLength": 200,
                "columnType": "varchar",
                "goType": "string",
                "goField": "title",
                "isPk": 0,
                "isIncrement": 0,
                "isRequired": 1,
                "isInsert": 1,
                "isEdit": 1,
                "isList": 1,
                "isQuery": 1,
                "queryType": "LIKE",
                "htmlType": "input",
                "dictType": "",
                "sort": 2,
            },
            {
                "id": 0,
                "tableId": 0,
                "columnName": "status",
                "columnComment": "Status",
                "columnLength": 0,
                "columnType": "tinyint",
                "goType": "int",
                "goField": "status",
                "isPk": 0,
                "isIncrement": 0,
                "isRequired": 0,
                "isInsert": 1,
                "isEdit": 1,
                "isList": 1,
                "isQuery": 1,
                "queryType": "=",
                "htmlType": "input",
                "dictType": "",
                "sort": 3,
            },
        ]

    if not isinstance(configured_columns, list):
        raise PlanError("codegen.columns must be a list")

    columns = [primary_key_column()]
    for index, item in enumerate(configured_columns, start=2):
        columns.append(configured_column(item, index))
    return columns


def primary_key_column() -> dict[str, Any]:
    return {
        "id": 0,
        "tableId": 0,
        "columnName": "id",
        "columnComment": "ID",
        "columnLength": 0,
        "columnType": "bigint",
        "goType": "uint",
        "goField": "id",
        "isPk": 1,
        "isIncrement": 1,
        "isRequired": 0,
        "isInsert": 0,
        "isEdit": 0,
        "isList": 1,
        "isQuery": 0,
        "queryType": "=",
        "htmlType": "input",
        "dictType": "",
        "sort": 1,
    }


def configured_column(item: Any, sort: int) -> dict[str, Any]:
    if not isinstance(item, dict):
        raise PlanError("codegen.columns entries must be objects")
    column_name = require_text(item, "columnName")
    go_field = require_text(item, "goField")
    html_type = require_text(item, "htmlType")
    defaults = HTML_DEFAULTS.get(html_type)
    if defaults is None:
        raise PlanError(f"unsupported codegen column htmlType: {html_type}")
    dict_type = item.get("dictType", "")
    if not isinstance(dict_type, str):
        raise PlanError("codegen column dictType must be a string")

    return {
        "id": 0,
        "tableId": 0,
        "columnName": column_name,
        "columnComment": item.get("columnComment") or column_name.replace("_", " ").title(),
        "columnLength": int(item.get("columnLength", defaults["columnLength"])),
        "columnType": item.get("columnType") or defaults["columnType"],
        "goType": item.get("goType") or defaults["goType"],
        "goField": go_field,
        "isPk": 0,
        "isIncrement": 0,
        "isRequired": int(item.get("isRequired", 0)),
        "isInsert": int(item.get("isInsert", 1)),
        "isEdit": int(item.get("isEdit", 1)),
        "isList": int(item.get("isList", 1)),
        "isQuery": int(item.get("isQuery", 1)),
        "queryType": item.get("queryType") or defaults["queryType"],
        "htmlType": html_type,
        "dictType": dict_type,
        "sort": int(item.get("sort", sort)),
    }


def require_text(data: dict[str, Any], key: str) -> str:
    value = data.get(key)
    if not isinstance(value, str) or not value.strip():
        raise PlanError(f"{key} must be a non-empty string")
    return value


def legacy_to_makeadmin_column(column: dict[str, Any]) -> dict[str, Any]:
    return {
        "id": column["id"],
        "tableId": column["tableId"],
        "columnName": column["columnName"],
        "columnComment": column["columnComment"],
        "columnType": column["columnType"],
        "columnLength": column["columnLength"],
        "goType": column["goType"],
        "goField": column["goField"],
        "jsonField": column["goField"],
        "isPk": column["isPk"],
        "isIncrement": column["isIncrement"],
        "isRequired": column["isRequired"],
        "isInsert": column["isInsert"],
        "isEdit": column["isEdit"],
        "isList": column["isList"],
        "isQuery": column["isQuery"],
        "queryType": column["queryType"],
        "htmlType": column["htmlType"],
        "dictType": column["dictType"],
        "sort": column["sort"],
    }


def print_summary(plan: dict[str, Any]) -> None:
    table = plan["makeadmin"]["table"]
    columns = plan["makeadmin"]["columns"]
    print("# Module codegen dry-run plan")
    print(plan["warning"])
    print()
    print("## Codegen Table")
    print(f"Source table: {table['sourceTable']}")
    print(f"Entity: {table['entityName']}")
    print(f"Module: {table['moduleName']}")
    print(f"Function: {table['functionName']}")
    print(f"Template: {table['templateType']}")
    print(f"Generate type: {table['generateType']}")
    print()
    print("## Codegen Columns")
    for column in columns:
        print(
            f"- {column['columnName']} goType={column['goType']} goField={column['goField']} "
            f"list={column['isList']} query={column['isQuery']} html={column['htmlType']}"
        )
    print()
    print("## JSON Preview")
    print(json.dumps(plan, ensure_ascii=False, indent=2))


def main() -> int:
    args = parse_args()
    manifest_path = Path(args.manifest)
    if not manifest_path.is_absolute():
        manifest_path = ROOT / manifest_path
    manifest_path = manifest_path.resolve()
    if not manifest_path.is_file():
        raise PlanError(f"manifest not found: {args.manifest}")

    manifest = load_manifest(manifest_path)
    plan = build_plan(manifest_path, manifest, args)
    if args.format == "json":
        print(json.dumps(plan, ensure_ascii=False, indent=2))
    else:
        print_summary(plan)
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except PlanError as exc:
        print(f"FAIL: {exc}")
        raise SystemExit(1)
