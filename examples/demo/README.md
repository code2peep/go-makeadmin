# Demo Module

`demo` 是 go-makeadmin 的标准后台 CRUD 示例模块，用来描述代码生成器输出后的接入约定。

## 模块约定

- 后端包名：`gencode`
- 模块名：`demo_article`
- 实体名：`DemoArticle`
- 表名：`ma_demo_article`
- 路由前缀：`/demo_article`
- 前端 API 目录：`admin/src/api/demo_article`
- 前端页面目录：`admin/src/views/demo_article`

## 后端路由

- `GET /demo_article/list`
- `GET /demo_article/detail`
- `POST /demo_article/add`
- `POST /demo_article/edit`
- `POST /demo_article/del`

## 权限标识

- `demo_article:list`
- `demo_article:add`
- `demo_article:edit`
- `demo_article:del`

## 生成器闭环

P2.9 已用 generator 单元测试验证 Go 模板闭环：

- 渲染 `model.go`、`schema.go`、`service.go`、`route.go`。
- 写入临时生成目录。
- 对临时生成包执行 `go test .`。
- 测试结束后清理临时生成目录。

该示例目前只作为模块接入标准，不默认注册进运行时路由，不创建表，不写入菜单或权限种子。
