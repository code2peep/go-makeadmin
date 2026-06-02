# P2 Module Install Plan

更新时间：2026-06-02

## 目标

P2.16 建立模块安装 dry-run 编排，把 manifest 校验、菜单权限注册 SQL、角色授权 SQL 和 runtime 开关集中输出为一份可审阅计划。

本阶段只输出计划，不连接数据库，不执行写入。

## 命令

不包含角色授权：

```bash
python3 scripts/module-install-plan.py --manifest examples/demo/manifest.json
```

包含角色授权预览：

```bash
python3 scripts/module-install-plan.py \
  --manifest examples/demo/manifest.json \
  --tenant-id 0 \
  --role-id 1
```

默认 manifest：

```bash
python3 scripts/module-install-plan.py
```

## 输出内容

安装计划包含：

- manifest 基础信息。
- 后端路由和权限映射。
- 前端 API 和页面路径。
- `ma_permission` / `ma_menu` / `ma_menu_permission` 注册 SQL。
- 可选 `ma_role_permission` 角色授权 SQL。
- runtime 环境变量提示。

## 安全规则

- 脚本不连接数据库。
- 脚本不执行 SQL。
- 注册 SQL 复用 `scripts/module-registry-plan.py` 的生成逻辑。
- 角色授权 SQL 复用 `scripts/module-role-grant-plan.py` 的生成逻辑。
- `--role-id` 不提供默认值；未传时不生成角色授权 SQL。
- `--tenant-id` 默认 `0`，但必须是非负整数。
- 本阶段不开放 apply/write 模式。

## 不在 P2.16 做

- 不执行注册写入。
- 不执行角色授权写入。
- 不创建、修改或迁移 schema。
- 不修改 `.env` 或生产配置。
- 不启动服务或做生产部署。

## 验证

P2.16 需要通过：

```bash
python3 -m py_compile scripts/check-module-manifests.py scripts/module-registry-plan.py scripts/module-role-grant-plan.py scripts/module-install-plan.py
python3 scripts/module-install-plan.py --manifest examples/demo/manifest.json
python3 scripts/module-install-plan.py --manifest examples/demo/manifest.json --tenant-id 0 --role-id 1
python3 scripts/module-install-plan.py --role-id 0
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
