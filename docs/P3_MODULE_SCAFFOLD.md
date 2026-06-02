# P3 Module Scaffold

更新时间：2026-06-02

## 目标

P3.1 提供标准业务模块脚手架入口，让一个新模块先拥有可验证的 manifest、README 和模块生命周期命令，再继续接入后端与前端代码生成。

## 命令

dry-run 预览：

```bash
python3 scripts/module-scaffold.py \
  --module billing_invoice \
  --entity BillingInvoice \
  --table ma_billing_invoice \
  --requires-schema \
  --dry-run
```

只输出 manifest JSON：

```bash
python3 scripts/module-scaffold.py \
  --module billing_invoice \
  --entity BillingInvoice \
  --table ma_billing_invoice \
  --requires-schema \
  --print-manifest
```

写入脚手架文件：

```bash
python3 scripts/module-scaffold.py \
  --module billing_invoice \
  --entity BillingInvoice \
  --table ma_billing_invoice \
  --requires-schema
```

写入后会创建：

```text
examples/billing_invoice/manifest.json
examples/billing_invoice/README.md
```

如果目标目录已存在，脚手架会失败，不覆盖已有文件。

## 默认生成内容

manifest 默认包含：

- `GET /<module>/list`
- `GET /<module>/detail`
- `POST /<module>/add`
- `POST /<module>/edit`
- `POST /<module>/del`
- `admin/src/api/<module>.ts`
- `admin/src/views/<module>/index.vue`
- `admin/src/views/<module>/edit.vue`
- `<module>:list`
- `<module>:detail`
- `<module>:add`
- `<module>:edit`
- `<module>:del`

默认菜单：

- 父级：`dev_tools`
- 路径：`/<module>`
- 名称：`<module>.index`
- 组件：`<module>/index`
- 可见性：隐藏

## 验证

脚手架生成内容会复用 `scripts/check-module-manifests.py` 做 manifest 校验，并调用模块注册、角色授权、卸载 SQL 生成器确认 manifest 能进入 P2 生命周期工具。

默认 no-db 验证已接入：

```bash
scripts/check-module-tools-no-db.sh
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```

## 边界

- 脚手架不创建业务 schema。
- 脚手架不连接数据库。
- 脚手架不读取或修改 `.env`。
- 脚手架不生成真实后端/前端代码文件；P3.2 已通过 `scripts/check-module-codegen.sh` 打通 codegen 验证联动。
- `requiresSchema=true` 的模块安装 apply 仍会按 P2 边界失败，不自动建表。
