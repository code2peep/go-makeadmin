# P4 Status

更新时间：2026-06-03

## 当前阶段

P4：可见后台与人工测试闭环。

P4 从 P3 冻结面继续推进，重点是把已完成的底座能力放到后台可见页面中验收，让管理端不只是能构建，也能被人工打开、点击和判断。

## P4.1 当前落地

框架工作台已产品化：

- 后端 `GET /api/common/index/console` 从旧业务假指标改为框架状态数据。
- 管理端工作台不再显示访问量、销售额、订单量等蓝本假数据。
- 工作台展示当前阶段、数据库边界、认证权限、模块闭环和本地验收状态。
- 工作台新增人工测试入口，直接进入代码生成器、菜单权限、角色、管理员、组织部门和网站信息。
- 人工测试入口使用当前动态菜单真实路由：`/code`、`/menu`、`/role`、`/admin`、`/department`、`/information`。
- 已通过浏览器人工验证登录、工作台、代码生成器和 manifest 预览弹窗。

详见 `docs/P4_VISIBLE_ADMIN_WORKBENCH.md`。

## P4.1 验收标准

- `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./admin/service/common ./admin/routers/common` 通过。
- `cd admin && npm run type-check` 通过。
- `cd admin && npm run build` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 浏览器能打开 `http://127.0.0.1:5173/workbench`。
- 管理端能登录并看到框架工作台。
- 工作台“代码生成器”入口能跳转到 `/code`。
- 代码生成器“Manifest 预览”能打开弹窗并生成 demo manifest 预览。
- 不改数据库 schema、不改 `.env`、不新增权限 SQL、不连接真实 zyai 业务库。

## P4.1 验收结果

- 已通过 `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./admin/service/common ./admin/routers/common`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `cd admin && npm run build`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 已通过浏览器人工验证 `http://127.0.0.1:5173/workbench`。
- 已通过浏览器人工验证本地 admin 登录。
- 已通过浏览器人工验证工作台“代码生成器”入口跳转到 `/code`。
- 已通过浏览器人工验证代码生成器“Manifest 预览”弹窗和 demo manifest 预览生成。
- 本阶段没有改数据库 schema、没有改 `.env`、没有新增权限 SQL、没有连接真实 zyai 业务库。

## 下一步

## P4.2 当前落地

模块中心页面骨架已建立：

- 新增 `admin/src/views/dev_tools/module/index.vue`。
- 工作台新增 `模块中心` 人工测试入口。
- P1 seed 在 `开发工具` 下新增 `模块中心` 菜单。
- P1 seed 新增 `module:center:view` 权限。
- 当前本地 `go_makeadmin` 开发库已同步补入模块中心菜单和权限，用于人工验证。
- 模块中心复用 `module-manifest-preview.vue`，可以打开 manifest 预览并生成 demo manifest 预览结果。

详见 `docs/P4_MODULE_CENTER.md`。

## P4.2 验收标准

- `cd admin && npm run type-check` 通过。
- `cd admin && npm run build` 通过。
- `./scripts/check-p1-seed.sh` 通过，确认当前本地库菜单和权限 seed 合法。
- 使用临时库运行 `scripts/init-p1-db.sh` 通过，确认 `sql/p1.seed.sql` 可完整导入。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 浏览器能打开 `http://127.0.0.1:5173/module`。
- 模块中心能打开 `Manifest 预览` 弹窗并生成 demo manifest 预览。
- 不改数据库 schema、不改 `.env`、不连接真实 zyai 业务库。

## P4.2 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `cd admin && npm run build`。
- 已通过 `./scripts/check-p1-seed.sh`，当前本地库 permission seed count 为 80，menu seed count 为 23。
- 已通过临时库 `scripts/init-p1-db.sh` 完整初始化，并确认临时库清理后残留数为 0。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 已通过浏览器人工验证 `http://127.0.0.1:5173/module`。
- 已通过浏览器人工验证模块中心 `Manifest 预览` 弹窗和 demo manifest 预览生成。
- 本阶段没有改数据库 schema、没有改 `.env`、没有连接真实 zyai 业务库。

## 下一步

P4.3：模块中心预览结果页面化。建议把 manifest 预览结果从弹窗状态抽到模块中心主页面，形成模块详情、安装计划、代码预览和审计预览的稳定工作区。

## P4.3 当前落地

模块中心预览结果已页面化：

- 模块中心新增 manifest 输入区，支持仓库路径和 JSON body 两种来源。
- 生成预览后，页面内展示模块来源、实体、表名、功能名、模板和运行时开关。
- 页面内展示字段明细，包括数据库字段、Go 字段、Go 类型、表单、查询和字典。
- `安装计划` 从当前页面预览结果打开 registry、role grant、install 和 uninstall SQL。
- `代码预览` 从当前页面预览结果打开后端和前端生成代码。
- 内置 `Demo Article` 清单的 `预览` 会直接写入 manifest 路径并生成页面内预览。

详见 `docs/P4_MODULE_CENTER_INLINE_PREVIEW.md`。

## P4.3 验收标准

- `cd admin && npm run type-check` 通过。
- `cd admin && npm run build` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 浏览器能打开 `http://127.0.0.1:5173/module`。
- 模块中心能在页面内生成 demo manifest 预览。
- 模块中心能从页面内预览打开安装计划。
- 模块中心能从页面内预览打开代码预览。
- 不改数据库 schema、不改 `.env`、不连接真实 zyai 业务库。

## P4.3 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `cd admin && npm run build`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- 已通过浏览器人工验证 `http://127.0.0.1:5173/module`。
- 已通过浏览器人工验证页面内生成 `DemoArticle`、`ma_demo_article` 和 `MAKEADMIN_ENABLE_DEMO_MODULE=1`。
- 已通过浏览器人工验证 `安装计划` 弹窗展示 `registry.sql`、`role_grant.sql`、`install.sql`、`uninstall.sql`。
- 已通过浏览器人工验证 `代码预览` 弹窗展示 `gocode/model.go`、`gocode/route.go`、`gocode/schema.go`、`api.ts`、`index.vue`、`edit.vue`。
- `npm run build` 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有改数据库 schema、没有改 `.env`、没有连接真实 zyai 业务库。

## 下一步

P4.4：模块中心安装/卸载 apply 结果内嵌化。建议把安装执行、卸载执行、apply 结果摘要和审计预览从弹窗迁到模块中心页面状态。

## P4.4 当前落地

模块中心 apply 结果已内嵌化：

- 模块中心复用已有 install/uninstall apply API，不新增后端接口。
- 预览结果区新增 `确认模块`、`安装写入`、`Schema 风险`、`删除确认`。
- 预览结果区新增 `安装执行` 和 `卸载执行`。
- 安装和卸载结果以内嵌 tabs 展示。
- apply 结果复用 `module-manifest-apply-result.vue`，包含状态、环境变量、权限、快照、检查项和审计预览。
- manifest 输入变化时会清空旧预览和旧 apply 结果，避免展示过期结果。

详见 `docs/P4_MODULE_CENTER_APPLY_RESULT.md`。

## P4.4 验收标准

- `cd admin && npm run type-check` 通过。
- `cd admin && npm run build` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 浏览器能打开 `http://127.0.0.1:5173/module`。
- 模块中心能在页面内生成 demo manifest 预览。
- 勾选 `安装写入` 后，`安装执行` 能在页面内展示 apply 结果。
- 安装结果能展开 `审计预览`。
- 勾选 `删除确认` 后，`卸载执行` 能在页面内展示 apply 结果。
- 不开启写入 env、不实际写入数据库、不改数据库 schema、不改 `.env`、不连接真实 zyai 业务库。

## P4.4 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `cd admin && npm run build`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- 已通过浏览器人工验证 `http://127.0.0.1:5173/module`。
- 已通过浏览器人工验证 demo manifest 预览生成。
- 已通过浏览器人工验证 `安装执行` 返回页面内 `安装结果`。
- 已通过浏览器人工验证 `安装结果` 的 `审计预览` 展开。
- 已通过浏览器人工验证 `卸载执行` 返回页面内 `卸载结果`。
- 本地未开启 `MAKEADMIN_ALLOW_MODULE_INSTALL_APPLY=1` 和 `MAKEADMIN_ALLOW_MODULE_UNINSTALL_APPLY=1`，安装和卸载 apply 均返回门禁阻断结果，未访问数据库。
- `npm run build` 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有改数据库 schema、没有改 `.env`、没有连接真实 zyai 业务库。

## 下一步

P4.5：模块中心安装状态探测与测试清单。建议增加模块当前安装状态、门禁条件状态和人工测试步骤，让模块中心成为早期框架验收的主操作台。
