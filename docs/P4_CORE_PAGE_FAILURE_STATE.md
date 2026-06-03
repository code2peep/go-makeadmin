# P4 Core Page Failure State

更新时间：2026-06-03

## 目标

P4.9 补核心页面请求失败后的状态复位，避免手写 loading 的页面在接口失败时卡住。

## 当前落地

- `admin/src/views/permission/menu/index.vue`
  - `getLists` 改为 `try/catch/finally`。
  - 请求失败时清空列表。
  - 无论成功或失败都复位 `loading`。
- `admin/src/views/organization/department/index.vue`
  - `getLists` 改为 `try/catch/finally`。
  - 请求失败时清空列表。
  - 无论成功或失败都复位 `loading`。

## 人工测试

本阶段已在浏览器完成：

- 打开 `/menu`，页面显示 `菜单名称`。
- 打开 `/department`，页面显示 `部门名称`。
- 两个页面均未出现 403、404 或空白页。

## 当前边界

P4.9 不做：

- 不改全局请求拦截器。
- 不给每个 usePaging 页面额外加本地错误弹窗，避免和现有全局错误提示重复。
- 不模拟接口失败写浏览器测试。
- 不新增数据库 schema。
- 不修改 `.env`。
- 不连接真实 zyai 业务库。

## 下一步

P4.10 建议进入 P4 冻结前的总体验收：梳理 P4 已完成内容、剩余缺口、人工测试账号和本地启动说明，判断是否可以从“可见后台”进入 P5。
