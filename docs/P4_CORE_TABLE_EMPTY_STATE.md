# P4 Core Table Empty State

更新时间：2026-06-03

## 目标

P4.8 收敛核心后台表格的空态文案，让菜单、角色、管理员、部门和日志页面在无数据时显示统一可读的状态。

## 当前落地

以下表格新增 `empty-text="暂无数据"`：

- `admin/src/views/permission/menu/index.vue`
- `admin/src/views/permission/role/index.vue`
- `admin/src/views/permission/admin/index.vue`
- `admin/src/views/organization/department/index.vue`
- `admin/src/views/setting/system/journal.vue`

## 人工测试

本阶段已在浏览器完成：

- 打开 `/menu`，页面显示 `菜单名称`。
- 打开 `/role`，页面显示 `管理员人数`。
- 打开 `/admin`，页面显示 `管理员账号`。
- 打开 `/department`，页面显示 `部门名称`。
- 打开 `/journal`，页面显示 `访问链接`。
- 上述页面均未出现 403、404 或空白页。

## 当前边界

P4.8 不做：

- 不改接口错误处理。
- 不改分页逻辑。
- 不改数据查询条件。
- 不新增数据库 schema。
- 不修改 `.env`。
- 不连接真实 zyai 业务库。

## 下一步

P4.9 建议继续补接口失败态和加载态提示，让核心页面在请求失败时有明确反馈，而不是只依赖默认异常提示。
