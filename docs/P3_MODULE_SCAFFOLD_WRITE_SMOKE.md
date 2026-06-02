# P3 Module Scaffold Write Smoke

更新时间：2026-06-02

## 目标

P3.3 验证 `module-scaffold` 的非 dry-run 写入路径：脚手架必须真实写出 `manifest.json` 和 `README.md`，并且写出的 manifest 能进入安装计划、卸载计划和 codegen 联动验证。

本阶段不写数据库，不创建业务 schema，不写入真实 `examples/<module>/` 目录。

## 命令

```bash
MAKEADMIN_ALLOW_MODULE_SCAFFOLD_WRITE=1 scripts/check-module-scaffold-write-smoke.sh
```

缺少 `MAKEADMIN_ALLOW_MODULE_SCAFFOLD_WRITE=1` 时，脚本会在写文件前失败。

## 验证内容

脚本会：

- 在 `.cache/module-scaffold-smoke/<timestamp>/examples/` 下创建临时模块目录。
- 调用 `scripts/module-scaffold.py --examples-root <cache examples dir>` 非 dry-run 写入文件。
- 检查 `manifest.json` 和 `README.md` 已存在。
- 校验 `manifest.json` 是合法 JSON。
- 执行 `module-install-plan.py --manifest <generated manifest> --tenant-id 0 --role-id 1` dry-run。
- 执行 `module-uninstall-plan.py --manifest <generated manifest>` dry-run。
- 执行 `scripts/check-module-codegen.sh --manifest <generated manifest>`。
- 检查 README 引用了生成 manifest 的实际路径。

## 边界

- 写入位置在 `.cache/` 下，已被 Git 忽略。
- 不删除仓库文件或目录。
- 不连接数据库。
- 不执行 SQL 写入或删除。
- 不创建业务 schema。
- 不读取或修改 `.env`。
