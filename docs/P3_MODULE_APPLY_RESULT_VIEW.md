# P3 Module Apply Result View

更新时间：2026-06-03

## 目标

P3.19 将管理端模块 manifest 安装、卸载结果视图提取为局部组件。

本阶段只调整前端组件结构，不改变页面行为，不改变接口协议，不改变后端写入门禁，不创建或修改数据库 schema。

## 组件

新增：

```text
admin/src/views/dev_tools/components/module-manifest-apply-result.vue
```

组件 props：

- `result`：`ModuleManifestApplyResult`。
- `fallbackTitle`：当前操作的兜底标题。

组件负责渲染：

- 结果 alert。
- 状态、模块、来源、环境变量。
- 操作类型、路由、权限编码。
- 执行前后快照。
- 门禁检查列表。

## 父组件变化

`module-manifest-preview.vue` 现在只保留：

- 安装结果 tab 入口。
- 卸载结果 tab 入口。
- apply 请求、状态和结果数据管理。

安装、卸载结果的重复模板和快照 helper 已移入新组件。

## 不在 P3.19 做

- 不改变按钮状态规则。
- 不改变错误归一化规则。
- 不修改请求 URL、method 或字段名。
- 不修改后端响应结构。
- 不修改后端写入门禁。
- 不创建、修改或迁移业务 schema。
- 不读取或修改 `.env`、密钥、CI/CD 或生产配置。
- 不新增权限 SQL。

## 验收

P3.19 需要通过：

```bash
cd admin && npm run type-check
scripts/check-module-tools-no-db.sh
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
