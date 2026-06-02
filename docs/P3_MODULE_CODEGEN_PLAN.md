# P3 Module Codegen Plan

更新时间：2026-06-02

## 目标

P3.4 将模块 manifest 转换为后台生成器配置预览，让脚手架产物能对接 `ma_codegen_table`、`ma_codegen_column` 和旧 `/gen/*` 兼容字段。

本阶段只做 dry-run 预览，不写数据库，不创建业务 schema，不修改旧 `/gen/*` 接口形状。

## 命令

生成可读预览：

```bash
python3 scripts/module-codegen-plan.py --manifest examples/demo/manifest.json
```

生成 JSON：

```bash
python3 scripts/module-codegen-plan.py \
  --manifest examples/demo/manifest.json \
  --tenant-id 0 \
  --format json
```

`--manifest` 可以传仓库内 manifest，也可以传脚手架临时输出的 manifest；命令只读取文件并生成预览。

使用脚手架输出校验：

```bash
scripts/check-module-codegen-plan.sh
```

使用脚手架实际写入产物校验：

```bash
MAKEADMIN_ALLOW_MODULE_SCAFFOLD_WRITE=1 scripts/check-module-scaffold-write-smoke.sh
```

该命令会把 smoke 产物写入 `.cache/module-scaffold-smoke/<timestamp>/examples/<module>/`，并使用生成的 `manifest.json` 生成 `codegen-plan.json`。

## 输出内容

预览包含两套形状：

- `makeadmin.table`：对应 `ma_codegen_table`。
- `makeadmin.columns`：对应 `ma_codegen_column`。
- `legacy.genTable`：对应旧 `GenTable` / `/gen/*` 表级字段。
- `legacy.genTableColumns`：对应旧 `GenTableColumn` / `/gen/*` 列级字段。

默认 CRUD 列配置：

- `id`：主键、自增、`uint`。
- `title`：必填、插入、编辑、列表、模糊查询。
- `status`：插入、编辑、列表、等值查询。

如果 manifest 声明了 `codegen.columns`，预览会使用 `id` 加 manifest 列配置，并保留 `htmlType`、`dictType`、`queryType` 等生成器字段。未显式提供的列类型会按 `htmlType` 使用保守默认值。

## 边界

- 不连接数据库。
- 不写入 `ma_codegen_table` 或 `ma_codegen_column`。
- 不创建业务 schema。
- 不读取或修改 `.env`。
- 不修改旧 `/gen/*` 请求/响应字段。
- 树表和子表预留字段继续放入 `options` JSON，不新增 SQL 列。
- 脚手架写入产物 smoke 只写 `.cache/` 下的临时模块文件，`.cache/` 不进入提交。
