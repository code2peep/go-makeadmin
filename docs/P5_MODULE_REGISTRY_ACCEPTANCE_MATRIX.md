# P5 Module Registry Acceptance Matrix

更新时间：2026-06-03

## 目标

P5.11 把模块 registry 的本地 API smoke、后端路由契约、全量 no-db 验证和登录后页面验收拆成稳定矩阵，避免后续模块接入时只靠临时记忆判断是否通过。

本阶段不新增后端 API，不新增数据库写入，不修改登录凭证。

## 自动验收矩阵

| 场景 | 入口 | 证明点 | 当前自动化 |
| --- | --- | --- | --- |
| 默认 registry 服务契约 | `scripts/check-module-registry-smoke.sh` | 未开启 broken fixture 时只返回 `Demo Article`，且 `manifestStatus=passed` | 已覆盖 |
| broken fixture 服务契约 | `scripts/check-module-registry-smoke.sh` | 开启 `MAKEADMIN_ENABLE_BROKEN_MODULE_REGISTRY_FIXTURE=1` 后返回 `Broken Manifest Fixture`，且不影响 Demo Article | 已覆盖 |
| registry 路由响应契约 | `scripts/check-module-registry-smoke.sh` | `GET /api/gen/moduleRegistry` handler 返回 `code=200` 和前端依赖字段 | 已覆盖 |
| no-db 全量链路 | `./scripts/verify-no-db.sh` | runtime residue、模块工具、Go test、前端 type-check/build、npm audit 均通过 | 已覆盖 |
| 前端类型契约 | `cd admin && npm run type-check` | 模块中心消费 `manifestStatus`、`manifestMessage`、`manifestChecks` 类型稳定 | 已覆盖 |

## 人工验收矩阵

| 场景 | 操作 | 期望结果 | 当前状态 |
| --- | --- | --- | --- |
| 默认模块中心 | 登录后台后打开 `/module` | 显示 `P5.11`、`Demo Article`、校验 `已通过`、可打开 `明细` | 待登录复验 |
| broken fixture 页面态 | 用 `MAKEADMIN_ENABLE_BROKEN_MODULE_REGISTRY_FIXTURE=1` 启动 API 后打开 `/module` | 显示 `Broken Manifest Fixture`，校验为 `异常` | 待登录复验 |
| 异常筛选 | 在 broken fixture 页面态点击 `异常` 筛选 | 只展示异常模块或包含异常模块 | 待登录复验 |
| 校验明细弹窗 | 点击 broken fixture 的 `明细` | 弹窗展示失败检查项和失败说明 | 待登录复验 |
| Demo Article 入口 | 点击 Demo Article 的 `打开` | 进入 `/demo/article`，页面可见 | 待登录复验 |

## 本地命令

自动验收：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
scripts/check-module-registry-smoke.sh
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```

页面验收需要先登录后台。当前本机没有 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD` 环境变量，不能自动完成登录后截图验收。

默认 registry 页面：

```bash
MAKEADMIN_ENABLE_DEMO_MODULE=1 \
MAKEADMIN_ALLOW_MODULE_INSTALL_APPLY=1 \
MAKEADMIN_ALLOW_MODULE_UNINSTALL_APPLY=1 \
./scripts/dev-api.sh
```

broken fixture 页面：

```bash
MAKEADMIN_ENABLE_DEMO_MODULE=1 \
MAKEADMIN_ENABLE_BROKEN_MODULE_REGISTRY_FIXTURE=1 \
MAKEADMIN_ALLOW_MODULE_INSTALL_APPLY=1 \
MAKEADMIN_ALLOW_MODULE_UNINSTALL_APPLY=1 \
./scripts/dev-api.sh
```

## 验收结果

- 已通过 `scripts/check-module-registry-smoke.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 登录后页面验收待重新登录后执行。

## 保留边界

P5.11 不做：

- 不伪造登录 token。
- 不修改 `.env` 或管理员密码。
- 不新增 registry 写入接口。
- 不新增数据库 schema。
- 不处理 PTLM 业务模块。
