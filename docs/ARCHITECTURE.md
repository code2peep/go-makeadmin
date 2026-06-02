# go-makeadmin 架构说明

## 定位

`go-makeadmin` 是 `makeadmin` 的 Go 版本后台基础框架，面向多项目复用。项目早期参考 LikeAdmin Go 的 MIT 代码，当前 P1 已冻结，默认运行链路已切到 `ma_*` 自研系统表，旧 `la_*` 代码只作为蓝本参考保留。

## 当前结构

```text
go-makeadmin/
├── admin/      # Vue 管理端
├── frontend/   # LikeAdmin 原发布目录，保留为蓝本资料
├── legacy/     # 历史来源和迁移参考资料
├── public/     # 静态资源和上传目录
├── server/     # Go API 服务
└── sql/        # P1 ma_* 初始化 SQL 和蓝本 SQL
```

## P1 运行模型

- 默认开发库是 `go_makeadmin`，系统表使用 `ma_*` 前缀。
- 已接管链路包括登录、管理员、角色、菜单、部门、岗位、设置、字典、文件、日志、公共首页和代码生成器。
- 后台认证中间件只接受 `makeadmin:token:*` 新 token，不再回退旧 `backstage:*` token。
- 操作审计固定写入 `ma_audit_log`，登录日志写入 `ma_login_log`。
- `server/admin/service/*` 和 `server/model/{system,setting,common}` 暂作参考代码保留，不作为 P1 核心运行兜底。
- `sql/install.sql` 和 `sql/install.core.sql` 是蓝本初始化资料，不代表框架默认 schema。

## P1 冻结标准

- `./scripts/verify-no-db.sh` 通过。
- `./scripts/check-p1-seed.sh` 通过。
- 本地 disposable P1 库上 `scripts/p1-smoke.py` 覆盖矩阵通过。
- 新增 P1 功能不得直接读写 `la_*` 表。
- 文档默认入口必须指向 P1 独立库和 `ma_*` 模型。
- P1 最终状态记录在 `docs/P1_FINAL_STATUS.md`。

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

## P2 优先级

1. 认证模型升级为 JWT + Redis session state。
2. 建立租户上下文 middleware 和租户数据边界。
3. 将角色数据范围落到查询层。
4. 让代码生成器生成可运行模块。
5. 建立 `examples/demo` 标准模块，验证框架扩展约定。
