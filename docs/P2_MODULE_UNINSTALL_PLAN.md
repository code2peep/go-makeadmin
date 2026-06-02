# P2 Module Uninstall Plan

更新时间：2026-06-02

## 目标

P2.19 建立模块卸载/回滚 dry-run 计划，用来生成 manifest 对应菜单、权限和角色授权的清理 SQL 预览。

本阶段只输出 SQL，不连接数据库，不执行删除。

当前更新：P2.20 已定义卸载写入边界；未来开放 apply 前必须遵守 `docs/P2_MODULE_UNINSTALL_APPLY_BOUNDARY.md`。

## 命令

```bash
python3 scripts/module-uninstall-plan.py --manifest examples/demo/manifest.json
```

默认 manifest：

```bash
python3 scripts/module-uninstall-plan.py
```

## 生成内容

SQL 预览包含：

- 清理 `ma_role_permission` 中引用模块权限的授权。
- 清理 `ma_menu_permission` 中引用模块菜单或模块权限的关联。
- 清理 `ma_menu` 中模块菜单节点。
- 清理 `ma_permission` 中模块权限。

## 安全规则

- 脚本不连接数据库。
- 脚本不执行 SQL。
- 清理顺序先关联表、后主表。
- SQL 只使用 manifest 中声明的权限 code 和菜单 routeName。
- 本阶段不开放 apply/write 模式。

## 不在 P2.19 做

- 不执行删除。
- 不删除前端文件或后端代码。
- 不关闭 runtime 环境变量。
- 不创建、修改或迁移 schema。
- 不修改 `.env` 或生产配置。

## 验证

P2.19 需要通过：

```bash
python3 -m py_compile scripts/check-module-manifests.py scripts/module-registry-plan.py scripts/module-uninstall-plan.py
python3 scripts/check-module-manifests.py
python3 scripts/module-uninstall-plan.py --manifest examples/demo/manifest.json
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
