# P5 Module Center Apply

更新时间：2026-06-03

## 目标

P5.2 验证 Demo Article 不只可以通过脚本安装，也可以在后台模块中心完成真实 install/uninstall apply。

本阶段仍只针对本地 `go_makeadmin` 开发库，不创建业务 schema。

## 运行时

本地 API 需要临时开启：

```bash
MAKEADMIN_ENABLE_DEMO_MODULE=1 \
MAKEADMIN_ALLOW_MODULE_INSTALL_APPLY=1 \
MAKEADMIN_ALLOW_MODULE_UNINSTALL_APPLY=1 \
./scripts/dev-api.sh
```

不要写入 `.env`。

## 后台入口

模块中心：

```text
http://127.0.0.1:5173/module
```

Demo Article 页面：

```text
http://127.0.0.1:5173/demo/article
```

## 管理端改动

- 模块中心内置模块清单标记为 `P5.2`。
- Demo Article 状态从 `可预览` 调整为 `可安装`。
- Demo Article 增加页面入口 `/demo/article`。
- 内置模块清单操作区增加 `打开`，可以从模块中心进入 demo 页面。

## 人工测试

建议顺序：

1. 打开 `/module`。
2. 点击 Demo Article 的 `预览`。
3. 勾选 `安装写入`。
4. 点击 `安装执行`，确认结果为 `applied`。
5. 勾选 `删除确认`。
6. 点击 `卸载执行`，确认结果为 `applied`。
7. 再次勾选 `安装写入` 并点击 `安装执行`，把 demo 模块恢复为可见状态。
8. 点击 Demo Article 的 `打开`，确认进入 `/demo/article`。
9. 在 Demo Article 页面点击 `运行时详情`，确认 `runtimeRegistered=true`。

## 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `cd admin && npm run build`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过浏览器人工验证模块中心显示 `P5.2`、`可安装` 和 `/demo/article`。
- 已通过浏览器人工验证安装 apply、卸载 apply、重新安装 apply，结果均为 `applied`。
- 卸载 apply 后权限、菜单、菜单权限、角色权限快照回到 0。
- 重新安装 apply 后 demo 模块恢复为本地可见状态。
- 已通过浏览器人工验证模块中心 `打开` 可进入 `/demo/article`。
- 已通过浏览器人工验证 Demo Article 页面 `运行时详情` 显示 `module article` 和 `runtimeRegistered true`。

## 保留边界

P5.2 不做：

- 不创建或迁移 `ma_demo_article`。
- 不把写入 env 写进 `.env`。
- 不连接 zyai 真实业务库。
- 不把 demo 模块默认开启到生产运行时。
