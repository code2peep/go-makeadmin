# P2 Module Runtime Registry

更新时间：2026-06-02

## 目标

P2.15 建立模块运行时注册闭环，让模块可以从 manifest、SQL 预览推进到后端路由可控挂载。

本阶段先接入 demo article 模块，用来验证路由注册边界，不创建 demo 表，不执行数据库写入。

## 开关

demo runtime 模块默认关闭。

启用方式：

```bash
MAKEADMIN_ENABLE_DEMO_MODULE=1
```

未设置该环境变量时，`server/modules/routers` 不返回任何 demo 路由。

## 运行时路由

启用后挂载以下后端 API：

```text
GET  /api/article/list
GET  /api/article/detail
POST /api/article/add
POST /api/article/edit
POST /api/article/del
```

这些路径和 `examples/demo/manifest.json` 中 `backend.routes` 保持一致。

## 安全规则

- demo runtime 模块默认关闭。
- demo article 路由仍使用 `middleware.TokenAuth()`。
- 未登录访问会返回 token empty 响应，不会绕过认证。
- list/detail 为只读示例响应。
- add/edit/del 返回只读失败响应。
- 不创建、修改或迁移 schema。
- 不连接或写入真实 zyai 业务库。

## 不在 P2.15 做

- 不创建 `ma_demo_article` 表。
- 不把 demo 菜单设为可见。
- 不自动给角色授权。
- 不开放无认证 demo API。
- 不做前端页面动态注册。

## 验证

P2.15 需要通过：

```bash
cd server
GOCACHE=/private/tmp/go-makeadmin-gocache go test ./modules/routers
GOCACHE=/private/tmp/go-makeadmin-gocache go test ./...
cd ..
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
