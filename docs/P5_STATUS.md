# P5 Status

更新时间：2026-06-03

## 当前阶段

P5：示例模块真实安装与后台菜单可见闭环。

P5 从 P4 可见后台冻结面继续推进，方向继续偏底座、生成器、安装卸载闭环。第一步先让一个 demo 模块从 manifest、前端产物、后端运行时、菜单权限写入到管理端可见页面全部打通。

## P5.1 当前落地

Demo Article 可见模块已建立：

- `examples/demo/manifest.json` 菜单改为 `visible=true`。
- `examples/demo/manifest.json` 运行时标记为 `runtimeRegistered=true`。
- demo 数据库菜单路径对齐为 `/dev_tools/demo/article`。
- demo 管理端运行路由为 `/demo/article`。
- 新增 `admin/src/api/article.ts`。
- 新增 `admin/src/views/article/index.vue`。
- 新增 `admin/src/views/article/edit.vue`。
- 前端动态路由允许加载 `admin/src/views/article/**/*.vue`。
- 新增 `scripts/check-demo-module-visible.sh`，显式写入本地 demo 菜单、权限和角色授权。

详见 `docs/P5_DEMO_MODULE_VISIBLE.md`。

## P5.1 验收标准

- `python3 scripts/check-module-manifests.py` 通过。
- `MAKEADMIN_ALLOW_DEMO_MODULE_VISIBLE_WRITE=1 scripts/check-demo-module-visible.sh` 通过。
- `cd admin && npm run type-check` 通过。
- `cd admin && npm run build` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 本地 API 使用 `MAKEADMIN_ENABLE_DEMO_MODULE=1` 启动。
- 浏览器能打开 `http://127.0.0.1:5173/demo/article`。
- 页面显示 `Demo Article`、`article / ma_demo_article` 和 `只读示例`。
- 页面不是 403、404 或空白页。
- 不创建 `ma_demo_article` 表、不改数据库 schema、不改 `.env`、不连接真实 zyai 业务库。

## P5.1 验收结果

- 已通过 `python3 scripts/check-module-manifests.py`。
- 已通过 `bash -n scripts/check-demo-module-visible.sh`。
- 已通过 `MAKEADMIN_ALLOW_DEMO_MODULE_VISIBLE_WRITE=1 scripts/check-demo-module-visible.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `cd admin && npm run build`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 本地 API 已用 `MAKEADMIN_ENABLE_DEMO_MODULE=1 ./scripts/dev-api.sh` 启动，并确认注册 `/api/article/list`、`/api/article/detail`、`/api/article/add`、`/api/article/edit`、`/api/article/del`。
- 已通过浏览器人工验证 `http://127.0.0.1:5173/demo/article`。
- 已通过浏览器人工验证页面显示 `Demo Article`、`article / ma_demo_article`、`P5.1` 和 `只读示例`。
- 已通过浏览器人工验证页面空表格显示 `暂无数据`。
- 已通过浏览器人工验证 `运行时详情` 显示 `module=article` 和 `runtimeRegistered=true`。
- `npm run build` 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有创建 `ma_demo_article` 表、没有改数据库 schema、没有改 `.env`、没有连接真实 zyai 业务库。

## 下一步

P5.2：示例模块安装/卸载后台操作闭环。建议让模块中心在开启本地写入 env 后能对 Demo Article 执行真实安装和卸载，并在页面内显示 applied、快照和审计预览。

## P5.2 当前落地

模块中心真实 apply 验收入口已补：

- 模块中心内置模块清单阶段标记更新为 `P5.2`。
- Demo Article 状态从 `可预览` 调整为 `可安装`。
- Demo Article 增加页面入口 `/demo/article`。
- 内置模块清单操作区增加 `打开`，可以从模块中心进入 demo 页面。
- 明确本地 API 需要临时开启 `MAKEADMIN_ENABLE_DEMO_MODULE=1`、`MAKEADMIN_ALLOW_MODULE_INSTALL_APPLY=1`、`MAKEADMIN_ALLOW_MODULE_UNINSTALL_APPLY=1`。

详见 `docs/P5_MODULE_CENTER_APPLY.md`。

## P5.2 验收标准

- `cd admin && npm run type-check` 通过。
- `cd admin && npm run build` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 本地 API 用 demo runtime 和 install/uninstall apply env 启动。
- 浏览器能打开 `/module`。
- 模块中心内置模块清单显示 `P5.2`、`可安装` 和 `/demo/article`。
- Demo Article 在模块中心可执行安装 apply，结果为 `applied`。
- Demo Article 在模块中心可执行卸载 apply，结果为 `applied`。
- 重新安装后能从模块中心打开 `/demo/article`。
- 不创建 `ma_demo_article` 表、不改数据库 schema、不改 `.env`、不连接真实 zyai 业务库。

## P5.2 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `cd admin && npm run build`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 本地 API 已用 `MAKEADMIN_ENABLE_DEMO_MODULE=1 MAKEADMIN_ALLOW_MODULE_INSTALL_APPLY=1 MAKEADMIN_ALLOW_MODULE_UNINSTALL_APPLY=1 ./scripts/dev-api.sh` 启动。
- 已通过浏览器人工验证 `http://127.0.0.1:5173/module`。
- 已通过浏览器人工验证模块中心显示 `P5.2`、`可安装` 和 `/demo/article`。
- 已通过浏览器人工验证 Demo Article 安装 apply，结果为 `applied`。
- 已通过浏览器人工验证 Demo Article 卸载 apply，结果为 `applied`，权限、菜单、菜单权限、角色权限快照回到 0。
- 已通过浏览器人工验证重新安装 apply，结果为 `applied`，并把本地 demo 模块恢复为可见状态。
- 已通过浏览器人工验证模块中心 `打开` 可进入 `http://127.0.0.1:5173/demo/article`。
- 已通过浏览器人工验证 Demo Article 页面显示 `Demo Article` 和 `只读示例`。
- 已通过浏览器人工验证 `运行时详情` 显示 `module article` 和 `runtimeRegistered true`。
- `npm run build` 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有创建 `ma_demo_article` 表、没有改数据库 schema、没有改 `.env`、没有连接真实 zyai 业务库。

## 下一步

P5.3：模块中心安装状态回读。建议让模块中心读取当前 demo 模块的安装状态，直接显示已安装、未安装、权限数、菜单可见和运行时状态，减少人工判断。

## P5.3 当前落地

模块中心安装状态回读已建立：

- 新增 `POST /api/gen/previewCode/status` 只读接口。
- 状态接口复用 manifest 加载和本地库目标校验。
- 状态响应包含 `status`、`snapshot`、`expected`、`missing`、`runtimeEnv`、`runtimeEnabled`、`runtimeRegistered` 和 `menuVisible`。
- 模块中心内置模块清单阶段标记更新为 `P5.3`。
- 模块中心增加 `刷新状态`、安装状态、安装快照和运行时状态列。
- 进入模块中心自动回读状态，安装 apply 或卸载 apply 后自动刷新状态。

详见 `docs/P5_MODULE_STATUS_READBACK.md`。

## P5.3 验收标准

- `go test ./generator/service/gen ./generator/routers/gen` 通过。
- `cd admin && npm run type-check` 通过。
- `cd admin && npm run build` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 本地 API 用 demo runtime 和 install/uninstall apply env 启动。
- 浏览器能打开 `/module`。
- 模块中心显示 `P5.3`、`已安装`、`权限 5/5`、`角色授权 5/5`、`已开启` 和 `MAKEADMIN_ENABLE_DEMO_MODULE=1`。
- 卸载 apply 后模块中心状态变为 `未安装`，快照回到 `0/5`。
- 重新安装 apply 后模块中心状态恢复为 `已安装`，快照恢复为 `5/5`。
- 不创建 `ma_demo_article` 表、不改数据库 schema、不改 `.env`、不连接真实 zyai 业务库。

## P5.3 验收结果

- 已通过 `go test ./generator/service/gen ./generator/routers/gen`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `cd admin && npm run build`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 本地 API 已用 `MAKEADMIN_ENABLE_DEMO_MODULE=1 MAKEADMIN_ALLOW_MODULE_INSTALL_APPLY=1 MAKEADMIN_ALLOW_MODULE_UNINSTALL_APPLY=1 ./scripts/dev-api.sh` 启动，并确认注册 `/api/gen/previewCode/status`。
- 已通过浏览器人工验证模块中心显示 `P5.3`、`已安装`、`权限 5/5`、`角色授权 5/5`、`已开启` 和 `MAKEADMIN_ENABLE_DEMO_MODULE=1`。
- 已通过浏览器人工验证卸载 apply 后状态变为 `未安装`，快照显示 `权限 0/5` 和 `角色授权 0/5`。
- 已通过浏览器人工验证重新安装 apply 后状态恢复为 `已安装`，快照显示 `权限 5/5` 和 `角色授权 5/5`。
- `npm run build` 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有创建 `ma_demo_article` 表、没有改数据库 schema、没有改 `.env`、没有连接真实 zyai 业务库。

## 下一步

P5.4：模块中心安装状态筛选与多模块兼容。建议把状态回读从 Demo Article 单行扩展为多模块列表能力，支持 installed/partial/uninstalled 筛选，并为后续非 demo 模块接入统一展示。

## P5.4 当前落地

模块中心状态筛选与汇总已建立：

- 模块中心内置模块清单阶段标记更新为 `P5.4`。
- 增加 `全部`、`已安装`、`部分安装`、`未安装`、`异常` 状态筛选。
- 增加总数、已安装、部分、未安装、异常汇总。
- 模块列表改为使用筛选后的 `filteredModules`。
- 无匹配模块时显示 `暂无匹配模块`。

详见 `docs/P5_MODULE_STATUS_FILTERS.md`。

## P5.4 验收标准

- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 浏览器能打开 `/module`。
- 模块中心显示 `P5.4`、状态汇总和 Demo Article 的 `已安装` 状态。
- `未安装` 筛选显示空态。
- `全部` 筛选恢复 Demo Article 行。

## P5.4 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过浏览器人工验证模块中心显示 `P5.4`、状态汇总和 Demo Article 的 `已安装` 状态。
- 已通过浏览器人工验证 `未安装` 筛选显示 `暂无匹配模块`。
- 已通过浏览器人工验证切回 `全部` 后恢复 Demo Article 行，并显示 `权限 5/5` 和 `角色授权 5/5`。

## 下一步

P5.5：模块 registry 后端只读列表。建议把内置模块清单从前端硬编码迁移为后端只读 registry 输出，为多模块和后续真实业务模块接入做准备。

## P5.5 当前落地

模块 registry 后端只读列表已建立：

- 新增 `GET /api/gen/moduleRegistry`。
- 后端 registry 当前返回 Demo Article 的 `module`、`manifest`、`table`、`runtime`、`entry` 和产品状态。
- 模块中心进入页面时先读取 registry，再逐项回读安装状态。
- 模块中心内置模块清单阶段标记更新为 `P5.5`。
- 前端不再把 Demo Article 的清单元数据写死在 `modules` 初始数据里。

详见 `docs/P5_MODULE_REGISTRY_READONLY.md`。

## P5.5 验收标准

- `go test ./generator/service/gen ./generator/routers/gen` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 本地 API 注册 `/api/gen/moduleRegistry`。
- 浏览器能打开 `/module`。
- 模块中心显示 `P5.5`、`Demo Article`、`examples/demo/manifest.json`、`MAKEADMIN_ENABLE_DEMO_MODULE=1`。
- 模块中心状态回读仍显示 `已安装`、`权限 5/5` 和 `角色授权 5/5`。

## P5.5 验收结果

- 已通过 `go test ./generator/service/gen ./generator/routers/gen`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 本地 API 已用 `MAKEADMIN_ENABLE_DEMO_MODULE=1 MAKEADMIN_ALLOW_MODULE_INSTALL_APPLY=1 MAKEADMIN_ALLOW_MODULE_UNINSTALL_APPLY=1 ./scripts/dev-api.sh` 启动，并确认注册 `/api/gen/moduleRegistry`。
- 已通过浏览器人工验证模块中心显示 `P5.5`、`Demo Article`、`examples/demo/manifest.json`、`MAKEADMIN_ENABLE_DEMO_MODULE=1`。
- 已通过浏览器人工验证状态回读仍显示 `已安装`、`权限 5/5` 和 `角色授权 5/5`。
- `npm run build` 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。

## 下一步

P5.6：模块 registry 与 manifest 校验联动。建议让 registry 项在后端启动或请求时校验 manifest 文件存在、模块名一致、入口字段完整，并把校验结果带回模块中心。

## P5.6 当前落地

模块 registry 与 manifest 校验联动已建立：

- `GET /api/gen/moduleRegistry` 返回项新增 `manifestStatus`、`manifestMessage` 和 `manifestChecks`。
- registry 返回时会校验 manifest 可读取、结构有效、`module` 一致、`table` 一致、runtime hint 一致、管理端入口合法、菜单路由字段完整。
- 单个 registry 项校验失败时仍返回该项，并标记为 `failed`，避免模块中心整页空白。
- 模块中心新增 `校验` 列。
- 模块中心 `异常` 筛选同时包含 registry 校验异常和安装状态异常。

详见 `docs/P5_MODULE_REGISTRY_MANIFEST_CHECK.md`。

## P5.6 验收标准

- `go test ./generator/service/gen ./generator/routers/gen` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 浏览器重新登录后，模块中心显示 `P5.6`、Demo Article 校验 `已通过`。

## P5.6 验收结果

- 已通过 `go test ./generator/service/gen ./generator/routers/gen`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 本地 API 已用 `MAKEADMIN_ENABLE_DEMO_MODULE=1 MAKEADMIN_ALLOW_MODULE_INSTALL_APPLY=1 MAKEADMIN_ALLOW_MODULE_UNINSTALL_APPLY=1 ./scripts/dev-api.sh` 启动，并确认注册 `/api/gen/moduleRegistry`。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此模块中心页面人工复验需要重新登录后执行。

## 下一步

P5.7：模块 registry 校验明细展示。建议在模块中心为校验列增加明细弹窗，展示每个 `manifestChecks` 项，方便以后多模块接入时快速定位 manifest 与 registry 的差异。

## P5.7 当前落地

模块 registry 校验明细展示已接入模块中心：

- 模块中心顶部阶段标识更新为 `P5.7`。
- `校验` 列新增带图标的 `明细` 操作。
- 明细弹窗展示模块名、manifest 路径、整体校验状态和整体说明。
- 明细弹窗展示 `manifestChecks` 检查项列表。
- 检查项状态映射为 `通过`、`异常`、`阻断` 或原始状态。
- 没有检查项时明细按钮不可点击。

详见 `docs/P5_MODULE_REGISTRY_CHECK_DETAIL.md`。

## P5.7 验收标准

- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 浏览器重新登录后，模块中心显示 `P5.7`，点击 Demo Article 的 `明细` 后可看到 manifest 校验检查项。

## P5.7 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此模块中心页面人工复验需要重新登录后执行。

## 下一步

P5.8：模块 registry 校验异常 fixture 与页面异常态验收。建议增加一个受环境变量控制的本地异常 registry 项，用来验证模块中心 `异常` 筛选、校验明细弹窗和后端单项失败不中断列表的产品闭环。

## P5.8 当前落地

模块 registry 异常 fixture 已接入：

- 新增环境变量开关 `MAKEADMIN_ENABLE_BROKEN_MODULE_REGISTRY_FIXTURE=1`。
- 默认情况下 `GET /api/gen/moduleRegistry` 仍只返回 `Demo Article`。
- 打开开关后额外返回 `Broken Manifest Fixture`。
- 异常 fixture 使用合法但不存在的 `examples/demo/missing/manifest.json`。
- 后端会把异常项标记为 `manifestStatus=failed`，同时保留正常 Demo Article。
- 模块中心阶段标识更新为 `P5.8`。
- 模块中心复用已有 `异常` 筛选、`校验` 列和 `明细` 弹窗展示失败检查项。

详见 `docs/P5_MODULE_REGISTRY_FAILURE_FIXTURE.md`。

## P5.8 验收标准

- `cd server && go test ./generator/service/gen` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 显式设置 `MAKEADMIN_ENABLE_BROKEN_MODULE_REGISTRY_FIXTURE=1` 后，模块中心可看到异常 fixture 并能打开校验明细。

## P5.8 验收结果

- 已通过 `cd server && go test ./generator/service/gen`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此模块中心页面人工复验需要重新登录后执行。

## 下一步

P5.9：模块 registry 后端契约测试与 CLI smoke。建议用脚本直接请求或调用 registry，分别验证默认清单和 broken fixture 清单，形成不用登录后台也能验收模块中心数据契约的 smoke。

## P5.9 当前落地

模块 registry CLI smoke 已接入：

- 新增 `scripts/check-module-registry-smoke.sh`。
- smoke 验证默认清单只返回 `Demo Article`，且 manifest 校验通过。
- smoke 验证 `MAKEADMIN_ENABLE_BROKEN_MODULE_REGISTRY_FIXTURE=1` 下返回 `Broken Manifest Fixture`。
- smoke 验证异常 fixture 校验失败，但不影响 Demo Article。
- `scripts/check-module-tools-no-db.sh` 已接入 registry smoke。
- `./scripts/verify-no-db.sh` 会覆盖 registry smoke。
- 模块中心阶段标识更新为 `P5.9`。

详见 `docs/P5_MODULE_REGISTRY_SMOKE.md`。

## P5.9 验收标准

- `scripts/check-module-registry-smoke.sh` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。

## P5.9 验收结果

- 已通过 `scripts/check-module-registry-smoke.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。

## 下一步

P5.10：模块中心登录后人工验收与页面截图记录。建议在你重新登录后台后，分别用默认 registry 和 broken fixture registry 做一次可视化验收，确认 P5.6-P5.9 的页面闭环。

## P5.10 当前落地

模块 registry 路由响应契约测试已接入：

- 新增 `server/generator/routers/gen/module_registry_test.go`。
- 路由测试验证默认响应 `code=200`。
- 路由测试验证默认 registry 只返回 `Demo Article`，且 `manifestStatus=passed`。
- 路由测试验证开启 `MAKEADMIN_ENABLE_BROKEN_MODULE_REGISTRY_FIXTURE=1` 后返回 `Broken Manifest Fixture`。
- 路由测试验证 broken fixture `manifestStatus=failed`，且不影响 Demo Article。
- `scripts/check-module-registry-smoke.sh` 已接入路由响应测试。
- 模块中心阶段标识更新为 `P5.10`。

详见 `docs/P5_MODULE_REGISTRY_ROUTE_CONTRACT.md`。

## P5.10 验收标准

- `cd server && go test ./generator/routers/gen` 通过。
- `scripts/check-module-registry-smoke.sh` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。

## P5.10 验收结果

- 已通过 `cd server && go test ./generator/routers/gen`。
- 已通过 `scripts/check-module-registry-smoke.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 下一步

P5.11：模块 registry 本地 API smoke 文档化。建议把默认 registry、broken fixture registry、路由契约和登录后页面验收拆成一个稳定的人工/自动验收矩阵。

## P5.11 当前落地

模块 registry 自动/人工验收矩阵已建立：

- 新增 `docs/P5_MODULE_REGISTRY_ACCEPTANCE_MATRIX.md`。
- 自动矩阵覆盖默认 registry 服务契约。
- 自动矩阵覆盖 broken fixture 服务契约。
- 自动矩阵覆盖 `/api/gen/moduleRegistry` 路由响应契约。
- 自动矩阵覆盖 `./scripts/verify-no-db.sh` 全量 no-db 链路。
- 人工矩阵拆出默认模块中心、broken fixture 页面态、异常筛选、校验明细弹窗和 Demo Article 入口。
- 模块中心阶段标识更新为 `P5.11`。

详见 `docs/P5_MODULE_REGISTRY_ACCEPTANCE_MATRIX.md`。

## P5.11 验收标准

- `scripts/check-module-registry-smoke.sh` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 登录后页面验收矩阵有清晰执行入口和预期结果。

## P5.11 验收结果

- 已通过 `scripts/check-module-registry-smoke.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 下一步

P5.12：模块中心页面验收辅助状态。建议在模块中心增加一个轻量的验收提示区，展示当前 registry 来源、broken fixture 是否开启、自动 smoke 命令和登录后人工验收入口，让页面本身更适合框架交付验收。

## P5.12 当前落地

模块中心页面验收辅助状态已接入：

- `内置模块清单` 区域新增只读状态条。
- 状态条展示 registry 来源 `/api/gen/moduleRegistry`。
- 状态条展示当前 registry 模块数量。
- 状态条展示 registry 校验异常数量。
- 状态条根据返回模块判断 `Broken Fixture` 是否开启。
- 状态条展示自动 smoke 命令 `check-module-registry-smoke.sh`。
- 状态条展示人工验收入口 `/module`。
- 模块中心阶段标识更新为 `P5.12`。

详见 `docs/P5_MODULE_CENTER_ACCEPTANCE_STATUS.md`。

## P5.12 验收标准

- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 登录后模块中心可看到 P5.12 验收辅助状态条。

## P5.12 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 下一步

P5.13：模块中心 registry 状态空态与错误态收敛。建议给 registry 读取失败和空清单加入更明确的页面状态，保证模块中心在 API 异常时仍能给出可执行的排查入口。

## P5.13 当前落地

模块中心 registry 状态空态与错误态已收敛：

- registry 读取失败时展示 `Registry 读取失败` alert。
- 失败 alert 带上错误信息和 `scripts/check-module-registry-smoke.sh`。
- registry 读取成功但模块列表为空时展示 `Registry 暂无模块` alert。
- 表格空态根据当前状态切换为 `registry 读取失败`、`registry 暂无模块` 或 `暂无匹配模块`。
- registry 读取失败时不继续逐项读取安装状态。
- 模块中心阶段标识更新为 `P5.13`。

详见 `docs/P5_MODULE_CENTER_REGISTRY_STATES.md`。

## P5.13 验收标准

- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 登录后模块中心在 registry 失败或空清单时有明确页面状态。

## P5.13 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 下一步

P5.14：模块中心 registry 状态单元测试。建议把 registry 错误态、空态、broken fixture 状态条等前端纯逻辑抽成可测试函数，降低后续页面调整时的回归风险。

## P5.14 当前落地

模块中心 registry 状态纯逻辑已抽出：

- 新增 `admin/src/views/dev_tools/module/registry-state.ts`。
- 抽出 registry 失败数量计算。
- 抽出 broken fixture 是否存在判断。
- 抽出 registry 空态判断。
- 抽出 registry 错误详情文案。
- 抽出表格 empty text 选择。
- 抽出验收辅助状态条 rows 构造。
- `admin/src/views/dev_tools/module/index.vue` 改为只负责请求、状态持有和渲染。
- 模块中心阶段标识更新为 `P5.14`。

当前前端没有 Vitest/Jest 等单测框架，P5.14 不新增测试依赖，先用现有 `vue-tsc` 和 no-db 全量验证覆盖 helper 类型契约。

详见 `docs/P5_MODULE_CENTER_REGISTRY_STATE_HELPER.md`。

## P5.14 验收标准

- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- registry 状态计算从页面内联逻辑收敛到独立 helper。

## P5.14 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 下一步

P5.15：模块中心 registry 状态 helper smoke。建议增加一个不依赖测试框架的 TypeScript 编译期 fixture 文件，覆盖 helper 的默认、broken fixture、错误态和空态输入输出形状。

## P5.15 当前落地

模块中心 registry 状态 helper 编译期 fixture 已接入：

- 新增 `admin/src/views/dev_tools/module/registry-state.fixture.ts`。
- fixture 覆盖默认 registry module 输入。
- fixture 覆盖 broken fixture registry module 输入。
- fixture 覆盖空 registry 状态输入。
- fixture 覆盖失败 registry 状态输入。
- fixture 覆盖验收辅助 rows、失败数量、broken fixture 判断、空态判断、empty text 和错误详情输出形状。
- 模块中心阶段标识更新为 `P5.15`。

详见 `docs/P5_MODULE_CENTER_REGISTRY_STATE_FIXTURE.md`。

## P5.15 验收标准

- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- fixture 位于 `admin/src/**/*` 并被 `admin/tsconfig.json` include 规则覆盖。

## P5.15 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 下一步

P5.16：模块中心 registry 文档索引收敛。建议把 P5.5-P5.15 的 registry 相关文档整理成一个短索引，降低后续继续做模块市场/真实模块接入时的查找成本。

## P5.16 当前落地

模块中心 registry 文档索引已收敛：

- 新增 `docs/P5_MODULE_REGISTRY_DOC_INDEX.md`。
- 索引串联 P5.5-P5.15 registry 文档阅读顺序。
- 索引列出后端 service、路由测试、管理端页面、helper、fixture 和 smoke 脚本入口。
- 索引保留默认 registry、broken fixture 和登录后人工验收启动方式。
- README 增加 P5.16 文档入口。

详见 `docs/P5_MODULE_REGISTRY_DOC_INDEX.md`。

## P5.16 验收标准

- `git diff --check` 通过。
- `scripts/check-module-registry-smoke.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- README 和 P5 状态文档都能定位到 P5.16 索引。

## P5.16 验收结果

- 已通过 `git diff --check`。
- 已通过 `scripts/check-module-registry-smoke.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 下一步

P5.17：模块中心登录后人工验收清单压缩。建议把当前反复出现的登录后复验项整理成一份短 checklist，后续用户登录后可以按一个入口完成默认 registry、broken fixture、异常筛选和明细弹窗复验。

## P5.17 当前落地

模块中心登录后人工验收清单已压缩到页面：

- 模块中心内置模块清单区域新增 registry checklist。
- checklist 覆盖默认 registry、broken fixture、异常筛选、校验明细和 Demo 入口。
- 新增 `buildRegistryManualChecklistRows()` helper。
- `registry-state.fixture.ts` 覆盖默认、broken fixture 和空 registry checklist。
- 模块中心阶段标识更新为 `P5.17`。
- README 和 registry 文档索引增加 P5.17 入口。

详见 `docs/P5_MODULE_CENTER_MANUAL_CHECKLIST.md`。

## P5.17 验收标准

- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 登录后模块中心显示 `P5.17` 和 registry checklist。

## P5.17 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 下一步

P5.18：模块中心 registry UI 截图验收准备。建议补一个本地浏览器验收说明或脚本输出，明确默认 registry 和 broken fixture 两种启动方式下页面应出现的关键文字，等你登录后就能一次性人工确认。

## P5.18 当前落地

模块中心 registry UI 契约检查已建立：

- 新增 `scripts/check-module-center-ui-contract.sh`。
- 脚本检查模块中心阶段标识、registry checklist、helper、fixture 和 P5.18 文档关键字。
- `scripts/check-module-tools-no-db.sh` 接入 `Module tools: module center UI contract`。
- 模块中心阶段标识更新为 `P5.18`。
- README 和 registry 文档索引增加 P5.18 入口。

详见 `docs/P5_MODULE_CENTER_UI_CONTRACT.md`。

## P5.18 验收标准

- `scripts/check-module-center-ui-contract.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 登录后模块中心可按默认 registry 和 broken fixture 关键字完成一次人工确认。

## P5.18 验收结果

- 已通过 `scripts/check-module-center-ui-contract.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 下一步

P5.19：模块中心 registry 进入真实模块接入前的冻结清单。建议整理当前 registry、manifest、UI、smoke 和人工验收边界，判断 P5 是否可以从 Demo Article 过渡到第二个示例模块或模块市场雏形。

## P5.19 当前落地

模块中心 registry 冻结清单已建立：

- 新增 `docs/P5_MODULE_REGISTRY_FREEZE_CHECKLIST.md`。
- 清单覆盖后端 registry、manifest 校验、broken fixture、smoke、路由契约、模块中心 UI、helper/fixture 和人工验收边界。
- 自动验收基线明确为 registry smoke、UI contract 和 no-db 全量验证。
- 判断当前自动验收基线足够进入第二个示例模块接入。
- README 和 registry 文档索引增加 P5.19 入口。

详见 `docs/P5_MODULE_REGISTRY_FREEZE_CHECKLIST.md`。

## P5.19 验收标准

- `scripts/check-module-registry-smoke.sh` 通过。
- `scripts/check-module-center-ui-contract.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。

## P5.19 验收结果

- 已通过 `scripts/check-module-registry-smoke.sh`。
- 已通过 `scripts/check-module-center-ui-contract.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 下一步

P5.20：第二个示例模块 registry 接入。建议新增一个轻量只读示例模块，用来验证 registry 和模块中心不是只服务 Demo Article 单模块。

## P5.20 当前落地

第二个只读示例模块 `Demo Notice` 已接入：

- 新增 `examples/demo_notice/manifest.json`。
- 新增 `admin/src/api/demoNotice.ts`。
- 新增 `admin/src/views/demo_notice/index.vue`。
- 管理端动态路由显式纳入 `admin/src/views/demo_notice/**/*.vue`。
- registry 增加 `MAKEADMIN_ENABLE_DEMO_NOTICE_MODULE=1` 开关。
- registry smoke 增加 Demo Notice service 和 route 契约。
- UI contract 增加 Demo Notice 和 P5.20 阶段标识检查。
- 模块中心阶段标识更新为 `P5.20`。

详见 `docs/P5_SECOND_DEMO_MODULE.md`。

## P5.20 验收标准

- `python3 scripts/check-module-manifests.py` 通过。
- `scripts/check-module-registry-smoke.sh` 通过。
- `scripts/check-module-center-ui-contract.sh` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。

## P5.20 验收结果

- 已通过 `python3 scripts/check-module-manifests.py`。
- 已通过 `scripts/check-module-registry-smoke.sh`。
- 已通过 `scripts/check-module-center-ui-contract.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 已生成 `demo_notice` chunk。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 下一步

P5.21：第二个示例模块安装计划和页面入口验收。建议验证 `Demo Notice` 的 manifest 预览、安装计划、模块中心打开入口和未注册运行时提示。

## P5.21 当前落地

`Demo Notice` 安装计划和页面入口验收已建立：

- 新增 `scripts/check-demo-notice-module.sh`。
- 脚本覆盖 manifest 校验、service preview/install plan、registry、route handler 和前端入口契约。
- `scripts/check-module-tools-no-db.sh` 接入 `Module tools: demo notice module contract`。
- `Demo Notice` 页面阶段标识更新为 `P5.21`。
- 模块中心阶段标识更新为 `P5.21`。
- `PreviewDemoNoticeManifestIncludesInstallPlan` 覆盖安装计划关键 SQL 和无运行时 gate 提示。
- README 增加 P5.21 文档入口。

详见 `docs/P5_DEMO_NOTICE_ACCEPTANCE.md`。

## P5.21 验收标准

- `scripts/check-demo-notice-module.sh` 通过。
- `scripts/check-module-registry-smoke.sh` 通过。
- `scripts/check-module-center-ui-contract.sh` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。

## P5.21 验收结果

- 已通过 `scripts/check-demo-notice-module.sh`。
- 已通过 `scripts/check-module-registry-smoke.sh`。
- 已通过 `scripts/check-module-center-ui-contract.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 已生成 `demo_notice` chunk。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 下一步

P5.22：第二个示例模块模块中心状态回读说明。建议把 `Demo Notice` 在未注册运行时和未安装状态下的模块中心状态文案整理成稳定契约，避免后续多模块状态显示混乱。

## P5.22 当前落地

`Demo Notice` 模块中心状态回读契约已建立：

- 新增 `buildModuleRuntimeStatus()` helper。
- 模块中心运行时状态判断从页面内联逻辑迁入 helper。
- helper 覆盖 `未注册`、`未开启` 和 `已开启` 三种状态。
- `registry-state.fixture.ts` 增加 `demoNoticeRuntimeStatus` 和 `demoArticleRuntimeStatus`。
- `scripts/check-demo-notice-module.sh` 增加 helper、未注册文案和 fixture 检查。
- 模块中心阶段标识更新为 `P5.22`。
- README 和 registry 文档索引增加 P5.22 入口。

详见 `docs/P5_DEMO_NOTICE_STATUS_CONTRACT.md`。

## P5.22 验收标准

- `scripts/check-demo-notice-module.sh` 通过。
- `scripts/check-module-center-ui-contract.sh` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。

## P5.22 验收结果

- 已通过 `scripts/check-demo-notice-module.sh`。
- 已通过 `scripts/check-module-center-ui-contract.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 已生成 `demo_notice` chunk。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 下一步

P5.23：模块中心多模块筛选验收。建议验证 `Demo Article` 与 `Demo Notice` 同时存在时，全部、未安装、异常筛选和验收状态条的数量统计保持稳定。

## P5.23 当前落地

模块中心多模块筛选验收已建立：

- 新增 `filterRegistryModules()` helper。
- 新增 `buildModuleStatusSummary()` helper。
- 新增 `isRegistryModuleFailed()` helper。
- 模块中心 `filteredModules` 和 `moduleStatusSummary` 改为调用 helper。
- `registry-state.fixture.ts` 增加 `article + demo_notice` 双模块状态。
- 新增 `scripts/check-module-center-filter-contract.sh`。
- `scripts/check-module-tools-no-db.sh` 接入 `Module tools: module center filter contract`。
- 模块中心阶段标识更新为 `P5.23`。
- README 和 registry 文档索引增加 P5.23 入口。

详见 `docs/P5_MODULE_CENTER_MULTI_FILTERS.md`。

## P5.23 验收标准

- `scripts/check-module-center-filter-contract.sh` 通过。
- `scripts/check-demo-notice-module.sh` 通过。
- `scripts/check-module-center-ui-contract.sh` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。

## P5.23 验收结果

- 已通过 `scripts/check-module-center-filter-contract.sh`。
- 已通过 `scripts/check-demo-notice-module.sh`。
- 已通过 `scripts/check-module-center-ui-contract.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 已生成 `demo_notice` chunk。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 下一步

P5.24：模块中心多模块人工验收入口收敛。建议把 `MAKEADMIN_ENABLE_DEMO_NOTICE_MODULE=1` 下模块中心应出现的关键字、筛选状态和打开入口整理成一个登录后 checklist。

## P5.24 当前落地

模块中心多模块人工验收入口已收敛：

- registry checklist 新增 `多模块` 项。
- registry checklist 新增 `Demo Notice` 项。
- checklist 明确展示 `MAKEADMIN_ENABLE_DEMO_NOTICE_MODULE=1` 和 `/demo/notice`。
- `registry-state.fixture.ts` 增加 `multiRegistryModules`、`multiRegistryState` 和 `multiChecklistRows`。
- 新增 `scripts/check-module-center-manual-checklist.sh`。
- `scripts/check-module-tools-no-db.sh` 接入 `Module tools: module center manual checklist contract`。
- 模块中心阶段标识更新为 `P5.24`。
- README 和 registry 文档索引增加 P5.24 入口。

详见 `docs/P5_MODULE_CENTER_MANUAL_MULTI_CHECKLIST.md`。

## P5.24 验收标准

- `scripts/check-module-center-manual-checklist.sh` 通过。
- `scripts/check-module-center-filter-contract.sh` 通过。
- `scripts/check-demo-notice-module.sh` 通过。
- `scripts/check-module-center-ui-contract.sh` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。

## P5.24 验收结果

- 已通过 `scripts/check-module-center-manual-checklist.sh`。
- 已通过 `scripts/check-module-center-filter-contract.sh`。
- 已通过 `scripts/check-demo-notice-module.sh`。
- 已通过 `scripts/check-module-center-ui-contract.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 已生成 `demo_notice` chunk。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 下一步

P5.25：P5 模块中心多模块冻结验收。建议对 Demo Article、Demo Notice、registry smoke、模块中心 UI contract、登录后人工 checklist 和 no-db 链路做 P5 最终状态归档。

## P5.25 当前落地

P5 模块中心多模块冻结验收已建立：

- 新增 `docs/P5_FINAL_STATUS.md`。
- 新增 `scripts/check-p5-module-center-freeze.sh`。
- `scripts/check-module-tools-no-db.sh` 接入 `Module tools: P5 module center freeze contract`。
- P5 最终状态覆盖 Demo Article、Demo Notice、registry smoke、UI contract、filter contract、manual checklist 和 no-db 链路。
- 模块中心阶段标识更新为 `P5.25`。
- README 和 registry 文档索引增加 P5 最终状态入口。

详见 `docs/P5_FINAL_STATUS.md`。

## P5.25 验收标准

- `scripts/check-p5-module-center-freeze.sh` 通过。
- `scripts/check-module-registry-smoke.sh` 通过。
- `scripts/check-module-center-ui-contract.sh` 通过。
- `scripts/check-module-center-filter-contract.sh` 通过。
- `scripts/check-module-center-manual-checklist.sh` 通过。
- `scripts/check-demo-notice-module.sh` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。

## P5.25 验收结果

- 已通过 `scripts/check-p5-module-center-freeze.sh`。
- 已通过 `scripts/check-module-registry-smoke.sh`。
- 已通过 `scripts/check-module-center-ui-contract.sh`。
- 已通过 `scripts/check-module-center-filter-contract.sh`。
- 已通过 `scripts/check-module-center-manual-checklist.sh`。
- 已通过 `scripts/check-demo-notice-module.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 已生成 `demo_notice` chunk。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 登录后人工 checklist 当前未执行；本地未提供 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，不修改 `.env` 或管理员密码。

## 下一步

P6.1：模块中心产品化入口。建议开始把当前模块中心从验收/工具界面推进为产品化模块市场入口，优先整理模块详情、安装向导和状态说明。
