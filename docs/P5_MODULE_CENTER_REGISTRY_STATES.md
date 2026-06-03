# P5 Module Center Registry States

更新时间：2026-06-03

## 目标

P5.13 收敛模块中心 registry 读取失败和空清单状态，让 API 异常或清单为空时页面仍能给出明确反馈和本地排查入口。

本阶段不新增后端接口，不新增数据库写入，不修改登录凭证。

## 管理端改动

模块中心新增 registry 状态处理：

- registry 读取失败时展示 `Registry 读取失败` alert。
- 失败 alert 描述里带上错误信息和 `scripts/check-module-registry-smoke.sh`。
- registry 读取成功但模块列表为空时展示 `Registry 暂无模块` alert。
- 表格空态根据当前状态切换为 `registry 读取失败`、`registry 暂无模块` 或 `暂无匹配模块`。
- registry 读取失败时不继续逐项读取安装状态。
- 模块中心阶段标识更新为 `P5.13`。

## 数据来源

页面状态只来自当前请求结果：

- `listModuleRegistry()` 成功或失败。
- 当前 registry rows 数量。
- 当前筛选后的模块列表。

## 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 登录后页面验收待重新登录后执行。

## 保留边界

P5.13 不做：

- 不伪造登录 token。
- 不修改 `.env` 或管理员密码。
- 不新增 registry 写入接口。
- 不新增数据库 schema。
- 不处理 PTLM 业务模块。
