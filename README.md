# go-makeadmin

`go-makeadmin` 是 `makeadmin` 的 Go 版本后台基础框架，目标是沉淀一套可复用、可商业使用、可迁移到多个项目的通用管理后台底座。

本项目早期参考 LikeAdmin Go 的 MIT 开源实现，保留合法来源说明；当前已进入 P2，默认运行链路面向 `ma_*` 自研系统表和独立开发库 `go_makeadmin`。

## 当前阶段

P2：框架能力增强。

当前目标：

- 默认使用 `go_makeadmin` 独立库，不和 zyai 业务库混用。
- 核心后台链路使用 `ma_*` 表：登录、菜单、权限、设置、字典、文件、日志和代码生成器。
- 旧 `la_*` SQL、模型和服务只作为蓝本参考保留，不再作为 P1 核心运行兜底。
- P2 当前任务：认证模型、多租户上下文、数据权限、模块生成闭环。

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

本地启动链路见 `docs/LOCAL_DEV.md`。P1 最终状态见 `docs/P1_FINAL_STATUS.md`。P2 状态见 `docs/P2_STATUS.md`。P2 最终状态见 `docs/P2_FINAL_STATUS.md`。P2 模块使用指南见 `docs/P2_MODULE_GUIDE.md`。P3 状态见 `docs/P3_STATUS.md`。P3.1 模块脚手架见 `docs/P3_MODULE_SCAFFOLD.md`。P3.2 模块 codegen 联动见 `docs/P3_MODULE_CODEGEN_LINK.md`。P3.3 模块脚手架写入 smoke 见 `docs/P3_MODULE_SCAFFOLD_WRITE_SMOKE.md`。P3.4 模块生成器配置预览见 `docs/P3_MODULE_CODEGEN_PLAN.md`。P3.5 模块生成器写入边界见 `docs/P3_MODULE_CODEGEN_APPLY_BOUNDARY.md`。P3.6 模块生成器受控写入见 `docs/P3_MODULE_CODEGEN_APPLY.md`。P3.7 模块生成器回读闭环见 `docs/P3_MODULE_CODEGEN_READBACK.md`。P3.8 后台 manifest 预览见 `docs/P3_MODULE_MANIFEST_PREVIEW.md`。P3.9 安装计划联动见 `docs/P3_MODULE_INSTALL_PLAN_PREVIEW.md`。P3.10 后台安装写入门禁见 `docs/P3_MODULE_INSTALL_APPLY_BOUNDARY.md`。P3.11 后台安装受控写入见 `docs/P3_MODULE_INSTALL_APPLY.md`。P3.12 后台卸载写入门禁见 `docs/P3_MODULE_UNINSTALL_APPLY_BOUNDARY.md`。P3.13 后台卸载受控删除见 `docs/P3_MODULE_UNINSTALL_APPLY.md`。P3.14 模块安装卸载页面闭环见 `docs/P3_MODULE_APPLY_UI_CLOSURE.md`。P3.15 模块操作摘要与审计规划见 `docs/P3_MODULE_APPLY_AUDIT_SUMMARY.md`。P3.16 模块 apply 页面状态收敛见 `docs/P3_MODULE_APPLY_UI_STATE.md`。P3.17 模块 manifest 前端 API 类型见 `docs/P3_MODULE_MANIFEST_API_TYPES.md`。P3.18 模块 apply 错误归一化见 `docs/P3_MODULE_MANIFEST_APPLY_ERROR.md`。P3.19 模块 apply 结果视图见 `docs/P3_MODULE_APPLY_RESULT_VIEW.md`。P3.20 模块 apply 审计 DTO 见 `docs/P3_MODULE_APPLY_AUDIT_DTO.md`。P2.1 认证模型见 `docs/P2_AUTH_MODEL.md`。P2.2 租户上下文见 `docs/P2_TENANT_CONTEXT.md`。P2.3 数据权限见 `docs/P2_DATA_SCOPE.md`。P2.4 租户切换见 `docs/P2_TENANT_SWITCH.md`。P2.5 租户迁移策略见 `docs/P2_TENANT_MIGRATION.md`。P2.6 租户初始化 dry-run 见 `docs/P2_TENANT_INIT_PLAN.md`。P2.7 租户初始化写入门禁见 `docs/P2_TENANT_INIT_APPLY_GUARD.md`。P2.8 租户初始化受控写入见 `docs/P2_TENANT_INIT_APPLY.md`。P2.9 代码生成器闭环见 `docs/P2_CODEGEN_CLOSURE.md`。P2.10 前端生成模板闭环见 `docs/P2_FRONTEND_CODEGEN_CLOSURE.md`。P2.11 模块注册清单见 `docs/P2_MODULE_REGISTRY.md`。P2.12 注册 SQL dry-run 见 `docs/P2_MODULE_REGISTRY_SQL_PLAN.md`。P2.13 模块注册受控写入见 `docs/P2_MODULE_REGISTRY_APPLY.md`。P2.14 角色授权 dry-run 见 `docs/P2_MODULE_ROLE_GRANT_PLAN.md`。P2.15 运行时模块注册见 `docs/P2_MODULE_RUNTIME_REGISTRY.md`。P2.16 模块安装计划见 `docs/P2_MODULE_INSTALL_PLAN.md`。P2.17 安装写入边界见 `docs/P2_MODULE_INSTALL_APPLY_BOUNDARY.md`。P2.18 模块安装受控写入见 `docs/P2_MODULE_INSTALL_APPLY.md`。P2.19 模块卸载计划见 `docs/P2_MODULE_UNINSTALL_PLAN.md`。P2.20 卸载写入边界见 `docs/P2_MODULE_UNINSTALL_APPLY_BOUNDARY.md`。P2.21 模块卸载受控删除见 `docs/P2_MODULE_UNINSTALL_APPLY.md`。P2.22 生命周期 smoke 见 `docs/P2_MODULE_LIFECYCLE_SMOKE.md`。P2.23 模块验证边界见 `docs/P2_MODULE_VERIFY_BOUNDARY.md`。P2.24 模块文档收敛见 `docs/P2_MODULE_GUIDE.md`。P2.25 P2 冻结状态见 `docs/P2_FINAL_STATUS.md`。

## 授权

本项目保留 MIT License。第三方来源见 `NOTICE.md`。
