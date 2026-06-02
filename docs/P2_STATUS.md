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

## 下一步

P2.2：多租户上下文 middleware。
