# P1 登录切换方案

更新时间：2026-05-31

## 目标

将后台认证从旧 `la_*` 链路切到 `ma_*` 链路，同时保持前端现有登录、菜单和权限接口形状基本不变。

## 当前阶段

P1.6 已先完成新框架内部登录应用服务。P1.7 已接入现有 API 形状，但保留旧链路兜底。

已具备能力：

- 从 `ma_admin` 读取账号。
- 使用 bcrypt 校验 P1 新账号密码。
- 兼容校验旧 MD5+salt 迁移账号。
- 构建管理员身份、角色、权限和菜单。
- 生成 token。
- 写入 Redis 会话。
- 更新 `ma_admin.last_login_ip` 和 `ma_admin.last_login_time`。
- 写入 `ma_login_log`。
- `/api/system/login` 在检测到可用 `ma_admin` 时走新链路，否则继续走旧 `la_*`。
- `/api/system/admin/self` 和 `/api/system/menu/route` 在检测到 `makeadmin:*` token 时走新链路，否则继续走旧链路。

## 切换顺序

1. 准备独立 P1 开发库，导入 `sql/p1.schema.sql`。可使用 `scripts/init-p1-db.sh`。
2. 由安装器或一次性本地命令生成 admin bcrypt hash，再导入最小种子。可使用 `ADMIN_PASSWORD=... scripts/init-p1-db.sh`。
3. 增加 `/api/system/login` 的 `ma_*` 实现适配层，响应仍返回 `{ token }`。已完成。
4. 增加 `ma_*` TokenAuth 中间件，读取新 Redis session，不再依赖旧 `la_system_auth_admin` 缓存结构。已完成。
5. 切 `/api/system/admin/self` 和 `/api/system/menu/route` 到 `ma_*`。已完成新 token 分流。
6. 切角色、菜单、组织、日志等管理页。
7. 冻结旧 `la_*` 运行路径，仅保留迁移脚本和来源说明。

## 不在本阶段做

- 不直接修改当前 P0 默认库连接。
- 不直接删除老 middleware。
- 不做旧账号登录成功后的 hash 自动升级写回。
- 不删除旧 `admin/service/system`。

## 待实现适配

- P1 安装脚本：生成 bcrypt admin hash 并替换 `sql/p1.seed.sql` 占位符。
- `makeadmin` 当前管理员接口：补全组织、岗位和角色展示字段。
- `makeadmin` 菜单路由接口：继续收敛 route path 和 permission code 命名。
- `makeadmin` token middleware：补 Redis identity cache，减少每次请求查库。
