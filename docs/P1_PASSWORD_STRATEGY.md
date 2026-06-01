# P1 密码策略

更新时间：2026-05-31

## 结论

P1 新账号统一使用 bcrypt。`ma_admin.password_hash` 保存 bcrypt 结果，`ma_admin.password_salt` 对新账号保持空字符串，只为旧 `la_*` 的 MD5+salt 一次性迁移兼容保留。

## 原则

- 仓库不保存默认明文密码。
- 仓库不保存可用于登录的默认密码 hash。
- 初始化 SQL 只允许出现安装时替换占位符。
- 管理员初始密码必须在安装命令或一次性本地变量中生成。
- 登录服务可以校验旧 MD5+salt，但新写入必须使用 bcrypt。

## 新账号格式

```text
password_hash: bcrypt hash generated at install time
password_salt: ""
```

bcrypt hash 自带盐，不再额外写 `password_salt`。

## 旧账号兼容

P1 允许校验旧 LikeAdmin 账号的 MD5+salt：

```text
password_hash: md5(plain_password + legacy_salt)
password_salt: legacy_salt
```

该兼容只用于迁移期认证。登录成功后的自动升级写回要等 P1 登录写路径确定后再做，不能在只读认证阶段提前修改数据。

## 当前实现

- `server/makeadmin/security/password.go`：bcrypt 生成、bcrypt 校验、旧 MD5+salt 校验、安装占位符拒绝。
- `server/makeadmin/service/auth.go`：`AuthenticateByUsername` 完成账号密码校验和身份构建；`Login` 完成 token、Redis session、登录日志和最后登录信息写入。
