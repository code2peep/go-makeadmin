# P3 Status

更新时间：2026-06-03

## 当前阶段

P3：业务模块产品化。

P3 从 P2 冻结面继续推进，重点是把 codegen、manifest、模块安装/卸载和验证命令串成可重复的模块开发体验。

## P3.1 当前落地

业务模块脚手架入口已建立：

- 新增 `scripts/module-scaffold.py`。
- 脚手架默认在 `examples/<module>/` 下生成 `manifest.json` 和 `README.md`。
- `manifest.json` 覆盖标准 CRUD 后端路由、前端 API、前端页面、菜单节点、权限元数据、runtime 状态和 schema 需求。
- `README.md` 输出模块约定、标准验证命令和后续 codegen 接入说明。
- 脚手架默认不覆盖已存在模块目录。
- `--dry-run` 模式只打印生成内容，不写文件、不连接数据库。
- 脚手架会校验生成 manifest，并确认它能进入注册、角色授权和卸载 SQL 生成器。
- `scripts/check-module-tools-no-db.sh` 已接入脚手架 dry-run 验证。

详见 `docs/P3_MODULE_SCAFFOLD.md`。

## P3.1 验收标准

- `python3 -m py_compile scripts/module-scaffold.py` 通过。
- `python3 scripts/module-scaffold.py --module billing_invoice --entity BillingInvoice --table ma_billing_invoice --requires-schema --dry-run` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不创建业务 schema、不执行数据库写入或删除、不读取或修改 `.env`、不连接真实 zyai 业务库。

## P3.1 验收结果

- 已通过 `python3 -m py_compile scripts/module-scaffold.py scripts/check-module-manifests.py scripts/module-registry-plan.py scripts/module-role-grant-plan.py scripts/module-uninstall-plan.py`。
- 已通过 `python3 scripts/module-scaffold.py --module billing_invoice --entity BillingInvoice --table ma_billing_invoice --requires-schema --dry-run`。
- 已通过 `python3 scripts/check-module-manifests.py`。
- 已通过 `scripts/check-module-tools-no-db.sh`，且 no-db guard 已执行脚手架 dry-run。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有创建业务 schema、没有执行数据库写入或删除、没有读取或修改 `.env`、没有连接真实 zyai 业务库。

## P3.2 当前落地

脚手架输出与代码生成器验证链路已打通：

- `scripts/module-scaffold.py` 新增 `--print-manifest`，用于输出纯 manifest JSON。
- 新增 `scripts/check-module-codegen.sh`。
- 新增 `TestGeneratedCrudCodeMatchesModuleManifest`，通过 `MAKEADMIN_CODEGEN_MANIFEST` 读取 manifest。
- 测试会使用 manifest 渲染 Go model、schema、service、route 模板，并编译临时生成的 Go 包。
- 测试会渲染 Vue API、列表页和编辑页模板，并检查 route URL 与 add/edit/del 权限和 manifest 对齐。
- `scripts/check-module-tools-no-db.sh` 已接入 `scripts/check-module-codegen.sh`。

详见 `docs/P3_MODULE_CODEGEN_LINK.md`。

## P3.2 验收标准

- `scripts/check-module-codegen.sh` 通过。
- `cd server && MAKEADMIN_CODEGEN_MANIFEST=<manifest> go test ./generator -run TestGeneratedCrudCodeMatchesModuleManifest -count=1` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不创建业务 schema、不执行数据库写入或删除、不读取或修改 `.env`、不连接真实 zyai 业务库。

## P3.2 验收结果

- 已通过 `scripts/check-module-codegen.sh`。
- 已通过临时 manifest 驱动的 `MAKEADMIN_CODEGEN_MANIFEST=<manifest> GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator -run TestGeneratedCrudCodeMatchesModuleManifest -count=1`。
- 已通过 `python3 -m py_compile scripts/module-scaffold.py`。
- 已通过 `bash -n scripts/check-module-codegen.sh scripts/check-module-tools-no-db.sh`。
- 已通过 `scripts/check-module-tools-no-db.sh`，且 no-db guard 已执行 `scaffold codegen link`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有创建业务 schema、没有执行数据库写入或删除、没有读取或修改 `.env`、没有连接真实 zyai 业务库。

## P3.3 当前落地

模块脚手架生成文件受控写入 smoke 已建立：

- `scripts/module-scaffold.py` 新增 `--examples-root`，允许把模块目录写入仓库内指定 examples root。
- 新增 `scripts/check-module-scaffold-write-smoke.sh`。
- smoke 必须显式设置 `MAKEADMIN_ALLOW_MODULE_SCAFFOLD_WRITE=1`。
- smoke 写入 `.cache/module-scaffold-smoke/<timestamp>/examples/<module>/manifest.json` 和 `README.md`。
- smoke 校验生成 manifest、安装计划 dry-run、卸载计划 dry-run 和 codegen 联动。
- `scripts/check-module-tools-no-db.sh` 已覆盖 smoke 写入门禁失败检查，默认 no-db 不执行写文件 smoke。

详见 `docs/P3_MODULE_SCAFFOLD_WRITE_SMOKE.md`。

## P3.3 验收标准

- `scripts/check-module-scaffold-write-smoke.sh` 失败，且错误说明没有写文件、没有访问数据库。
- `MAKEADMIN_ALLOW_MODULE_SCAFFOLD_WRITE=1 scripts/check-module-scaffold-write-smoke.sh` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不创建业务 schema、不执行数据库写入或删除、不读取或修改 `.env`、不连接真实 zyai 业务库。

## P3.3 验收结果

- 已通过 `scripts/check-module-scaffold-write-smoke.sh` 失败门禁；失败文案明确没有写文件且没有访问数据库。
- 已通过 `MAKEADMIN_ALLOW_MODULE_SCAFFOLD_WRITE=1 scripts/check-module-scaffold-write-smoke.sh`。
- smoke 已在 `.cache/module-scaffold-smoke/<timestamp>/examples/<module>/` 下写出 `manifest.json` 和 `README.md`。
- 已通过生成 manifest 的 JSON 校验、安装计划 dry-run、卸载计划 dry-run 和 `scripts/check-module-codegen.sh --manifest <generated manifest>`。
- 已通过 `python3 -m py_compile scripts/module-scaffold.py`。
- 已通过 `bash -n scripts/check-module-scaffold-write-smoke.sh scripts/check-module-codegen.sh scripts/check-module-tools-no-db.sh`。
- 已通过 `scripts/check-module-tools-no-db.sh`，且 no-db guard 已执行 filesystem write gate。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已确认 `.cache/` 为 Git 忽略目录，smoke 产物不会进入提交。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有创建业务 schema、没有执行数据库写入或删除、没有读取或修改 `.env`、没有连接真实 zyai 业务库。

## P3.4 当前落地

模块脚手架产物与后台生成器配置预览已打通：

- 新增 `scripts/module-codegen-plan.py`。
- 新增 `scripts/check-module-codegen-plan.sh`。
- `module-codegen-plan.py` 读取 manifest 并输出 `ma_codegen_table` / `ma_codegen_column` 对应的 `makeadmin` 配置预览。
- 同时输出旧 `GenTable` / `GenTableColumn` 对应的 `legacy` 配置预览，保留 `/gen/*` 兼容面。
- 默认 CRUD 列配置包含 `id`、`title`、`status`。
- manifest 声明 `codegen.columns` 时会使用 `id` 加自定义列配置，并保留 `htmlType`、`dictType`、`queryType`。
- `scripts/check-module-scaffold-write-smoke.sh` 已使用实际写入的 `.cache/.../manifest.json` 生成 `codegen-plan.json` 并断言表名和列配置。
- `scripts/check-module-tools-no-db.sh` 已接入 codegen plan 验证。

详见 `docs/P3_MODULE_CODEGEN_PLAN.md`。

## P3.4 验收标准

- `python3 scripts/module-codegen-plan.py --manifest examples/demo/manifest.json --format json` 通过。
- `scripts/check-module-codegen-plan.sh` 通过。
- `MAKEADMIN_ALLOW_MODULE_SCAFFOLD_WRITE=1 scripts/check-module-scaffold-write-smoke.sh` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不写 `ma_codegen_*`、不创建业务 schema、不执行数据库写入或删除、不读取或修改 `.env`、不连接真实 zyai 业务库。

## P3.4 验收结果

- 已通过 `python3 -m py_compile scripts/module-codegen-plan.py scripts/module-scaffold.py scripts/check-module-manifests.py`。
- 已通过 `bash -n scripts/check-module-codegen-plan.sh scripts/check-module-tools-no-db.sh`。
- 已通过 `python3 scripts/module-codegen-plan.py --manifest examples/demo/manifest.json --format json | python3 -m json.tool >/dev/null`。
- 已通过 `scripts/check-module-codegen-plan.sh`，覆盖默认列和 manifest `codegen.columns` 自定义列配置。
- 已通过 `MAKEADMIN_ALLOW_MODULE_SCAFFOLD_WRITE=1 scripts/check-module-scaffold-write-smoke.sh`，实际写入产物已生成 `codegen-plan.json` 并完成表名和列配置断言。
- 已通过 `scripts/check-module-tools-no-db.sh`，且 no-db guard 已执行 scaffold codegen plan。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有写入 `ma_codegen_*`、没有创建业务 schema、没有执行数据库写入或删除、没有读取或修改 `.env`、没有连接真实 zyai 业务库。

## P3.5 当前落地

模块生成器配置写入边界已建立：

- `scripts/module-codegen-plan.py` 新增保留的 `--apply` 写入入口。
- 写入必须显式设置 `MAKEADMIN_ALLOW_MODULE_CODEGEN_WRITE=1`。
- 写入必须显式确认 `--confirm-module <module>` 和 `--confirm-source-table <table>`。
- 写入必须显式传入 `--confirm-sync-columns`，用于确认未来同步 `ma_codegen_column` 时可能删除 stale 列配置。
- P3.5 当前即使全部门禁满足，也会在数据库访问前失败；真正写入执行器不在本阶段开放。
- 新增 `scripts/check-module-codegen-apply-boundary.sh`。
- `scripts/check-module-tools-no-db.sh` 已接入 codegen apply boundary 验证。

详见 `docs/P3_MODULE_CODEGEN_APPLY_BOUNDARY.md`。

## P3.5 验收标准

- `scripts/check-module-codegen-apply-boundary.sh` 通过。
- `python3 scripts/module-codegen-plan.py --apply` 失败，且错误说明没有访问数据库。
- `MAKEADMIN_ALLOW_MODULE_CODEGEN_WRITE=1 python3 scripts/module-codegen-plan.py --apply --confirm-module article --confirm-source-table ma_demo_article --confirm-sync-columns` 失败，且错误说明没有访问数据库。
- `scripts/check-module-tools-no-db.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不写 `ma_codegen_*`、不创建业务 schema、不执行数据库写入或删除、不读取或修改 `.env`、不连接业务项目数据库。

## P3.5 验收结果

- 已通过 `python3 -m py_compile scripts/module-codegen-plan.py scripts/module-scaffold.py scripts/check-module-manifests.py`。
- 已通过 `bash -n scripts/check-module-codegen-apply-boundary.sh scripts/check-module-tools-no-db.sh scripts/check-module-codegen-plan.sh`。
- 已通过 `scripts/check-module-codegen-apply-boundary.sh`，覆盖缺少环境变量、缺少确认模块、缺少确认表名、缺少列同步确认和全部门禁满足但执行器未开放的失败路径。
- 已通过 `scripts/check-module-tools-no-db.sh`，且 no-db guard 已执行 codegen apply boundary。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有写入 `ma_codegen_*`、没有创建业务 schema、没有执行数据库写入或删除、没有读取或修改 `.env`、没有连接业务项目数据库。

## P3.6 当前落地

模块生成器配置本地受控写入已建立：

- `scripts/module-codegen-plan.py --apply` 已开放本地受控写入。
- 写入必须满足 `MAKEADMIN_ALLOW_MODULE_CODEGEN_WRITE=1`、`--confirm-module`、`--confirm-source-table` 和 `--confirm-sync-columns`。
- apply 会在单事务内写入或更新 `ma_codegen_table` 和 `ma_codegen_column`。
- live 表配置按 `tenant_id + table_name + delete_time=0` 幂等。
- 列配置按 `table_id + column_name` 幂等 upsert。
- 已存在 live 表配置必须与 manifest 的 `module_name`、`business_name`、`entity_name` 对齐，否则停止，不覆盖。
- `--confirm-sync-columns` 会删除同一 `table_id` 下 manifest 已移除的 stale 列配置。
- 新增 `scripts/check-module-codegen-apply-smoke.sh`。
- `scripts/check-module-tools-no-db.sh` 已覆盖 codegen apply smoke 缺环境变量门禁。

详见 `docs/P3_MODULE_CODEGEN_APPLY.md`。

## P3.6 验收标准

- `python3 -m py_compile scripts/module-codegen-plan.py` 通过。
- `bash -n scripts/check-module-codegen-apply-boundary.sh scripts/check-module-codegen-apply-smoke.sh scripts/check-module-tools-no-db.sh` 通过。
- `scripts/check-module-codegen-apply-boundary.sh` 通过。
- `MAKEADMIN_ALLOW_MODULE_CODEGEN_SMOKE_WRITE=1 scripts/check-module-codegen-apply-smoke.sh` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- smoke 清理后 `ma_demo_article` live codegen 表配置残留为 0。
- 不创建业务 schema、不读取或修改 `.env`、不连接业务项目数据库。

## P3.6 验收结果

- 已通过 `python3 -m py_compile scripts/module-codegen-plan.py`。
- 已通过 `bash -n scripts/check-module-codegen-apply-boundary.sh scripts/check-module-codegen-apply-smoke.sh scripts/check-module-tools-no-db.sh`。
- 已通过 `scripts/check-module-codegen-apply-boundary.sh`。
- 已通过 `MAKEADMIN_ALLOW_MODULE_CODEGEN_SMOKE_WRITE=1 scripts/check-module-codegen-apply-smoke.sh`。
- smoke 已完成第一次 apply、stale 列插入、第二次 apply 同步删除 stale 列和最终清理。
- 已确认清理后 `tenant_id=0 + ma_demo_article + delete_time=0` live codegen 表配置残留为 0。
- 已通过 `scripts/check-module-tools-no-db.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有创建业务 schema、没有读取或修改 `.env`、没有连接业务项目数据库。

## P3.7 当前落地

生成器配置回读与模板生成闭环已建立：

- 新增 `TestCodegenConfigReadbackAndTemplateGenerationSmoke`。
- 新增 `scripts/check-module-codegen-readback-smoke.sh`。
- smoke 必须显式设置 `MAKEADMIN_ALLOW_MODULE_CODEGEN_READBACK_WRITE=1`。
- smoke 会先确认本地没有 `tenant_id=0 + ma_demo_article + delete_time=0` live codegen 表配置。
- smoke 复用 P3.6 的受控 apply，把 demo manifest 写入 `ma_codegen_table` 和 `ma_codegen_column`。
- 测试会通过生成器服务回读 `List`、`Detail`、`PreviewCode` 和 `DownloadCode`。
- `Detail` 验证旧 `/gen/*` 兼容响应字段仍保持可用。
- `PreviewCode` 验证 Go 和 Vue 模板可以由回读配置渲染。
- `DownloadCode` 验证 zip 中包含按模块目录组织的生成文件。
- smoke 完成后清理本次列配置并软删本次表配置，确认 live 残留为 0。
- `scripts/check-module-tools-no-db.sh` 已覆盖 readback smoke 缺环境变量门禁。

详见 `docs/P3_MODULE_CODEGEN_READBACK.md`。

## P3.7 验收标准

- `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestCodegenConfigReadbackAndTemplateGenerationSmoke|TestCodegenTableLegacyConversionPreservesOldFields|TestCodegenColumnLegacyConversionPreservesOldFields' -count=1` 通过。
- `bash -n scripts/check-module-codegen-readback-smoke.sh scripts/check-module-tools-no-db.sh` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `MAKEADMIN_ALLOW_MODULE_CODEGEN_READBACK_WRITE=1 scripts/check-module-codegen-readback-smoke.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- smoke 清理后 `ma_demo_article` live codegen 表配置残留为 0。
- 不创建业务 schema、不读取或修改 `.env`、不连接业务项目数据库。

## P3.7 验收结果

- 已通过 `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestCodegenConfigReadbackAndTemplateGenerationSmoke|TestCodegenTableLegacyConversionPreservesOldFields|TestCodegenColumnLegacyConversionPreservesOldFields' -count=1`。
- 已通过 `bash -n scripts/check-module-codegen-readback-smoke.sh scripts/check-module-tools-no-db.sh`。
- 已通过 `scripts/check-module-tools-no-db.sh`，且 no-db guard 已执行 readback smoke 缺环境变量门禁。
- 已通过 `MAKEADMIN_ALLOW_MODULE_CODEGEN_READBACK_WRITE=1 scripts/check-module-codegen-readback-smoke.sh`。
- 已确认清理后 `tenant_id=0 + ma_demo_article + delete_time=0` live codegen 表配置残留为 0。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有创建业务 schema、没有读取或修改 `.env`、没有连接业务项目数据库。

## P3.8 当前落地

后台生成器页面与模块 manifest 的操作闭环已建立：

- 新增 `POST /gen/previewCode`。
- 新接口复用现有 `gen:previewCode` 权限面，不新增权限 SQL。
- 旧 `GET /gen/previewCode?id=...` 保持不变。
- 新接口支持仓库内 `manifest.json` 路径和 inline manifest JSON。
- `manifestPath` 限制在仓库内，且文件名必须为 `manifest.json`。
- 新接口返回 manifest 摘要、兼容 `/gen/detail` 的生成器配置和模板代码预览。
- 新增 `PreviewModuleManifest` 服务方法和 no-db 单元测试。
- 新增管理端 `Manifest 预览` 弹窗。
- 弹窗支持仓库路径模式、JSON 模式、配置摘要、字段表格和代码预览。
- 新增 `scripts/check-module-manifest-preview.sh`。
- `scripts/check-module-tools-no-db.sh` 已接入 manifest preview 验证。

详见 `docs/P3_MODULE_MANIFEST_PREVIEW.md`。

## P3.8 验收标准

- `scripts/check-module-manifest-preview.sh` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/... ./generator/service/gen -run 'TestPreviewModuleManifest|TestGeneratedCrudCodeMatchesModuleManifest|TestCodegenConfigReadbackAndTemplateGenerationSmoke|TestCodegenTableLegacyConversionPreservesOldFields|TestCodegenColumnLegacyConversionPreservesOldFields' -count=1` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不写 `ma_codegen_*`、不生成文件、不创建业务 schema、不读取或修改 `.env`、不新增权限 SQL。

## P3.8 验收结果

- 已通过 `scripts/check-module-manifest-preview.sh`。
- 已通过 `scripts/check-module-tools-no-db.sh`。
- 已通过 `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/... ./generator/service/gen -run 'TestPreviewModuleManifest|TestGeneratedCrudCodeMatchesModuleManifest|TestCodegenConfigReadbackAndTemplateGenerationSmoke|TestCodegenTableLegacyConversionPreservesOldFields|TestCodegenColumnLegacyConversionPreservesOldFields' -count=1`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有写 `ma_codegen_*`、没有生成文件、没有创建业务 schema、没有读取或修改 `.env`、没有新增权限 SQL。

## P3.9 当前落地

模块 manifest 预览结果的安装计划联动已建立：

- `POST /gen/previewCode` 返回新增 `plan`。
- `plan` 包含 `registrySql`、`roleGrantSql`、`installSql`、`uninstallSql` 和 `runtimeHint`。
- 安装计划 SQL 由 Go 服务内存生成，不在请求中调用 Python 脚本。
- `registrySql` 对应 manifest 菜单和权限注册预览。
- `roleGrantSql` 对应指定租户和角色的权限授权预览。
- `installSql` 合并注册 SQL 和角色授权 SQL。
- `uninstallSql` 只按 manifest 声明的权限 code 和菜单 routeName 生成清理预览。
- 管理端 `Manifest 预览` 弹窗新增租户 ID、角色 ID 输入。
- 管理端预览结果新增租户、角色、运行时提示展示。
- 管理端新增 `安装计划` 按钮，使用代码预览弹窗展示 SQL。
- 新增 `scripts/check-module-install-plan-preview.sh`。
- `scripts/check-module-tools-no-db.sh` 已接入 install plan preview 验证。

详见 `docs/P3_MODULE_INSTALL_PLAN_PREVIEW.md`。

## P3.9 验收标准

- `scripts/check-module-install-plan-preview.sh` 通过。
- `scripts/check-module-manifest-preview.sh` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestPreviewModuleManifest|TestPreviewModuleManifestIncludesInstallPlan' -count=1` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不执行安装 SQL、不执行卸载 SQL、不写数据库、不创建业务 schema、不读取或修改 `.env`、不新增权限 SQL。

## P3.9 验收结果

- 已通过 `scripts/check-module-install-plan-preview.sh`。
- 已通过 `scripts/check-module-manifest-preview.sh`。
- 已通过 `scripts/check-module-tools-no-db.sh`。
- 已通过 `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestPreviewModuleManifest|TestPreviewModuleManifestIncludesInstallPlan' -count=1`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有执行安装 SQL、没有执行卸载 SQL、没有写数据库、没有创建业务 schema、没有读取或修改 `.env`、没有新增权限 SQL。

## 下一步

## P3.10 当前落地

后台模块安装写入门禁与确认参数已建立：

- 新增 `PUT /gen/previewCode`。
- 新接口复用现有 `gen:previewCode` 权限面，不新增权限 SQL。
- 新增 `ApplyModuleManifestInstall` 服务方法。
- 写入门禁必须显式设置 `MAKEADMIN_ALLOW_MODULE_INSTALL_APPLY=1`。
- 确认参数覆盖 `confirmModule`、`confirmTenantId`、`confirmRoleId`、`confirmInstall` 和 `confirmSchemaRisk`。
- `manifest.requiresSchema=true` 时必须显式确认 schema 风险。
- P3.10 即使全部门禁满足，也会在安装执行器处阻断，并明确返回 `no database access was attempted`。
- 门禁响应返回 manifest 摘要、租户、角色、安装计划和结构化检查结果。
- 管理端 `Manifest 预览` 弹窗新增写入确认控件和 `写入门禁` 按钮。
- 新增 `scripts/check-module-install-apply-boundary.sh`。
- `scripts/check-module-tools-no-db.sh` 已接入 install apply boundary 验证。

详见 `docs/P3_MODULE_INSTALL_APPLY_BOUNDARY.md`。

## P3.10 验收标准

- `scripts/check-module-install-apply-boundary.sh` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestPreviewModuleManifest|TestModuleManifestInstallApplyGate' -count=1` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不执行安装 SQL、不执行卸载 SQL、不写数据库、不创建业务 schema、不读取或修改 `.env`、不新增权限 SQL。

## P3.10 验收结果

- 已通过 `scripts/check-module-install-apply-boundary.sh`。
- 已通过 `scripts/check-module-tools-no-db.sh`，且 no-db guard 已执行 install apply boundary。
- 已通过 `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestPreviewModuleManifest|TestModuleManifestInstallApplyGate' -count=1`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有执行安装 SQL、没有执行卸载 SQL、没有写数据库、没有创建业务 schema、没有读取或修改 `.env`、没有新增权限 SQL。

## 下一步

## P3.11 当前落地

后台模块安装受控本地写入 smoke 已建立：

- `PUT /gen/previewCode` 在 P3.10 门禁基础上开放本地受控写入。
- 写入仍必须显式设置 `MAKEADMIN_ALLOW_MODULE_INSTALL_APPLY=1`。
- 写入前会校验数据库目标必须是本地 `go_makeadmin`。
- 写入前后返回 `permissions`、`menus`、`menuPermissions`、`rolePermissions` 快照。
- 写入在单事务内完成。
- 权限按 `ma_permission.code` 幂等。
- 菜单按 `ma_menu.route_name + delete_time=0` 幂等。
- 菜单权限按 `menu_id + permission_id` 幂等。
- 角色授权按 `tenant_id + role_id + permission_id` 幂等。
- 目标角色不存在时跳过角色授权，不影响菜单和权限注册。
- 新增 `TestModuleManifestInstallApplyLocalSmoke`。
- 新增 `scripts/check-module-install-apply-smoke.sh`。
- `scripts/check-module-tools-no-db.sh` 已覆盖 install apply smoke 缺环境变量门禁。

详见 `docs/P3_MODULE_INSTALL_APPLY.md`。

## P3.11 验收标准

- `scripts/check-module-install-apply-boundary.sh` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `MAKEADMIN_ALLOW_MODULE_INSTALL_SMOKE_WRITE=1 scripts/check-module-install-apply-smoke.sh` 通过。
- `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestModuleManifestInstallApplyGate|TestModuleManifestInstallApplyLocalSmoke' -count=1` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- smoke 清理后 demo article 注册行和角色授权残留为 0。
- 不执行卸载接口、不创建业务 schema、不读取或修改 `.env`、不新增权限 SQL、不连接业务项目数据库。

## P3.11 验收结果

- 已通过 `scripts/check-module-install-apply-boundary.sh`。
- 已通过 `scripts/check-module-tools-no-db.sh`，且 no-db guard 已覆盖 install apply smoke 缺环境变量门禁。
- 已通过 `MAKEADMIN_ALLOW_MODULE_INSTALL_SMOKE_WRITE=1 scripts/check-module-install-apply-smoke.sh`。
- smoke 已完成第一次安装、第二次幂等安装、最终清理和残留检查。
- 已通过 `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestModuleManifestInstallApplyGate|TestModuleManifestInstallApplyLocalSmoke' -count=1`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有执行卸载接口、没有创建业务 schema、没有读取或修改 `.env`、没有新增权限 SQL、没有连接业务项目数据库。

## P3.12 当前落地

后台模块卸载写入门禁与确认参数已建立：

- 新增 `DELETE /gen/previewCode`。
- 新接口复用现有 `gen:previewCode` 权限面，不新增权限 SQL。
- 新增 `ApplyModuleManifestUninstall` 服务方法。
- 写入门禁必须显式设置 `MAKEADMIN_ALLOW_MODULE_UNINSTALL_APPLY=1`。
- 确认参数覆盖 `confirmModule` 和 `confirmDelete`。
- 门禁响应返回 manifest 摘要、卸载 SQL 预览和结构化检查结果。
- P3.12 即使全部门禁满足，也会在卸载执行器处阻断，并明确返回 `no database access was attempted`。
- 新增 `scripts/check-module-uninstall-apply-boundary.sh`。
- `scripts/check-module-tools-no-db.sh` 已接入 uninstall apply boundary 验证。

详见 `docs/P3_MODULE_UNINSTALL_APPLY_BOUNDARY.md`。

## P3.12 验收标准

- `scripts/check-module-uninstall-apply-boundary.sh` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestModuleManifestUninstallApplyGate' -count=1` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不执行卸载 SQL、不删除数据库行、不创建业务 schema、不读取或修改 `.env`、不新增权限 SQL。

## P3.12 验收结果

- 已通过 `scripts/check-module-uninstall-apply-boundary.sh`。
- 已通过 `scripts/check-module-tools-no-db.sh`，且 no-db guard 已执行 uninstall apply boundary。
- 已通过 `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestModuleManifestUninstallApplyGate' -count=1`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有执行卸载 SQL、没有删除数据库行、没有创建业务 schema、没有读取或修改 `.env`、没有新增权限 SQL。

## P3.13 当前落地

后台模块卸载受控本地删除 smoke 已建立：

- `DELETE /gen/previewCode` 在 P3.12 门禁基础上开放本地受控删除。
- 删除仍必须显式设置 `MAKEADMIN_ALLOW_MODULE_UNINSTALL_APPLY=1`。
- 删除前会校验数据库目标必须是本地 `go_makeadmin`。
- 删除前后返回 `permissions`、`menus`、`menuPermissions`、`rolePermissions` 快照。
- 删除在单事务内完成。
- 删除顺序为角色授权、菜单权限、菜单、权限。
- 删除范围只来自 manifest 的 permission codes 和 `menu.routeName`。
- 新增 `TestModuleManifestUninstallApplyLocalSmoke`。
- 新增 `scripts/check-module-uninstall-apply-smoke.sh`。
- `scripts/check-module-tools-no-db.sh` 已覆盖 uninstall apply smoke 缺环境变量门禁。

详见 `docs/P3_MODULE_UNINSTALL_APPLY.md`。

## P3.13 验收标准

- `scripts/check-module-uninstall-apply-boundary.sh` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `MAKEADMIN_ALLOW_MODULE_UNINSTALL_SMOKE_WRITE=1 scripts/check-module-uninstall-apply-smoke.sh` 通过。
- `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestModuleManifestUninstallApplyGate|TestModuleManifestUninstallApplyLocalSmoke' -count=1` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- smoke 清理后 demo article 注册行和角色授权残留为 0。
- 不删除业务表、不创建业务 schema、不读取或修改 `.env`、不新增权限 SQL、不连接业务项目数据库。

## P3.13 验收结果

- 已通过 `scripts/check-module-uninstall-apply-boundary.sh`。
- 已通过 `scripts/check-module-tools-no-db.sh`，且 no-db guard 已覆盖 uninstall apply smoke 缺环境变量门禁。
- 已通过 `MAKEADMIN_ALLOW_MODULE_UNINSTALL_SMOKE_WRITE=1 scripts/check-module-uninstall-apply-smoke.sh`。
- smoke 已完成安装种子、第一次卸载、第二次 no-op 卸载和最终残留检查。
- 已通过 `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestModuleManifestUninstallApplyGate|TestModuleManifestUninstallApplyLocalSmoke' -count=1`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有删除业务表、没有创建业务 schema、没有读取或修改 `.env`、没有新增权限 SQL、没有连接业务项目数据库。

## P3.14 当前落地

后台模块安装/卸载门禁结果页面闭环已建立：

- 管理端 `Manifest 预览` 弹窗新增卸载执行入口。
- 前端新增 `applyModuleManifestUninstall` API。
- 安装执行结果进入安装结果 tab。
- 卸载执行结果进入卸载结果 tab。
- 成功和失败响应都会在页面落地，不只依赖 toast。
- 页面展示后端返回的状态、模块、来源、写入环境变量。
- 页面展示门禁检查列表。
- 页面展示安装/卸载执行前后的权限、菜单、菜单权限和角色授权快照。
- 安装和卸载继续复用 `gen:previewCode` 权限面。

详见 `docs/P3_MODULE_APPLY_UI_CLOSURE.md`。

## P3.14 验收标准

- `cd admin && npm run type-check` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不新增后端写入规则、不创建业务 schema、不读取或修改 `.env`、不新增权限 SQL。

## P3.14 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `scripts/check-module-tools-no-db.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有新增后端写入规则、没有创建业务 schema、没有读取或修改 `.env`、没有新增权限 SQL。

## 下一步

## P3.15 当前落地

模块安装/卸载操作摘要与审计规划已建立：

- 安装 apply 响应新增 `summary`。
- 卸载 apply 响应新增 `summary`。
- `summary` 覆盖操作类型、模块、实体、来源表、菜单路由、权限编码、schema 风险、本地库范围和 runtime 提示。
- 成功和失败响应都会返回 `summary`。
- 管理端安装结果和卸载结果新增操作类型、路由名和权限编码标签。
- 文档规划后续审计日志模型，但本阶段不创建审计表。

详见 `docs/P3_MODULE_APPLY_AUDIT_SUMMARY.md`。

## P3.15 验收标准

- `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestModuleManifestInstallApplyGate|TestModuleManifestUninstallApplyGate' -count=1` 通过。
- `cd admin && npm run type-check` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不创建审计表、不修改数据库 schema、不读取或修改 `.env`、不新增权限 SQL。

## P3.15 验收结果

- 已通过 `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestModuleManifestInstallApplyGate|TestModuleManifestUninstallApplyGate' -count=1`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `scripts/check-module-tools-no-db.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有创建审计表、没有修改数据库 schema、没有读取或修改 `.env`、没有新增权限 SQL。

## P3.16 当前落地

模块 apply 页面权限与按钮状态已收敛：

- `Manifest 预览` 弹窗内安装、卸载继续复用 `gen:previewCode` 权限面。
- 安装计划和代码预览只允许基于当前输入对应的 preview 打开。
- 安装执行按钮必须满足当前 preview、确认模块匹配、安装写入确认和 schema 风险确认。
- 卸载执行按钮必须满足当前 preview、确认模块匹配和删除确认。
- preview、安装、卸载请求期间会锁住对应按钮状态。
- apply 请求期间会锁住确认输入和确认勾选项。
- 修改来源模式、manifest 路径、manifest JSON、作者、租户或角色后会清空旧 preview 和旧 apply 结果。
- 重新生成 preview 会清空旧安装结果和旧卸载结果。
- 发起安装会清空旧安装结果，发起卸载会清空旧卸载结果。

详见 `docs/P3_MODULE_APPLY_UI_STATE.md`。

## P3.16 验收标准

- `cd admin && npm run type-check` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不修改后端写入门禁、不新增后端接口、不创建业务 schema、不读取或修改 `.env`、不新增权限 SQL。

## P3.16 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `scripts/check-module-tools-no-db.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- 本阶段没有修改后端写入门禁、没有新增后端接口、没有创建业务 schema、没有读取或修改 `.env`、没有新增权限 SQL。

## P3.17 当前落地

模块 manifest 前端 API 类型已收敛：

- `admin/src/api/tools/code.ts` 新增 manifest preview、install apply、uninstall apply 请求参数类型。
- `admin/src/api/tools/code.ts` 新增 manifest preview、apply result、plan、summary、check、snapshot 响应类型。
- `previewModuleManifest` 返回 `Promise<ModuleManifestPreviewResult>`。
- `applyModuleManifestInstall` 返回 `Promise<ModuleManifestApplyResult>`。
- `applyModuleManifestUninstall` 返回 `Promise<ModuleManifestApplyResult>`。
- `Manifest 预览` 弹窗使用 API 类型替换核心 `any`。
- 快照表格行使用本地 `SnapshotRow` 类型。
- 失败响应仍保留兜底对象，兼容 request 拦截器 reject 后端 `data` 的现有行为。

详见 `docs/P3_MODULE_MANIFEST_API_TYPES.md`。

## P3.17 验收标准

- `cd admin && npm run type-check` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不修改请求 URL、method 或字段名，不修改 request 封装，不修改后端响应结构，不创建业务 schema，不读取或修改 `.env`，不新增权限 SQL。

## P3.17 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `scripts/check-module-tools-no-db.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- 本阶段没有修改请求 URL、method 或字段名，没有修改 request 封装，没有修改后端响应结构，没有创建业务 schema，没有读取或修改 `.env`，没有新增权限 SQL。

## P3.18 当前落地

模块 manifest apply 错误响应已归一化：

- `admin/src/api/tools/code.ts` 新增 `normalizeModuleManifestApplyError`。
- 后端业务失败返回的 apply result 会保留 `summary`、`checks`、`before`、`after` 等结构化字段。
- 后端 result 缺少 `message` 时使用当前操作 fallback message。
- 后端 result 缺少 `checks` 时补为空数组。
- `Error` 对象会转换为 `{ message, checks: [] }`。
- 字符串错误会转换为 `{ message, checks: [] }`。
- 其他未知错误会使用 fallback message。
- `Manifest 预览` 弹窗安装、卸载 catch 统一调用归一化 helper，不再直接强转 `unknown error`。

详见 `docs/P3_MODULE_MANIFEST_APPLY_ERROR.md`。

## P3.18 验收标准

- `cd admin && npm run type-check` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不修改请求 URL、method 或字段名，不修改后端响应结构，不修改后端写入门禁，不创建业务 schema，不读取或修改 `.env`，不新增权限 SQL。

## P3.18 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `scripts/check-module-tools-no-db.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- 本阶段没有修改请求 URL、method 或字段名，没有修改后端响应结构，没有修改后端写入门禁，没有创建业务 schema，没有读取或修改 `.env`，没有新增权限 SQL。

## P3.19 当前落地

模块 manifest apply 前端结果视图已提取：

- 新增 `admin/src/views/dev_tools/components/module-manifest-apply-result.vue`。
- 新组件接收 `ModuleManifestApplyResult` 和 `fallbackTitle`。
- 新组件负责渲染结果 alert、状态摘要、操作摘要、权限编码、执行前后快照和门禁检查列表。
- `module-manifest-preview.vue` 只保留安装结果 tab、卸载结果 tab、apply 请求和状态管理。
- 安装、卸载结果重复模板已移入局部组件。
- 父组件中由本次提取产生的快照 helper 和权限标签样式已清理。

详见 `docs/P3_MODULE_APPLY_RESULT_VIEW.md`。

## P3.19 验收标准

- `cd admin && npm run type-check` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不改变按钮状态规则，不改变错误归一化规则，不修改请求 URL、method 或字段名，不修改后端响应结构，不修改后端写入门禁，不创建业务 schema，不读取或修改 `.env`，不新增权限 SQL。

## P3.19 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `scripts/check-module-tools-no-db.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- 本阶段没有改变按钮状态规则，没有改变错误归一化规则，没有修改请求 URL、method 或字段名，没有修改后端响应结构，没有修改后端写入门禁，没有创建业务 schema，没有读取或修改 `.env`，没有新增权限 SQL。

## P3.20 当前落地

模块 manifest apply 审计事件 DTO 草图已建立：

- 新增 `server/generator/schemas/resp/module_manifest_audit.go`。
- Go DTO 覆盖审计操作人、审计范围和审计事件。
- 审计事件草图包含事件 ID、操作类型、manifest 来源、manifest 摘要、操作摘要、执行范围、执行状态、门禁环境变量、检查结果、执行前后快照、操作人、请求时间和完成时间。
- `admin/src/api/tools/code.ts` 新增对应 TypeScript DTO。
- TypeScript DTO 字段按后端 JSON 命名，供后续审计页面、接口草图和结果详情复用。
- 当前 DTO 不接入 `PUT /gen/previewCode` 或 `DELETE /gen/previewCode`。
- 当前 DTO 不接入数据库模型、数据库迁移、审计写入服务或菜单权限 SQL。

详见 `docs/P3_MODULE_APPLY_AUDIT_DTO.md`。

## P3.20 验收标准

- `gofmt -w server/generator/schemas/resp/module_manifest_audit.go` 通过。
- `cd admin && npm run type-check` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不创建审计表、不修改数据库 schema、不接入写库、不修改接口响应、不读取或修改 `.env`、不新增权限 SQL。

## P3.20 验收结果

- 已通过 `gofmt -w server/generator/schemas/resp/module_manifest_audit.go`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `scripts/check-module-tools-no-db.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- 本阶段没有创建审计表、没有修改数据库 schema、没有接入写库、没有修改接口响应、没有读取或修改 `.env`、没有新增权限 SQL。

## P3.21 当前落地

模块 manifest apply 审计事件构造器 dry-run 已建立：

- 新增 `server/generator/service/gen/module_manifest_apply_audit.go`。
- 新增 `buildModuleManifestInstallAuditEvent`。
- 新增 `buildModuleManifestUninstallAuditEvent`。
- 构造器把现有 install/uninstall apply result 组合成 `ModuleManifestApplyAuditEventResp`。
- 构造器输入包含 event ID、actor、requestedAt 和 completedAt。
- 构造器不生成 event ID、不读取数据库、不读取当前用户。
- 新增 `server/generator/service/gen/module_manifest_apply_audit_test.go`。
- 单测覆盖安装和卸载审计事件映射。

详见 `docs/P3_MODULE_APPLY_AUDIT_BUILDER.md`。

## P3.21 验收标准

- `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestBuildModuleManifest.*AuditEvent' -count=1` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不接入接口、不创建审计表、不修改数据库 schema、不写库、不读取或修改 `.env`、不新增权限 SQL。

## P3.21 验收结果

- 已通过 `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestBuildModuleManifest.*AuditEvent' -count=1`。
- 已通过 `scripts/check-module-tools-no-db.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- 本阶段没有接入接口、没有创建审计表、没有修改数据库 schema、没有写库、没有读取或修改 `.env`、没有新增权限 SQL。

## P3.22 当前落地

模块 manifest apply 审计事件前端 dry-run 预览已建立：

- `admin/src/api/tools/code.ts` 新增 `buildModuleManifestApplyAuditPreview`。
- 前端构造器基于 `ModuleManifestApplyResult` 生成 `ModuleManifestApplyAuditEventResult`。
- `scope.tenantId` 和 `scope.roleId` 优先来自 apply result，其次来自 plan。
- 默认 `eventId` 为 `preview`。
- 默认 `actor.type` 为 `frontend-preview`。
- 默认时间使用前端本地 ISO 时间，仅作为 dry-run 展示。
- `module-manifest-apply-result.vue` 新增 `审计预览` 操作。
- 点击后展示格式化 JSON 代码块。

详见 `docs/P3_MODULE_APPLY_AUDIT_PREVIEW.md`。

## P3.22 验收标准

- `cd admin && npm run type-check` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不调用后端、不新增接口、不写库、不创建审计表、不读取当前登录用户、不读取或修改 `.env`、不新增权限 SQL。

## P3.22 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `scripts/check-module-tools-no-db.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- 本阶段没有调用后端、没有新增接口、没有写库、没有创建审计表、没有读取当前登录用户、没有读取或修改 `.env`、没有新增权限 SQL。

## P3.23 当前落地

模块 manifest apply 审计预览已调整为摘要优先：

- `admin/src/api/tools/code.ts` 新增 `buildModuleManifestApplyAuditPreviewSummary`。
- 审计预览摘要包含操作、模块、状态、路由、权限数量、检查项数量、执行前快照总数、执行后快照总数、数据库范围和操作人类型。
- `module-manifest-apply-result.vue` 点击 `审计预览` 后先展示摘要。
- 完整 JSON 需要点击 `JSON` 后二次展开。
- 关闭审计预览时同步收起 JSON。

详见 `docs/P3_MODULE_APPLY_AUDIT_PREVIEW_SUMMARY.md`。

## P3.23 验收标准

- `cd admin && npm run type-check` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不改变后端、不新增接口、不调用后端、不写库、不创建审计表、不读取当前登录用户、不读取或修改 `.env`、不新增权限 SQL。

## P3.23 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `scripts/check-module-tools-no-db.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- 本阶段没有改变后端、没有新增接口、没有调用后端、没有写库、没有创建审计表、没有读取当前登录用户、没有读取或修改 `.env`、没有新增权限 SQL。

## P3.24 当前落地

模块 manifest apply 结果空态文案已收敛：

- 权限编码为空时显示 `无权限编码`。
- 没有执行前后快照时显示 `无快照`。
- 没有门禁检查项时显示 `无检查项`。
- 空态文案统一使用低强调样式。

详见 `docs/P3_MODULE_APPLY_RESULT_EMPTY_STATE.md`。

## P3.24 验收标准

- `cd admin && npm run type-check` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不改变按钮状态规则、不改变错误归一化规则、不改变审计预览构造规则、不改变后端、不新增接口、不写库、不创建审计表、不读取或修改 `.env`、不新增权限 SQL。

## P3.24 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `scripts/check-module-tools-no-db.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有改变按钮状态规则、没有改变错误归一化规则、没有改变审计预览构造规则、没有改变后端、没有新增接口、没有写库、没有创建审计表、没有读取或修改 `.env`、没有新增权限 SQL。

## 下一步

## P3.25 当前落地

P3 阶段冻结验收与最终状态文档已补齐：

- 新增 `docs/P3_FINAL_STATUS.md`。
- 汇总 P3.1-P3.24 的脚手架、codegen、manifest 后台闭环、管理端结果视图和审计 dry-run 范围。
- 明确 P3 冻结验收命令、保留边界、不覆盖范围和 P4 入口。
- README 已补充 P3 最终状态入口，并把当前阶段从 P2 更新为 P4 入口前状态。

## P3.25 验收标准

- `scripts/check-module-tools-no-db.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- `git diff --check` 通过。
- `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache` 通过。
- 不改后端、不改前端运行逻辑、不新增接口、不写库、不创建 schema、不读取或修改 `.env`、不新增权限 SQL。

## P3.25 验收结果

- 已通过 `scripts/check-module-tools-no-db.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有改后端、没有改前端运行逻辑、没有新增接口、没有写库、没有创建 schema、没有读取或修改 `.env`、没有新增权限 SQL。

## 下一步

P4.1：可见后台与人工测试闭环。建议启动本地 API 和管理端，确认登录、工作台、菜单、系统管理、代码生成器和模块 manifest 预览页面能被人工打开和操作。
