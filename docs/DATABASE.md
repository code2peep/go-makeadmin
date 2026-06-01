# go-makeadmin 数据库边界

## 当前状态

P0 阶段不连接 zyai 现有业务库做写操作。

当前本机已按 P0.7 初始化一次性开发库 `go_makeadmin`，用于验证登录、菜单、权限和工作台链路。LikeAdmin 原始 SQL 保留在 `sql/install.sql`，当前只作为 P0 蓝本初始化数据，P1 会切换为 `ma_*` 自研系统表。

P1 阶段已明确：`go-makeadmin` 使用独立数据库，不和 zyai 业务库混用；独立库最终命名沿用 `go_makeadmin`，采用一次性初始化或迁移方式，不保留长期双写。

## 默认开发库

当前默认 DSN 指向：

```text
root:@tcp(127.0.0.1:3306)/go_makeadmin?charset=utf8mb4&parseTime=True&loc=Local
```

这是开发默认值。当前本机数据库已创建，其他环境仍需按 `docs/DB_INIT_PLAN.md` 初始化。

## 完整联调前置条件

启动 API 并验证登录链路前，需要先满足：

- MySQL 服务可连接。
- 开发库 `go_makeadmin` 已存在。
- 当前 P0 蓝本需要 `la_*` 系统表和初始化数据。
- Redis 服务可连接。

可先运行只读检查：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/check-services.sh
```

该脚本只检查服务和数据库是否存在，不创建数据库、不导入 SQL、不写入 Redis。

本地烟测账号：

```text
P0 蓝本：
username: admin
password: 123456

P1 ma_*：
username: admin
password: 由 scripts/init-p1-db.sh 执行时的 ADMIN_PASSWORD 决定
```

P0.7 的一次性本地库方案见：

```text
docs/DB_INIT_PLAN.md
```

初始化后可运行只读种子检查：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/check-db-seed.sh
```

## 命名规划

系统表建议使用：

```text
ma_*
```

P1 的系统表草案见：

```text
docs/P1_SCHEMA_PLAN.md
sql/p1.schema.sql
sql/p1.seed.sql
```

P1 独立库种子只读检查：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/check-p1-seed.sh
```

业务表由具体项目决定。例如 zyai 迁移时，优先兼容现有业务表：

```text
tenants
organizations
students
screening_projects
```

## 迁移原则

- 开发阶段可以用显式 SQL 或迁移工具。
- 生产阶段不依赖隐式 AutoMigrate。
- 系统表和业务表分离。
- 不复用 Gin-Vue-Admin 的 `sys_*` 表，避免授权和结构耦合。
- 从 GVA 迁移用户、角色、菜单、权限前，必须先出映射表和回滚方案。

## 红线

以下操作必须先确认：

- 创建新数据库。
- 执行 `sql/install.sql`。
- 修改现有 zyai 业务表结构。
- 向 zyai 业务库写入数据。
- 从 GVA `sys_*` 表迁移数据。
