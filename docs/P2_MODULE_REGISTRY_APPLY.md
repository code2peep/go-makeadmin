# P2 Module Registry Apply

更新时间：2026-06-02

## 目标

P2.13 开放模块注册 SQL 的本地受控写入模式，用来把模块 manifest 中声明的菜单和权限注册到本地 `go_makeadmin` 开发库。

写入模式只处理：

- `ma_permission` 缺失权限。
- `ma_menu` 缺失菜单节点。
- `ma_menu_permission` 菜单与主权限关联。

不处理：

- 角色授权。
- 运行时路由注册。
- demo 表创建。
- schema 变更。

## 命令

dry-run 预览：

```bash
python3 scripts/module-registry-plan.py --manifest examples/demo/manifest.json
```

受控写入：

```bash
MAKEADMIN_ALLOW_MODULE_REGISTRY_WRITE=1 \
python3 scripts/module-registry-plan.py \
  --manifest examples/demo/manifest.json \
  --confirm-module article \
  --apply
```

默认数据库连接参数：

- host：`127.0.0.1`
- port：`3306`
- user：`root`
- database：`go_makeadmin`

可通过命令参数或 `MYSQL_HOST`、`MYSQL_PORT`、`MYSQL_USER`、`MYSQL_DATABASE`、`MYSQL_PASSWORD` 环境变量覆盖。

## 写入门禁

`--apply` 必须同时满足：

- 环境变量 `MAKEADMIN_ALLOW_MODULE_REGISTRY_WRITE=1`。
- 显式传入 `--confirm-module <module>`。
- `--confirm-module` 必须等于 manifest 中的 `module`。

缺少环境变量或确认参数时，脚本会在数据库访问前失败。

## 写入规则

- 写入 SQL 来自 dry-run 同一套生成逻辑。
- 权限按 `ma_permission.code` 幂等插入。
- 菜单按 `ma_menu.route_name` 和 `delete_time=0` 幂等插入。
- 菜单权限关联按 `menu_id + permission_id` 幂等插入。
- 菜单父级通过 manifest `menu.parent` 对应的 `ma_menu.route_name` 查找。
- 父级菜单缺失时，菜单会以 `parent_id=0` 写入；后续模块安装器再统一处理父级策略。

## 本地 smoke

P2.13 已用本地 `go_makeadmin` 执行一次写入 smoke：

- 写入前确认 `article:*` 权限、`demo.article` 菜单和关联残留为 0。
- 第一次 apply 写入 5 条权限、1 条菜单、1 条菜单权限关联。
- 第二次 apply 后计数仍为 5 条权限、1 条菜单、1 条菜单权限关联，确认幂等。
- 清理 demo article 注册行。
- 清理后残留计数为 0。

## 验收

P2.13 需要通过：

```bash
python3 -m py_compile scripts/check-module-manifests.py scripts/module-registry-plan.py
python3 scripts/check-module-manifests.py
python3 scripts/module-registry-plan.py --apply
MAKEADMIN_ALLOW_MODULE_REGISTRY_WRITE=1 python3 scripts/module-registry-plan.py --apply
python3 scripts/module-registry-plan.py --manifest examples/demo/manifest.json
MAKEADMIN_ALLOW_MODULE_REGISTRY_WRITE=1 python3 scripts/module-registry-plan.py --manifest examples/demo/manifest.json --confirm-module article --apply
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
