# P2 Module Registry

更新时间：2026-06-02

## 目标

P2.11 建立模块注册清单规范，让后续模块接入时统一描述后端路由、前端页面、菜单节点、权限和初始化数据需求。

本阶段只做规范和校验，不执行数据库写入。

## Manifest 位置

每个示例模块必须提供：

```text
examples/<module>/manifest.json
```

## 必填字段

- `version`：当前为 `1`。
- `module`：模块名，对应生成器模块名。
- `entity`：实体名。
- `table`：表名。
- `backendPackage`：生成 Go 包名。
- `backend.routes`：后端路由列表。
- `frontend.api`：前端 API 文件路径。
- `frontend.views`：前端页面路径列表。
- `menu`：菜单节点声明。
- `permissions`：权限元数据列表。
- `runtimeRegistered`：是否默认注册进运行时。
- `requiresSchema`：是否需要 schema。

## 路由规则

`backend.routes` 每项必须包含：

- `method`：`GET`、`POST`、`PUT`、`PATCH`、`DELETE` 之一。
- `path`：以 `/` 开头。
- `permission`：必须存在于 `permissions`。

## 权限规则

权限编码采用二段或三段冒号格式，例如：

```text
article:list
system:admin:list
```

每个权限必须包含：

- `code`
- `name`
- `module`
- `resource`
- `action`

## 菜单规则

`menu` 必须包含：

- `code`
- `parent`
- `type`：`catalog`、`page`、`button` 之一。
- `name`
- `routePath`
- `routeName`
- `component`
- `permission`
- `visible`
- `sort`

`menu.permission` 必须存在于 `permissions`。

## 校验

```bash
python3 scripts/check-module-manifests.py
```

校验内容包括：

- manifest JSON 可解析。
- 必填字段存在且类型正确。
- 后端路由不重复。
- 路由权限和菜单权限均已声明。
- 前端路径位于 `admin/src/api` 或 `admin/src/views`。
- 权限编码格式正确。

## 不在 P2.11 做

- 不把 manifest 自动写入数据库。
- 不创建、修改或迁移 schema。
- 不写入菜单、权限或角色授权种子。
- 不默认注册 demo 运行时路由。
