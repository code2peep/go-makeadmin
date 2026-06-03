# P5 Module Registry Freeze Checklist

更新时间：2026-06-03

## 目标

P5.19 在继续接入第二个示例模块或模块市场雏形前，对当前 registry、manifest、模块中心 UI、自动 smoke 和人工验收边界做一次冻结判断。

本阶段不新增代码行为，不写数据库，不处理真实业务模块。

## 冻结判断

| 范围 | 当前状态 | 证明入口 | 判断 |
| --- | --- | --- | --- |
| 后端只读 registry | 已具备 | `docs/P5_MODULE_REGISTRY_READONLY.md` | 可冻结 |
| manifest 校验 | 已具备 | `docs/P5_MODULE_REGISTRY_MANIFEST_CHECK.md` | 可冻结 |
| broken fixture | 已具备 | `docs/P5_MODULE_REGISTRY_FAILURE_FIXTURE.md` | 可冻结 |
| registry smoke | 已具备 | `scripts/check-module-registry-smoke.sh` | 可冻结 |
| 路由契约 | 已具备 | `server/generator/routers/gen/module_registry_test.go` | 可冻结 |
| 模块中心状态条 | 已具备 | `docs/P5_MODULE_CENTER_ACCEPTANCE_STATUS.md` | 可冻结 |
| registry 失败态/空态 | 已具备 | `docs/P5_MODULE_CENTER_REGISTRY_STATES.md` | 可冻结 |
| helper/fixture | 已具备 | `admin/src/views/dev_tools/module/registry-state.ts` | 可冻结 |
| 页面 UI 契约 | 已具备 | `scripts/check-module-center-ui-contract.sh` | 可冻结 |
| 登录后截图验收 | 待登录 | `docs/P5_MODULE_CENTER_UI_CONTRACT.md` | 不阻断继续开发 |

## 自动验收基线

继续接入新模块前，至少保持以下命令通过：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
scripts/check-module-registry-smoke.sh
scripts/check-module-center-ui-contract.sh
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```

## 人工验收基线

登录后台后打开：

```text
http://127.0.0.1:5173/module
```

默认 registry 页面重点看：

- `P5.18`
- `默认 Registry`
- `Demo 入口`
- `/demo/article`

broken fixture 页面重点看：

- `Broken Fixture`
- `已开启`
- `异常筛选`
- `校验明细`

当前本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，不做自动登录。

## 进入下一步的判断

当前自动验收基线已经足够支撑继续接入第二个示例模块。第二个示例模块应满足：

- manifest 可被 registry 读取。
- 菜单入口能在模块中心展示。
- 页面可以作为只读示例打开。
- 不新增真实业务表。
- 不写 `.env`。
- 不连接 zyai 业务库。

## 验收结果

- 已通过 `scripts/check-module-registry-smoke.sh`。
- 已通过 `scripts/check-module-center-ui-contract.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 下一步

P5.20：第二个示例模块 registry 接入。建议新增一个轻量只读示例模块，用来验证 registry 和模块中心不是只服务 Demo Article 单模块。
