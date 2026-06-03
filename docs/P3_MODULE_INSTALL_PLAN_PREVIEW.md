# P3 Module Install Plan Preview

更新时间：2026-06-03

## 目标

P3.9 将 P3.8 的 manifest 预览结果继续联动到模块安装计划，让后台生成器页面可以只读查看模块注册、角色授权、安装和卸载 SQL。

本阶段只生成 SQL 预览，不连接数据库，不执行 SQL，不修改 `ma_menu`、`ma_permission`、`ma_role_permission` 或 `ma_codegen_*`。

## 后端返回

`POST /gen/previewCode` 的返回新增 `plan`：

```json
{
  "plan": {
    "tenantId": 0,
    "roleId": 1,
    "registrySql": "...",
    "roleGrantSql": "...",
    "installSql": "...",
    "uninstallSql": "...",
    "runtimeHint": "MAKEADMIN_ENABLE_DEMO_MODULE=1"
  }
}
```

`registrySql` 对应模块菜单和权限注册预览。

`roleGrantSql` 对应把模块权限授予指定角色的预览。

`installSql` 是注册 SQL 和角色授权 SQL 的组合预览。

`uninstallSql` 是按 manifest 声明的权限 code 和菜单 routeName 生成的清理预览。

## 管理端入口

`Manifest 预览` 弹窗新增：

- 租户 ID 输入，默认 `0`。
- 角色 ID 输入，默认 `1`。
- 预览结果展示租户、角色和运行时提示。
- `安装计划` 按钮，用代码预览弹窗展示：
  - `registry.sql`
  - `role_grant.sql`
  - `install.sql`
  - `uninstall.sql`

## 边界

- 不执行安装计划 SQL。
- 不执行卸载计划 SQL。
- 不写数据库。
- 不创建或迁移业务 schema。
- 不读取或修改 `.env`。
- 不新增权限 SQL 或 seed 数据。
- 不改变旧 `GET /gen/previewCode` 请求和响应。

## 验收

P3.9 需要通过：

```bash
scripts/check-module-install-plan-preview.sh
scripts/check-module-manifest-preview.sh
scripts/check-module-tools-no-db.sh
cd admin && npm run type-check
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
