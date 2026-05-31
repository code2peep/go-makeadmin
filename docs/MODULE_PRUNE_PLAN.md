# P0.9 模块裁剪清单

## 结论

`go-makeadmin` 的基础框架应保留通用后台能力，裁掉或隔离 LikeAdmin 蓝本里的业务演示模块。

当前菜单种子已经只暴露通用后台模块，前端动态路由也已限制为核心模块。P0.11 已将前端业务演示源码迁入 `legacy/likeadmin-demo`，保留作蓝本参考，但不再留在核心 `admin/src` 源码树。

## 保留为核心框架

这些模块有 Go 后端路由或属于后台基础能力，P0 继续保留：

- 工作台：`workbench`
- 权限管理：管理员、角色、菜单
- 组织管理：部门、岗位
- 系统设置：网站信息、备案、协议、存储、字典、系统环境、缓存、日志
- 素材管理：素材中心、图库、上传
- 开发工具：代码生成器

对应后端范围：

```text
server/admin/routers/common
server/admin/routers/monitor
server/admin/routers/setting
server/admin/routers/system
server/generator
```

对应核心 SQL 表：

```text
la_album
la_album_cate
la_dict_data
la_dict_type
la_gen_table
la_gen_table_column
la_system_auth_admin
la_system_auth_dept
la_system_auth_menu
la_system_auth_perm
la_system_auth_post
la_system_auth_role
la_system_config
la_system_log_login
la_system_log_operate
```

## 裁剪或隔离候选

这些模块当前不适合作为通用后台核心：

- 文章内容：`admin/src/views/article`、`admin/src/api/article.ts`
- C 端用户：`admin/src/views/consumer`、`admin/src/api/consumer.ts`
- 微信/公众号/H5 渠道：`admin/src/views/channel`、`admin/src/api/channel`
- 装修器：`admin/src/views/decoration`、`admin/src/api/decoration.ts`
- 消息通知/短信：`admin/src/views/message`、`admin/src/api/message.ts`
- 搜索设置：`admin/src/views/setting/search`、`admin/src/api/setting/search.ts`
- 用户注册登录设置：`admin/src/views/setting/user`、`admin/src/api/setting/user.ts`

对应蓝本 SQL 表：

```text
la_article
la_article_category
la_article_collect
la_decorate_page
la_decorate_tabbar
la_hot_search
la_notice_setting
la_official_reply
la_user
la_user_auth
```

## 裁剪顺序

1. 先确认菜单种子不暴露候选模块。
2. 再移除工作台、配置、种子数据里的可见业务运营入口。
3. 新建最小核心初始化 SQL，只保留核心框架表和数据。
4. 前端源码按模块删除或迁入 `legacy` 隔离区。
5. 最后删除对应无用 API 封装、图片和类型。

## P0.10 隔离方式

当前先不删除候选模块源码，只在动态路由层设置核心视图白名单：

```text
admin/src/router/index.ts
```

动态路由只允许加载：

```text
workbench
permission
organization
material
setting/dict
setting/storage
setting/system
setting/website
dev_tools
```

这样可以保留蓝本源码作对照，同时确保文章、C 端用户、渠道、装修、消息等演示模块不会进入基础框架的动态路由集合。

P0.10 已验证：

- `./scripts/verify-no-db.sh` 通过。
- `admin/dist` 未发现业务演示路由字符串或 `LikeAdmin` 品牌字符串。
- 工作台接口已从 `version.channel` 改为 `version.links`。
- 管理端 `/workbench` 浏览器烟测通过，页面不再出现“服务支持”和 `LikeAdmin`。

## P0.11 源码隔离结果

业务演示源码和 API 封装已迁入：

```text
legacy/likeadmin-demo
```

迁移范围：

```text
admin/src/views/article
admin/src/views/consumer
admin/src/views/channel
admin/src/views/decoration
admin/src/views/message
admin/src/views/setting/search
admin/src/views/setting/user
admin/src/api/article.ts
admin/src/api/consumer.ts
admin/src/api/channel
admin/src/api/decoration.ts
admin/src/api/message.ts
admin/src/api/setting/search.ts
admin/src/api/setting/user.ts
admin/src/components/link
```

`legacy/README.md` 和 `legacy/likeadmin-demo/README.md` 已定义目录约定：`legacy` 内容只作参考，不允许被运行时代码导入。

## P0.12 配置残留清理

已从核心运行链路清理：

- 网站设置的 `shopName`、`shopLogo` 字段。
- 管理端网站信息页的“前台设置”“商城名称”“商城LOGO”。
- 后端免权限白名单里的 `article:cate:all`。
- `sql/install.core.sql` 中的 `website.shopName`、`website.shopLogo` 种子。
- 仅装修器演示使用的 `admin/src/components/link`。

更详细的 SQL 残留记录见：

```text
docs/P0_SQL_RESIDUE_AUDIT.md
```

## 约束

- `NOTICE.md` 和授权来源文档必须保留 LikeAdmin 来源说明。
- 不把业务演示模块改造成半成品核心能力。
- 裁剪源码前必须跑 `npm run type-check` 和 `npm run build`。
- 裁剪 SQL 前必须先保留原始 `sql/install.sql`，新增核心 SQL 文件，不直接丢失来源蓝本。
