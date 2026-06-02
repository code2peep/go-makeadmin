# P1 Acceptance Checklist

## 目标

P1 的目标是把 `go-makeadmin` 从 LikeAdmin Go 蓝本过渡为可继续自研的 Go 后台框架底座。当前默认运行链路应面向独立开发库 `go_makeadmin` 和 `ma_*` 系统表。

## 当前默认入口

- 初始化 P1 独立库：`ADMIN_PASSWORD='your-local-admin-password' ./scripts/init-p1-db.sh`
- 只读检查 MySQL 和 Redis：`./scripts/check-services.sh`
- 只读检查 P1 种子：`./scripts/check-p1-seed.sh`
- 启动 API：`./scripts/dev-api.sh`
- 启动管理端：`./scripts/dev-admin.sh`
- P1 运行残留守卫：`./scripts/check-runtime-residue.sh`
- 不触库全量验证：`./scripts/verify-no-db.sh`
- P1 HTTP smoke：`P1_SMOKE_ALLOW_WRITE=1 P1_SMOKE_ADMIN_PASSWORD='your-local-admin-password' python3 scripts/p1-smoke.py`

## P1 已接管链路

- 认证登录：`ma_admin`、`ma_role`、`makeadmin:token:*`
- 管理员：`ma_admin`、`ma_admin_profile`、`ma_admin_role`、`ma_admin_org`
- 角色和权限：`ma_role`、`ma_permission`、`ma_role_permission`
- 菜单：`ma_menu`、`ma_menu_permission`
- 部门和岗位：`ma_org_unit`、`ma_position`
- 设置：`ma_setting`
- 字典：`ma_dict_type`、`ma_dict_item`
- 文件：`ma_file_category`、`ma_file`
- 日志：`ma_login_log`、`ma_audit_log`
- 代码生成器：`ma_codegen_table`、`ma_codegen_column`
- 公共首页：从 `ma_setting` 返回配置和控制台信息

## P1 验收命令

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/verify-no-db.sh
./scripts/check-runtime-residue.sh
./scripts/check-services.sh
./scripts/check-p1-seed.sh
```

本地 disposable P1 库写入 smoke：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
P1_SMOKE_ALLOW_WRITE=1 \
P1_SMOKE_BASE_URL=http://127.0.0.1:8000/api \
P1_SMOKE_ADMIN_PASSWORD='your-local-admin-password' \
python3 scripts/p1-smoke.py
```

## P1 不覆盖范围

- 不迁移 zyai 真实业务库。
- 不把旧 `la_*` 表作为默认运行模型。
- 不删除 `legacy`、旧模型、旧服务和蓝本 SQL。
- 不处理生产部署、远程对象存储、CI/CD 和发布流程。
- 不把 Redis token 作为长期认证最终方案。

## P1 完成标准

- 默认文档入口不再指向 P0 蓝本启动链路。
- 后台核心接口不再回退旧 `la_*` 服务。
- P1 运行残留守卫已接入 `verify-no-db`。
- `verify-no-db`、P1 seed 检查和 P1 HTTP smoke 均可通过。
- P1 运行残留边界已记录在 `docs/P1_RUNTIME_RESIDUE_AUDIT.md`。
- P2 入口任务已明确。

## P2 入口任务

1. P2.1：认证模型升级，设计并落地 JWT + Redis session state。
2. P2.2：多租户上下文 middleware，明确 `tenant_id` 来源、默认值和越权处理。
3. P2.3：数据权限查询约束，把角色数据范围落到 repository/service 查询层。
4. P2.4：代码生成器闭环，生成可编译、可挂载、可验证的示例模块。
5. P2.5：新增 `examples/demo` 标准模块，沉淀模块目录、路由、权限和前端页面约定。
