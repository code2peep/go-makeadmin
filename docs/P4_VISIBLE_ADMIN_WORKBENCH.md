# P4 Visible Admin Workbench

更新时间：2026-06-03

## 目标

P4.1 把 `go-makeadmin` 的后台首页从旧业务指标改成基础框架工作台。

本阶段目标是让本地登录后的第一屏能直接回答三个问题：

- 当前框架做到哪个阶段。
- 底座能力有哪些。
- 人工测试从哪些入口开始。

## 后端接口

`GET /api/common/index/console` 当前返回：

- `version`：版本和技术栈。
- `framework`：阶段、数据库、表前缀、认证模型和模块生命周期。
- `milestones`：P1-P4 阶段状态。
- `validation`：本地验收状态。

接口不写库，不新增表，不读取 `.env`。

## 管理端工作台

`admin/src/views/workbench/index.vue` 当前展示：

- 框架标题和阶段标签。
- 核心后台、认证权限、模块闭环和当前版本。
- P1-P4 阶段状态。
- 本地验收状态。
- 人工测试入口。

人工测试入口使用真实动态菜单路由：

```text
/code
/menu
/role
/admin
/department
/information
```

## 人工验证

本阶段已在浏览器完成：

- 打开 `http://127.0.0.1:5173`。
- 本地 admin 登录。
- 打开新版工作台。
- 从工作台进入 `/code`。
- 打开 `Manifest 预览` 弹窗。
- 使用 `examples/demo/manifest.json` 生成预览成功。

## 当前边界

P4.1 不做：

- 不新增模块中心独立页面。
- 不新增审计落库。
- 不新增数据库 schema。
- 不修改 `.env`。
- 不新增权限 SQL。
- 不连接真实 zyai 业务库。

## 下一步

P4.2 建议做模块中心页面骨架，把 P3 的 manifest 预览、安装计划、安装执行、卸载执行和审计预览从代码生成器弹窗中抽出，形成后台可见的模块管理入口。
