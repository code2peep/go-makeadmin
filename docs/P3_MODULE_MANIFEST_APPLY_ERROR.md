# P3 Module Manifest Apply Error

更新时间：2026-06-03

## 目标

P3.18 统一管理端模块 manifest 安装、卸载 apply 的错误响应形态。

本阶段只调整前端错误结果归一化，不改变接口协议，不改变后端写入门禁，不创建或修改数据库 schema。

## 归一化入口

`admin/src/api/tools/code.ts` 新增：

```ts
normalizeModuleManifestApplyError(error, fallbackMessage)
```

该 helper 输出 `ModuleManifestApplyResult`，用于页面结果 tab 渲染。

## 处理规则

- 后端业务失败返回的 apply result 会保留 `summary`、`checks`、`before`、`after` 等结构化字段。
- 后端 result 缺少 `message` 时使用当前操作的 fallback message。
- 后端 result 缺少 `checks` 时补为空数组。
- `Error` 对象使用 `error.message`。
- 字符串错误直接作为 `message`。
- 其他未知错误使用 fallback message。

## 页面接入

`Manifest 预览` 弹窗现在在安装和卸载 catch 中调用同一个归一化 helper，不再直接把 `unknown error` 强转成 apply result。

## 不在 P3.18 做

- 不修改请求 URL、method 或字段名。
- 不修改后端响应结构。
- 不修改后端写入门禁。
- 不创建、修改或迁移业务 schema。
- 不读取或修改 `.env`、密钥、CI/CD 或生产配置。
- 不新增权限 SQL。

## 验收

P3.18 需要通过：

```bash
cd admin && npm run type-check
scripts/check-module-tools-no-db.sh
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
