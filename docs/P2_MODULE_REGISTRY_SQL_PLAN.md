# P2 Module Registry SQL Plan

更新时间：2026-06-02

## 目标

P2.12 根据模块 manifest 生成菜单和权限初始化 SQL 预览。

本阶段只输出 SQL，不连接数据库，不执行写入。

当前更新：P2.13 已在该 dry-run 生成逻辑基础上开放本地受控写入模式；现行写入规则见 `docs/P2_MODULE_REGISTRY_APPLY.md`。

## 命令

```bash
python3 scripts/module-registry-plan.py --manifest examples/demo/manifest.json
```

默认 manifest：

```bash
python3 scripts/module-registry-plan.py
```

## 生成内容

SQL 预览包含：

- `ma_permission` 缺失权限插入。
- `ma_menu` 缺失菜单节点插入。
- `ma_menu_permission` 菜单和主权限关联插入。

脚本会先复用 `scripts/check-module-manifests.py` 校验 manifest。

## 安全规则

- 脚本不连接数据库。
- 脚本不执行 SQL。
- SQL 使用 `WHERE NOT EXISTS` 防止重复插入。
- 菜单父级通过 manifest `menu.parent` 对应的 `ma_menu.route_name` 查找。
- 本阶段不写 `ma_role_permission`，不自动给角色授权。

## 不在 P2.12 做

- 不执行 SQL。
- 不开放 apply/write 模式。
- 不创建、修改或迁移 schema。
- 不写菜单、权限、角色授权种子。
- 不注册 demo 运行时路由。

## 验证

P2.12 需要通过：

```bash
python3 -m py_compile scripts/check-module-manifests.py scripts/module-registry-plan.py
python3 scripts/check-module-manifests.py
python3 scripts/module-registry-plan.py --manifest examples/demo/manifest.json
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
