# Examples

本目录存放 go-makeadmin 框架示例模块，用来验证模块扩展约定，不存放一次性业务代码。

## 目录约定

- 每个示例一个子目录，目录名使用小写短横线或小写下划线。
- 每个示例必须包含 `README.md`，说明模块目标、后端路由、前端页面、权限标识和验证方式。
- 必须包含 `manifest.json`，用结构化方式描述模块元信息、后端路由、前端路径、菜单节点和权限。
- 示例代码必须能从文档追溯到框架约定或生成器输出。
- 示例不得依赖真实业务库、生产配置、`.env` 或私有密钥。

## 校验

```bash
python3 scripts/check-module-manifests.py
```

## 脚手架

```bash
python3 scripts/module-scaffold.py \
  --module billing_invoice \
  --entity BillingInvoice \
  --table ma_billing_invoice \
  --requires-schema \
  --dry-run
```

去掉 `--dry-run` 后会创建 `examples/billing_invoice/manifest.json` 和 `examples/billing_invoice/README.md`；如果目录已存在会失败，不覆盖已有文件。

## 当前示例

- `demo/`：标准后台 CRUD 模块约定，用于 P2.9/P2.10 代码生成器闭环。
