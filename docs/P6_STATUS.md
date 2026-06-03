# P6 Status

## P6.1：模块中心产品化入口

P6.1 从已冻结的 P5 多模块 registry 底座出发，把模块中心从验收工具台推进为后台产品入口。

## P6.1 当前落地

模块中心产品化入口已建立：

- 模块中心顶部新增 `模块市场`、`模块详情` 和 `安装向导` 三个产品化区域。
- `模块市场` 基于 registry 与安装状态回读展示模块数、已安装、待安装和异常数量。
- `模块详情` 跟随当前选中模块展示 manifest、运行时和入口，并提供预览、安装计划和打开动作。
- `安装向导` 基于当前选中模块展示 manifest 校验、安装计划、写入状态和页面入口。
- 页面同时保留 `P6.1` 当前阶段标识和 `P5.25` 冻结标识，避免 P5 历史 contract 漂移。
- 新增 `buildModuleMarketRows()` 和 `buildModuleInstallWizardRows()`，让产品化状态由纯 TypeScript helper 输出。
- `registry-state.fixture.ts` 增加 P6.1 helper 编译期 fixture。
- 新增 `scripts/check-module-center-product-entry.sh` 并接入 no-db 模块工具验证。

## P6.1 验收标准

- `scripts/check-module-center-product-entry.sh` 通过。
- `scripts/check-module-center-ui-contract.sh` 通过。
- `scripts/check-module-center-filter-contract.sh` 通过。
- `scripts/check-module-center-manual-checklist.sh` 通过。
- `scripts/check-p5-module-center-freeze.sh` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。

## P6.1 验收结果

- 已通过 `scripts/check-module-center-product-entry.sh`。
- 已通过 `scripts/check-module-tools-no-db.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 已生成更新后的 `module` chunk 和 `demo_notice` chunk。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器人工复验仍受登录态限制；当前不修改 `.env` 或管理员密码。

## 下一步

P6.2：模块详情弹窗或详情页。建议把 P6.1 的选中模块摘要推进为可打开的详情视图，承载 manifest 校验明细、安装计划预览和运行状态说明。

## P6.2：模块详情弹窗

P6.2 把 P6.1 的选中模块摘要推进为可打开的模块详情弹窗，让模块市场入口具备更完整的产品查看面。

## P6.2 当前落地

- 模块中心新增 `模块详情` 弹窗。
- 顶部详情面板新增 `详情` 按钮。
- 模块清单表格的 `详情` 操作会选中当前模块并打开弹窗。
- 弹窗展示模块名称、标识、manifest、表名、安装状态、运行时、快照和页面入口。
- 弹窗复用 `buildModuleInstallWizardRows()` 展示安装向导。
- 弹窗直接展示当前模块的 manifest 校验检查项。
- 页面同时保留 `P6.2`、`P6.1` 和 `P5.25` 阶段标识，保持历史 contract 可回放。
- 新增 `scripts/check-module-center-detail-dialog.sh` 并接入 no-db 模块工具验证。

## P6.2 验收标准

- `scripts/check-module-center-detail-dialog.sh` 通过。
- `scripts/check-module-center-product-entry.sh` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。

## P6.2 验收结果

- 已通过 `scripts/check-module-center-detail-dialog.sh`。
- 已通过 `scripts/check-module-center-product-entry.sh`。
- 已通过 `scripts/check-module-tools-no-db.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 已生成更新后的 `module` chunk 和 `demo_notice` chunk。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器人工复验仍受登录态限制；当前停留在 `http://127.0.0.1:5173/login?redirect=/module`，不修改 `.env` 或管理员密码。

## 方向纠偏

P6.1 和 P6.2 已经让模块中心具备模块市场、模块详情和安装向导雏形，但继续推进模块市场不是当前最直接目标。

当前项目目标应回到通用管理后台框架：

- 让后台第一眼可见，登录后能看到工作台、系统管理、设置、字典、文件、日志和代码生成器。
- 让 AI 后续生成业务模块时有稳定约定，而不是依赖复杂模块市场功能。
- 模块中心保留为开发工具和 manifest 验收入口，不继续扩展成完整商业模块商店。
- 下一阶段优先做框架可用性、生成器模板、基础页面体验和本地人工测试链路。

## 下一步

P6.3：通用后台框架交付面校准。建议停止继续扩展模块市场，转而整理后台首页、基础菜单、核心 CRUD 页面、代码生成器入口和 AI 业务模块生成约定，确保这个框架适合后续 vibe coding 具体业务功能。

## P6.3：通用后台框架交付面校准

P6.3 停止继续扩展模块市场，把工作台重新校准为通用后台框架首页。

## P6.3 当前落地

- 工作台阶段从 `P4.10 P4 冻结验收` 更新为 `P6.3 通用后台框架交付面`。
- 工作台顶部标签更新为 `P5 已冻结` + 当前 P6 阶段。
- 工作台能力卡把 `模块闭环` 收敛为 `AI CRUD scaffold + codegen`。
- 验收状态把 `模块中心` 降级为 `开发工具`，主线强调 `核心页面入口` 和 `通用框架交付`。
- 人工测试入口优先展示代码生成器、菜单、角色、管理员、组织、设置、缓存和日志；模块中心保留为最后的开发工具入口。
- 后端 `GET /api/common/index/console` 返回同样的 P6.3 工作台数据，避免刷新后回到旧 P4 文案。
- 新增 `scripts/check-framework-workbench-contract.sh` 并接入 `scripts/verify-no-db.sh`。

## P6.3 验收标准

- `scripts/check-framework-workbench-contract.sh` 通过。
- `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./admin/service/common ./admin/routers/common` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。

## P6.3 验收结果

- 已通过 `scripts/check-framework-workbench-contract.sh`。
- 已通过 `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./admin/service/common ./admin/routers/common`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 已生成更新后的 `workbench` chunk。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器人工复验仍受登录态限制；当前不修改 `.env` 或管理员密码。

## 下一步

P6.4：生成器业务模块模板体验校准。建议检查代码生成器和 manifest 脚手架输出，让 AI 后续生成具体业务功能时默认得到可用的列表、搜索、编辑、状态字段和 API 文件。
