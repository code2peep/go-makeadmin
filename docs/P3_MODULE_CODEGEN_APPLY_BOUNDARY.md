# P3 Module Codegen Apply Boundary

更新时间：2026-06-02

## 目标

P3.5 定义模块 manifest 写入后台生成器配置的受控边界，为后续把 `module-codegen-plan.py` 从 dry-run 预览升级为本地可控写入做准备。

本阶段只开放 `--apply` 门禁检查，不执行数据库写入，不创建业务 schema，不修改旧 `/gen/*` 接口形状。

当前更新：P3.6 已按该边界开放本地受控写入模式；现行写入规则见 `docs/P3_MODULE_CODEGEN_APPLY.md`。

## 写入门禁

未来 `module-codegen-plan.py --apply` 必须同时满足：

- 环境变量 `MAKEADMIN_ALLOW_MODULE_CODEGEN_WRITE=1`。
- 显式传入 `--confirm-module <module>`。
- `--confirm-module` 必须等于 manifest 中的 `module`。
- 显式传入 `--confirm-source-table <table>`。
- `--confirm-source-table` 必须等于 manifest 中的 `table`。
- 显式传入 `--confirm-sync-columns`。

缺少环境变量或确认参数时，脚本必须在数据库访问前失败。

P3.5 阶段即使满足全部门禁，也会在数据库访问前失败。P3.6 起，该命令在本地开发库受控写入：

```bash
MAKEADMIN_ALLOW_MODULE_CODEGEN_WRITE=1 \
python3 scripts/module-codegen-plan.py \
  --manifest examples/demo/manifest.json \
  --tenant-id 0 \
  --apply \
  --confirm-module article \
  --confirm-source-table ma_demo_article \
  --confirm-sync-columns
```

## 数据库范围

未来写入默认只允许本地开发库：

- host：`127.0.0.1`
- port：`3306`
- user：`root`
- database：`go_makeadmin`

可通过命令参数或 `MYSQL_HOST`、`MYSQL_PORT`、`MYSQL_USER`、`MYSQL_DATABASE`、`MYSQL_PASSWORD` 环境变量覆盖。

执行器不得读取 `.env` 猜测数据库密码，不得连接业务项目数据库。

## 写入顺序

未来 apply 必须在单事务内完成：

1. 校验 manifest。
2. 生成 codegen plan。
3. 按 `tenant_id + table_name + delete_time=0` 查找 live `ma_codegen_table`。
4. 如果不存在，插入 `ma_codegen_table`。
5. 如果已存在，只有 `module_name`、`business_name`、`entity_name` 与 manifest 对齐时才允许更新生成器配置字段。
6. 写入或更新 manifest 期望的 `ma_codegen_column`。
7. 在 `--confirm-sync-columns` 存在时，删除同一 `table_id` 下不再由 manifest 声明的列配置。

事务内任一步失败，全部回滚。

## 幂等规则

- 表配置按 `tenant_id + table_name + delete_time=0` 幂等。
- 列配置按 `table_id + column_name` 幂等。
- 重复执行同一 manifest 不应增加重复行。
- 已存在但归属不一致的 live 表配置必须停止，不覆盖。
- `ma_codegen_table` 删除仍走 `delete_time` 软删。
- `ma_codegen_column` 没有软删字段；删除 stale 列必须依赖 `--confirm-sync-columns`，且只限同一 codegen table。

## schema 边界

- 不创建、修改或迁移业务 schema。
- 不扫描业务数据库推断字段。
- `requiresSchema=true` 只表示模块业务表需要另行准备；codegen 配置写入不会代替建表。
- 树表和子表预留字段继续放入 `ma_codegen_table.options` JSON，不新增 SQL 列。

## 验证

P3.5 需要通过：

```bash
scripts/check-module-codegen-apply-boundary.sh
scripts/check-module-tools-no-db.sh
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
