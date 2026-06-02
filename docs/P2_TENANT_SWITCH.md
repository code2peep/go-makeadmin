# P2 Tenant Switch

更新时间：2026-06-02

## 目标

P2.4 开放租户成员校验和后端租户切换入口，让后台可以在已有 `ma_tenant` / `ma_tenant_member` 表基础上选择租户上下文。

本阶段不新增数据库表，不导入租户数据，不做前端租户切换 UI。

## 规则

登录阶段：

- 未携带 `X-Tenant-ID` 时继续使用默认租户 `0`。
- 携带 `X-Tenant-ID` 时解析为目标租户。
- `tenant_id=0` 是 P1/P2 兼容默认上下文，不要求 `ma_tenant_member`。
- 非 `0` 租户必须满足：
  - `ma_tenant.id` 存在、启用、未软删除。
  - `ma_tenant_member` 中存在当前管理员的启用成员关系。

认证后请求：

- JWT `tenantId` 仍是可信租户来源。
- 请求头 `X-Tenant-ID` 可以携带，但必须和 JWT `tenantId` 一致。
- 中间件重建身份时会重新校验租户成员关系；成员或租户失效后，旧 token 不能继续访问该租户。

租户切换：

- `GET /system/tenant/list` 返回当前管理员可访问租户。
- `POST /system/tenant/switch` 接收 `tenantId`，校验后签发目标租户的新 JWT 和 Redis session。
- 切换接口不隐式吊销旧 token，多端/单端策略后续单独设计。

## 接口

`GET /system/tenant/list`

返回：

```json
[
  {
    "id": 0,
    "code": "default",
    "name": "默认租户",
    "memberType": "default",
    "isCurrent": 1
  }
]
```

`POST /system/tenant/switch`

请求：

```json
{
  "tenantId": 0
}
```

返回：

```json
{
  "token": "new.jwt.token",
  "tenantId": 0
}
```

登录响应也新增 `tenantId` 字段；旧前端只读取 `token` 不受影响。

## 不在 P2.4 做

- 不新增租户、租户成员或租户设置表。
- 不写入或迁移租户数据。
- 不做前端切换菜单。
- 不改变前端 token header 名称。
- 不做切换时强制吊销旧 token。
- 不设计租户套餐、到期时间、配额和计费能力。

## 验证

P2.4 已通过：

- `python3 -m py_compile scripts/p1-smoke.py`
- `GOCACHE=/private/tmp/go-makeadmin-gocache go test ./...`
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`
- `./scripts/check-services.sh`
- `./scripts/check-p1-seed.sh`

完整 P1 HTTP smoke 仍只在本机提供一次性 `P1_SMOKE_ADMIN_PASSWORD` 或 `ADMIN_PASSWORD` 时运行。
