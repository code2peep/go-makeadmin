# P4 Module Center

更新时间：2026-06-03

## 目标

P4.2 建立后台模块中心页面骨架，把 P3 的 manifest 预览、安装计划、安装执行、卸载执行和审计预览能力从代码生成器页面中抽出一个独立入口。

## 当前落地

- 新增 `admin/src/views/dev_tools/module/index.vue`。
- 工作台人工测试入口新增 `模块中心`。
- P1 seed 在 `开发工具` 下新增 `模块中心` 菜单。
- P1 seed 新增 `module:center:view` 权限。
- 当前本地 `go_makeadmin` 开发库已同步补入模块中心菜单和权限，用于人工验证。

## 页面能力

模块中心当前展示：

- `Manifest`：从仓库路径或 JSON 读取 manifest。
- `Codegen`：把 manifest 转换为生成器配置预览。
- `Install`：查看安装计划和本地受控安装 apply。
- `Audit`：查看 apply 结果摘要和审计 dry-run 预览。
- 内置 `Demo Article` 清单，指向 `examples/demo/manifest.json`。

## 人工验证

本阶段已在浏览器完成：

- 打开 `http://127.0.0.1:5173/module`。
- 模块中心页面显示 `Demo Article`。
- 点击 `Manifest 预览` 打开弹窗。
- 使用 `examples/demo/manifest.json` 生成预览成功。

## 当前边界

P4.2 不做：

- 不把 manifest 预览结果内嵌成模块中心主页面状态。
- 不新增审计落库。
- 不新增数据库 schema。
- 不修改 `.env`。
- 不连接真实 zyai 业务库。

## 下一步

P4.3 建议把模块中心的 manifest 预览结果从弹窗状态抽到页面状态，形成模块详情/安装计划/代码预览/审计预览的稳定工作区。
