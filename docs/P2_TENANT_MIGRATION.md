# P2 Tenant Migration Strategy

更新时间：2026-06-02

## 目标

P2.5 明确租户级设置、文件和日志从默认租户 `0` 走向多租户运行时的迁移策略。

本阶段只沉淀策略和防回流检查，不执行数据库写入，不修改 schema，不迁移真实数据。

## 当前状态

设置：

- `ma_setting` 已按 `tenant_id + setting_group + setting_key` 隔离。
- 网站信息、备案、协议、存储设置、公共首页配置和控制台信息都从 request context 读取 `tenant_id`。
- 非 `0` 租户如果没有自己的设置行，不会自动回退读取租户 `0`。

文件：

- `ma_file_category` 和 `ma_file` 已按 `tenant_id` 隔离。
- 文件分类、文件列表、移动、重命名、删除和上传元数据都使用当前租户。
- 物理上传目录仍由当前 storage driver 决定；P2.5 不改变物理文件路径规则。

日志：

- 登录日志写入登录时选择的租户。
- 操作日志从 request context 写入当前租户。
- 登录日志和操作日志列表按当前租户查询，并叠加 P2.3 数据权限。
- 旧 `tenant_id=0` 日志保留在默认租户，不自动复制到新租户。

## 迁移原则

- 默认租户 `0` 是框架基线，不是业务租户模板的永久替代品。
- 新租户上线前必须准备自己的设置和文件分类基线。
- 迁移命令必须显式指定源租户和目标租户。
- 迁移命令必须可 dry-run，先输出将写入的记录数量和 key 列表。
- 涉及密钥、token、云存储 `secretKey` / `accessKey` 的配置不默认复制。
- 不跨项目库迁移，不连接 zyai 业务库写数据。
- 不把真实密钥写进文档、脚本参数日志或 Git commit。

## 设置迁移策略

建议按 setting group 分类处理：

| setting_group | 策略 |
| --- | --- |
| `website` 基础展示项 | 可从租户 `0` 复制到新租户，缺失则补齐 |
| `website.copyright` | 可从租户 `0` 复制到新租户 |
| `protocol` | 可从租户 `0` 复制到新租户 |
| `storage.default` | 可复制默认 driver 名称 |
| `storage.local` | 可复制本地公开配置 |
| `storage.qiniu` / `storage.aliyun` / `storage.qcloud` | 不默认复制密钥；只允许显式 opt-in 或重新配置 |

后续命令建议：

```text
makeadmin tenant seed-settings --from-tenant=0 --to-tenant=2 --dry-run
makeadmin tenant seed-settings --from-tenant=0 --to-tenant=2 --apply --copy-secret=false
```

执行要求：

- `--dry-run` 为默认模式。
- `--apply` 必须显式传入。
- 默认只插入目标租户缺失 key，不覆盖已有 key。
- 覆盖已有 key 必须额外传 `--overwrite`。

## 文件迁移策略

文件分类：

- 新租户应创建基础分类，例如图片和视频根分类。
- 可以从租户 `0` 复制分类结构，但要重新生成目标租户内唯一的分类 ID。
- 分类 code 可以沿用，但唯一性必须按目标租户校验。

文件元数据：

- 不默认复制 `ma_file` 元数据。
- 只有在明确要把某批公共素材作为租户模板素材时，才按白名单复制。
- 复制文件元数据前必须确认物理文件仍可访问。

物理文件：

- P2.5 不迁移物理文件。
- 后续如要做租户级物理路径隔离，应在 storage driver 层增加 tenant-aware path 规则，并先验证旧文件 URL 兼容。

后续命令建议：

```text
makeadmin tenant seed-file-categories --from-tenant=0 --to-tenant=2 --dry-run
makeadmin tenant seed-file-categories --from-tenant=0 --to-tenant=2 --apply
```

## 日志迁移策略

登录日志和操作日志不建议做跨租户复制：

- `tenant_id=0` 的历史日志代表默认租户历史上下文。
- 新租户从启用后开始自然写入自己的日志。
- 如果需要审计视图合并，应在查询层做只读聚合，不修改日志原始归属。

后续可以考虑：

- 为超级管理员增加跨租户只读审计视图。
- 为租户管理员保持仅当前租户日志视图。

## 防回流检查

`scripts/check-runtime-residue.sh` 已增加检查：

- 设置、文件和日志的 service/repository 不允许硬编码 `GlobalTenantID`。
- 这些链路必须继续从调用参数或 request context 读取租户。

## 验证

P2.5 已通过：

- `bash -n scripts/check-runtime-residue.sh scripts/verify-no-db.sh scripts/check-services.sh scripts/check-p1-seed.sh`
- `python3 -m py_compile scripts/p1-smoke.py`
- `./scripts/check-runtime-residue.sh`
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`
- `./scripts/check-services.sh`
- `./scripts/check-p1-seed.sh`

## 不在 P2.5 做

- 不执行租户数据复制。
- 不新增迁移表或 schema。
- 不读取或修改 `.env`。
- 不复制真实云存储密钥。
- 不迁移物理上传文件。
- 不连接 zyai 业务库写数据。

## 下一步建议

P2.6 可以做租户初始化命令的 dry-run 版本：只读取源租户配置，生成将要写入的计划和 SQL 预览，不直接执行写库。
