# go-makeadmin

`go-makeadmin` 是 `makeadmin` 的 Go 版本后台基础框架，目标是沉淀一套可复用、可商业使用、可迁移到多个项目的通用管理后台底座。

本项目早期参考 LikeAdmin Go 的 MIT 开源实现，保留合法来源说明；当前已进入 P1 收口阶段，默认运行链路面向 `ma_*` 自研系统表和独立开发库 `go_makeadmin`。

## 当前阶段

P1：Go 后台框架底座收口。

当前目标：

- 默认使用 `go_makeadmin` 独立库，不和 zyai 业务库混用。
- 核心后台链路使用 `ma_*` 表：登录、菜单、权限、设置、字典、文件、日志和代码生成器。
- 旧 `la_*` SQL、模型和服务只作为蓝本参考保留，不再作为 P1 核心运行兜底。
- P1 收口后进入 P2：认证模型、多租户上下文、数据权限、模块生成闭环。

## 技术栈

后端：

- Go
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

## 目录

```text
go-makeadmin/
├── CLAUDE.md
├── LICENSE
├── NOTICE.md
├── README.md
├── admin/
├── docs/
├── frontend/
├── legacy/
├── public/
├── scripts/
├── server/
└── sql/
```

## 验证

不触库验证用于编译、类型检查、前端构建和依赖审计：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/verify-no-db.sh
```

P1 独立库种子只读检查：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/check-p1-seed.sh
```

P1 HTTP smoke 只允许在本地 disposable P1 库执行写操作：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
P1_SMOKE_ALLOW_WRITE=1 \
P1_SMOKE_BASE_URL=http://127.0.0.1:8000/api \
P1_SMOKE_ADMIN_PASSWORD='your-local-admin-password' \
python3 scripts/p1-smoke.py
```

本地启动链路见 `docs/LOCAL_DEV.md`。P1 验收清单见 `docs/P1_ACCEPTANCE_CHECKLIST.md`。P1 数据模型设计见 `docs/P1_SCHEMA_PLAN.md`。

## 授权

本项目保留 MIT License。第三方来源见 `NOTICE.md`。
