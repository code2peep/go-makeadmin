# P2 Module Guide

更新时间：2026-06-02

## 目标

本文是 P2 模块生命周期能力的统一入口，串联模块 manifest、注册 SQL、角色授权、安装、卸载、runtime 开关和验证命令。

P2 冻结状态见 `docs/P2_FINAL_STATUS.md`。

## 模块清单

模块清单放在：

```text
examples/<module>/manifest.json
```

当前 demo 模块：

```text
examples/demo/manifest.json
```

校验命令：

```bash
python3 scripts/check-module-manifests.py
```

## 注册预览

菜单和权限注册 dry-run：

```bash
python3 scripts/module-registry-plan.py --manifest examples/demo/manifest.json
```

角色授权 dry-run：

```bash
python3 scripts/module-role-grant-plan.py \
  --manifest examples/demo/manifest.json \
  --tenant-id 0 \
  --role-id 1
```

完整安装计划 dry-run：

```bash
python3 scripts/module-install-plan.py \
  --manifest examples/demo/manifest.json \
  --tenant-id 0 \
  --role-id 1
```

这些命令不连接数据库，不执行写入。

## 安装写入

只允许在本地 `go_makeadmin` 开发库受控执行：

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

安装写入只处理：

- `ma_permission`
- `ma_menu`
- `ma_menu_permission`
- 可选 `ma_role_permission`

不会创建 schema，不会修改 runtime 环境变量。

## 卸载预览

卸载 dry-run：

```bash
python3 scripts/module-uninstall-plan.py --manifest examples/demo/manifest.json
```

该命令只输出 SQL，不连接数据库，不执行删除。

## 卸载删除

只允许在本地 `go_makeadmin` 开发库受控执行：

```bash
MAKEADMIN_ALLOW_MODULE_UNINSTALL_WRITE=1 \
python3 scripts/module-uninstall-plan.py \
  --manifest examples/demo/manifest.json \
  --confirm-module article \
  --confirm-delete \
  --apply
```

卸载删除只处理：

- `ma_role_permission`
- `ma_menu_permission`
- `ma_menu`
- `ma_permission`

不会删除前端文件、后端代码、schema、runtime 环境变量或 codegen 元数据。

## runtime 开关

demo runtime 模块默认关闭。

启用方式：

```bash
MAKEADMIN_ENABLE_DEMO_MODULE=1
```

该开关只影响后端运行时路由是否挂载，不会由安装或卸载脚本自动修改。

## 验证

默认不触库验证：

```bash
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```

`verify-no-db` 已包含模块工具 no-db guard：

```bash
scripts/check-module-tools-no-db.sh
```

本地写库生命周期 smoke：

```bash
MAKEADMIN_ALLOW_MODULE_LIFECYCLE_WRITE=1 \
scripts/check-module-lifecycle-smoke.sh
```

生命周期 smoke 会安装 demo article、检查计数、卸载 demo article、确认残留为 0，并验证二次卸载 no-op。

## 红线边界

- 不连接真实 zyai 业务库执行写操作。
- 不读取或修改 `.env`。
- 不创建、修改或迁移 schema。
- 不删除文件或目录。
- 不部署生产环境。

## 详细文档

- `docs/P2_MODULE_REGISTRY.md`
- `docs/P2_MODULE_REGISTRY_SQL_PLAN.md`
- `docs/P2_MODULE_REGISTRY_APPLY.md`
- `docs/P2_MODULE_ROLE_GRANT_PLAN.md`
- `docs/P2_MODULE_RUNTIME_REGISTRY.md`
- `docs/P2_MODULE_INSTALL_PLAN.md`
- `docs/P2_MODULE_INSTALL_APPLY_BOUNDARY.md`
- `docs/P2_MODULE_INSTALL_APPLY.md`
- `docs/P2_MODULE_UNINSTALL_PLAN.md`
- `docs/P2_MODULE_UNINSTALL_APPLY_BOUNDARY.md`
- `docs/P2_MODULE_UNINSTALL_APPLY.md`
- `docs/P2_MODULE_LIFECYCLE_SMOKE.md`
- `docs/P2_MODULE_VERIFY_BOUNDARY.md`
- `docs/P2_FINAL_STATUS.md`
