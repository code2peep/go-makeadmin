# P2 Module Role Grant Plan

更新时间：2026-06-02

## 目标

P2.14 建立模块注册后的角色授权 dry-run 计划，用来把 manifest 中声明的模块权限授权给指定租户下的指定角色。

本阶段只输出 SQL，不连接数据库，不执行写入。

## 命令

```bash
python3 scripts/module-role-grant-plan.py \
  --manifest examples/demo/manifest.json \
  --tenant-id 0 \
  --role-id 1
```

默认 manifest：

```bash
python3 scripts/module-role-grant-plan.py --role-id 1
```

## 生成内容

SQL 预览包含：

- `ma_role_permission` 缺失授权插入。
- 按 `manifest.permissions[*].code` 查找权限 ID。
- 按 `tenant_id + role_id` 校验目标角色存在、启用且未软删除。
- 按 `tenant_id + role_id + permission_id` 防止重复授权。

## 安全规则

- 脚本不连接数据库。
- 脚本不执行 SQL。
- `--role-id` 必须显式传入，不提供默认角色。
- `--tenant-id` 默认 `0`，但必须是非负整数。
- `--role-id` 必须是正整数。
- 权限不存在、权限禁用或角色不可用时，生成 SQL 执行后不会插入对应授权。

## 不在 P2.14 做

- 不开放 apply/write 模式。
- 不自动授权 super_admin 或任意默认角色。
- 不创建角色。
- 不创建、修改或迁移 schema。
- 不注册 demo 运行时路由。

## 验证

P2.14 需要通过：

```bash
python3 -m py_compile scripts/check-module-manifests.py scripts/module-role-grant-plan.py
python3 scripts/check-module-manifests.py
python3 scripts/module-role-grant-plan.py --role-id 1
python3 scripts/module-role-grant-plan.py --role-id 0
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
