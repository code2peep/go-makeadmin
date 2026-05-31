# go-makeadmin 架构说明

## 定位

`go-makeadmin` 是 `makeadmin` 的 Go 版本后台基础框架，面向多项目复用。P0 阶段参考 LikeAdmin Go 的 MIT 代码，先完成现代化升级和最小可运行验证；后续逐步重构为自研架构。

## 当前结构

```text
go-makeadmin/
├── admin/      # Vue 管理端
├── frontend/   # LikeAdmin 原发布目录，P0 阶段保留
├── public/     # 静态资源和上传目录
├── server/     # Go API 服务
└── sql/        # LikeAdmin 原始初始化 SQL
```

## P0 已确定方向

- 后端 module 改为 `go-makeadmin`。
- 前端包名改为 `go-makeadmin-admin`。
- 默认响应结构未来向 `code=0/msg/data` 收敛，P0 暂不批量重写业务响应。
- MySQL 和 Redis 改为懒加载，避免包加载时强制连接基础设施。
- 默认 `npm run build` 只构建 `admin/dist`，发布复制改为 `npm run build:release`。

## 后续重构方向

后端目标结构：

```text
server/
├── cmd/api/
├── internal/
│   ├── config/
│   ├── database/
│   ├── http/
│   ├── middleware/
│   ├── permission/
│   ├── tenant/
│   ├── storage/
│   └── audit/
└── modules/
    ├── system/
    └── demo/
```

前端目标结构：

```text
admin/src/
├── api/
├── layouts/
├── modules/
├── router/
├── stores/
└── utils/
```

## 必须保留的框架能力

- 登录、登出、会话管理。
- 菜单、按钮、API 权限。
- 角色数据范围。
- 多租户上下文。
- 文件上传抽象。
- 操作审计日志。
- 后台 CRUD 模块约定。

## 不继承的能力

- 不继承 Gin-Vue-Admin 框架层代码。
- 不沿用 LikeAdmin 的强全局初始化模式。
- 不把 Redis token 作为长期认证模型，后续改为 JWT + Redis 会话状态。
