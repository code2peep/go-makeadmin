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
