# P3 Module Codegen Link

更新时间：2026-06-02

## 目标

P3.2 将 `module-scaffold` 输出的 manifest 与现有 Go/Vue 代码生成器模板验证链路打通。

本阶段不改旧 `/gen/*` 接口形状，不写数据库，不生成持久业务文件。它只验证脚手架输出能被 codegen 模板消费。

## 命令

```bash
scripts/check-module-codegen.sh
```

验证指定 manifest：

```bash
scripts/check-module-codegen.sh --manifest .cache/module-scaffold-smoke/examples/billing_invoice/manifest.json
```

该脚本会：

- 调用 `scripts/module-scaffold.py --print-manifest` 生成临时 manifest。
- 校验 manifest 是合法 JSON。
- 设置 `MAKEADMIN_CODEGEN_MANIFEST`。
- 运行 `go test ./generator -run TestGeneratedCrudCodeMatchesModuleManifest -count=1`。

## 验证内容

Go 测试会读取 manifest，并使用 manifest 的：

- `module`
- `entity`
- `table`
- `menu.name`
- `backend.routes`
- `permissions`

然后执行：

- 渲染 Go model、schema、service、route 模板。
- 渲染 Vue API、列表页和编辑页模板。
- 检查生成的后端 route 是否覆盖 manifest 的 route。
- 检查生成的前端 API URL 是否覆盖 manifest 的 route。
- 检查列表页按钮权限是否覆盖 manifest 的 add、edit、del 权限。
- 写入临时 Go 生成目录并执行 `go test .`，确认生成后端代码可编译。

## 默认 no-db 接入

`scripts/check-module-tools-no-db.sh` 已调用 `scripts/check-module-codegen.sh`。

因此默认验证会覆盖：

```bash
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```

## 边界

- 不连接数据库。
- 不创建业务 schema。
- 不读取或修改 `.env`。
- 不写入 `examples/<module>/`。
- 不写入 `admin/src/api/` 或 `admin/src/views/`。
- 只在临时目录写入 manifest 和 Go 生成代码，并在脚本结束时清理。
- 使用 `--manifest` 时只读取指定 manifest，不清理调用方提供的文件。
