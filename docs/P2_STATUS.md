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

## P2.7 当前落地

租户初始化 apply/write 模式的安全门禁已建立：

- `scripts/tenant-init-plan.py` 新增预留 `--apply` 参数。
- 当前传入 `--apply` 会在任何数据库访问前失败。
- dry-run SQL 预览模式保持不变。
- 明确后续真正开放写入需要额外 `MAKEADMIN_ALLOW_TENANT_INIT_WRITE=1` 和 `--confirm-to-tenant <id>` 双门禁。
- 明确未来写入必须在单事务内完成，只插入缺失 setting 和文件分类，不覆盖目标租户已有配置。
- 明确本阶段不执行 SQL、不修改 schema、不修改 `.env`。

详见 `docs/P2_TENANT_INIT_APPLY_GUARD.md`。

## P2.7 验收标准

- `python3 -m py_compile scripts/tenant-init-plan.py scripts/p1-smoke.py` 通过。
- `python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 999999 --apply` 失败，且错误说明没有访问数据库。
- `python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 999999 --sql-only` 继续通过。
- `./scripts/verify-no-db.sh` 通过。
- `./scripts/check-services.sh` 通过。
- `./scripts/check-p1-seed.sh` 通过。
- 不执行数据写入、不修改 schema、不读取或修改 `.env`。

## P2.7 验收结果

- 已通过 `python3 -m py_compile scripts/tenant-init-plan.py scripts/p1-smoke.py`。
- 已通过 `python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 999999 --apply` 失败门禁；失败文案明确没有访问数据库。
- 已通过 `python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 999999 --sql-only`；dry-run SQL 预览继续可用。
- 已通过 `bash -n scripts/check-runtime-residue.sh scripts/verify-no-db.sh scripts/check-services.sh scripts/check-p1-seed.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `./scripts/check-services.sh` 和 `./scripts/check-p1-seed.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有执行租户数据写入、schema 变更或 `.env` 修改。

## P2.8 当前落地

租户初始化 apply/write 模式已开放为本地受控写入：

- `scripts/tenant-init-plan.py --apply` 现在支持执行缺失初始化行写入。
- 写入必须同时满足 `MAKEADMIN_ALLOW_TENANT_INIT_WRITE=1` 和 `--confirm-to-tenant <id>`。
- 缺少环境变量或确认参数时，脚本会在数据库访问前失败。
- apply 前会校验目标租户已存在、启用且未软删除。
- 写入在单个事务中执行。
- 只插入目标租户缺失的 `ma_setting` 和 `ma_file_category` 行。
- 已存在 setting key 和文件分类 code 只跳过，不覆盖、不更新。
- `storage` JSON 配置里的 `secretKey` / `accessKey` 默认清空，除非显式传入 `--copy-secret`。
- 本阶段不创建租户、不创建租户成员、不迁移 `ma_file` 元数据、不迁移物理上传文件、不修改 schema。

详见 `docs/P2_TENANT_INIT_APPLY.md`。

## P2.8 验收标准

- `python3 -m py_compile scripts/tenant-init-plan.py scripts/p1-smoke.py` 通过。
- `python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 999999 --apply` 失败，且错误说明没有访问数据库。
- `MAKEADMIN_ALLOW_TENANT_INIT_WRITE=1 python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 999999 --apply` 失败，且错误说明缺少确认参数。
- `python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 999999 --sql-only` 继续通过。
- 对本地 `go_makeadmin` 执行一次临时租户写入 smoke，并清理临时行。
- `./scripts/verify-no-db.sh` 通过。
- `./scripts/check-services.sh` 通过。
- `./scripts/check-p1-seed.sh` 通过。
- 不修改 schema、不读取或修改 `.env`、不连接真实 zyai 业务库。

## P2.8 验收结果

- 已通过 `python3 -m py_compile scripts/tenant-init-plan.py scripts/p1-smoke.py`。
- 已通过 `python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 999999 --apply` 失败门禁；失败文案明确没有访问数据库。
- 已通过 `MAKEADMIN_ALLOW_TENANT_INIT_WRITE=1 python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 999999 --apply` 失败门禁；失败文案明确没有访问数据库。
- 已通过 `python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 999999 --sql-only`；dry-run SQL 预览继续可用。
- 已对本地 `go_makeadmin` 临时租户 `990028` 完成写入 smoke：第一次插入 12 条 setting 和 2 条文件分类，第二次插入 0 条并跳过已有 12 条 setting 和 2 条文件分类。
- 已校验临时租户 storage 密钥字段没有非空复制。
- 已清理临时租户、setting 和文件分类行，清理后残留计数为 0。
- 已通过 `bash -n scripts/check-runtime-residue.sh scripts/verify-no-db.sh scripts/check-services.sh scripts/check-p1-seed.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `./scripts/check-services.sh` 和 `./scripts/check-p1-seed.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有修改 schema、没有读取或修改 `.env`、没有连接真实 zyai 业务库。

## P2.9 当前落地

代码生成器 Go 输出闭环已建立：

- 修复 `route.go.tpl` 中 list handler 使用错误接收者的问题。
- `schema.go.tpl` 只在需要 `core.*` 类型时导入 `go-makeadmin/core`，避免生成无用导入。
- `service.go.tpl` 只在需要 URL 绝对化时导入 `go-makeadmin/util`，避免生成无用导入。
- `service.go.tpl` 的 `Detail` / `Del` 主键参数类型跟随生成主键类型。
- `EditReq` 始终包含主键字段，保证编辑逻辑可编译。
- 新增 `server/generator/tpl_test.go`，渲染 CRUD Go 模板到临时目录并执行 `go test .` 编译生成包。
- 新增 `examples/README.md` 和 `examples/demo/`，沉淀标准 CRUD 模块接入约定。

详见 `docs/P2_CODEGEN_CLOSURE.md`。

## P2.9 验收标准

- `GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator ./generator/service/gen` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache go test ./...` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不创建 demo 表、不写菜单或权限种子、不默认注册运行时路由。

## P2.9 验收结果

- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator ./generator/service/gen`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache go test ./...`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有创建 demo 表、没有写菜单或权限种子、没有默认注册运行时路由。

## P2.10 当前落地

前端生成模板闭环已建立：

- 新增 `scripts/check-codegen-frontend.sh`。
- 新增 env-gated 测试 `TestGeneratedCrudFrontendCodeTypeChecks`，默认后端测试不触发 Node。
- 显式运行脚本时，测试会渲染 `api.ts`、`index.vue`、`edit.vue`。
- 测试临时写入 `admin/src/api/article.ts` 和 `admin/src/views/article/`。
- 测试执行 `npm run type-check`，验证生成前端代码符合当前 admin TypeScript/Vue 约定。
- 测试结束后清理临时生成文件。
- `examples/demo` 已补充前端生成模板验证说明。

详见 `docs/P2_FRONTEND_CODEGEN_CLOSURE.md`。

## P2.10 验收标准

- `./scripts/check-codegen-frontend.sh` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 临时生成的 `admin/src/api/article.ts` 和 `admin/src/views/article/` 不残留。

## P2.10 验收结果

- 已通过 `bash -n scripts/check-codegen-frontend.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator -run TestGeneratedCrudFrontendCodeTypeChecks`；默认未设置环境变量时前端 type-check 测试会跳过。
- 已通过 `./scripts/check-codegen-frontend.sh`；脚本临时生成前端 API 和页面并执行 `npm run type-check`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已确认临时生成的 `admin/src/api/article.ts` 和 `admin/src/views/article/` 没有残留。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。

## P2.11 当前落地

模块注册清单规范已建立：

- `examples/<module>/manifest.json` 作为模块注册清单。
- manifest 统一描述后端路由、前端 API、前端页面、菜单节点、权限元数据、运行时注册状态和 schema 需求。
- `examples/demo/manifest.json` 已升级为结构化注册清单。
- 新增 `scripts/check-module-manifests.py`，校验 examples 下所有 manifest。
- 新增 `docs/P2_MODULE_REGISTRY.md`，记录字段、路由、权限、菜单和校验规则。
- 本阶段不写数据库、不创建 demo 表、不写菜单或权限种子。

详见 `docs/P2_MODULE_REGISTRY.md`。

## P2.11 验收标准

- `python3 -m py_compile scripts/check-module-manifests.py` 通过。
- `python3 scripts/check-module-manifests.py` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不执行数据库写入、不修改 schema、不默认注册 demo 运行时路由。

## P2.11 验收结果

- 已通过 `python3 -m py_compile scripts/check-module-manifests.py`。
- 已通过 `python3 scripts/check-module-manifests.py`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有执行数据库写入、没有修改 schema、没有默认注册 demo 运行时路由。

## P2.12 当前落地

模块注册 SQL dry-run 已建立：

- 新增 `scripts/module-registry-plan.py`。
- 脚本读取 manifest 并复用 `scripts/check-module-manifests.py` 做结构校验。
- 脚本只输出 SQL 预览，不连接数据库、不执行 SQL。
- SQL 预览包含 `ma_permission`、`ma_menu`、`ma_menu_permission`。
- SQL 使用 `WHERE NOT EXISTS` 防止重复插入。
- 菜单父级通过 manifest `menu.parent` 对应的 `ma_menu.route_name` 查找。
- 本阶段不写 `ma_role_permission`，不自动给角色授权。

详见 `docs/P2_MODULE_REGISTRY_SQL_PLAN.md`。

## P2.12 验收标准

- `python3 -m py_compile scripts/check-module-manifests.py scripts/module-registry-plan.py` 通过。
- `python3 scripts/check-module-manifests.py` 通过。
- `python3 scripts/module-registry-plan.py --manifest examples/demo/manifest.json` 通过，只输出 SQL。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不执行数据库写入、不修改 schema、不默认注册 demo 运行时路由。

## P2.12 验收结果

- 已通过 `python3 -m py_compile scripts/check-module-manifests.py scripts/module-registry-plan.py`。
- 已通过 `python3 scripts/check-module-manifests.py`。
- 已通过 `python3 scripts/module-registry-plan.py --manifest examples/demo/manifest.json`；生成 76 行 SQL 预览，没有执行 SQL。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有执行数据库写入、没有修改 schema、没有默认注册 demo 运行时路由。

## P2.13 当前落地

模块注册 apply/write 模式已开放为本地受控写入：

- `scripts/module-registry-plan.py --apply` 现在支持执行 manifest 注册 SQL。
- 写入必须同时满足 `MAKEADMIN_ALLOW_MODULE_REGISTRY_WRITE=1` 和 `--confirm-module <module>`。
- 缺少环境变量或确认参数时，脚本会在数据库访问前失败。
- 写入 SQL 复用 dry-run 同一套生成逻辑。
- 写入内容限定为 `ma_permission`、`ma_menu`、`ma_menu_permission`。
- 权限、菜单、菜单权限关联均使用缺失插入方式保证幂等。
- 本阶段不写 `ma_role_permission`，不自动给角色授权。
- 本阶段不创建 demo 表、不注册 demo 运行时路由、不修改 schema。

详见 `docs/P2_MODULE_REGISTRY_APPLY.md`。

## P2.13 验收标准

- `python3 -m py_compile scripts/check-module-manifests.py scripts/module-registry-plan.py` 通过。
- `python3 scripts/check-module-manifests.py` 通过。
- `python3 scripts/module-registry-plan.py --apply` 失败，且错误说明没有访问数据库。
- `MAKEADMIN_ALLOW_MODULE_REGISTRY_WRITE=1 python3 scripts/module-registry-plan.py --apply` 失败，且错误说明没有访问数据库。
- `python3 scripts/module-registry-plan.py --manifest examples/demo/manifest.json` 继续通过，只输出 SQL。
- 对本地 `go_makeadmin` 执行一次 demo article 注册写入 smoke，并清理临时行。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不修改 schema、不读取或修改 `.env`、不连接真实 zyai 业务库。

## P2.13 验收结果

- 已通过 `python3 -m py_compile scripts/check-module-manifests.py scripts/module-registry-plan.py`。
- 已通过 `python3 scripts/check-module-manifests.py`。
- 已通过 `python3 scripts/module-registry-plan.py --apply` 失败门禁；失败文案明确没有访问数据库。
- 已通过 `MAKEADMIN_ALLOW_MODULE_REGISTRY_WRITE=1 python3 scripts/module-registry-plan.py --apply` 失败门禁；失败文案明确没有访问数据库。
- 已通过 `python3 scripts/module-registry-plan.py --manifest examples/demo/manifest.json`；dry-run SQL 预览继续可用。
- 已对本地 `go_makeadmin` 完成 demo article 注册写入 smoke：第一次 apply 后得到 5 条权限、1 条菜单、1 条菜单权限关联。
- 已第二次执行 apply，计数仍为 5 条权限、1 条菜单、1 条菜单权限关联，确认幂等。
- 已清理 demo article 注册行，清理后残留计数为 0。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有修改 schema、没有读取或修改 `.env`、没有连接真实 zyai 业务库。

## P2.14 当前落地

模块注册角色授权 dry-run 已建立：

- 新增 `scripts/module-role-grant-plan.py`。
- 脚本读取 manifest 并复用 `scripts/check-module-manifests.py` 做结构校验。
- 脚本只输出 SQL 预览，不连接数据库、不执行 SQL。
- SQL 预览包含 `ma_role_permission` 缺失授权插入。
- 授权目标必须显式传入 `--role-id`，不提供默认角色。
- `--tenant-id` 默认 `0`，但必须是非负整数。
- SQL 通过 `ma_permission.code` 查找权限 ID。
- SQL 校验目标 `ma_role` 存在、启用且未软删除。
- SQL 使用 `tenant_id + role_id + permission_id` 防止重复授权。

详见 `docs/P2_MODULE_ROLE_GRANT_PLAN.md`。

## P2.14 验收标准

- `python3 -m py_compile scripts/check-module-manifests.py scripts/module-role-grant-plan.py` 通过。
- `python3 scripts/check-module-manifests.py` 通过。
- `python3 scripts/module-role-grant-plan.py --role-id 1` 通过，只输出 SQL。
- `python3 scripts/module-role-grant-plan.py --role-id 0` 失败，且不访问数据库。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不执行数据库写入、不修改 schema、不默认授权任何角色。

## P2.14 验收结果

- 已通过 `python3 -m py_compile scripts/check-module-manifests.py scripts/module-role-grant-plan.py`。
- 已通过 `python3 scripts/check-module-manifests.py`。
- 已通过 `python3 scripts/module-role-grant-plan.py --role-id 1`；生成 `ma_role_permission` SQL 预览，没有执行 SQL。
- 已通过 `python3 scripts/module-role-grant-plan.py --role-id 0` 失败校验；命令在参数解析阶段失败，没有访问数据库。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有执行数据库写入、没有修改 schema、没有默认授权任何角色。

## P2.15 当前落地

模块运行时注册闭环已建立：

- 新增 `server/modules/routers` 作为运行时模块路由入口。
- 新增 `MAKEADMIN_ENABLE_DEMO_MODULE=1` 环境变量开关。
- 默认不挂载 demo runtime 模块。
- 开启后主路由会追加 `server/modules/routers.InitRouters()` 返回的模块路由。
- 新增 demo article 路由，路径与 `examples/demo/manifest.json` 的 `backend.routes` 保持一致。
- demo article 路由仍使用 `middleware.TokenAuth()`。
- 未登录访问 `/api/article/list` 返回 token empty 响应，确认没有绕过认证。
- demo article 写操作返回只读失败响应。
- 本阶段不创建 demo 表、不写数据库、不自动授权角色。

详见 `docs/P2_MODULE_RUNTIME_REGISTRY.md`。

## P2.15 验收标准

- `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./modules/routers` 通过。
- `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./...` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不执行数据库写入、不修改 schema、不开放无认证 demo API。

## P2.15 验收结果

- 已通过 `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./modules/routers`。
- 已通过 `cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./...`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有执行数据库写入、没有修改 schema、没有开放无认证 demo API。

## P2.16 当前落地

模块安装 dry-run 编排已建立：

- 新增 `scripts/module-install-plan.py`。
- 脚本读取 manifest 并复用已有 manifest 校验。
- 脚本只输出安装计划，不连接数据库、不执行 SQL。
- 安装计划包含 manifest 摘要、后端路由、前端路径、注册 SQL、可选角色授权 SQL 和 runtime 开关。
- 注册 SQL 复用 `scripts/module-registry-plan.py` 的生成逻辑。
- 角色授权 SQL 复用 `scripts/module-role-grant-plan.py` 的生成逻辑。
- 未传 `--role-id` 时不生成角色授权 SQL。
- runtime 提示当前输出 `MAKEADMIN_ENABLE_DEMO_MODULE=1`。

详见 `docs/P2_MODULE_INSTALL_PLAN.md`。

## P2.16 验收标准

- `python3 -m py_compile scripts/check-module-manifests.py scripts/module-registry-plan.py scripts/module-role-grant-plan.py scripts/module-install-plan.py` 通过。
- `python3 scripts/module-install-plan.py --manifest examples/demo/manifest.json` 通过，只输出计划。
- `python3 scripts/module-install-plan.py --manifest examples/demo/manifest.json --tenant-id 0 --role-id 1` 通过，只输出计划。
- `python3 scripts/module-install-plan.py --role-id 0` 失败，且不访问数据库。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不执行数据库写入、不修改 schema、不修改 `.env`。

## P2.16 验收结果

- 已通过 `python3 -m py_compile scripts/check-module-manifests.py scripts/module-registry-plan.py scripts/module-role-grant-plan.py scripts/module-install-plan.py`。
- 已通过 `python3 scripts/module-install-plan.py --manifest examples/demo/manifest.json`；生成安装计划，没有执行 SQL。
- 已通过 `python3 scripts/module-install-plan.py --manifest examples/demo/manifest.json --tenant-id 0 --role-id 1`；生成包含角色授权 SQL 的安装计划，没有执行 SQL。
- 已通过 `python3 scripts/module-install-plan.py --role-id 0` 失败校验；命令在参数解析阶段失败，没有访问数据库。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有执行数据库写入、没有修改 schema、没有修改 `.env`。

## P2.17 当前落地

模块安装 apply/write 边界已定义：

- 新增 `docs/P2_MODULE_INSTALL_APPLY_BOUNDARY.md`。
- 明确未来安装写入必须使用 `MAKEADMIN_ALLOW_MODULE_INSTALL_WRITE=1`。
- 明确未来写入必须显式传入 `--confirm-module <module>`。
- 如果传入 `--role-id`，还必须显式传入 `--confirm-role-id <id>`。
- 明确执行器不得读取 `.env` 猜测数据库密码。
- 明确未来写入必须在单事务内完成。
- 明确本阶段不自动建表，`requiresSchema=true` 时直接失败。
- 明确 runtime 开关只输出提示，不修改 `.env` 或系统环境变量。
- 明确未来本地 smoke 写入前必须确认 demo article 注册行不存在，避免误删已有数据。

详见 `docs/P2_MODULE_INSTALL_APPLY_BOUNDARY.md`。

## P2.17 验收标准

- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不开放新的写入命令。
- 不执行数据库写入、不修改 schema、不修改 `.env`。

## P2.17 验收结果

- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有开放新的写入命令。
- 本阶段没有执行数据库写入、没有修改 schema、没有修改 `.env`。

## P2.18 当前落地

模块安装 apply/write 模式已开放为本地受控写入：

- `scripts/module-install-plan.py --apply` 现在支持执行模块安装写入。
- 写入必须满足 `MAKEADMIN_ALLOW_MODULE_INSTALL_WRITE=1` 和 `--confirm-module <module>`。
- 如果传入 `--role-id`，必须同时传入匹配的 `--confirm-role-id <id>`。
- 缺少环境变量或确认参数时，脚本会在数据库访问前失败。
- manifest `requiresSchema=true` 时，脚本会在数据库访问前失败，不自动建表。
- 写入在单个事务中执行。
- 写入内容包含 `ma_permission`、`ma_menu`、`ma_menu_permission` 和可选 `ma_role_permission`。
- runtime 开关只输出提示，不修改 `.env` 或系统环境变量。

详见 `docs/P2_MODULE_INSTALL_APPLY.md`。

## P2.18 验收标准

- `python3 -m py_compile scripts/check-module-manifests.py scripts/module-registry-plan.py scripts/module-role-grant-plan.py scripts/module-install-plan.py` 通过。
- `python3 scripts/module-install-plan.py --apply` 失败，且错误说明没有访问数据库。
- `MAKEADMIN_ALLOW_MODULE_INSTALL_WRITE=1 python3 scripts/module-install-plan.py --apply` 失败，且错误说明没有访问数据库。
- `MAKEADMIN_ALLOW_MODULE_INSTALL_WRITE=1 python3 scripts/module-install-plan.py --apply --confirm-module article --role-id 1` 失败，且错误说明没有访问数据库。
- `python3 scripts/module-install-plan.py --manifest examples/demo/manifest.json --tenant-id 0 --role-id 1` 继续通过，只输出计划。
- 对本地 `go_makeadmin` 执行一次 demo article 安装写入 smoke，并清理临时行。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不修改 schema、不读取或修改 `.env`、不连接真实 zyai 业务库。

## P2.18 验收结果

- 已通过 `python3 -m py_compile scripts/check-module-manifests.py scripts/module-registry-plan.py scripts/module-role-grant-plan.py scripts/module-install-plan.py`。
- 已通过 `python3 scripts/module-install-plan.py --apply` 失败门禁；失败文案明确没有访问数据库。
- 已通过 `MAKEADMIN_ALLOW_MODULE_INSTALL_WRITE=1 python3 scripts/module-install-plan.py --apply` 失败门禁；失败文案明确没有访问数据库。
- 已通过 `MAKEADMIN_ALLOW_MODULE_INSTALL_WRITE=1 python3 scripts/module-install-plan.py --apply --confirm-module article --role-id 1` 失败门禁；失败文案明确没有访问数据库。
- 已通过 `python3 scripts/module-install-plan.py --manifest examples/demo/manifest.json --tenant-id 0 --role-id 1`；dry-run 安装计划继续可用。
- 已对本地 `go_makeadmin` 完成 demo article 安装写入 smoke：第一次 apply 后得到 5 条权限、1 条菜单、1 条菜单权限关联、5 条角色授权。
- 已第二次执行 apply，计数仍为 5 条权限、1 条菜单、1 条菜单权限关联、5 条角色授权，确认幂等。
- 已清理 demo article 注册行和角色授权，清理后残留计数为 0。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有修改 schema、没有读取或修改 `.env`、没有连接真实 zyai 业务库。

## P2.19 当前落地

模块卸载/回滚 dry-run 已建立：

- 新增 `scripts/module-uninstall-plan.py`。
- 脚本读取 manifest 并复用已有 manifest 校验。
- 脚本只输出 SQL 预览，不连接数据库、不执行删除。
- SQL 预览包含 `ma_role_permission`、`ma_menu_permission`、`ma_menu`、`ma_permission` 清理语句。
- 清理顺序先关联表、后主表。
- SQL 只使用 manifest 声明的权限 code 和菜单 routeName。
- 本阶段不删除前端文件、后端代码、runtime 环境变量或 schema。

详见 `docs/P2_MODULE_UNINSTALL_PLAN.md`。

## P2.19 验收标准

- `python3 -m py_compile scripts/check-module-manifests.py scripts/module-registry-plan.py scripts/module-uninstall-plan.py` 通过。
- `python3 scripts/check-module-manifests.py` 通过。
- `python3 scripts/module-uninstall-plan.py --manifest examples/demo/manifest.json` 通过，只输出 SQL。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不执行数据库删除、不删除文件、不修改 schema、不修改 `.env`。

## P2.19 验收结果

- 已通过 `python3 -m py_compile scripts/check-module-manifests.py scripts/module-registry-plan.py scripts/module-uninstall-plan.py`。
- 已通过 `python3 scripts/check-module-manifests.py`。
- 已通过 `python3 scripts/module-uninstall-plan.py --manifest examples/demo/manifest.json`；生成卸载 SQL 预览，没有执行 SQL。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有执行数据库删除、没有删除文件、没有修改 schema、没有修改 `.env`。

## P2.20 当前落地

模块卸载 apply/write 边界已定义：

- 新增 `docs/P2_MODULE_UNINSTALL_APPLY_BOUNDARY.md`。
- 明确未来卸载删除必须使用 `MAKEADMIN_ALLOW_MODULE_UNINSTALL_WRITE=1`。
- 明确未来删除必须显式传入 `--confirm-module <module>`。
- 明确未来删除必须显式传入 `--confirm-delete`。
- 明确执行器不得读取 `.env` 猜测数据库密码。
- 明确未来删除前必须输出四类快照计数。
- 明确未来删除必须在单事务内完成。
- 明确未来删除只按 manifest 权限 code 和菜单 routeName 执行。
- 明确不删除前端文件、后端代码、runtime 环境变量、schema 或 codegen 元数据。

详见 `docs/P2_MODULE_UNINSTALL_APPLY_BOUNDARY.md`。

## P2.20 验收标准

- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不开放新的删除命令。
- 不执行数据库删除、不删除文件、不修改 schema、不修改 `.env`。

## P2.20 验收结果

- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有开放新的删除命令。
- 本阶段没有执行数据库删除、没有删除文件、没有修改 schema、没有修改 `.env`。

## P2.21 当前落地

模块卸载 apply/write 模式已开放为本地受控删除：

- `scripts/module-uninstall-plan.py --apply` 现在支持执行模块卸载删除。
- 删除必须满足 `MAKEADMIN_ALLOW_MODULE_UNINSTALL_WRITE=1`、`--confirm-module <module>` 和 `--confirm-delete`。
- 缺少环境变量或确认参数时，脚本会在数据库访问前失败。
- 删除前输出权限、菜单、菜单权限、角色授权四类快照计数。
- 四类计数均为 0 时报告 no-op，不执行删除。
- 删除在单个事务中执行。
- 删除内容限定为 `ma_role_permission`、`ma_menu_permission`、`ma_menu`、`ma_permission`。
- 本阶段不删除前端文件、后端代码、runtime 环境变量、schema 或 codegen 元数据。

详见 `docs/P2_MODULE_UNINSTALL_APPLY.md`。

## P2.21 验收标准

- `python3 -m py_compile scripts/check-module-manifests.py scripts/module-registry-plan.py scripts/module-uninstall-plan.py` 通过。
- `python3 scripts/module-uninstall-plan.py --apply` 失败，且错误说明没有访问数据库。
- `MAKEADMIN_ALLOW_MODULE_UNINSTALL_WRITE=1 python3 scripts/module-uninstall-plan.py --apply` 失败，且错误说明没有访问数据库。
- `MAKEADMIN_ALLOW_MODULE_UNINSTALL_WRITE=1 python3 scripts/module-uninstall-plan.py --apply --confirm-module article` 失败，且错误说明没有访问数据库。
- `python3 scripts/module-uninstall-plan.py --manifest examples/demo/manifest.json` 继续通过，只输出 SQL。
- 对本地 `go_makeadmin` 执行一次 demo article 安装后卸载 smoke。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不修改 schema、不读取或修改 `.env`、不连接真实 zyai 业务库。

## P2.21 验收结果

- 已通过 `python3 -m py_compile scripts/check-module-manifests.py scripts/module-registry-plan.py scripts/module-uninstall-plan.py`。
- 已通过 `python3 scripts/module-uninstall-plan.py --apply` 失败门禁；失败文案明确没有访问数据库。
- 已通过 `MAKEADMIN_ALLOW_MODULE_UNINSTALL_WRITE=1 python3 scripts/module-uninstall-plan.py --apply` 失败门禁；失败文案明确没有访问数据库。
- 已通过 `MAKEADMIN_ALLOW_MODULE_UNINSTALL_WRITE=1 python3 scripts/module-uninstall-plan.py --apply --confirm-module article` 失败门禁；失败文案明确没有访问数据库。
- 已通过 `python3 scripts/module-uninstall-plan.py --manifest examples/demo/manifest.json`；dry-run 卸载 SQL 预览继续可用。
- 已对本地 `go_makeadmin` 完成 demo article 安装后卸载 smoke：安装后得到 5 条权限、1 条菜单、1 条菜单权限关联、5 条角色授权。
- 已执行卸载，卸载后权限、菜单、菜单权限和角色授权计数均为 0。
- 已第二次执行卸载，脚本报告 no-op。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有修改 schema、没有读取或修改 `.env`、没有连接真实 zyai 业务库。

## P2.22 当前落地

模块生命周期 smoke 脚本已建立：

- 新增 `scripts/check-module-lifecycle-smoke.sh`。
- 脚本必须显式设置 `MAKEADMIN_ALLOW_MODULE_LIFECYCLE_WRITE=1`。
- 缺少环境变量时，脚本会在数据库访问前失败。
- 脚本写入前确认 demo article 注册行和授权残留为 0。
- 脚本调用 `module-install-plan.py --apply` 完成安装。
- 脚本校验安装后得到 5 条权限、1 条菜单、1 条菜单权限关联、5 条角色授权。
- 脚本调用 `module-uninstall-plan.py --apply` 完成卸载。
- 脚本校验卸载后残留为 0。
- 脚本第二次调用卸载，确认 no-op 幂等。

详见 `docs/P2_MODULE_LIFECYCLE_SMOKE.md`。

## P2.22 验收标准

- `bash -n scripts/check-module-lifecycle-smoke.sh` 通过。
- `scripts/check-module-lifecycle-smoke.sh` 失败，且错误说明没有访问数据库。
- `MAKEADMIN_ALLOW_MODULE_LIFECYCLE_WRITE=1 scripts/check-module-lifecycle-smoke.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不修改 schema、不读取或修改 `.env`、不连接真实 zyai 业务库。

## P2.22 验收结果

- 已通过 `bash -n scripts/check-module-lifecycle-smoke.sh`。
- 已通过 `scripts/check-module-lifecycle-smoke.sh` 失败门禁；失败文案明确没有访问数据库。
- 已通过 `MAKEADMIN_ALLOW_MODULE_LIFECYCLE_WRITE=1 scripts/check-module-lifecycle-smoke.sh`；脚本完成 install、uninstall、二次 no-op 和残留检查。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有修改 schema、没有读取或修改 `.env`、没有连接真实 zyai 业务库。

## P2.23 当前落地

模块工具验证边界已建立：

- 新增 `scripts/check-module-tools-no-db.sh`。
- `scripts/verify-no-db.sh` 已接入模块工具 no-db guard。
- no-db guard 覆盖 Python 语法、shell 语法、manifest 校验、模块 dry-run 预览和写入门禁失败检查。
- no-db guard 不连接数据库、不执行写入或删除。
- `scripts/check-module-lifecycle-smoke.sh` 继续保持显式写库 smoke，不进入默认 no-db 验证。

详见 `docs/P2_MODULE_VERIFY_BOUNDARY.md`。

## P2.23 验收标准

- `bash -n scripts/verify-no-db.sh scripts/check-module-tools-no-db.sh scripts/check-module-lifecycle-smoke.sh` 通过。
- `scripts/check-module-tools-no-db.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 不把写库 smoke 放入 `verify-no-db.sh`。
- 不修改 schema、不读取或修改 `.env`、不连接真实 zyai 业务库。

## P2.23 验收结果

- 已通过 `bash -n scripts/verify-no-db.sh scripts/check-module-tools-no-db.sh scripts/check-module-lifecycle-smoke.sh`。
- 已通过 `scripts/check-module-tools-no-db.sh`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`，且 verify 链路已执行模块工具 no-db guard。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有把写库 smoke 放入 `verify-no-db.sh`。
- 本阶段没有修改 schema、没有读取或修改 `.env`、没有连接真实 zyai 业务库。

## P2.24 当前落地

模块文档入口已收敛：

- 新增 `docs/P2_MODULE_GUIDE.md`。
- README 已补充 P2 模块使用指南入口。
- 模块使用指南覆盖 manifest 校验、注册预览、角色授权预览、安装计划、安装写入、卸载预览、卸载删除、runtime 开关、默认 no-db 验证、本地生命周期 smoke 和红线边界。
- P2 模块生命周期的详细文档继续保留，统一由 `docs/P2_MODULE_GUIDE.md` 汇总进入。

## P2.24 验收标准

- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- `git diff --check` 通过。
- 不修改 schema、不执行数据库写入或删除、不读取或修改 `.env`、不连接真实 zyai 业务库。

## P2.24 验收结果

- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `git diff --check`。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有修改 schema、没有执行数据库写入或删除、没有读取或修改 `.env`、没有连接真实 zyai 业务库。

## P2.25 当前落地

P2 冻结状态文档已建立：

- 新增 `docs/P2_FINAL_STATUS.md`。
- README 已补充 P2 最终状态入口。
- `docs/P2_MODULE_GUIDE.md` 已指向 P2 最终状态文档。
- 最终状态文档汇总 P2 已完成范围、冻结验收范围、已知验证噪音、保留边界、不覆盖范围和 P3 入口。

## P2.25 验收标准

- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- `MAKEADMIN_ALLOW_MODULE_LIFECYCLE_WRITE=1 scripts/check-module-lifecycle-smoke.sh` 通过。
- `git diff --check` 通过。
- 不修改 schema、不读取或修改 `.env`、不连接真实 zyai 业务库。

## P2.25 验收结果

- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 已通过 `MAKEADMIN_ALLOW_MODULE_LIFECYCLE_WRITE=1 scripts/check-module-lifecycle-smoke.sh`。
- 已通过 `git diff --check`。
- 已通过 `git check-ignore server/.env admin/.env.development admin/node_modules admin/dist frontend public/admin public/assets`。
- 生命周期 smoke 安装后计数为 5 条权限、1 条菜单、1 条菜单权限关联、5 条角色授权；卸载后残留计数为 0，二次卸载为 no-op。
- 冻结验收后再次查询 demo article 残留计数为 0。
- `verify-no-db` 中前端 build 仍输出 Rolldown 对 `node_modules/@vueuse/core/dist/index.js` 的 `/* #__PURE__ */` annotation warning；当前退出码为 0，不影响验收。
- 本阶段没有修改 schema、没有读取或修改 `.env`、没有连接真实 zyai 业务库。

## 下一步

P3.1：业务模块脚手架与产品化模板。该任务把 P2 的 codegen、manifest 和模块生命周期能力串成标准模块开发体验。
