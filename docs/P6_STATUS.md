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

## P6.4：生成器业务模块模板体验校准

P6.4 把 P6 主线从复杂模块市场继续收敛到通用后台框架的业务生成底座。

## P6.4 当前落地

- 前端生成模板 contract 明确覆盖列表、搜索、重置、新增、编辑、删除、分页、字典展示和弹窗表单。
- 生成器前端 smoke 使用独立 `codegen_template_smoke` 临时模块，不再和仓库已有 `article` 示例模块冲突。
- 默认 manifest codegen 字段仍保持 `id/title/status`，但 `status` 改为 `common_status` 字典单选，符合常见后台业务起步形态。
- 新增 `scripts/check-codegen-business-template-contract.sh`，验证模板关键能力和默认 codegen 字段。
- `scripts/check-module-tools-no-db.sh` 接入业务模板 contract 和生成前端 type-check。

## P6.4 验收标准

- `scripts/check-codegen-business-template-contract.sh` 通过。
- `scripts/check-codegen-frontend.sh` 通过。
- `scripts/check-module-codegen-plan.sh` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。

## P6.4 验收结果

- 已通过 `scripts/check-codegen-business-template-contract.sh`。
- 已通过 `scripts/check-codegen-frontend.sh`。
- 已通过 `scripts/check-module-codegen-plan.sh`。
- 已通过 `scripts/check-module-tools-no-db.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。

## 下一步

P6.5：通用后台人工可见链路验收。建议启动本地 API 和管理端，按“登录、工作台、系统管理、设置、字典、文件、日志、代码生成器”顺序做一次真实浏览器验收，重点找阻塞后台第一眼可用的问题。

## P6.5：通用后台人工可见链路验收

P6.5 从“框架能不能被真实打开”出发，验收本地服务、P1 种子、核心 API 和浏览器登录页。

## P6.5 当前落地

- 本地管理端 `http://127.0.0.1:5173` 可访问。
- 本地 Go API `http://127.0.0.1:18000/api` 可访问，`8000` 当前由 nginx 占用，不作为本轮主验收入口。
- `scripts/check-services.sh` 确认 MySQL、Redis 和 `go_makeadmin` 可用。
- `scripts/check-p1-seed.sh` 确认 `ma_*` 表、admin、super_admin、权限、菜单、设置、字典和文件分类种子可用。
- 本地开发默认账号已统一为 `admin / 123456`；`scripts/init-p1-db.sh` 会在未指定 `ADMIN_PASSWORD` 时生成该密码的 bcrypt hash。
- P6.4 默认 `common_status` 字典已补进 P1 seed，生成器默认状态字段能直接拿到 `1=启用`、`0=禁用`。
- 使用 `admin / 123456` 验证登录 API 返回 token。
- 浏览器已使用 `admin / 123456` 正常登录，并按 redirect 进入 `/module`。
- 浏览器已打开核心后台页面：工作台、管理员、角色、菜单、部门、字典、存储、系统日志、代码生成器和模块中心。

## P6.5 验收标准

- `scripts/check-services.sh` 通过。
- `scripts/check-p1-seed.sh` 通过。
- `scripts/check-local-dev-login-contract.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 核心只读 API 返回成功。
- 浏览器能使用 `admin / 123456` 登录。
- 登录后能打开工作台、系统管理、设置、字典、日志、代码生成器和模块中心页面。

## P6.5 验收结果

- 已通过 `scripts/check-services.sh`。
- 已通过 `scripts/check-p1-seed.sh`。
- 已通过 `scripts/check-local-dev-login-contract.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过核心只读 API 检查：工作台 `code=200`，自身信息 `admin`，菜单 `24` 条，租户 `1` 条，字典类型 `5` 条，`common_status` 字典项 `2` 条，代码生成器列表、网站设置、存储设置和登录日志均返回成功。
- 已通过 `admin / 123456` 登录 API 检查。
- 已通过浏览器登录检查，登录后进入 `http://127.0.0.1:5173/module`。
- 已通过浏览器核心页面检查：`/workbench`、`/admin`、`/role`、`/menu`、`/department`、`/dict`、`/storage`、`/journal`、`/code`、`/module` 均可打开，未发现 404 或 API 错误。
- 登录页截图保存到 `/tmp/go-makeadmin-p65-login.png`，工作台截图保存到 `/tmp/go-makeadmin-p65-workbench.png`；核心页面验收快照保存到 `/tmp/go-makeadmin-p65-actual-pages.json`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。

## 下一步

P6.6：菜单层级与工作台入口体验修整。建议把当前扁平化实际路由和菜单层级展示校准清楚，让后台侧边栏更像通用管理系统，而不是一串顶级页面。

## P6.6：菜单层级与目录空白页修整

P6.6 修复核心后台菜单被拍平成顶级页面的问题，并让目录路由自动进入第一个真实子页面。

## P6.6 当前落地

- 修复 `util.ArrayUtil.ListToTree`：父节点没有预置 `children` 字段时也会正确挂载子节点。
- `GET /api/system/menu/route` 返回的权限管理、组织管理、系统设置、开发工具等菜单恢复层级结构。
- 前端目录路由增加 redirect：直接访问 `/setting`、`/permission`、`/dev_tools` 等目录路径时，会跳到第一个可见子页面；嵌套目录使用完整绝对路径，避免 `/setting/website` 被错误跳到 `/information`。
- 新增 `scripts/check-menu-tree-contract.sh`，验证菜单树结构和目录 redirect 契约。
- `scripts/verify-no-db.sh` 已接入菜单树契约检查。

## P6.6 验收标准

- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/check-menu-tree-contract.sh` 通过。
- `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./util ./makeadmin/adapter` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 浏览器打开 `/setting` 不再停留在空白目录页。

## P6.6 验收结果

- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/check-menu-tree-contract.sh`。
- 已通过 `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./util ./makeadmin/adapter`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过本地菜单接口检查：`/api/system/menu/route` 顶层菜单为 `6` 个，`系统设置` 子菜单为 `3` 个，`开发工具` 子菜单为 `4` 个。
- 已通过浏览器目录路由检查：`/setting` -> `/setting/website/information`，`/permission` -> `/permission/admin`，`/organization` -> `/organization/department`，`/dev_tools` -> `/dev_tools/dict`，`/setting/website` -> `/setting/website/information`。
- 修复后截图保存到 `/tmp/go-makeadmin-p66-setting-fixed.png`。

## 下一步

P6.7：后台第一屏和侧边栏人工体验验收。建议继续用浏览器走一遍登录、展开菜单、点击目录和子页面、刷新目录路径，修掉真实使用中还会看到的空白页、缺组件页或入口命名问题。

## P6.7：后台第一屏和侧边栏人工体验验收

P6.7 把“菜单能点开”从单点修复推进到核心菜单全量可见验收，重点防止页面组件缺失、目录刷新 404 和第一眼控制台噪音。

## P6.7 当前落地

- 浏览器已直达验收核心页面：工作台、管理员、角色、菜单、部门、岗位、素材管理、网站信息、网站备案、政策协议、存储设置、系统环境、系统缓存、系统日志、字典管理、代码生成器、模块中心和 Demo Article。
- `素材管理` 页面不是路由空白，当前为空数据状态，可见图片/视频切换、分组、本地上传和空态分页。
- 修复 `resetRouter()` 无条件移除动态根路由导致的 Vue Router warning。
- 新增 `scripts/check-admin-route-components.py`，从 `sql/p1.seed.sql` 解析核心菜单并检查页面组件是否存在。
- `scripts/verify-no-db.sh` 已接入后台路由组件契约检查。

## P6.7 验收标准

- `scripts/check-admin-route-components.py` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 浏览器核心页面直达不出现空白页或 404。

## P6.7 验收结果

- 已通过 `scripts/check-admin-route-components.py`，P1 seed 当前覆盖 `17` 个核心页面和 `6` 个目录。
- 已通过浏览器核心页面直达验收，18 个当前本地可见页面均非空白且非 404。
- 当前本地库比 P1 seed 多出的 `Demo Article` 来自模块安装链路，已在浏览器直达验收中覆盖。

## 下一步

P6.8：素材管理空态和上传入口体验修整。建议把素材页从“空数据能打开”继续调整到更像通用后台的文件管理入口，包括空态提示、上传入口、分组默认态和失败态说明。
