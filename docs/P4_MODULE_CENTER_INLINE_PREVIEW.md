# P4 Module Center Inline Preview

更新时间：2026-06-03

## 目标

P4.3 把模块中心的 manifest 预览从弹窗操作改为页面内工作区，让模块详情、字段、安装计划和代码预览能在同一个后台入口完成。

## 当前落地

- `admin/src/views/dev_tools/module/index.vue` 新增 manifest 输入区。
- 支持从仓库路径或 JSON body 生成预览。
- 页面内展示预览结果：来源、实体、表名、功能名、模板和运行时开关。
- 页面内展示字段明细：数据库字段、Go 字段、Go 类型、表单类型、查询类型和字典类型。
- `安装计划` 按钮打开当前预览的 registry、role grant、install、uninstall SQL。
- `代码预览` 按钮打开当前预览的后端和前端生成代码。
- 内置模块清单中的 `预览` 会写入该模块 manifest 路径，并生成页面内预览。

## 人工测试

本阶段已在浏览器完成：

- 打开 `http://127.0.0.1:5173/module`。
- 点击 `生成预览`，页面内显示 `DemoArticle`、`ma_demo_article` 和 `MAKEADMIN_ENABLE_DEMO_MODULE=1`。
- 点击 `安装计划`，代码预览弹窗显示 `registry.sql`、`role_grant.sql`、`install.sql`、`uninstall.sql`。
- 点击 `代码预览`，代码预览弹窗显示 `gocode/model.go`、`gocode/route.go`、`gocode/schema.go`、`api.ts`、`index.vue`、`edit.vue`。

## 当前边界

P4.3 不做：

- 不执行安装或卸载写入。
- 不新增数据库 schema。
- 不修改 `.env`。
- 不连接真实 zyai 业务库。

## 下一步

P4.4 建议把模块安装、卸载和 apply 结果摘要也内嵌到模块中心页面状态，减少用户在多个弹窗之间切换。
