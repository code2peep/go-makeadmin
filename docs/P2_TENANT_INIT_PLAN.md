# P2 Tenant Init Plan

更新时间：2026-06-02

## 目标

P2.6 提供租户初始化命令的 dry-run 版本，用来从源租户读取设置和文件分类，生成目标租户的初始化计划与 SQL 预览。

本阶段不执行 SQL，不写数据库，不修改 schema。

## 脚本

```bash
python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 2
```

默认读取：

- `ma_setting` 中 `website`、`protocol`、`storage` 三类设置。
- `ma_file_category` 中未软删除的文件分类。

默认输出：

- 将插入的 setting key。
- 已存在会跳过的 setting key。
- 将插入的文件分类 code。
- 已存在会跳过的文件分类 code。
- SQL 预览。

只输出 SQL：

```bash
python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 2 --sql-only
```

## 安全规则

- 脚本只查询数据库，不执行输出的 SQL。
- `--apply` 是 P2.7 预留写入入口，当前会在数据库访问前失败。
- `--from-tenant` 和 `--to-tenant` 必须不同。
- `--to-tenant` 不会被自动创建。
- 目标租户已存在的 setting key 和文件分类 code 默认跳过。
- `storage` JSON 配置里的 `secretKey` 和 `accessKey` 默认清空。
- 只有显式传入 `--copy-secret` 才保留 `secretKey` / `accessKey`。

数据库连接参数来自：

- `MYSQL_HOST`，默认 `127.0.0.1`
- `MYSQL_PORT`，默认 `3306`
- `MYSQL_USER`，默认 `root`
- `MYSQL_PASSWORD`，可选
- `MYSQL_DATABASE`，默认 `go_makeadmin`

脚本通过本机 `mysql` client 读取数据；不会读取 `.env`。

## 输出示例

```text
Tenant init dry-run: from=0 to=2
Settings: insert=12 skip_existing=0
File categories: insert=2 skip_existing=0

-- Dry-run SQL preview. Review manually before applying; this script did not execute it.
SET @tenant_id = 2;
SET @now = UNIX_TIMESTAMP();
```

## 不在 P2.6 做

- 不执行 SQL。
- 不开放 `--apply` 写入模式。
- 不创建租户或租户成员。
- 不覆盖目标租户已有配置。
- 不迁移 `ma_file` 文件元数据。
- 不迁移物理上传文件。
- 不复制真实密钥，除非调用者显式传 `--copy-secret` 查看预览。

## 验证

P2.6 需要通过：

- `python3 -m py_compile scripts/tenant-init-plan.py`
- `python3 scripts/tenant-init-plan.py --help`
- `python3 scripts/tenant-init-plan.py --from-tenant 0 --to-tenant 999999 --sql-only`
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`
- `./scripts/check-services.sh`
- `./scripts/check-p1-seed.sh`

后续写入门禁见 `docs/P2_TENANT_INIT_APPLY_GUARD.md`。
