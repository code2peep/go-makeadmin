# P5 Demo Module Visible

更新时间：2026-06-03

## 目标

P5.1 把 Demo Article 从“只存在于 manifest 和代码预览”推进到“本地后台可见页面”。

本阶段只针对本地 `go_makeadmin` 开发库，写入：

- `ma_permission`
- `ma_menu`
- `ma_menu_permission`
- `ma_role_permission`

不创建 `ma_demo_article` 表，不迁移业务 schema。

## 模块入口

管理端入口：

```text
http://127.0.0.1:5173/demo/article
```

数据库 `ma_menu.route_path` 保持为 `/dev_tools/demo/article`，前端动态路由会按父级菜单剥离后生成运行路由 `/demo/article`。

后端运行时需要显式开启：

```bash
MAKEADMIN_ENABLE_DEMO_MODULE=1 ./scripts/dev-api.sh
```

前端真实产物：

```text
admin/src/api/article.ts
admin/src/views/article/index.vue
admin/src/views/article/edit.vue
```

## 安装验证

P5.1 新增本地写入验证：

```bash
MAKEADMIN_ALLOW_DEMO_MODULE_VISIBLE_WRITE=1 scripts/check-demo-module-visible.sh
```

脚本会：

1. 校验 demo manifest。
2. 校验 demo manifest 菜单可见、运行时已注册、路径和组件正确。
3. 校验 demo article 前端文件存在。
4. 对本地 `go_makeadmin` 先执行 demo article 卸载 no-op 或清理。
5. 再执行 demo article 安装。
6. 校验 5 条权限、1 条可见菜单、1 条菜单权限、5 条角色授权。

脚本保留安装后的 demo article 菜单和权限，用于浏览器人工测试。

## 页面行为

Demo Article 页面当前是只读示例：

- 列表接口 `GET /article/list` 返回空分页。
- 详情接口 `GET /article/detail` 返回运行时状态。
- 新增、编辑、删除入口不写业务表。
- 后端写接口仍返回 `demo module is read-only`。

## 保留边界

P5.1 不做：

- 不创建 `ma_demo_article` 表。
- 不生成真实业务 CRUD schema。
- 不修改 `.env` 或系统环境变量。
- 不迁移 zyai 真实业务库。
- 不把 demo 模块默认开启到生产运行时。
