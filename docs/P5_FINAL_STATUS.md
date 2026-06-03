# P5 Final Status

更新时间：2026-06-03

## 结论

P5 模块中心多模块底座已冻结：`Demo Article` 和 `Demo Notice` 已覆盖 registry、manifest 校验、模块中心状态条、筛选统计、登录后人工 checklist、前端 helper fixture 和 no-db 自动验收。

本阶段仍不把 P5 视为真实业务模块市场；它是进入下一阶段产品化模块中心前的多模块底座验收点。

## 已冻结范围

| 范围 | 状态 | 证明入口 |
| --- | --- | --- |
| Demo Article 可见模块 | 已冻结 | `docs/P5_DEMO_MODULE_VISIBLE.md` |
| Demo Article install/uninstall apply | 已冻结 | `docs/P5_MODULE_CENTER_APPLY.md` |
| 模块状态回读和筛选 | 已冻结 | `docs/P5_MODULE_STATUS_READBACK.md`、`docs/P5_MODULE_STATUS_FILTERS.md` |
| 后端 registry 只读接口 | 已冻结 | `docs/P5_MODULE_REGISTRY_READONLY.md` |
| manifest 校验和异常 fixture | 已冻结 | `docs/P5_MODULE_REGISTRY_MANIFEST_CHECK.md`、`docs/P5_MODULE_REGISTRY_FAILURE_FIXTURE.md` |
| registry smoke 和 route contract | 已冻结 | `scripts/check-module-registry-smoke.sh` |
| 模块中心 UI contract | 已冻结 | `scripts/check-module-center-ui-contract.sh` |
| Demo Notice 第二示例模块 | 已冻结 | `docs/P5_SECOND_DEMO_MODULE.md` |
| Demo Notice 安装计划和入口 | 已冻结 | `scripts/check-demo-notice-module.sh` |
| Demo Notice 未注册运行时状态 | 已冻结 | `docs/P5_DEMO_NOTICE_STATUS_CONTRACT.md` |
| 多模块筛选统计 | 已冻结 | `scripts/check-module-center-filter-contract.sh` |
| 多模块登录后 checklist | 已冻结 | `scripts/check-module-center-manual-checklist.sh` |
| no-db 全量验证 | 已冻结 | `./scripts/verify-no-db.sh` |

## 自动验收基线

P5 冻结后至少保持以下命令通过：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
scripts/check-module-registry-smoke.sh
scripts/check-module-center-ui-contract.sh
scripts/check-module-center-filter-contract.sh
scripts/check-module-center-manual-checklist.sh
scripts/check-demo-notice-module.sh
scripts/check-p5-module-center-freeze.sh
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```

## 登录后人工验收基线

默认 Demo Article 页面：

```bash
MAKEADMIN_ENABLE_DEMO_MODULE=1 \
MAKEADMIN_ALLOW_MODULE_INSTALL_APPLY=1 \
MAKEADMIN_ALLOW_MODULE_UNINSTALL_APPLY=1 \
./scripts/dev-api.sh
```

多模块页面：

```bash
MAKEADMIN_ENABLE_DEMO_MODULE=1 \
MAKEADMIN_ENABLE_DEMO_NOTICE_MODULE=1 \
MAKEADMIN_ALLOW_MODULE_INSTALL_APPLY=1 \
MAKEADMIN_ALLOW_MODULE_UNINSTALL_APPLY=1 \
./scripts/dev-api.sh
```

登录后台后打开：

```text
http://127.0.0.1:5173/module
```

多模块关键字：

- `P5.25`
- `Demo Article`
- `Demo Notice`
- `多模块`
- `MAKEADMIN_ENABLE_DEMO_NOTICE_MODULE=1`
- `/demo/article`
- `/demo/notice`
- `未安装`
- `未注册`

当前本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此不自动登录、不伪造 token、不修改管理员密码。

## 保留边界

P5 冻结后仍不做：

- 不新增真实模块市场 UI。
- 不新增 `ma_demo_notice` 表。
- 不注册 `/demo_notice/*` 后端运行时路由。
- 不修改 `.env` 或管理员密码。
- 不连接或迁移 zyai 业务库。
- 不处理 PTLM 业务模块。

## 下一阶段建议

P6.1：模块中心产品化入口。建议从当前已冻结的 registry 多模块底座出发，开始整理模块市场/模块详情页/模块安装向导的产品化界面，而不是继续扩展临时 demo。

## 验收结果

- 已通过 `scripts/check-p5-module-center-freeze.sh`。
- 已通过 `scripts/check-module-registry-smoke.sh`。
- 已通过 `scripts/check-module-center-ui-contract.sh`。
- 已通过 `scripts/check-module-center-filter-contract.sh`。
- 已通过 `scripts/check-module-center-manual-checklist.sh`。
- 已通过 `scripts/check-demo-notice-module.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 已生成 `demo_notice` chunk。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 登录后人工 checklist 当前未执行；本地未提供 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，不修改 `.env` 或管理员密码。
