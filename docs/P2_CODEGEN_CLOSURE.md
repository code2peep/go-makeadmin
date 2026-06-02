# P2 Codegen Closure

更新时间：2026-06-02

## 目标

P2.9 让代码生成器的 Go 输出具备可验证闭环：生成后的基础 CRUD 后端代码至少能编译，并且有标准示例模块说明后续挂载约定。

## 当前落地

- 修复 `route.go.tpl` 中 list handler 使用错误接收者的问题。
- `schema.go.tpl` 只在生成字段使用 `core.*` 类型时导入 `go-makeadmin/core`。
- `service.go.tpl` 只在需要 URL 绝对化时导入 `go-makeadmin/util`。
- `service.go.tpl` 的 `Detail` / `Del` 主键参数类型跟随生成主键类型。
- `EditReq` 始终包含主键字段，保证编辑逻辑可编译。
- 新增 `server/generator/tpl_test.go`，渲染 CRUD Go 模板到临时目录并执行 `go test .` 编译生成包。
- 新增 `examples/README.md` 和 `examples/demo/`，沉淀标准 CRUD 模块接入约定。

## 不在 P2.9 做

- 不默认把 demo 模块注册进运行时路由。
- 不创建 `ma_demo_article` 表。
- 不写入菜单或权限种子。
- 不生成或提交真实业务模块代码。
- 不扩大前端生成模板改造范围。

## 验证

P2.9 需要通过：

```bash
cd server
GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator ./generator/service/gen
GOCACHE=/private/tmp/go-makeadmin-gocache go test ./...
cd ..
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
