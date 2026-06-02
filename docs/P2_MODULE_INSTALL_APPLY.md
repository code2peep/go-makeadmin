# P2 Module Install Apply

更新时间：2026-06-02

## 目标

P2.18 开放模块安装计划的本地受控写入模式，用来把 manifest 中声明的菜单、权限和可选角色授权写入本地 `go_makeadmin` 开发库。

写入模式只处理：

- `ma_permission` 缺失权限。
- `ma_menu` 缺失菜单节点。
- `ma_menu_permission` 菜单与主权限关联。
- 可选 `ma_role_permission` 角色授权。

不处理：

- schema 创建或迁移。
- runtime 开关修改。
- `.env` 或系统环境变量修改。
- 生产部署。

## 命令

dry-run 预览：

```bash
python3 scripts/module-install-plan.py --manifest examples/demo/manifest.json
```

受控写入，不包含角色授权：

```bash
MAKEADMIN_ALLOW_MODULE_INSTALL_WRITE=1 \
python3 scripts/module-install-plan.py \
  --manifest examples/demo/manifest.json \
  --confirm-module article \
  --apply
```

受控写入，包含角色授权：

```bash
MAKEADMIN_ALLOW_MODULE_INSTALL_WRITE=1 \
python3 scripts/module-install-plan.py \
  --manifest examples/demo/manifest.json \
  --tenant-id 0 \
  --role-id 1 \
  --confirm-module article \
  --confirm-role-id 1 \
  --apply
```

## 写入门禁

`--apply` 必须同时满足：

- 环境变量 `MAKEADMIN_ALLOW_MODULE_INSTALL_WRITE=1`。
- 显式传入 `--confirm-module <module>`。
- `--confirm-module` 必须等于 manifest 中的 `module`。
- manifest `requiresSchema` 必须为 `false`。
- 如果传入 `--role-id`，必须同时传入匹配的 `--confirm-role-id <id>`。

缺少环境变量或确认参数时，脚本会在数据库访问前失败。

## 写入规则

- 写入在单个事务中执行。
- 写入 SQL 复用 dry-run 的注册和角色授权 SQL。
- 权限按 `ma_permission.code` 幂等插入。
- 菜单按 `ma_menu.route_name + delete_time=0` 幂等插入。
- 菜单权限按 `menu_id + permission_id` 幂等插入。
- 角色授权按 `tenant_id + role_id + permission_id` 幂等插入。
- 已存在行不覆盖、不更新。
- runtime 开关只输出提示，不修改配置。

## 本地 smoke

P2.18 已用本地 `go_makeadmin` 执行一次写入 smoke：

- 写入前确认 demo article 注册行和 role grant 残留为 0。
- 第一次 apply 后得到 5 条权限、1 条菜单、1 条菜单权限关联、5 条角色授权。
- 第二次 apply 后计数仍为 5 条权限、1 条菜单、1 条菜单权限关联、5 条角色授权，确认幂等。
- 清理 demo article 注册行和 role grant。
- 清理后残留计数为 0。

## 验收

P2.18 需要通过：

```bash
python3 -m py_compile scripts/check-module-manifests.py scripts/module-registry-plan.py scripts/module-role-grant-plan.py scripts/module-install-plan.py
python3 scripts/module-install-plan.py --apply
MAKEADMIN_ALLOW_MODULE_INSTALL_WRITE=1 python3 scripts/module-install-plan.py --apply
MAKEADMIN_ALLOW_MODULE_INSTALL_WRITE=1 python3 scripts/module-install-plan.py --apply --confirm-module article --role-id 1
python3 scripts/module-install-plan.py --manifest examples/demo/manifest.json --tenant-id 0 --role-id 1
MAKEADMIN_ALLOW_MODULE_INSTALL_WRITE=1 python3 scripts/module-install-plan.py --manifest examples/demo/manifest.json --tenant-id 0 --role-id 1 --confirm-module article --confirm-role-id 1 --apply
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
