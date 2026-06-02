# P2.1 Auth Model

## 目标

P2.1 将后台认证从 P1 的纯 Redis opaque token 升级为 JWT + Redis session state：

- JWT 是客户端携带的访问凭证。
- Redis session state 是服务端会话有效性和登出吊销来源。
- 数据库 schema 不变。
- 不恢复旧 `backstage:*` token 链路。

## Token 结构

JWT 使用 HS256 签名，签名密钥来自 `config.Config.Secret`。

Payload 字段：

- `sid`：随机 session id，Redis 会话主键。
- `adminId`：管理员 ID。
- `tenantId`：租户 ID，当前仍为 `0`。
- `iat`：签发时间。
- `exp`：过期时间。
- `iss`：固定为 `go-makeadmin`。

客户端仍通过 HTTP header `token` 传递 JWT，保持前端接口形状不变。

## Redis session state

登录成功后写入：

```text
MakeAdmin:makeadmin:session:<sid> = <adminId>
MakeAdmin:makeadmin:session:set:<adminId> contains <sid>
```

其中 `MakeAdmin:` 来自 `config.Config.RedisPrefix`。

中间件校验顺序：

1. 读取 header `token`。
2. 解析并校验 JWT 签名、issuer 和过期时间。
3. 使用 JWT `sid` 查询 Redis session state。
4. 校验 Redis 中的 admin id 和 JWT `adminId` 一致。
5. 从 `ma_admin` 重建实时身份、角色和权限。
6. 执行权限判断。

JWT 过期或 Redis session state 不存在时，均视为 token 无效。

## 登出

`/system/logout` 仍要求有效 token。登出时解析 JWT，删除对应 Redis session state。JWT 本身无服务端存储，删除 Redis state 后即使 JWT 尚未过期也不能继续访问接口。

## TTL 策略

- 默认 TTL 仍为 `7200` 秒。
- JWT `exp` 固定为签发时间加 TTL。
- Redis session state 可在请求中续期，但续期不会超过 JWT 剩余有效时间。
- `makeadmin:session:set:<adminId>` 会跟随 session TTL 过期，避免只保存已过期 sid。

## 不在 P2.1 做

- 不增加数据库 session 表。
- 不做 refresh token。
- 不做多端设备列表和批量踢出。
- 不修改前端 header 名称。
- 不处理生产密钥轮换。

## 后续任务

- P2.2：多租户上下文 middleware。
- P2.3：数据权限查询约束。
- 后续认证增强：session 设备信息、批量吊销、密钥轮换和 refresh token。

## P2.1 落地状态

- 已实现标准库 HS256 JWT 签发和解析。
- 已将 Redis session state 从旧 token key 改为 `makeadmin:session:<sid>`。
- 已将认证中间件改为 JWT 校验、Redis session state 校验、实时身份重建。
- 已保持前端 header `token` 和登录响应 `token` 字段不变。
- 已补 JWT 签名、篡改和过期测试。
