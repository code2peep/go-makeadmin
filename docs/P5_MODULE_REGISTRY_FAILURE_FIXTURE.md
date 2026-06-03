# P5 Module Registry Failure Fixture

更新时间：2026-06-03

## 目标

P5.8 增加一个本地可控的 registry 异常 fixture，用来验证模块中心异常筛选、校验明细弹窗和后端单项失败不中断列表。

本阶段不新增数据库写入，不新增 schema，不默认启用异常模块。

## 后端改动

`GET /api/gen/moduleRegistry` 在默认情况下仍只返回 `Demo Article`。

当显式设置以下环境变量时，会额外返回一个异常 registry 项：

```bash
MAKEADMIN_ENABLE_BROKEN_MODULE_REGISTRY_FIXTURE=1
```

异常项：

- `name=Broken Manifest Fixture`
- `module=broken_fixture`
- `manifest=examples/demo/missing/manifest.json`
- `runtime=MAKEADMIN_ENABLE_BROKEN_MODULE_REGISTRY_FIXTURE=1`

该 manifest 路径合法但文件不存在，因此会被 `manifestStatus=failed` 标记，并进入 `manifestChecks` 明细。

## 管理端改动

模块中心阶段标识更新为 `P5.8`。

P5.8 复用 P5.6/P5.7 已有能力：

- `异常` 筛选会包含 broken fixture。
- `校验` 列会显示 `异常`。
- `明细` 弹窗会展示失败检查项。

## 验收结果

- 已通过 `cd server && go test ./generator/service/gen`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前登录态已过期；页面人工复验需要重新登录后执行。

## 保留边界

P5.8 不做：

- 不默认启用 broken fixture。
- 不新增 registry 写入接口。
- 不新增数据库 schema。
- 不创建异常 manifest 文件。
- 不处理 PTLM 业务模块。
