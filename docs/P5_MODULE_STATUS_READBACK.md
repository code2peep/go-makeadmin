# P5 Module Status Readback

更新时间：2026-06-03

## 目标

P5.3 让模块中心直接读取 Demo Article 的真实安装状态，减少靠人工对比 apply 结果和数据库快照的判断成本。

本阶段只读本地 `go_makeadmin` 开发库，不创建业务 schema，不写 `.env`。

## 后端接口

新增只读接口：

```text
POST /api/gen/previewCode/status
```

请求参数复用模块 manifest 预览参数：

- `manifestPath`
- `manifestBody`
- `tenantId`
- `roleId`
- `authorName`

响应包含：

- `status`：`installed`、`partial`、`uninstalled`、`blocked` 或 `failed`。
- `snapshot`：当前权限、菜单、菜单权限、角色授权数量。
- `expected`：manifest 期望的权限、菜单、菜单权限、角色授权数量。
- `missing`：当前快照相对期望快照缺失的数量。
- `runtimeEnv`：例如 `MAKEADMIN_ENABLE_DEMO_MODULE`。
- `runtimeEnabled`：当前运行时 env 是否开启。
- `runtimeRegistered`：manifest 是否声明已注册运行时。
- `menuVisible`：数据库菜单是否可见。

状态接口只读，但仍要求目标库是本地 `go_makeadmin`，防止误连业务库做状态探测。

## 管理端改动

模块中心内置模块清单新增：

- `刷新状态`。
- `安装` 状态列。
- `快照` 列。
- `运行时状态` 列。

进入模块中心时会自动读取状态。安装 apply 或卸载 apply 后会自动刷新状态。

## 验收结果

- 已通过 `go test ./generator/service/gen ./generator/routers/gen`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `cd admin && npm run build`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 本地 API 已用 `MAKEADMIN_ENABLE_DEMO_MODULE=1 MAKEADMIN_ALLOW_MODULE_INSTALL_APPLY=1 MAKEADMIN_ALLOW_MODULE_UNINSTALL_APPLY=1 ./scripts/dev-api.sh` 启动。
- 已通过浏览器人工验证 `/module` 显示 `P5.3`、`已安装`、`权限 5/5`、`角色授权 5/5`、`已开启` 和 `MAKEADMIN_ENABLE_DEMO_MODULE=1`。
- 已通过浏览器人工验证卸载 apply 后状态变为 `未安装`，快照显示 `权限 0/5` 和 `角色授权 0/5`。
- 已通过浏览器人工验证重新安装 apply 后状态恢复为 `已安装`，快照显示 `权限 5/5` 和 `角色授权 5/5`。
- `npm run build` 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。

## 保留边界

P5.3 不做：

- 不创建或迁移 `ma_demo_article`。
- 不把写入 env 写进 `.env`。
- 不连接 zyai 真实业务库。
- 不处理 PTLM 业务模块。
