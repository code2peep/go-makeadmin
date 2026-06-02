# P2 Status

更新时间：2026-06-02

## 当前阶段

P2：框架能力增强。

P2 从 P1 冻结底座继续推进，不再扩大 P1 范围。P2 的重点是把框架能力从“能运行的后台底座”推进到“可复用的多项目后台框架”。

## P2.1 当前落地

认证模型已从 P1 的纯 Redis opaque token 升级为 JWT + Redis session state：

- 登录返回 JWT，前端继续通过 header `token` 传递。
- JWT 使用 HS256，签名密钥来自 `config.Config.Secret`。
- JWT payload 包含 `sid`、`adminId`、`tenantId`、`iat`、`exp`、`iss`。
- Redis 保存 `sid -> adminId` 的 session state，用于服务端有效性校验和登出吊销。
- Redis session set 会跟随 session TTL 过期，避免积累已过期 sid。
- 中间件先校验 JWT，再查 Redis session state，再从 `ma_admin` 重建实时身份和权限。
- 登出删除 Redis session state，JWT 即使尚未过期也不能继续访问。
- `scripts/check-runtime-residue.sh` 已增加旧 makeadmin opaque token key 回流检查。

## P2.1 验收标准

- `go test ./...` 通过。
- `./scripts/check-runtime-residue.sh` 通过。
- `./scripts/verify-no-db.sh` 通过。
- P1 HTTP smoke 继续通过，证明前端 token header 形状兼容。

## P2.1 验收结果

- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache go test ./...`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `./scripts/check-services.sh` 和 `./scripts/check-p1-seed.sh`。
- 已在本地 18082 临时 API 上通过 `scripts/p1-smoke.py`，包含 JWT claims 和 logout 校验。
- smoke 后 live 残留计数为 0，临时 API 已停止。

## 不在 P2.1 做

- 不增加数据库 session 表。
- 不做 refresh token。
- 不做多端设备列表和批量踢出。
- 不修改前端 header 名称。
- 不处理生产密钥轮换。

## P2.2 当前落地

多租户上下文 middleware 已建立：

- 默认租户仍为 `tenant_id=0`。
- 登录阶段不开放非 `0` 租户切换。
- 认证后请求以 JWT `tenantId` 为可信来源。
- 可选 header `X-Tenant-ID` 必须与 JWT `tenantId` 一致，否则直接无权限。
- 认证中间件会把租户上下文写入 request context 和 gin context。
- makeadmin adapter 中的租户隔离链路已从硬编码 `GlobalTenantID` 改为读取上下文。
- 操作日志写入当前请求租户。
- P1 HTTP smoke 增加 `X-Tenant-ID` mismatch guard。

详见 `docs/P2_TENANT_CONTEXT.md`。

## P2.2 验收标准

- `go test ./...` 通过。
- `./scripts/check-runtime-residue.sh` 通过。
- `./scripts/verify-no-db.sh` 通过。
- P1 HTTP smoke 继续通过，并包含租户 header 越权校验。

## P2.2 验收结果

- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache go test ./...`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `./scripts/check-runtime-residue.sh`，并新增 adapter/middleware 禁止硬编码 `GlobalTenantID` 的防回流检查。
- 已通过 `./scripts/check-services.sh` 和 `./scripts/check-p1-seed.sh`。
- 已在本地 18082 临时 API 上通过最小 HTTP guard：合法 JWT + Redis session 携带错误 `X-Tenant-ID` 返回 `403`。
- 临时 Redis session key 已清理，临时 API 已停止。
- 完整 `scripts/p1-smoke.py` 写入 smoke 因本机未提供 `P1_SMOKE_ADMIN_PASSWORD` 或 `ADMIN_PASSWORD` 环境变量未运行；脚本已补租户 mismatch guard 并通过 `python3 -m py_compile`。

## P2.3 当前落地

数据权限查询约束已接入核心列表型查询：

- 认证身份会根据当前租户角色解析 `ma_data_scope` / `ma_role_data_scope`。
- 超级管理员获得当前租户内全部数据权限。
- 普通管理员支持 `all`、`self`、`org`、`org_tree`、`custom_org`。
- 多个数据范围取并集，`all` 优先。
- 无角色或无有效数据范围时保守回退为 `self`。
- 数据范围无法解析出本人或组织 ID 时回退为 `NoAccess`。
- Adapter 会把认证身份同步写入 request context，供查询入口读取 `DataScope`。
- Adapter 缺失 request identity 时按 `NoAccess` 处理，避免异常上下文放开列表查询。
- 管理员列表、登录日志、操作日志已在 repository 层应用数据范围约束。

详见 `docs/P2_DATA_SCOPE.md`。

## P2.3 验收标准

- `go test ./...` 通过。
- `./scripts/verify-no-db.sh` 通过。
- `./scripts/check-runtime-residue.sh` 通过。
- `./scripts/check-services.sh` 通过。
- `./scripts/check-p1-seed.sh` 通过。
- P1 HTTP smoke 如本机提供一次性密码变量则继续通过；未提供时不读取 `.env` 猜测密码。

## P2.3 验收结果

- 已通过 `bash -n scripts/check-runtime-residue.sh scripts/verify-no-db.sh scripts/check-services.sh scripts/check-p1-seed.sh`。
- 已通过 `python3 -m py_compile scripts/p1-smoke.py`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache go test ./...`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `./scripts/check-services.sh` 和 `./scripts/check-p1-seed.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 完整 `scripts/p1-smoke.py` 写入 smoke 因本机未提供 `P1_SMOKE_ADMIN_PASSWORD` 或 `ADMIN_PASSWORD` 环境变量未运行。

## 下一步

P2.4：租户成员和租户切换入口。
