# P5 Module Registry Document Index

更新时间：2026-06-03

## 目标

P5.16 将 P5.5-P5.24 的模块 registry 文档收敛成一个短索引，后续继续做模块市场、真实模块接入或后台人工验收时，先从这里定位入口。

本阶段只整理文档入口，不新增接口、不改管理端行为、不写数据库。

## 推荐阅读顺序

| 顺序 | 文档 | 用途 |
| --- | --- | --- |
| 1 | `docs/P5_MODULE_REGISTRY_READONLY.md` | 后端只读 registry 接口起点 |
| 2 | `docs/P5_MODULE_REGISTRY_MANIFEST_CHECK.md` | registry 与 manifest 校验规则 |
| 3 | `docs/P5_MODULE_REGISTRY_CHECK_DETAIL.md` | 模块中心校验明细展示 |
| 4 | `docs/P5_MODULE_REGISTRY_FAILURE_FIXTURE.md` | broken fixture 异常模块入口 |
| 5 | `docs/P5_MODULE_REGISTRY_SMOKE.md` | 默认和异常 registry 自动 smoke |
| 6 | `docs/P5_MODULE_REGISTRY_ROUTE_CONTRACT.md` | `/api/gen/moduleRegistry` 路由响应契约 |
| 7 | `docs/P5_MODULE_REGISTRY_ACCEPTANCE_MATRIX.md` | 自动和人工验收矩阵 |
| 8 | `docs/P5_MODULE_CENTER_ACCEPTANCE_STATUS.md` | 模块中心验收辅助状态条 |
| 9 | `docs/P5_MODULE_CENTER_REGISTRY_STATES.md` | registry 失败态和空态 |
| 10 | `docs/P5_MODULE_CENTER_REGISTRY_STATE_HELPER.md` | 前端 registry 状态 helper |
| 11 | `docs/P5_MODULE_CENTER_REGISTRY_STATE_FIXTURE.md` | TypeScript 编译期 fixture |
| 12 | `docs/P5_MODULE_CENTER_MANUAL_CHECKLIST.md` | 模块中心登录后人工验收清单 |
| 13 | `docs/P5_MODULE_CENTER_UI_CONTRACT.md` | 模块中心 UI 文案契约检查 |
| 14 | `docs/P5_MODULE_REGISTRY_FREEZE_CHECKLIST.md` | 进入真实模块接入前的冻结判断 |
| 15 | `docs/P5_SECOND_DEMO_MODULE.md` | 第二个只读示例模块接入 |
| 16 | `docs/P5_DEMO_NOTICE_ACCEPTANCE.md` | Demo Notice 安装计划和页面入口验收 |
| 17 | `docs/P5_DEMO_NOTICE_STATUS_CONTRACT.md` | Demo Notice 未注册运行时状态契约 |
| 18 | `docs/P5_MODULE_CENTER_MULTI_FILTERS.md` | 多模块筛选与统计契约 |
| 19 | `docs/P5_MODULE_CENTER_MANUAL_MULTI_CHECKLIST.md` | 多模块登录后人工验收入口 |

## 代码入口

| 层级 | 文件 | 说明 |
| --- | --- | --- |
| 后端 service | `server/generator/service/gen/module_registry.go` | registry 清单、manifest 校验和 broken fixture 组装 |
| 后端 route test | `server/generator/routers/gen/module_registry_test.go` | `/api/gen/moduleRegistry` 路由响应契约 |
| 管理端页面 | `admin/src/views/dev_tools/module/index.vue` | 模块中心 registry 展示、筛选、明细和验收状态条 |
| 管理端 helper | `admin/src/views/dev_tools/module/registry-state.ts` | registry 页面状态纯逻辑 |
| 管理端 fixture | `admin/src/views/dev_tools/module/registry-state.fixture.ts` | helper 编译期 smoke 输入输出形状 |
| 自动 smoke | `scripts/check-module-registry-smoke.sh` | 默认 registry、broken fixture、路由契约和 manifest 工具验证 |

## 当前自动验收入口

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
scripts/check-module-registry-smoke.sh
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```

前端 helper/fixture 的最小验证入口：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin/admin
npm run type-check
```

## 登录后人工验收入口

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

登录后台后打开：

```text
http://127.0.0.1:5173/module
```

当前页面复验仍需要真实登录态；本项目不伪造登录 token，不修改管理员密码，不修改 `.env`。

## 后续接入判断

新增真实模块时，先检查：

- manifest 是否能通过 `docs/P5_MODULE_REGISTRY_MANIFEST_CHECK.md` 的字段规则。
- registry 是否能通过 `scripts/check-module-registry-smoke.sh`。
- 模块中心是否能在默认、异常筛选、明细弹窗和空态下保持可读。
- `docs/P5_MODULE_REGISTRY_ACCEPTANCE_MATRIX.md` 的人工验收项是否有新增场景。

## 保留边界

P5.16 不做：

- 不新增 registry 写入接口。
- 不新增模块市场 UI。
- 不修改后台登录流程。
- 不修改 `.env` 或管理员密码。
- 不新增数据库 schema。
- 不处理 PTLM 业务模块。
