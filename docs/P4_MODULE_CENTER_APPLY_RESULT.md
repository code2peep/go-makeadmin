# P4 Module Center Apply Result

更新时间：2026-06-03

## 目标

P4.4 把模块安装、卸载和 apply 结果摘要迁入模块中心页面状态，减少用户在 manifest 预览弹窗、安装计划和结果弹窗之间切换。

## 当前落地

- `admin/src/views/dev_tools/module/index.vue` 复用模块 install/uninstall apply API。
- 预览结果区新增写入确认表单：
  - `确认模块`
  - `安装写入`
  - `Schema 风险`
  - `删除确认`
- 预览结果区新增 `安装执行` 和 `卸载执行` 操作。
- apply 结果以内嵌 tabs 展示：
  - `安装结果`
  - `卸载结果`
- apply 结果复用 `module-manifest-apply-result.vue`，包含状态、环境变量、权限、快照、检查项和审计预览。
- 当 manifest 输入变化时，页面会清空旧预览和旧 apply 结果，避免把过期 manifest 的结果继续展示为当前状态。

## 人工测试

本阶段已在浏览器完成：

- 打开 `http://127.0.0.1:5173/module`。
- 生成 `examples/demo/manifest.json` 预览。
- 勾选 `安装写入` 后点击 `安装执行`。
- 页面内展示 `安装结果`，当前本地未开启 `MAKEADMIN_ALLOW_MODULE_INSTALL_APPLY=1`，返回阻断结果且未访问数据库。
- 展开 `安装结果` 的 `审计预览`，能看到操作、模块、状态、路由、权限、检查、执行前和执行后摘要。
- 勾选 `删除确认` 后点击 `卸载执行`。
- 页面内展示 `卸载结果`，当前本地未开启 `MAKEADMIN_ALLOW_MODULE_UNINSTALL_APPLY=1`，返回阻断结果且未访问数据库。

## 当前边界

P4.4 不做：

- 不开启安装/卸载写入环境变量。
- 不实际写入或删除本地模块菜单权限。
- 不新增数据库 schema。
- 不修改 `.env`。
- 不连接真实 zyai 业务库。

## 下一步

P4.5 建议把模块中心的可见状态继续产品化：增加模块当前状态探测、安装条件提示和人工测试清单，让用户进入页面后能判断一个模块处于未安装、已安装、被阻断或需开启本地写入门禁。
