# P2 Tenant Init Apply

更新时间：2026-06-02

## 目标

P2.8 开放租户初始化命令的本地受控写入模式，用来把默认租户的基础设置和文件分类初始化到一个已存在的新租户。

写入模式只处理：

- `ma_setting` 中 `website`、`protocol`、`storage` 三类设置。
- `ma_file_category` 中未软删除的文件分类。

不处理：

- 租户创建。
- 租户成员创建。
- `ma_file` 文件元数据迁移。
- 物理上传文件迁移。
- schema 变更。

## 命令

dry-run 预览：

```bash
python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 2
```

受控写入：

```bash
MAKEADMIN_ALLOW_TENANT_INIT_WRITE=1 \
python3 scripts/tenant-init-plan.py \
  --from-tenant 0 \
  --to-tenant 2 \
  --confirm-to-tenant 2 \
  --apply
```

## 写入门禁

`--apply` 必须同时满足：

- 环境变量 `MAKEADMIN_ALLOW_TENANT_INIT_WRITE=1`。
- 显式传入 `--confirm-to-tenant <id>`。
- `--confirm-to-tenant` 必须等于 `--to-tenant`。
- `--from-tenant` 和 `--to-tenant` 必须不同且非负。
- 目标租户必须已存在、启用且未软删除。

缺少环境变量或确认参数时，脚本会在数据库访问前失败。

## 写入规则

- 写入在单个事务中执行。
- 只插入目标租户缺失的 setting key 和文件分类 code。
- 目标租户已有 setting key 和文件分类 code 只跳过，不覆盖、不更新。
- `storage` JSON 配置里的 `secretKey` / `accessKey` 默认清空。
- 只有显式传入 `--copy-secret` 才保留源租户存储密钥字段。
- apply 完成后输出插入数量、跳过数量和事务结果。

## 本地 smoke

P2.8 已用本地 `go_makeadmin` 执行一次写入 smoke：

- 创建临时租户 `990028`。
- 第一次 apply 插入 12 条 setting 和 2 条文件分类。
- 校验 storage 密钥字段没有非空复制。
- 第二次 apply 插入 0 条，跳过 12 条 setting 和 2 条文件分类。
- 清理临时租户、setting 和文件分类行。
- 清理后临时租户残留计数为 0。

## 验收

P2.8 需要通过：

```bash
python3 -m py_compile scripts/tenant-init-plan.py scripts/p1-smoke.py
python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 999999 --apply
MAKEADMIN_ALLOW_TENANT_INIT_WRITE=1 python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 999999 --apply
python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 999999 --sql-only
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
./scripts/check-services.sh
./scripts/check-p1-seed.sh
```
