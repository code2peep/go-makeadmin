# P3 Module Codegen Readback

更新时间：2026-06-03

## 目标

P3.7 验证模块生成器配置写入后的回读和模板生成闭环。

本阶段不新增生成器数据结构，不改变旧 `/gen/*` 响应形状，只验证 `ma_codegen_table` 和 `ma_codegen_column` 能被现有生成器服务读取，并驱动预览和下载模板生成。

## 命令

默认 Go 测试会跳过本地数据库 smoke：

```bash
cd server
go test ./generator/service/gen -run TestCodegenConfigReadbackAndTemplateGenerationSmoke -count=1
```

受控写入 smoke：

```bash
MAKEADMIN_ALLOW_MODULE_CODEGEN_READBACK_WRITE=1 \
scripts/check-module-codegen-readback-smoke.sh
```

## Smoke 流程

`scripts/check-module-codegen-readback-smoke.sh` 会执行以下步骤：

- 确认本地 `go_makeadmin` 中不存在 `tenant_id=0 + ma_demo_article + delete_time=0` live codegen 表配置。
- 复用 `scripts/module-codegen-plan.py --apply` 写入 demo manifest 对应的 `ma_codegen_table` 和 `ma_codegen_column`。
- 运行 `TestCodegenConfigReadbackAndTemplateGenerationSmoke`。
- 测试 `List` 能回读 `ma_demo_article`。
- 测试 `Detail` 保持旧 `/gen/*` 兼容响应字段，包括 `genTpl`、`genType`、`moduleName`、`functionName` 和列配置。
- 测试 `PreviewCode` 输出 Go model、schema、service、route 和 Vue api、list、edit 模板。
- 测试 `DownloadCode` 输出 zip，并包含按模块目录组织的模板文件。
- smoke 结束后删除本次列配置，并软删本次 live 表配置。
- 确认清理后 live codegen 表配置残留为 0。

## 边界

- 不创建、修改或迁移业务 schema。
- 不读取或修改 `.env`。
- 不连接业务项目数据库。
- 不删除文件或目录。
- 不修改旧 `/gen/*` 请求或响应字段。
- 不把临时生成文件写入源码目录。

## 验收

P3.7 需要通过：

```bash
cd server
GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestCodegenConfigReadbackAndTemplateGenerationSmoke|TestCodegenTableLegacyConversionPreservesOldFields|TestCodegenColumnLegacyConversionPreservesOldFields' -count=1

bash -n scripts/check-module-codegen-readback-smoke.sh scripts/check-module-tools-no-db.sh
scripts/check-module-tools-no-db.sh
MAKEADMIN_ALLOW_MODULE_CODEGEN_READBACK_WRITE=1 scripts/check-module-codegen-readback-smoke.sh
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
