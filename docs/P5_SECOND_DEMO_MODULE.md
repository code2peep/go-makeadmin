# P5 Second Demo Module

更新时间：2026-06-03

## 目标

P5.20 接入第二个只读示例模块 `Demo Notice`，验证模块 registry 和模块中心支持多模块，而不是只服务 `Demo Article` 单模块。

本阶段不新增数据库 schema，不改 `.env`，不注册后端运行时路由。

## 新增模块

```text
examples/demo_notice/manifest.json
admin/src/api/demoNotice.ts
admin/src/views/demo_notice/index.vue
```

`Demo Notice` 是前端只读示例：

- module：`demo_notice`
- table：`ma_demo_notice`
- entry：`/demo/notice`
- manifest：`examples/demo_notice/manifest.json`
- runtimeRegistered：`false`

## Registry 开关

默认 registry 仍只返回 `Demo Article`，避免破坏既有单模块 smoke。

开启第二个示例模块：

```bash
MAKEADMIN_ENABLE_DEMO_NOTICE_MODULE=1
```

开启后 `GET /api/gen/moduleRegistry` 会返回 `Demo Notice`。

## 验收标准

- `python3 scripts/check-module-manifests.py` 通过。
- `scripts/check-module-registry-smoke.sh` 通过。
- `scripts/check-module-center-ui-contract.sh` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 登录后在开启 `MAKEADMIN_ENABLE_DEMO_NOTICE_MODULE=1` 的 API 下，模块中心能看到 `Demo Notice` 和 `/demo/notice`。

## 验收结果

- 已通过 `python3 scripts/check-module-manifests.py`。
- 已通过 `scripts/check-module-registry-smoke.sh`。
- 已通过 `scripts/check-module-center-ui-contract.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 已生成 `demo_notice` chunk。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 保留边界

P5.20 不做：

- 不新增 `ma_demo_notice` 表。
- 不注册 `/demo_notice/*` 后端运行时路由。
- 不修改 `.env`。
- 不处理 PTLM 业务模块。
