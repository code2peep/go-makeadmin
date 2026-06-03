# P5 Module Center Acceptance Status

更新时间：2026-06-03

## 目标

P5.12 在模块中心页面增加验收辅助状态，让页面直接展示 registry 来源、模块数量、校验异常数量、broken fixture 状态、自动 smoke 命令和人工验收入口。

本阶段不新增后端接口，不新增数据库写入，不修改登录凭证。

## 管理端改动

模块中心 `内置模块清单` 区域新增只读状态条：

- `来源`：固定显示 `/api/gen/moduleRegistry`。
- `模块`：显示当前 registry 返回模块数量。
- `校验异常`：显示当前 registry 校验失败数量。
- `Broken Fixture`：根据返回模块中是否存在 `broken_fixture` 显示 `已开启` 或 `未开启`。
- `Smoke`：显示 `check-module-registry-smoke.sh`。
- `人工入口`：显示 `/module`。

模块中心阶段标识更新为 `P5.12`。

## 数据来源

状态条只消费当前页面已有数据：

- `listModuleRegistry()` 返回的 registry rows。
- 每行的 `module`、`manifestStatus` 和安装状态回读结果。

不新增 API，不读取 `.env`，不猜测服务端环境变量。

## 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 登录后页面验收待重新登录后执行。

## 保留边界

P5.12 不做：

- 不伪造登录 token。
- 不修改 `.env` 或管理员密码。
- 不新增 registry 写入接口。
- 不新增数据库 schema。
- 不处理 PTLM 业务模块。
