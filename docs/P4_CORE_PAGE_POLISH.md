# P4 Core Page Polish

更新时间：2026-06-03

## 目标

P4.7 对核心后台页面做低风险细节修整，优先处理第一眼可见的文案、格式和窄视口表格可读性问题。

## 当前落地

- `admin/src/styles/public.scss` 增加全局表格 cell 换行规则。
- `admin/src/styles/public.scss` 增加描述表 label 不换行规则。
- `admin/src/views/setting/website/information.vue` 将 `网站logo` 改为 `网站 Logo`。
- `admin/src/views/setting/website/information.vue` 将网站信息相关注释从备案口径改为网站口径。
- `admin/src/views/organization/department/index.vue` 修正跑偏的 import 缩进。

## 人工测试

本阶段已在浏览器完成：

- 打开 `http://127.0.0.1:5173/information`。
- 页面显示 `网站 Logo`，不再显示 `网站logo`。
- 打开 `http://127.0.0.1:5173/department`。
- 页面显示 `部门名称`、`部门状态`、`新增` 和 `展开/折叠`。
- 浏览器验证 `.el-table .cell` 具备可换行样式，当前计算值为 `break-word`。

## 当前边界

P4.7 不做：

- 不重构核心 CRUD 页面。
- 不改接口和数据结构。
- 不新增数据库 schema。
- 不修改 `.env`。
- 不连接真实 zyai 业务库。

## 下一步

P4.8 建议继续做核心页面空态和加载态收敛，优先让菜单、角色、管理员、部门、日志在空数据、加载中、接口失败时都有稳定可读的表现。
