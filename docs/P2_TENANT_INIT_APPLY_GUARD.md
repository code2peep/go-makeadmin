# P2 Tenant Init Apply Guard

更新时间：2026-06-02

## 目标

P2.7 只建立租户初始化写入模式的安全门禁和设计边界，不执行数据库写入。

当前 `scripts/tenant-init-plan.py` 已接受 `--apply` 参数，但会在任何数据库访问之前失败：

```bash
python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 2 --apply
```

预期结果：

```text
FAIL: --apply is intentionally disabled until DB write approval is granted; no database access was attempted
```

## 当前门禁

- `--apply` 只作为预留入口，当前不执行查询、不执行写入。
- dry-run 行为保持不变：不传 `--apply` 时只读取数据库并输出 SQL 预览。
- 本阶段不引入自动迁移、不创建租户、不创建租户成员。
- 本阶段不迁移 `ma_file` 元数据和物理上传文件。
- 本阶段不默认复制 `secretKey` / `accessKey`。

## 后续开放写入的必要条件

真正进入 apply/write 实现前，需要单独授权数据库写入，并满足这些门禁：

- 命令必须同时包含 `--apply` 和 `MAKEADMIN_ALLOW_TENANT_INIT_WRITE=1`。
- 命令必须包含 `--confirm-to-tenant <id>`，且值必须等于 `--to-tenant`。
- `--from-tenant` 和 `--to-tenant` 必须不同且非负。
- 目标租户必须已存在且处于启用状态；脚本不自动创建租户。
- 写入必须在单个事务中完成。
- 只插入目标租户缺失的 `ma_setting` 和 `ma_file_category` 行。
- 已存在的 setting key 和文件分类 code 不覆盖、不更新。
- `storage` 密钥字段默认清空；只有显式 `--copy-secret` 才允许写入预览中的密钥字段。
- apply 完成后必须输出插入数量、跳过数量和事务结果。

## 验收要求

P2.7 安全门禁需要证明：

- `--apply` 会失败。
- `--apply` 失败发生在数据库访问前。
- dry-run SQL 预览继续可用。
- 无数据库写入、无 schema 变更、无 `.env` 修改。

推荐验证：

```bash
python3 -m py_compile scripts/tenant-init-plan.py scripts/p1-smoke.py
python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 999999 --apply
python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 999999 --sql-only
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
./scripts/check-services.sh
./scripts/check-p1-seed.sh
```

## 不在 P2.7 做

- 不执行租户初始化 SQL。
- 不开放 `MAKEADMIN_ALLOW_TENANT_INIT_WRITE` 写入路径。
- 不对真实业务库或生产库做任何操作。
- 不修改 schema。
- 不修改 `.env` 或密钥。
