# P2 Module Uninstall Apply

更新时间：2026-06-02

## 目标

P2.21 开放模块卸载计划的本地受控删除模式，用来清理 manifest 中声明的菜单、权限和相关授权。

删除模式只处理：

- `ma_role_permission` 中引用模块权限的授权。
- `ma_menu_permission` 中引用模块菜单或模块权限的关联。
- `ma_menu` 中模块菜单节点。
- `ma_permission` 中模块权限。

不处理：

- 前端文件删除。
- 后端代码删除。
- runtime 环境变量关闭。
- schema 删除或迁移。
- codegen 元数据删除。

## 命令

dry-run 预览：

```bash
python3 scripts/module-uninstall-plan.py --manifest examples/demo/manifest.json
```

受控删除：

```bash
MAKEADMIN_ALLOW_MODULE_UNINSTALL_WRITE=1 \
python3 scripts/module-uninstall-plan.py \
  --manifest examples/demo/manifest.json \
  --confirm-module article \
  --confirm-delete \
  --apply
```

## 删除门禁

`--apply` 必须同时满足：

- 环境变量 `MAKEADMIN_ALLOW_MODULE_UNINSTALL_WRITE=1`。
- 显式传入 `--confirm-module <module>`。
- `--confirm-module` 必须等于 manifest 中的 `module`。
- 显式传入 `--confirm-delete`。

缺少环境变量或确认参数时，脚本会在数据库访问前失败。

## 删除规则

- 删除前输出权限、菜单、菜单权限、角色授权四类快照计数。
- 四类计数均为 0 时报告 no-op，不执行删除。
- 删除在单个事务中执行。
- 删除顺序先关联表、后主表。
- 只按 manifest 权限 code 和菜单 routeName 清理。
- 不删除 manifest 未声明的文件、代码、schema 或 codegen 元数据。

## 本地 smoke

P2.21 已用本地 `go_makeadmin` 执行一次安装后卸载 smoke：

- 写入前确认 demo article 注册行和授权残留为 0。
- 先用 `module-install-plan.py --apply` 安装 demo article。
- 安装后得到 5 条权限、1 条菜单、1 条菜单权限关联、5 条角色授权。
- 再用 `module-uninstall-plan.py --apply` 卸载 demo article。
- 卸载后权限、菜单、菜单权限和角色授权计数均为 0。
- 第二次执行 uninstall apply 报告 no-op。

## 验收

P2.21 需要通过：

```bash
python3 -m py_compile scripts/check-module-manifests.py scripts/module-registry-plan.py scripts/module-uninstall-plan.py
python3 scripts/module-uninstall-plan.py --apply
MAKEADMIN_ALLOW_MODULE_UNINSTALL_WRITE=1 python3 scripts/module-uninstall-plan.py --apply
MAKEADMIN_ALLOW_MODULE_UNINSTALL_WRITE=1 python3 scripts/module-uninstall-plan.py --apply --confirm-module article
python3 scripts/module-uninstall-plan.py --manifest examples/demo/manifest.json
MAKEADMIN_ALLOW_MODULE_UNINSTALL_WRITE=1 python3 scripts/module-uninstall-plan.py --manifest examples/demo/manifest.json --confirm-module article --confirm-delete --apply
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
