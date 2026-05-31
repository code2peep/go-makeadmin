# go-makeadmin

`go-makeadmin` 是 `makeadmin` 的 Go 版本后台基础框架，目标是沉淀一套可复用、可商业使用、可迁移到多个项目的通用管理后台底座。

本项目 P0 阶段参考 LikeAdmin Go 的 MIT 开源实现，先完成现代化升级和最小可运行验证，再逐步重构为自研框架能力。

## 当前阶段

P0：基础框架 POC。

当前目标：

- 保留 LikeAdmin Go 的合法 MIT 来源说明。
- 升级后端和前端依赖。
- 跑通管理后台登录、菜单、角色权限、上传。
- 建立后续多租户、数据权限、审计日志的框架边界。

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

不触库验证：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/verify-no-db.sh
```

后端：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin/server
go test ./...
```

前端：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin/admin
npm run build
```

本地启动链路见 `docs/LOCAL_DEV.md`。P1 数据模型设计见 `docs/P1_SCHEMA_PLAN.md`。

## 授权

本项目保留 MIT License。第三方来源见 `NOTICE.md`。
