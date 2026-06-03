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

P4.2：模块中心页面骨架。建议把 P3 的 manifest、安装计划、安装执行、卸载执行和审计预览能力从代码生成器弹窗中抽成独立模块中心入口，让后台用户不必先进代码生成器才能管理模块。
