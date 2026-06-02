# P2 Final Status

更新时间：2026-06-02

## 结论

P2 已冻结，可以作为 `go-makeadmin` 后续产品化框架能力的基础面进入 P3。

P2 的冻结目标是把 P1 的可运行后台底座推进为可复用的多项目后台框架能力。当前冻结面包括：

- JWT + Redis session state 认证模型。
- 请求级租户上下文和默认租户兼容策略。
- 核心列表查询的数据权限约束。
- 后端租户切换入口和租户初始化命令。
- Go 后端代码生成器闭环。
- Vue 前端生成模板闭环。
- manifest 驱动的模块注册、安装、卸载和 runtime 开关。
- 默认 no-db 验证链路和显式本地写库生命周期 smoke。

## 已完成范围

认证与租户：

- `docs/P2_AUTH_MODEL.md`
- `docs/P2_TENANT_CONTEXT.md`
- `docs/P2_DATA_SCOPE.md`
- `docs/P2_TENANT_SWITCH.md`
- `docs/P2_TENANT_MIGRATION.md`
- `docs/P2_TENANT_INIT_PLAN.md`
- `docs/P2_TENANT_INIT_APPLY_GUARD.md`
- `docs/P2_TENANT_INIT_APPLY.md`

代码生成器：

- `docs/P2_CODEGEN_CLOSURE.md`
- `docs/P2_FRONTEND_CODEGEN_CLOSURE.md`
- `scripts/check-codegen-frontend.sh`
- `examples/README.md`
- `examples/demo/`

模块生命周期：

- `examples/<module>/manifest.json` 作为模块清单约定。
- `scripts/check-module-manifests.py` 校验 manifest。
- `scripts/module-registry-plan.py` 生成和执行模块注册 SQL。
- `scripts/module-role-grant-plan.py` 生成角色授权 SQL。
- `scripts/module-install-plan.py` 生成和执行模块安装计划。
- `scripts/module-uninstall-plan.py` 生成和执行模块卸载计划。
- `scripts/check-module-lifecycle-smoke.sh` 执行本地写库生命周期 smoke。
- `scripts/check-module-tools-no-db.sh` 进入默认 no-db 验证链路。
- `docs/P2_MODULE_GUIDE.md` 作为模块能力统一入口。

## 冻结验收范围

P2.25 冻结验收需要通过：

```bash
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
MAKEADMIN_ALLOW_MODULE_LIFECYCLE_WRITE=1 scripts/check-module-lifecycle-smoke.sh
git diff --check
```

验收含义：

- `verify-no-db` 覆盖运行残留守卫、模块工具 no-db guard、Go test、前端 type-check、前端 build 和 npm audit。
- 模块工具 no-db guard 覆盖 manifest 校验、模块 dry-run 预览和写入门禁失败检查。
- 生命周期 smoke 在本地 `go_makeadmin` 开发库安装 demo article、检查计数、卸载 demo article、确认残留为 0，并验证二次卸载 no-op。

## 冻结验收结果

- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `MAKEADMIN_ALLOW_MODULE_LIFECYCLE_WRITE=1 scripts/check-module-lifecycle-smoke.sh`。
- 已通过 `git diff --check`。
- 已确认 `server/.env`、`admin/.env.development`、`admin/node_modules`、`admin/dist`、`frontend`、`public/admin`、`public/assets` 继续被 Git 忽略。
- 生命周期 smoke 安装后计数为 5 条权限、1 条菜单、1 条菜单权限关联、5 条角色授权。
- 生命周期 smoke 卸载后残留计数为 0，二次卸载为 no-op。
- 冻结验收后再次查询 demo article 残留计数为 0。
- 本阶段没有修改 schema、没有读取或修改 `.env`、没有连接真实 zyai 业务库。

## 已知验证噪音

前端构建会输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 中 `/* #__PURE__ */` 注释位置的 warning；当前不影响构建退出码。

## 保留边界

- demo runtime 模块默认关闭，必须显式设置 `MAKEADMIN_ENABLE_DEMO_MODULE=1` 才会挂载。
- 模块安装写入必须显式设置 `MAKEADMIN_ALLOW_MODULE_INSTALL_WRITE=1`。
- 模块卸载删除必须显式设置 `MAKEADMIN_ALLOW_MODULE_UNINSTALL_WRITE=1`。
- 模块生命周期 smoke 必须显式设置 `MAKEADMIN_ALLOW_MODULE_LIFECYCLE_WRITE=1`。
- manifest `requiresSchema=true` 时安装写入直接失败，不自动建表。
- runtime 开关只输出提示，不修改 `.env` 或系统环境变量。
- 卸载只删除 manifest 声明的菜单、权限和授权关联，不删除前端文件、后端代码、schema 或 codegen 元数据。

## 不覆盖范围

P2 冻结不代表以下内容已完成：

- 不迁移 zyai 真实业务库。
- 不连接真实 zyai 业务库执行写操作。
- 不提供生产部署、CI/CD、npm publish 或线上发布能力。
- 不提供图形化模块市场或后台模块安装页面。
- 不自动创建业务 schema。
- 不删除 `legacy/`、`frontend/` 或旧蓝本文件。
- 不处理多数据库、多 Redis namespace 或生产密钥轮换策略。

## P3 入口

下一步进入 P3.1：业务模块脚手架与产品化模板。

P3.1 的目标是把 P2 的 codegen、manifest 和模块生命周期能力串成标准模块开发体验：一个模块从 manifest 到后端路由、前端 API、前端页面、安装计划、卸载计划和本地验证都有明确入口。

P3 状态见 `docs/P3_STATUS.md`。P3.1 模块脚手架见 `docs/P3_MODULE_SCAFFOLD.md`。
