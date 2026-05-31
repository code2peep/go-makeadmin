# go-makeadmin 项目级开发规范

本项目是 `makeadmin` 的 Go 版本后台基础框架，目标是形成一套可复用、可商业使用、可迁移到多个项目的通用管理后台底座。

## 项目定位

- 名称：`go-makeadmin`
- 定位：Go 后端 + Vue 管理端的通用后台框架。
- 来源策略：参考 LikeAdmin Go 的 MIT 开源底座，保留合法来源说明，逐步重构为自研基础框架。
- 首个业务验证项目：`zyai` 中医舌诊体质筛查 SaaS。

## 第一原则

- 先定规范，再写代码。
- 先做可验证底座，再迁真实业务。
- 先保证权限和数据边界，再做页面体验。
- 不复制 Gin-Vue-Admin 框架层代码。
- 不在未确认数据库方案前连接或迁移真实业务库。

## 目录约定

- `server/`：Go 后端。
- `admin/`：Vue 管理后台。
- `docs/`：架构、数据库、模块和迁移文档。
- `scripts/`：本地开发、验证、构建脚本。
- `examples/`：框架示例模块。
- `NOTICE.md`：第三方开源来源与授权说明。

## 后端目标技术栈

- Go
- Gin
- Gorm
- MySQL
- Redis
- Zap
- Viper

## 前端目标技术栈

- Vue 3
- Vite
- TypeScript
- Element Plus
- Pinia
- Vue Router
- Axios

## P0 成功标准

- 后端依赖完成现代化升级。
- 前端依赖完成现代化升级。
- 本地能启动管理后台和 API 服务。
- 管理员登录、菜单、角色权限、上传能力可用。
- `go test ./...` 通过。
- `npm run build` 通过。

## 红线

以下操作必须先确认：

- 删除文件或目录。
- Git 回滚、reset、checkout 覆盖文件、清理工作树。
- 修改 `.env`、密钥、token、生产配置、CI/CD 配置。
- 创建、修改或迁移数据库 schema。
- 连接现有 zyai 业务数据库执行写操作。
- 安装新的全局依赖或修改系统配置。
- 部署、发布、npm publish、生产环境操作。

## 开发纪律

- 每次只改和当前阶段目标直接相关的文件。
- 不为了跑通而绕过认证、权限、租户隔离或审计。
- 外部依赖升级必须配套验证命令。
- 发现 LikeAdmin 原始代码不适合作为长期底座时，重写小核心，不继续堆补丁。

## 验证命令

不触库全量验证：

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
