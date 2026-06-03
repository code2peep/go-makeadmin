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
