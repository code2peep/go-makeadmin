# P1 Runtime Residue Audit

## 目标

P1.20 的目标是明确 `la_*` 残留边界：P1 运行默认值和已接管模块应面向 `ma_*`，旧 `la_*` 只作为蓝本资料、P0 本地验证脚本和过渡代码存在。P1.22 已冻结旧运行兜底，核心后台接口不再在运行时回退到旧 `la_*` 服务。

## 当前结论

- P1 默认表前缀应为 `ma_`。
- `DB_TABLE_PREFIX` 可通过环境变量覆盖，默认不再是 `la_`。
- 代码生成器 `/gen/*` 已在 P1.19 切到 `ma_codegen_table`、`ma_codegen_column`。
- `server/model/gen` 当前只作为代码生成模板 DTO 保留，不再作为代码生成器持久化模型。
- 后台核心路由和认证中间件不再接受旧 `backstage:*` token，也不再回退到旧 `server/admin/service/*` 运行链路。

## 已接管运行链路

- 登录、自身信息、菜单路由：走 `ma_admin`、`ma_role`、`ma_menu`、`ma_permission`。
- 设置：网站信息、备案、协议、存储配置走 `ma_setting`。
- 字典：字典类型和字典项走 `ma_dict_type`、`ma_dict_item`。
- 文件：图库分类和文件元数据走 `ma_file_category`、`ma_file`。
- 日志：登录日志和操作审计走 `ma_login_log`、`ma_audit_log`。
- 权限：角色、部门、岗位、管理员、菜单管理走对应 `ma_*` 表。
- 代码生成：生成表配置和列配置走 `ma_codegen_table`、`ma_codegen_column`。
- 公共首页配置和控制台信息走 `ma_setting`。

## 允许保留的 `la_*` 残留

- `sql/install.sql`、`sql/install.core.sql`：P0 蓝本来源和旧库初始化资料。
- `scripts/check-db-seed.sh`、`scripts/build-core-sql.sh`、`scripts/init-local-blueprint-db.sh`：P0 蓝本验证和构建脚本。
- `docs/P0_*`、`docs/DB_INIT_PLAN.md`、`docs/MODULE_PRUNE_PLAN.md`：历史决策与迁移资料。
- `server/model/system`、`server/model/setting`、`server/model/common`：旧蓝本模型，暂作参考代码保留，不再作为 P1 核心运行兜底。
- `server/admin/service/*` 旧服务：旧蓝本服务，暂作参考代码保留，不再作为 P1 核心运行兜底。

## 不再允许新增的残留

- 新 P1 功能直接读写 `la_*` 表。
- 新 P1 运行默认值使用 `la_` 前缀。
- 新增需要长期双写 `la_*` 和 `ma_*` 的接口。
- 新增 `la_*` 初始化种子作为框架默认数据。

## 旧兜底冻结状态

- P1 初始化脚本可稳定生成完整 `ma_*` 开发库：已完成。
- 登录、权限、菜单、设置、文件、日志、代码生成的 smoke 覆盖核心写操作：已完成，并在 P1.22 补充公共首页和日志查询。
- P0 蓝本库不再作为 P1 默认运行库：已完成。
- P1 运行残留守卫脚本已接入 `verify-no-db`，防止核心运行目录重新引用旧服务、旧模型、旧 token 和 `la_*` 表名。

## 自动守卫

`scripts/check-runtime-residue.sh` 会扫描 P1 核心运行目录：

- `server/admin/routers`
- `server/middleware`
- `server/makeadmin`
- `server/generator`

守卫范围：

- 禁止重新导入旧 `server/admin/service/{system,setting,common}`。
- 禁止重新导入旧 `server/model/{system,setting,common}`。
- 禁止重新使用旧 `backstage:*` Redis token key。
- 禁止路由和中间件重新通过 `.Available()` 分支做运行时 fallback。
- 禁止 P1 运行目录直接引用 `la_*` 表名。
- 禁止 P1 运行目录重新使用 `ConfigUtil` 设置 fallback。
- 禁止配置默认表前缀回到 `la_`。

旧蓝本源码、旧 SQL 和 P0 脚本仍按“允许保留的 `la_*` 残留”规则保留，不在守卫脚本的运行路径扫描范围内。

## 后续清理条件

删除或迁入 `legacy` 前，需要单独确认删除边界；当前阶段只冻结运行引用，不删除文件。
