# P2 Module Verify Boundary

更新时间：2026-06-02

## 目标

P2.23 明确模块工具验证边界，把不触库检查接入 `verify-no-db.sh`，把真实写库生命周期 smoke 保持为显式命令。

## no-db 验证

新增脚本：

```bash
scripts/check-module-tools-no-db.sh
```

该脚本会执行：

- Python 脚本语法检查。
- shell 脚本语法检查。
- module manifest 校验。
- 注册、角色授权、安装、卸载 dry-run 预览。
- 注册、安装、卸载和生命周期脚本的写入门禁失败检查。

该脚本不会连接数据库，不会执行写入或删除。

## verify-no-db 接入

`scripts/verify-no-db.sh` 已在运行残留守卫之后接入：

```bash
scripts/check-module-tools-no-db.sh
```

因此默认 no-db 验证现在覆盖模块工具的语法、manifest、dry-run 和写入门禁。

## 写库 smoke 边界

真实写库生命周期 smoke 不进入 `verify-no-db.sh`。

需要显式执行：

```bash
MAKEADMIN_ALLOW_MODULE_LIFECYCLE_WRITE=1 \
scripts/check-module-lifecycle-smoke.sh
```

这样可以保持默认验证不触库，同时保留本地开发库的完整生命周期回归能力。

## 不在 P2.23 做

- 不把写库 smoke 放入 `verify-no-db.sh`。
- 不修改 `.env`。
- 不连接真实 zyai 业务库。
- 不创建、修改或迁移 schema。

## 验收

P2.23 需要通过：

```bash
bash -n scripts/verify-no-db.sh scripts/check-module-tools-no-db.sh scripts/check-module-lifecycle-smoke.sh
scripts/check-module-tools-no-db.sh
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
