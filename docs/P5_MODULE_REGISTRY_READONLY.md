# P5 Module Registry Readonly

更新时间：2026-06-03

## 目标

P5.5 把模块中心的内置模块清单从前端硬编码迁移到后端只读 registry，为后续多模块接入提供统一来源。

本阶段不新增数据库写入，不创建业务 schema。

## 后端接口

新增只读接口：

```text
GET /api/gen/moduleRegistry
```

当前 registry 返回 Demo Article：

- `module=article`
- `manifest=examples/demo/manifest.json`
- `table=ma_demo_article`
- `runtime=MAKEADMIN_ENABLE_DEMO_MODULE=1`
- `entry=/demo/article`

## 管理端改动

模块中心进入页面时先读取 registry，再对返回的模块逐项调用 P5.3 状态接口：

```text
POST /api/gen/previewCode/status
```

前端不再把 Demo Article 的清单元数据写死在表格数据中。

## 验收结果

- 已通过 `go test ./generator/service/gen ./generator/routers/gen`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过浏览器人工验证 `/module` 显示 `P5.5`。
- 已通过浏览器人工验证模块中心显示 registry 返回的 `Demo Article`、`examples/demo/manifest.json`、`MAKEADMIN_ENABLE_DEMO_MODULE=1`。
- 已通过浏览器人工验证状态回读仍显示 `已安装`、`权限 5/5` 和 `角色授权 5/5`。
- `npm run build` 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。

## 保留边界

P5.5 不做：

- 不从数据库读取模块 registry。
- 不新增模块安装市场。
- 不创建或迁移 `ma_demo_article`。
- 不把写入 env 写进 `.env`。
- 不连接 zyai 真实业务库。
- 不处理 PTLM 业务模块。
