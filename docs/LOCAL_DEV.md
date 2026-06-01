# 本地开发启动链路

本文档定义 `go-makeadmin` P0 阶段的本地开发入口。默认先跑不触库验证，只有在 MySQL、Redis 和初始化数据明确准备好之后，才启动完整登录链路。

## 两条链路

### 1. 不触库验证

用于依赖升级、类型修复、前端构建和后端编译测试。该链路不连接 MySQL，不连接 Redis，不执行 SQL。

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/verify-no-db.sh
```

它会执行：

```bash
cd server && GOPROXY=https://goproxy.cn,direct go test ./...
cd admin && npm run type-check
cd admin && npm run build
cd admin && npm audit --audit-level=moderate
```

### 2. 完整联调启动

用于验证登录、菜单、角色权限、上传等真实 API 行为。该链路需要 MySQL、Redis 和系统初始化数据已经存在。

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/check-services.sh
```

检查通过后，分别启动后端和前端：

```bash
./scripts/dev-api.sh
```

```bash
./scripts/dev-admin.sh
```

默认访问：

```text
管理端：http://127.0.0.1:5173
API：http://127.0.0.1:8000/api
```

如果本机 `8000` 已被其他服务占用，调整：

```text
server/.env: SERVER_PORT
server/.env: PUBLIC_URL
admin/.env.development: VITE_API_PROXY_TARGET
```

当前本机开发环境使用：

```text
API：http://127.0.0.1:18000/api
管理端：http://127.0.0.1:5173
```

## MySQL 前置条件

当前默认后端 DSN：

```text
root:@tcp(127.0.0.1:3306)/go_makeadmin?charset=utf8mb4&parseTime=True&loc=Local
```

完整联调需要满足：

- MySQL 服务可连接。
- 数据库 `go_makeadmin` 已存在。
- 系统表已存在；P1 阶段优先使用 `ma_*` 表，旧 `la_*` 表只作为过渡兜底。
- 基础管理员、菜单、角色、网站配置等初始化数据已存在。

P0.7 已在当前本机初始化一次性开发库 `go_makeadmin`，用于烟测当前蓝本链路。

P0.7 的本地库初始化策略见：

```text
docs/DB_INIT_PLAN.md
```

新项目优先使用 P0.9 核心初始化 SQL：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
MAKEADMIN_ALLOW_SCHEMA_WRITE=1 ./scripts/init-local-core-db.sh --apply
```

## Redis 前置条件

当前默认 Redis：

```text
redis://127.0.0.1:6379/0
```

完整联调需要满足：

- Redis 服务可连接。
- 登录 token、管理员缓存、菜单权限缓存可以写入当前 Redis DB。
- Redis key 前缀为 `MakeAdmin:`，避免和其他项目混用。

## 示例配置

后端示例配置：

```text
server/.env.example
```

前端示例配置：

```text
admin/.env.development.example
admin/.env.production.example
```

本地开发时，前端默认请求自身域名下的 `/api`，由 Vite proxy 转发到 `http://127.0.0.1:8000`。如果后端端口变化，调整前端开发环境的 `VITE_API_PROXY_TARGET`。

本地烟测账号：

```text
P0 蓝本：
username: admin
password: 123456

P1 ma_*：
username: admin
password: 由 scripts/init-p1-db.sh 执行时的 ADMIN_PASSWORD 决定
```

## 启动顺序

1. 跑 `./scripts/verify-no-db.sh`，确认无数据库依赖的基础质量。
2. 准备 MySQL 和 Redis。
3. 如本机未初始化 P1 开发库，执行 `scripts/init-p1-db.sh` 并通过 `ADMIN_PASSWORD` 指定本地管理员密码。
4. 跑 `./scripts/check-services.sh`。
5. 跑 `./scripts/check-p1-seed.sh`，确认 `ma_*` 表和基础数据存在。
6. 分别启动 `./scripts/dev-api.sh` 和 `./scripts/dev-admin.sh`。
7. 浏览器打开 `http://127.0.0.1:5173`。

## 当前限制

- P1 已开始接管认证、权限、菜单和设置链路，旧 `la_*` 路径暂时保留为过渡兜底。
- `sql/install.sql` 是完整蓝本初始化 SQL，不是最终自研 schema。
- `sql/install.core.sql` 是 P0.9 生成的最小核心初始化 SQL，已剔除文章、C 端用户、渠道、装修等业务演示表。
- 未导入初始化 SQL 时，后端服务可能启动成功，但登录和菜单接口不能完成业务验证。
