# P3 Module Codegen Apply

更新时间：2026-06-03

## 目标

P3.6 开放模块生成器配置的本地受控写入，把 manifest 生成的 codegen plan 写入 `ma_codegen_table` 和 `ma_codegen_column`。

本阶段只写生成器配置，不创建业务 schema，不扫描业务数据库，不修改旧 `/gen/*` 请求和响应字段。

## 命令

dry-run 预览：

```bash
python3 scripts/module-codegen-plan.py --manifest examples/demo/manifest.json
```

受控写入：

```bash
MAKEADMIN_ALLOW_MODULE_CODEGEN_WRITE=1 \
python3 scripts/module-codegen-plan.py \
  --manifest examples/demo/manifest.json \
  --tenant-id 0 \
  --confirm-module article \
  --confirm-source-table ma_demo_article \
  --confirm-sync-columns \
  --apply
```

本地 smoke：

```bash
MAKEADMIN_ALLOW_MODULE_CODEGEN_SMOKE_WRITE=1 \
scripts/check-module-codegen-apply-smoke.sh
```

## 写入门禁

`--apply` 必须同时满足：

- 环境变量 `MAKEADMIN_ALLOW_MODULE_CODEGEN_WRITE=1`。
- manifest 必须在仓库内。
- 显式传入 `--confirm-module <module>`，且等于 manifest `module`。
- 显式传入 `--confirm-source-table <table>`，且等于 manifest `table`。
- 显式传入 `--confirm-sync-columns`。

缺少环境变量或确认参数时，脚本会在数据库访问前失败。

## 写入规则

- 写入在单个事务中执行。
- 表配置按 `tenant_id + table_name + delete_time=0` 查找 live 行。
- 不存在 live 行时插入 `ma_codegen_table`。
- 存在 live 行时，只有 `module_name`、`business_name`、`entity_name` 与 manifest 对齐才允许更新。
- 列配置按 `table_id + column_name` 幂等 upsert。
- `--confirm-sync-columns` 会删除同一 `table_id` 下 manifest 不再声明的 stale 列配置。
- 旧树表和子表字段继续写入 `ma_codegen_table.options` JSON。

## 本地 smoke

P3.6 已用本地 `go_makeadmin` 完成一次 codegen 写入 smoke：

- 写入前确认 `tenant_id=0 + ma_demo_article + delete_time=0` 不存在 live codegen 表配置。
- 第一次 apply 写入 1 条 `ma_codegen_table` 和 3 条 `ma_codegen_column`。
- 手工插入 1 条 `stale_column`。
- 第二次 apply 后列配置回到 `id,title,status`，确认 stale 列同步删除。
- smoke 结束后删除本次列配置，并软删本次 live 表配置。
- 清理后 live codegen 表配置残留为 0。

## 边界

- 不创建、修改或迁移业务 schema。
- 不读取或修改 `.env`。
- 不连接业务项目数据库。
- 不删除文件或目录。
- 不修改 runtime 环境变量。
- 不改变旧 `/gen/*` 兼容字段形状。

## 验收

P3.6 需要通过：

```bash
python3 -m py_compile scripts/module-codegen-plan.py
bash -n scripts/check-module-codegen-apply-boundary.sh scripts/check-module-codegen-apply-smoke.sh scripts/check-module-tools-no-db.sh
scripts/check-module-codegen-apply-boundary.sh
MAKEADMIN_ALLOW_MODULE_CODEGEN_SMOKE_WRITE=1 scripts/check-module-codegen-apply-smoke.sh
scripts/check-module-tools-no-db.sh
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
