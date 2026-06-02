# Demo Module

`demo` 是 go-makeadmin 的标准后台 CRUD 示例模块，用来描述代码生成器输出后的接入约定。

## 模块约定

- 后端包名：`gencode`
- 模块名：`article`
- 实体名：`DemoArticle`
- 表名：`ma_demo_article`
- 路由前缀：`/article`
- 前端 API：`admin/src/api/article.ts`
- 前端页面目录：`admin/src/views/article`

## 后端路由

- `GET /article/list`
- `GET /article/detail`
- `POST /article/add`
- `POST /article/edit`
- `POST /article/del`

## 权限标识

- `article:list`
- `article:detail`
- `article:add`
- `article:edit`
- `article:del`

## 注册清单

`manifest.json` 描述本模块接入后台所需的注册信息：

- 后端路由和权限映射。
- 前端 API 和页面路径。
- 菜单节点。
- 权限元数据。
- 是否默认注册运行时路由。
- 是否需要 schema。

## 生成器闭环

P2.9 已用 generator 单元测试验证 Go 模板闭环：

- 渲染 `model.go`、`schema.go`、`service.go`、`route.go`。
- 写入临时生成目录。
- 对临时生成包执行 `go test .`。
- 测试结束后清理临时生成目录。

P2.10 已用显式脚本验证前端模板闭环：

- 渲染 `api.ts`、`index.vue`、`edit.vue`。
- 临时写入 `admin/src/api/article.ts` 和 `admin/src/views/article/`。
- 执行 `npm run type-check`。
- 测试结束后清理临时生成文件。

该示例目前只作为模块接入标准，不默认注册进运行时路由，不创建表，不写入菜单或权限种子。
