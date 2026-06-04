# go-makeadmin

`go-makeadmin` 是一个轻量通用管理后台基础框架，面向“先搭好后台底座，再用 AI / vibe coding 快速生成具体业务功能”的开发方式。

它不是一个复杂的模块市场，也不是某个具体业务系统。当前首版目标已经收口为：登录可用、菜单可用、权限和基础设置可用、素材管理可用、代码生成器可作为业务 CRUD 起点。



- 默认本地管理员账号：`admin / 123456`。
- 独立开发库：默认使用 `go_makeadmin`，不和任何业务库混用。
- 自研核心表：默认使用 `ma_*` 系统表
- 核心菜单：工作台、权限管理、组织管理、素材管理、系统设置、开发工具。
- 核心页面：管理员、角色、菜单、部门、岗位、素材、网站信息、存储、系统环境、缓存、日志、字典、代码生成器、模块中心。
- 素材管理：支持全部、未分组、图片/视频上传入口和清晰空态。
- 代码生成器：保留为 AI 业务功能起步工具。
- 验证脚本：覆盖菜单树、路由组件、素材空态、代码生成模板、模块工具、Go 测试、前端类型检查、前端构建和 npm audit。


## 技术栈

后端：

- Go 1.26
- Gin
- Gorm
- MySQL
- Redis
- Zap
- Viper

前端：

- Vue 3
- Vite
- TypeScript
- Element Plus
- Pinia
- Vue Router
- Axios

## 功能概览

### 认证与权限

- 管理员登录和登出。
- 管理员信息、角色、菜单、权限字符。
- 菜单树和动态路由。
- 角色授权和数据范围基础模型。
- 租户上下文和租户切换基础能力。

### 通用后台页面

- 工作台：展示框架状态、核心入口和验证信息。
- 权限管理：管理员、角色、菜单。
- 组织管理：部门、岗位。
- 素材管理：图片/视频素材、分组、未分组、上传入口、空态。
- 系统设置：网站信息、备案、协议、存储设置。
- 系统维护：系统环境、Redis 缓存、系统日志。
- 开发工具：字典管理、代码生成器、模块中心。

### 代码生成器

代码生成器用于快速生成业务模块起步代码，适合后续用 AI 按业务表继续扩展：

- 数据表导入和同步。
- 字段配置和模板预览。
- 后端、前端列表页和表单模板。
- 默认业务模板包含列表、搜索、分页、新增、编辑、删除和状态字典。

### 模块中心

模块中心当前定位为开发工具和示例验收入口，不再继续扩展成复杂模块市场。它可以保留 manifest、预览、安装计划等开发期能力，但首版使用重点仍是通用后台和代码生成器。

## 项目结构

```text
go-makeadmin/
├── admin/                  # Vue 3 管理后台
├── server/                 # Go API 服务
├── sql/                    # ma_* schema 和 seed
├── scripts/                # 本地启动、初始化、验证脚本
├── docs/                   # 阶段文档和架构说明
├── examples/               # 示例模块
├── legacy/                 # 历史蓝本和来源隔离
├── public/                 # 公共静态目录
├── NOTICE.md               # 第三方来源说明
├── LICENSE
└── README.md
```

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/code2peep/go-makeadmin.git
cd go-makeadmin
```

### 2. 准备配置

```bash
cp server/.env.example server/.env
cp admin/.env.development.example admin/.env.development
```

默认配置：

```text
API: http://127.0.0.1:8000/api
Admin: http://127.0.0.1:5173
MySQL: root:@tcp(127.0.0.1:3306)/go_makeadmin
Redis: redis://127.0.0.1:6379/0
```

如果 `8000` 端口已被占用，修改：

```text
server/.env: SERVER_PORT
server/.env: PUBLIC_URL
admin/.env.development: VITE_API_PROXY_TARGET
```

### 3. 安装前端依赖

```bash
cd admin
npm ci
cd ..
```

### 4. 初始化本地数据库

确保 MySQL 可连接后执行：

```bash
./scripts/init-p1-db.sh
```

默认会创建或使用数据库 `go_makeadmin`，并导入 `sql/p1.schema.sql` 和 `sql/p1.seed.sql`。

如果本地 disposable 开发库需要重建：

```bash
INIT_P1_DROP=1 ./scripts/init-p1-db.sh
```

如需设置本地管理员密码：

```bash
ADMIN_PASSWORD='your-local-password' ./scripts/init-p1-db.sh
```

### 5. 检查服务和种子数据

```bash
./scripts/check-services.sh
./scripts/check-p1-seed.sh
```

### 6. 启动后端

```bash
./scripts/dev-api.sh
```

### 7. 启动前端

另开一个终端：

```bash
./scripts/dev-admin.sh
```

访问：

```text
http://127.0.0.1:5173
```

本地默认账号：

```text
username: admin
password: 123456
```

## 验证

不触库全量验证：

```bash
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```

它会执行：

- 运行时残留检查。
- 工作台、登录、菜单树、路由组件、素材空态契约检查。
- 模块工具和代码生成器 no-db 契约检查。
- `cd server && go test ./...`。
- `cd admin && npm run type-check`。
- `cd admin && npm run build`。
- `cd admin && npm audit --audit-level=moderate`。

前端 build 目前可能出现 `@vueuse/core` 的 Rolldown pure annotation warning；只要命令退出码为 0，不阻塞验证结论。

## 常用脚本

```bash
./scripts/verify-no-db.sh                       # 不触库全量验证
./scripts/init-p1-db.sh                         # 初始化 go_makeadmin 开发库
./scripts/check-services.sh                     # 检查 MySQL / Redis / 基础服务
./scripts/check-p1-seed.sh                      # 检查 ma_* 种子数据
./scripts/dev-api.sh                            # 启动 Go API
./scripts/dev-admin.sh                          # 启动 Vue 管理端
./scripts/check-admin-route-components.py       # 检查菜单页面组件
./scripts/check-material-empty-state-contract.sh # 检查素材页空态契约
```

## 数据库说明

当前默认数据库是 `go_makeadmin`，核心运行表统一使用 `ma_*`：

- `ma_admin`
- `ma_role`
- `ma_menu`
- `ma_permission`
- `ma_setting`
- `ma_dict_type`
- `ma_dict_data`
- `ma_file`
- `ma_file_category`
- `ma_login_log`
- `ma_codegen_table`
- `ma_codegen_column`

`sql/install.sql` 和 `sql/install.core.sql` 保留为历史蓝本或核心蓝本，不是当前默认初始化入口。当前默认入口是：

```text
sql/p1.schema.sql
sql/p1.seed.sql
```

## 开发原则

- 先保持底座轻量可用，再做具体业务。
- 具体业务功能优先通过代码生成器和现有页面模式生成。
- 不把默认账号、真实密钥、生产配置提交到仓库。
- 不在当前底座里继续扩展复杂模块市场。

## 文档


## 来源与授权

本项目保留 MIT License。第三方来源见 `NOTICE.md`。
