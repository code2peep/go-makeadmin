# P3 Module Apply UI State

更新时间：2026-06-03

## 目标

P3.16 收敛后台 `Manifest 预览` 弹窗的模块安装、卸载按钮状态和结果清理规则。

本阶段只调整管理端交互，不改变后端写入逻辑，不新增权限 SQL，不创建或修改数据库 schema。

## 按钮状态

弹窗内安装、卸载继续复用 `gen:previewCode` 权限面。

预览相关按钮：

- `安装计划`：必须已有当前输入对应的 preview，且没有 apply 请求执行中。
- `代码预览`：必须已有当前输入对应的 preview，且没有 apply 请求执行中。

安装执行按钮：

- 必须已有当前输入对应的 preview。
- `确认模块` 必须等于 manifest module。
- 必须勾选 `安装写入`。
- manifest 声明 `requiresSchema=true` 时必须勾选 `Schema 风险`。
- 没有 preview、安装或卸载请求执行中。

卸载执行按钮：

- 必须已有当前输入对应的 preview。
- `确认模块` 必须等于 manifest module。
- 必须勾选 `删除确认`。
- 没有 preview、安装或卸载请求执行中。

## 加载态

- 预览请求执行中会锁住 apply 按钮。
- 安装请求执行中只展示安装按钮 loading。
- 卸载请求执行中只展示卸载按钮 loading。
- apply 请求执行中会锁住确认输入和确认勾选项。

## 结果清理

- 重新生成 preview 后清空旧安装结果和旧卸载结果。
- 修改来源模式、manifest 路径、manifest JSON、作者、租户或角色后清空当前 preview 和旧 apply 结果。
- 发起安装时清空旧安装结果。
- 发起卸载时清空旧卸载结果。

## 不在 P3.16 做

- 不修改后端写入门禁。
- 不新增后端接口。
- 不创建、修改或迁移业务 schema。
- 不读取或修改 `.env`、密钥、CI/CD 或生产配置。
- 不新增权限 SQL。

## 验收

P3.16 需要通过：

```bash
cd admin && npm run type-check
scripts/check-module-tools-no-db.sh
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
