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
