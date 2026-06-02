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

## P2.4 当前落地

租户成员校验和后端租户切换入口已建立：

- 登录阶段允许解析 `X-Tenant-ID`，非 `0` 租户由 auth service 校验成员关系。
- `tenant_id=0` 保持 P1/P2 默认兼容上下文，不要求 `ma_tenant_member`。
- 非 `0` 租户必须存在启用租户和启用成员关系。
- 认证中间件每次重建身份都会重新校验租户成员关系，租户或成员失效后旧 token 不能继续访问该租户。
- 登录响应新增 `tenantId` 字段，旧前端只读取 `token` 不受影响。
- 新增 `GET /system/tenant/list`，返回当前管理员可访问租户。
- 新增 `POST /system/tenant/switch`，切换成功后签发目标租户的新 JWT 和 Redis session。
- P1 HTTP smoke 脚本补充默认租户列表和切换到默认租户检查；本阶段未提供一次性密码变量时只做脚本编译检查。

详见 `docs/P2_TENANT_SWITCH.md`。

## P2.4 验收标准

- `go test ./...` 通过。
- `./scripts/verify-no-db.sh` 通过。
- `./scripts/check-services.sh` 通过。
- `./scripts/check-p1-seed.sh` 通过。
- `scripts/p1-smoke.py` 语法编译通过。
- P1 HTTP smoke 如本机提供一次性密码变量则继续通过；未提供时不读取 `.env` 猜测密码。

## P2.4 验收结果

- 已通过 `python3 -m py_compile scripts/p1-smoke.py`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache go test ./...`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `./scripts/check-services.sh` 和 `./scripts/check-p1-seed.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 完整 `scripts/p1-smoke.py` 写入 smoke 因本机未提供 `P1_SMOKE_ADMIN_PASSWORD` 或 `ADMIN_PASSWORD` 环境变量未运行。

## P2.5 当前落地

租户级设置、文件和日志的迁移策略已沉淀：

- 明确 `ma_setting` 当前已按 `tenant_id + setting_group + setting_key` 隔离。
- 明确非 `0` 租户不会自动回退读取租户 `0` 设置，新租户上线前必须准备自己的设置基线。
- 明确 `ma_file_category` 和 `ma_file` 元数据已按 `tenant_id` 隔离，物理上传路径本阶段不改。
- 明确登录日志和操作日志不做跨租户复制，历史 `tenant_id=0` 日志保留默认租户归属。
- 明确云存储 `secretKey` / `accessKey` 不默认复制，必须显式 opt-in 或重新配置。
- `scripts/check-runtime-residue.sh` 增加设置、文件和日志 service/repository 不得硬编码 `GlobalTenantID` 的防回流检查。

详见 `docs/P2_TENANT_MIGRATION.md`。

## P2.5 验收标准

- `scripts/check-runtime-residue.sh` 通过，并包含租户级设置/文件/日志防回流检查。
- `./scripts/verify-no-db.sh` 通过。
- `./scripts/check-services.sh` 通过。
- `./scripts/check-p1-seed.sh` 通过。
- 不执行数据迁移、不修改 schema、不读取或修改 `.env`。

## P2.5 验收结果

- 已通过 `bash -n scripts/check-runtime-residue.sh scripts/verify-no-db.sh scripts/check-services.sh scripts/check-p1-seed.sh`。
- 已通过 `python3 -m py_compile scripts/p1-smoke.py`。
- 已通过 `./scripts/check-runtime-residue.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `./scripts/check-services.sh` 和 `./scripts/check-p1-seed.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有执行租户数据复制、schema 变更或 `.env` 修改。

## P2.6 当前落地

租户初始化命令 dry-run 版本已建立：

- 新增 `scripts/tenant-init-plan.py`，用于生成目标租户初始化 SQL 预览。
- 默认从源租户读取 `ma_setting` 的 `website`、`protocol`、`storage` 三类设置。
- 默认从源租户读取未软删除的 `ma_file_category` 文件分类。
- 目标租户已存在的 setting key 和文件分类 code 会跳过，不生成覆盖 SQL。
- `storage` JSON 配置中的 `secretKey` / `accessKey` 默认清空。
- 只有显式传入 `--copy-secret` 才保留云存储密钥字段在 SQL 预览里。
- 脚本只通过本机 `mysql` client 执行查询，不执行输出 SQL。
- 脚本不读取 `.env`，数据库连接参数来自 `MYSQL_*` 环境变量或命令参数默认值。

详见 `docs/P2_TENANT_INIT_PLAN.md`。

## P2.6 验收标准

- `python3 -m py_compile scripts/tenant-init-plan.py scripts/p1-smoke.py` 通过。
- `python3 scripts/tenant-init-plan.py --help` 通过。
- `python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 999999 --sql-only` 通过，且只输出 SQL 预览。
- `./scripts/verify-no-db.sh` 通过。
- `./scripts/check-services.sh` 通过。
- `./scripts/check-p1-seed.sh` 通过。
- 不执行数据写入、不修改 schema、不读取或修改 `.env`。

## P2.6 验收结果

- 已通过 `python3 -m py_compile scripts/tenant-init-plan.py scripts/p1-smoke.py`。
- 已通过 `python3 scripts/tenant-init-plan.py --help`。
- 已通过 `python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 999999 --sql-only`；本机预览生成 12 条 setting 和 2 条文件分类插入 SQL，没有执行 SQL。
- 已通过 `bash -n scripts/check-runtime-residue.sh scripts/verify-no-db.sh scripts/check-services.sh scripts/check-p1-seed.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `./scripts/check-services.sh` 和 `./scripts/check-p1-seed.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有执行租户数据写入、schema 变更或 `.env` 修改。

## 下一步

P2.7：租户初始化 apply/write 模式设计与安全门禁。该任务会触及数据库写入红线，进入前需要明确授权。
