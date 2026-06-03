# Demo Module

`demo` 是 go-makeadmin 的标准后台 CRUD 示例模块，用来描述代码生成器输出后的接入约定。

## 模块约定

- 后端包名：`gencode`
- 模块名：`article`
- 实体名：`DemoArticle`
- 表名：`ma_demo_article`
- 后端路由前缀：`/article`
- 管理端运行路由：`/demo/article`
- 数据库菜单路径：`/dev_tools/demo/article`
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

## P5.1 可见安装

P5.1 起，该示例可以作为本地可见模块安装到 `go_makeadmin`：

- manifest 菜单 `visible=true`。
- manifest 运行时 `runtimeRegistered=true`。
- 前端真实产物位于 `admin/src/api/article.ts` 和 `admin/src/views/article/`。
- 菜单安装后挂在开发工具下，访问路径为 `/demo/article`。
- 后端运行时仍需显式设置 `MAKEADMIN_ENABLE_DEMO_MODULE=1`。
- 示例模块不创建 `ma_demo_article` 表，列表接口返回空分页，写接口保持只读。

本地写入验证：

```bash
MAKEADMIN_ALLOW_DEMO_MODULE_VISIBLE_WRITE=1 scripts/check-demo-module-visible.sh
```

该脚本只写入 demo article 的菜单、权限、菜单权限和角色授权，不创建或迁移业务 schema。
