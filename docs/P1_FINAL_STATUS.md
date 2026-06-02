# P1 Final Status

更新时间：2026-06-02

## 结论

P1 已冻结，可以作为 `go-makeadmin` 后续自研框架的基础底座继续进入 P2。

P1 的默认运行目标是：

- 独立开发库：`go_makeadmin`
- 系统表前缀：`ma_*`
- 后台核心运行链路：`makeadmin` 新模型、repository、service、adapter
- 旧 `la_*` 蓝本：只作为参考资料保留，不再作为 P1 核心运行兜底

## 已完成范围

- P1 schema 和 seed：`sql/p1.schema.sql`、`sql/p1.seed.sql`
- P1 初始化脚本：`scripts/init-p1-db.sh`
- P1 种子只读检查：`scripts/check-p1-seed.sh`
- P1 HTTP smoke：`scripts/p1-smoke.py`
- 运行残留守卫：`scripts/check-runtime-residue.sh`
- 默认 no-db 验证链路：`scripts/verify-no-db.sh`
- 后台认证、菜单、权限、设置、字典、文件、日志、公共首页、代码生成器已切到 `ma_*`
- 核心后台接口不再回退旧 `server/admin/service/*` 运行链路
- 后台认证不再接受旧 `backstage:*` token
- 操作日志固定写入 `ma_audit_log`

## 冻结验收结果

本轮冻结验收在本地 `go_makeadmin` P1 开发库完成。

通过：

```bash
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
./scripts/check-services.sh
./scripts/check-p1-seed.sh
P1_SMOKE_ALLOW_WRITE=1 \
P1_SMOKE_BASE_URL=http://127.0.0.1:18082/api \
P1_SMOKE_ADMIN_PASSWORD='your-local-admin-password' \
python3 scripts/p1-smoke.py
```

验收记录：

- `verify-no-db` 通过，包含运行残留守卫、Go test、前端 type-check、前端 build、npm audit。
- `check-services` 通过，MySQL 和 Redis 可用。
- `check-p1-seed` 通过，live 种子计数为权限 79、菜单 22、设置 12、字典类型 4、字典项 14、文件分类 2。
- P1 HTTP smoke 通过，覆盖登录、自身信息、菜单路由、公共首页、登录日志、操作日志、角色、管理员、菜单、字典、文件上传和代码生成器。
- smoke 后 live 残留计数为 0：管理员、角色、菜单、字典、文件、代码生成和 `system:p1smoke:*` 权限均无 live 残留。
- 临时 API 使用 `SERVER_PORT=18082` 启动，验收后已停止。

已知验证噪音：

- 前端构建会输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 中 `/* #__PURE__ */` 注释位置的 warning；当前不影响构建退出码。

## 保留边界

P1 冻结不删除历史资料：

- `legacy/`
- `frontend/`
- `sql/install.sql`
- `sql/install.core.sql`
- `server/admin/service/*`
- `server/model/{system,setting,common}`

这些内容只作为蓝本参考、历史验证和后续迁移资料存在。新增 P1/P2 功能不得重新依赖旧运行链路。

## 不覆盖范围

- 不迁移 zyai 真实业务库。
- 不处理生产部署、CI/CD、npm publish 或线上发布。
- 不迁移旧 `la_*` 历史数据。
- 不删除旧蓝本文件。
- 不把 Redis token 作为长期最终认证模型。
- 不完成多租户隔离和数据权限查询约束，这些进入 P2。

## P2 入口

下一步进入 P2.1：认证模型升级。

P2.1 的目标是设计并落地 JWT + Redis session state，替换当前纯 Redis token 作为长期认证模型。该任务应先出认证模型文档，再改中间件、登录签发和 session 校验。
