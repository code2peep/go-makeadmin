# P5 Module Registry Manifest Check

更新时间：2026-06-03

## 目标

P5.6 让后端 registry 返回模块清单时同步校验 manifest，避免模块中心展示一个已经和 manifest 脱节的模块入口。

本阶段不新增数据库写入，不创建业务 schema。

## 后端改动

`GET /api/gen/moduleRegistry` 返回项新增：

- `manifestStatus`
- `manifestMessage`
- `manifestChecks`

当前校验项：

- manifest 能读取。
- manifest 结构通过现有校验。
- registry `module` 与 manifest `module` 一致。
- registry `table` 与 manifest `table` 一致。
- registry `runtime` 与 manifest runtime hint 一致。
- registry `entry` 是绝对管理端路由。
- manifest 菜单 `routeName`、`routePath`、`component` 完整。

如果单个模块校验失败，registry 仍返回该模块，并把该行标记为 `failed`，避免整个模块中心空白。

## 管理端改动

模块中心新增 `校验` 列：

- `已通过`：registry 与 manifest 一致。
- `异常`：registry 与 manifest 不一致或 manifest 无法读取。

`异常` 筛选同时包含 registry 校验异常和安装状态异常。

## 验收结果

- 已通过 `go test ./generator/service/gen ./generator/routers/gen`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，浏览器登录页 token 已过期，因此模块中心页面人工复验需要重新登录后执行。

## 保留边界

P5.6 不做：

- 不从数据库读取模块 registry。
- 不新增 registry 写入接口。
- 不创建或迁移 `ma_demo_article`。
- 不把写入 env 写进 `.env`。
- 不连接 zyai 真实业务库。
- 不处理 PTLM 业务模块。
