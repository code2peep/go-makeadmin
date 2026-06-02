# P2 Data Scope

更新时间：2026-06-02

## 目标

P2.3 把 P1 已建立的角色数据范围表落到运行时查询链路，先保护后台核心列表型查询：

- 管理员列表。
- 登录日志列表。
- 操作日志列表。

本阶段不改数据库 schema，不导入新数据，不调整前端数据范围配置页面。

## 数据来源

数据权限从当前认证身份解析：

- `ma_admin_role`：读取当前管理员在当前租户下的角色。
- `ma_role_data_scope`：读取角色绑定的数据范围。
- `ma_data_scope`：只读取当前租户、启用、未软删除的数据范围。
- `ma_admin_org`：读取管理员主组织。
- `ma_org_unit`：解析组织树范围。

超级管理员不读取角色数据范围，直接获得当前租户内的全部数据权限。

## 解析规则

运行时解析为 `repository.DataScopeFilter`：

- `all`：当前租户内不追加管理员维度约束。
- `self`：只允许当前管理员自己的数据。
- `org`：允许当前管理员主组织下的管理员数据。
- `org_tree`：允许当前管理员主组织及下级组织内的管理员数据。
- `custom_org`：从 `scope_value` JSON 解析组织 ID，支持数组或对象字段 `org_ids`、`orgIds`、`ids`。

多个数据范围取并集；只要包含 `all`，直接按全部数据处理。

保守默认：

- 没有角色时按 `self` 处理。
- 有角色但没有绑定有效数据范围时按 `self` 处理。
- 数据范围无法解析出本人或组织 ID 时按 `NoAccess` 处理。

## 查询约束

Repository 层统一使用 `applyDataScopeFilter` 追加约束：

- `self`：追加 `admin_id = current_admin_id`。
- `org` / `org_tree` / `custom_org`：追加 `admin_id IN (SELECT admin_id FROM ma_admin_org WHERE tenant_id = ? AND org_id IN ? ...)`。
- 同时命中本人和组织范围时使用 `OR`。
- `NoAccess`：追加 `1 = 0`。
- `all` 或未启用数据范围：不追加管理员维度约束。

Adapter 如果在受保护查询入口拿不到 request identity，会按 `NoAccess` 处理，避免异常上下文放开列表查询。

当前已接入：

- `repository.AdminRepository.ListAdmins`
- `repository.LogRepository.ListLoginLogs`
- `repository.LogRepository.ListAuditLogs`

Adapter 从 request context 读取认证身份中的 `DataScope`，并传入对应 service filter。

## 不在 P2.3 做

- 不新增或迁移数据库表。
- 不改 `.env` 或生产配置。
- 不做前端角色数据范围配置页面。
- 不做详情、编辑、删除等资源级动作授权。
- 不开放非 `tenant_id=0` 登录或租户切换。

资源级动作授权和租户切换入口应进入后续 P2 任务。

## 验证

P2.3 已通过：

- `GOCACHE=/private/tmp/go-makeadmin-gocache go test ./...`
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`
- `./scripts/check-runtime-residue.sh`
- `./scripts/check-services.sh`
- `./scripts/check-p1-seed.sh`

已通过 `python3 -m py_compile scripts/p1-smoke.py`。

完整 P1 HTTP smoke 因本机未提供一次性 `P1_SMOKE_ADMIN_PASSWORD` 或 `ADMIN_PASSWORD` 环境变量未运行；没有密码变量时不读取 `.env` 猜测密码。
