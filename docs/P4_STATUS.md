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

## P4.5 当前落地

模块中心新增状态探测和人工测试清单：

- `模块状态` 展示预览、安装、卸载、写入门禁和运行时。
- `人工测试清单` 展示 Manifest 预览、安装计划、代码预览、安装执行、卸载执行和审计预览。
- 打开安装计划后，清单标记为 `已打开`。
- 打开代码预览后，清单标记为 `已打开`，并显示生成文件数量。
- install/uninstall apply 返回后，清单显示 blocked/applied 等状态和结果说明。
- 写入门禁 env 以可换行 tag 展示。
- 运行时开关、描述 label 和表格长文本已做换行/宽度处理，避免窄视口下文字互相挤压。

详见 `docs/P4_MODULE_CENTER_STATUS_CHECKLIST.md`。

## P4.5 验收标准

- `cd admin && npm run type-check` 通过。
- `cd admin && npm run build` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 浏览器能打开 `http://127.0.0.1:5173/module`。
- 生成 demo manifest 预览后能看到 `模块状态`。
- 生成 demo manifest 预览后能看到 `人工测试清单`。
- 打开安装计划和代码预览后，清单状态能更新。
- 安装/卸载 apply 返回门禁阻断后，模块状态和清单能显示 blocked 与对应 env。
- 窄视口下写入门禁 env、运行时和描述 label 不出现明显挤压。
- 不开启写入 env、不实际写入数据库、不改数据库 schema、不改 `.env`、不连接真实 zyai 业务库。

## P4.5 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `cd admin && npm run build`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- 已通过浏览器人工验证 `http://127.0.0.1:5173/module`。
- 已通过浏览器人工验证 `模块状态` 和 `人工测试清单`。
- 已通过浏览器人工验证安装计划、代码预览、安装执行、卸载执行的清单状态更新。
- 已通过浏览器人工验证写入门禁 env tag、运行时和描述 label 的换行/宽度样式。
- 本地未开启 `MAKEADMIN_ALLOW_MODULE_INSTALL_APPLY=1` 和 `MAKEADMIN_ALLOW_MODULE_UNINSTALL_APPLY=1`，安装和卸载 apply 均返回门禁阻断结果，未访问数据库。
- `npm run build` 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有改数据库 schema、没有改 `.env`、没有连接真实 zyai 业务库。

## 下一步

P4.6：核心管理页可见验收。建议检查系统设置、菜单、角色、管理员、部门等后台基础页面，优先修正明显蓝本残留、空态和第一屏可见性问题。

## P4.6 当前落地

工作台新增核心页面验收矩阵：

- 后端 `GET /api/common/index/console` 返回 `corePages`。
- 工作台阶段更新为 `P4.6 核心管理页可见验收`。
- 工作台新增 `核心页面验收` 表格，包含页面、状态、范围、路由和入口。
- 工作台验收状态新增 `模块中心` 和 `核心页面入口`。
- 核心页面验收覆盖 `/menu`、`/role`、`/admin`、`/department`、`/information`、`/cache`、`/journal`。

详见 `docs/P4_CORE_ADMIN_VISIBLE_CHECK.md`。

## P4.6 验收标准

- `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./admin/service/common ./admin/routers/common` 通过。
- `cd admin && npm run type-check` 通过。
- `cd admin && npm run build` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 浏览器能打开 `http://127.0.0.1:5173/workbench`。
- 工作台能显示 `核心页面验收`。
- `/menu`、`/role`、`/admin`、`/department`、`/information`、`/cache`、`/journal` 均可打开且不是 403/404/空白页。
- 不改数据库 schema、不改 `.env`、不连接真实 zyai 业务库。

## P4.6 验收结果

- 已通过 `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./admin/service/common ./admin/routers/common`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `cd admin && npm run build`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets .cache`。
- 已通过浏览器人工验证 `http://127.0.0.1:5173/workbench` 显示 `核心页面验收` 和 `P4.6`。
- 已通过浏览器人工验证 `/menu` 显示 `菜单名称`。
- 已通过浏览器人工验证 `/role` 显示 `管理员人数`。
- 已通过浏览器人工验证 `/admin` 显示 `管理员账号`。
- 已通过浏览器人工验证 `/department` 显示 `部门名称`。
- 已通过浏览器人工验证 `/information` 显示 `网站名称`。
- 已通过浏览器人工验证 `/cache` 显示 `基本信息`。
- 已通过浏览器人工验证 `/journal` 显示 `访问链接`。
- `npm run build` 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有改数据库 schema、没有改 `.env`、没有连接真实 zyai 业务库。

## 下一步

P4.7：核心页面细节修整。建议优先处理管理员、部门、日志页面中的格式不一致、蓝本残留、空态和窄视口表格可读性。
