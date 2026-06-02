# P3 Status

更新时间：2026-06-02

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

## 下一步

P3.6：模块生成器配置本地受控写入。该任务在 P3.5 门禁基础上实现本地 `ma_codegen_*` 事务写入和幂等 smoke。
