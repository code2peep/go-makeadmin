# P0.7 本地开发库初始化方案

## 结论

P0.7 只为本地烟测准备一次性开发库方案，不把 LikeAdmin 的 `la_*` 表当成 `go-makeadmin` 的最终模型。

当前推荐：

- P0 登录、菜单、权限、上传烟测：临时使用 `sql/install.sql` 中的 `la_*` 蓝本表。
- P1 自研基础框架：重新设计 `ma_*` 系统表、认证、权限、多租户和数据范围模型。

不建议现在立刻把 `la_*` 批量改成 `ma_*`。原因是当前服务层、模型、菜单、权限、初始化数据仍然按 LikeAdmin 结构工作，单纯改表名前缀不能形成自研框架，反而会把 P0 验证和 P1 权限模型重写混在一起。

## 初始化边界

以下事情属于红线操作，必须单独确认后执行：

- 创建数据库。
- 执行 `sql/install.sql`。
- 导入、删除、覆盖任何表。
- 连接 zyai 业务库做写操作。

P0.7 只新增 dry-run 默认脚本和只读检查脚本，不实际执行 schema 操作。

## 一次性本地库策略

目标库：

```text
go_makeadmin
```

目标用途：

- 仅用于本机 P0 烟测。
- 仅验证当前蓝本链路是否可启动、可登录、可读取菜单和权限。
- 不作为最终生产 schema。

初始化 SQL：

```text
sql/install.sql
```

该 SQL 会创建完整蓝本 `la_*` 表并插入默认管理员、角色、菜单、配置等数据。P0.9 已新增最小核心 SQL：

```text
sql/install.core.sql
```

新项目优先使用核心 SQL，它只保留当前 Go 后端实际接入的后台基础表和种子数据。

默认管理员来自蓝本数据：

```text
username: admin
password: 123456
```

该账号只用于本机烟测。任何面向外部网络的环境都必须修改默认密码或重新初始化管理员。

## 脚本

### dry-run 初始化计划

默认不会写数据库：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/init-local-blueprint-db.sh
```

### 执行初始化

需要显式打开写入开关：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
MAKEADMIN_ALLOW_SCHEMA_WRITE=1 ./scripts/init-local-blueprint-db.sh --apply
```

如果目标库已存在且非空，脚本默认拒绝继续，避免覆盖已有表。确认为一次性开发库后，才允许额外设置：

```bash
MAKEADMIN_ALLOW_SCHEMA_WRITE=1 MAKEADMIN_ALLOW_NONEMPTY_DB=1 ./scripts/init-local-blueprint-db.sh --apply
```

### 执行核心初始化

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
MAKEADMIN_ALLOW_SCHEMA_WRITE=1 ./scripts/init-local-core-db.sh --apply
```

### 只读检查初始化结果

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/check-db-seed.sh
```

该脚本只读检查关键表和基础数据，不创建数据库、不写入数据。

## 环境变量

通用 MySQL 参数：

```text
MYSQL_HOST=127.0.0.1
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=
MYSQL_DATABASE=go_makeadmin
```

可按本机环境覆盖：

```bash
MYSQL_USER=makeadmin MYSQL_PASSWORD=your_password ./scripts/check-db-seed.sh
```

## P1 方向

P1 不应继续堆 LikeAdmin 表结构补丁，应从问题本质重新设计：

- 管理员、角色、权限、菜单的边界。
- 多租户和组织模型。
- 数据范围策略。
- 审计日志与操作日志。
- 初始化种子数据和迁移策略。
- 与 zyai 业务表的映射边界。
