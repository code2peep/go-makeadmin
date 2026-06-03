# P3 Module Apply UI Closure

更新时间：2026-06-03

## 目标

P3.14 将后台 `Manifest 预览` 弹窗与 P3.11/P3.13 的安装、卸载结构化结果闭环。

本阶段不改变后端写入边界，不新增数据库 schema，不新增权限 SQL。

## 管理端入口

`Manifest 预览` 弹窗现在支持：

- 安装执行。
- 卸载执行。
- 安装结果 tab。
- 卸载结果 tab。
- 门禁检查列表。
- 写入前后快照。
- 失败原因展示。

写入前后快照包含：

- 权限数。
- 菜单数。
- 菜单权限关联数。
- 角色授权数。

## API

前端新增：

```ts
applyModuleManifestUninstall(params)
```

对应：

```text
DELETE /gen/previewCode
```

安装仍使用：

```text
PUT /gen/previewCode
```

两个接口都复用 `gen:previewCode` 权限面。

## 页面行为

- 安装执行成功或失败都会写入安装结果 tab。
- 卸载执行成功或失败都会写入卸载结果 tab。
- 后端返回 `status=applied` 时展示成功态。
- 门禁失败、环境变量缺失或数据库目标不满足时展示 warning 态。
- 后端返回快照时展示执行前后计数。

## 不在 P3.14 做

- 不新增后端写入规则。
- 不创建、修改或迁移业务 schema。
- 不读取或修改 `.env`、密钥、CI/CD 或生产配置。
- 不新增权限 SQL。

## 验收

P3.14 需要通过：

```bash
cd admin && npm run type-check
scripts/check-module-tools-no-db.sh
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
