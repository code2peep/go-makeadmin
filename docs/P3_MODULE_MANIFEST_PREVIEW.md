# P3 Module Manifest Preview

更新时间：2026-06-03

## 目标

P3.8 将模块 manifest 与后台生成器页面打通，让管理端可以直接预览 manifest 对应的生成器配置和模板产物。

本阶段只做只读预览，不写 `ma_codegen_*`，不生成文件，不创建业务 schema，不修改旧 `GET /gen/previewCode` 行为。

## 后端接口

新增：

```http
POST /gen/previewCode
```

该接口复用现有 `gen:previewCode` 权限面。旧接口保持不变：

```http
GET /gen/previewCode?id=<table_id>
```

请求支持两种来源：

```json
{
  "manifestPath": "examples/demo/manifest.json",
  "authorName": "codepeep"
}
```

或：

```json
{
  "manifestBody": "{...manifest json...}",
  "authorName": "codepeep"
}
```

`manifestPath` 只允许读取仓库内的 `manifest.json`，避免后台预览接口变成任意文件读取入口。

## 返回内容

返回结构包含：

- `source`：manifest 来源。
- `manifest`：模块、实体、来源表和菜单名称摘要。
- `detail`：兼容旧 `/gen/detail` 的表配置和列配置。
- `code`：Go model、schema、service、route 和 Vue api、list、edit 的模板预览。
- `warning`：只读边界提示。

## 管理端入口

`admin/src/views/dev_tools/code/index.vue` 已新增 `Manifest 预览` 按钮。

弹窗支持：

- 仓库路径模式，默认 `examples/demo/manifest.json`。
- JSON 模式，粘贴 manifest 内容。
- 表配置摘要。
- 字段配置表格。
- 代码预览弹窗。

## 边界

- 不连接数据库。
- 不写 `ma_codegen_table` 或 `ma_codegen_column`。
- 不执行 `module-codegen-plan.py --apply`。
- 不生成或删除源码文件。
- 不读取或修改 `.env`。
- 不新增生成器权限 SQL。
- 不改变旧 `GET /gen/previewCode` 请求和响应。

## 验收

P3.8 需要通过：

```bash
scripts/check-module-manifest-preview.sh
scripts/check-module-tools-no-db.sh
cd admin && npm run type-check
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
