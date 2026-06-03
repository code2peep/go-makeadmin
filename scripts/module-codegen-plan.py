#!/usr/bin/env python3
"""Preview backend codegen table and column configuration from a module manifest."""

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
WRITE_ENV = "MAKEADMIN_ALLOW_MODULE_CODEGEN_WRITE"
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
    parser.add_argument("--apply", action="store_true", help=f"reserved codegen config write mode; requires {WRITE_ENV}=1")
    parser.add_argument("--confirm-module", help="required with --apply; must equal manifest module")
    parser.add_argument("--confirm-source-table", help="required with --apply; must equal manifest table")
    parser.add_argument("--confirm-sync-columns", action="store_true", help="required with --apply")
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
            "codegenApplyBoundary": (
                f"{WRITE_ENV}=1 python3 scripts/module-codegen-plan.py --manifest {source} "
                f"--tenant-id {args.tenant_id} --apply --confirm-module {module} "
                f"--confirm-source-table {manifest['table']} --confirm-sync-columns"
            ),
        },
    }


def validate_apply_gate(args: argparse.Namespace, manifest_path: Path, manifest: dict[str, Any]) -> None:
    if not manifest_path.is_relative_to(ROOT):
        raise PlanError("--apply requires manifest inside repository; no database access was attempted")
    if os.environ.get(WRITE_ENV) != "1":
        raise PlanError(f"--apply requires {WRITE_ENV}=1; no database access was attempted")
    if not args.confirm_module:
        raise PlanError("--apply requires --confirm-module; no database access was attempted")
    if args.confirm_module != manifest["module"]:
        raise PlanError("--confirm-module must equal manifest module; no database access was attempted")
    if not args.confirm_source_table:
        raise PlanError("--apply requires --confirm-source-table; no database access was attempted")
    if args.confirm_source_table != manifest["table"]:
        raise PlanError("--confirm-source-table must equal manifest table; no database access was attempted")
    if not args.confirm_sync_columns:
        raise PlanError("--apply requires --confirm-sync-columns; no database access was attempted")


def mysql_query(args: argparse.Namespace, query: str) -> str:
    command = mysql_command(args) + ["--batch", "--raw", "--skip-column-names", "--execute", query]
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


def mysql_exec(args: argparse.Namespace, query: str) -> None:
    command = mysql_command(args) + ["--execute", query]
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


def mysql_command(args: argparse.Namespace) -> list[str]:
    return [
        "mysql",
        "--host",
        args.mysql_host,
        "--port",
        str(args.mysql_port),
        "--user",
        args.mysql_user,
        "--database",
        args.mysql_database,
    ]


def apply_codegen_config(args: argparse.Namespace, plan: dict[str, Any]) -> None:
    assert_live_table_owner(args, plan)
    mysql_exec(args, build_apply_sql(plan))
    snapshot = apply_snapshot(args, plan)
    table = plan["makeadmin"]["table"]
    print(f"Module codegen apply: module={plan['module']} sourceTable={table['sourceTable']}")
    print(f"Table: id={snapshot['tableId']} tenantId={table['tenantId']}")
    print(f"Columns: count={snapshot['columnCount']} names={snapshot['columnNames']}")


def assert_live_table_owner(args: argparse.Namespace, plan: dict[str, Any]) -> None:
    table = plan["makeadmin"]["table"]
    raw = mysql_query(
        args,
        "\n".join(
            [
                "SELECT `id`, `module_name`, `business_name`, `entity_name`",
                "FROM `ma_codegen_table`",
                f"WHERE `tenant_id` = {int(table['tenantId'])}",
                f"  AND `table_name` = {sql_quote(table['sourceTable'])}",
                "  AND `delete_time` = 0",
                "LIMIT 2;",
            ]
        ),
    )
    rows = [row for row in raw.splitlines() if row.strip()]
    if not rows:
        return
    if len(rows) > 1:
        raise PlanError("multiple live codegen table rows found; no database writes were executed")
    parts = rows[0].split("\t")
    if len(parts) != 4:
        raise PlanError(f"unexpected codegen table owner output: {rows[0]}")
    _, module_name, business_name, entity_name = parts
    if (
        module_name != table["moduleName"]
        or business_name != table["businessName"]
        or entity_name != table["entityName"]
    ):
        raise PlanError("live codegen table belongs to another module or entity; no database writes were executed")


def build_apply_sql(plan: dict[str, Any]) -> str:
    table = plan["makeadmin"]["table"]
    columns = plan["makeadmin"]["columns"]
    expected_column_names = ", ".join(sql_quote(column["columnName"]) for column in columns)
    statements = [
        "START TRANSACTION;",
        "SET @now = UNIX_TIMESTAMP();",
        f"SET @tenant_id = {int(table['tenantId'])};",
        f"SET @source_table = {sql_quote(table['sourceTable'])};",
        table_insert_sql(table),
        table_update_sql(table),
    ]
    statements.extend(column_upsert_sql(column) for column in columns)
    statements.append(
        "\n".join(
            [
                "DELETE FROM `ma_codegen_column`",
                "WHERE `table_id` = @codegen_table_id",
                f"  AND `column_name` NOT IN ({expected_column_names});",
            ]
        )
    )
    statements.append("COMMIT;")
    return "\n\n".join(statements)


def table_insert_sql(table: dict[str, Any]) -> str:
    return "\n".join(
        [
            "SET @codegen_table_id = COALESCE((",
            "    SELECT `id` FROM `ma_codegen_table`",
            "    WHERE `tenant_id` = @tenant_id AND `table_name` = @source_table AND `delete_time` = 0",
            "    LIMIT 1",
            "), 0);",
            "",
            "INSERT INTO `ma_codegen_table`",
            "(`tenant_id`, `table_name`, `table_comment`, `module_name`, `package_name`, `business_name`,",
            " `entity_name`, `function_name`, `author_name`, `template_type`, `generate_type`, `generate_path`,",
            " `options`, `remark`, `create_time`, `update_time`, `delete_time`)",
            "SELECT",
            f"@tenant_id, @source_table, {sql_quote(table['tableComment'])}, {sql_quote(table['moduleName'])},",
            f"{sql_quote(table['packageName'])}, {sql_quote(table['businessName'])}, {sql_quote(table['entityName'])},",
            f"{sql_quote(table['functionName'])}, {sql_quote(table['authorName'])}, {sql_quote(table['templateType'])},",
            f"{sql_quote(table['generateType'])}, {sql_quote(table['generatePath'])}, {sql_quote(table['options'])},",
            f"{sql_quote(table['remark'])}, @now, @now, 0",
            "WHERE @codegen_table_id = 0;",
            "",
            "SET @codegen_table_id = COALESCE((",
            "    SELECT `id` FROM `ma_codegen_table`",
            "    WHERE `tenant_id` = @tenant_id AND `table_name` = @source_table AND `delete_time` = 0",
            "    LIMIT 1",
            "), 0);",
        ]
    )


def table_update_sql(table: dict[str, Any]) -> str:
    return "\n".join(
        [
            "UPDATE `ma_codegen_table`",
            "SET",
            f"  `table_comment` = {sql_quote(table['tableComment'])},",
            f"  `module_name` = {sql_quote(table['moduleName'])},",
            f"  `package_name` = {sql_quote(table['packageName'])},",
            f"  `business_name` = {sql_quote(table['businessName'])},",
            f"  `entity_name` = {sql_quote(table['entityName'])},",
            f"  `function_name` = {sql_quote(table['functionName'])},",
            f"  `author_name` = {sql_quote(table['authorName'])},",
            f"  `template_type` = {sql_quote(table['templateType'])},",
            f"  `generate_type` = {sql_quote(table['generateType'])},",
            f"  `generate_path` = {sql_quote(table['generatePath'])},",
            f"  `options` = {sql_quote(table['options'])},",
            f"  `remark` = {sql_quote(table['remark'])},",
            "  `update_time` = @now",
            "WHERE `id` = @codegen_table_id;",
        ]
    )


def column_upsert_sql(column: dict[str, Any]) -> str:
    return "\n".join(
        [
            "INSERT INTO `ma_codegen_column`",
            "(`table_id`, `column_name`, `column_comment`, `column_type`, `column_length`, `go_type`, `go_field`,",
            " `json_field`, `is_pk`, `is_increment`, `is_required`, `is_insert`, `is_edit`, `is_list`, `is_query`,",
            " `query_type`, `html_type`, `dict_type`, `sort`, `create_time`, `update_time`)",
            "VALUES",
            f"(@codegen_table_id, {sql_quote(column['columnName'])}, {sql_quote(column['columnComment'])},",
            f" {sql_quote(column['columnType'])}, {int(column['columnLength'])}, {sql_quote(column['goType'])},",
            f" {sql_quote(column['goField'])}, {sql_quote(column['jsonField'])}, {int(column['isPk'])},",
            f" {int(column['isIncrement'])}, {int(column['isRequired'])}, {int(column['isInsert'])},",
            f" {int(column['isEdit'])}, {int(column['isList'])}, {int(column['isQuery'])},",
            f" {sql_quote(column['queryType'])}, {sql_quote(column['htmlType'])}, {sql_quote(column['dictType'])},",
            f" {int(column['sort'])}, @now, @now)",
            "ON DUPLICATE KEY UPDATE",
            "  `column_comment` = VALUES(`column_comment`),",
            "  `column_type` = VALUES(`column_type`),",
            "  `column_length` = VALUES(`column_length`),",
            "  `go_type` = VALUES(`go_type`),",
            "  `go_field` = VALUES(`go_field`),",
            "  `json_field` = VALUES(`json_field`),",
            "  `is_pk` = VALUES(`is_pk`),",
            "  `is_increment` = VALUES(`is_increment`),",
            "  `is_required` = VALUES(`is_required`),",
            "  `is_insert` = VALUES(`is_insert`),",
            "  `is_edit` = VALUES(`is_edit`),",
            "  `is_list` = VALUES(`is_list`),",
            "  `is_query` = VALUES(`is_query`),",
            "  `query_type` = VALUES(`query_type`),",
            "  `html_type` = VALUES(`html_type`),",
            "  `dict_type` = VALUES(`dict_type`),",
            "  `sort` = VALUES(`sort`),",
            "  `update_time` = @now;",
        ]
    )


def apply_snapshot(args: argparse.Namespace, plan: dict[str, Any]) -> dict[str, Any]:
    table = plan["makeadmin"]["table"]
    raw = mysql_query(
        args,
        "\n".join(
            [
                "SELECT",
                "  COALESCE((",
                "      SELECT `id` FROM `ma_codegen_table`",
                f"      WHERE `tenant_id` = {int(table['tenantId'])}",
                f"        AND `table_name` = {sql_quote(table['sourceTable'])}",
                "        AND `delete_time` = 0",
                "      LIMIT 1",
                "  ), 0),",
                "  (",
                "      SELECT COUNT(*) FROM `ma_codegen_column` AS c",
                "      INNER JOIN `ma_codegen_table` AS t ON t.id = c.table_id",
                f"      WHERE t.`tenant_id` = {int(table['tenantId'])}",
                f"        AND t.`table_name` = {sql_quote(table['sourceTable'])}",
                "        AND t.`delete_time` = 0",
                "  ),",
                "  COALESCE((",
                "      SELECT GROUP_CONCAT(c.`column_name` ORDER BY c.`sort` SEPARATOR ',')",
                "      FROM `ma_codegen_column` AS c",
                "      INNER JOIN `ma_codegen_table` AS t ON t.id = c.table_id",
                f"      WHERE t.`tenant_id` = {int(table['tenantId'])}",
                f"        AND t.`table_name` = {sql_quote(table['sourceTable'])}",
                "        AND t.`delete_time` = 0",
                "  ), '');",
            ]
        ),
    )
    parts = raw.split("\t")
    if len(parts) != 3:
        raise PlanError(f"unexpected apply snapshot output: {raw}")
    return {"tableId": int(parts[0]), "columnCount": int(parts[1]), "columnNames": parts[2]}


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


def sql_quote(value: Any) -> str:
    text = str(value)
    return "'" + text.replace("\\", "\\\\").replace("'", "''") + "'"


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
    if args.apply:
        validate_apply_gate(args, manifest_path, manifest)
        apply_codegen_config(args, plan)
        return 0
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
