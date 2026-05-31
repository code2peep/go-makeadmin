# P0 SQL 残留审计

更新时间：2026-05-31

## 结论

P0 阶段保留两类 SQL：

- `sql/install.sql`：LikeAdmin 原始蓝本 SQL，作为来源参考保留，不作为最终自研 schema。
- `sql/install.core.sql`：`go-makeadmin` P0 最小核心初始化 SQL，由 `scripts/build-core-sql.sh` 生成，用于新项目快速初始化基础后台。

当前核心运行链路已经不再依赖文章、C 端用户、渠道、装修、消息通知、搜索设置、用户注册设置等业务演示模块。

## P0.12 已清理

- `sql/install.core.sql` 不再包含 `website.shopName` 和 `website.shopLogo`。
- 后端网站设置接口不再返回或保存 `shopName`、`shopLogo`。
- 管理端网站信息页面不再展示“前台设置”“商城名称”“商城LOGO”。
- 后端免权限列表移除 `article:cate:all`。
- 装修器专用的 `admin/src/components/link` 已迁入 `legacy/likeadmin-demo`。

## 仍保留的蓝本残留

这些残留是 P0 有意保留，不在本轮强行改名：

- 表前缀仍为 `la_*`。
- Go 模型仍对应当前 `la_*` 表结构。
- `sql/install.sql` 仍包含文章、用户、渠道、装修、通知、短信等蓝本业务表和演示数据。
- 当前本机开发库 `go_makeadmin` 可能仍有旧蓝本配置行；核心代码不再读取这些字段。

## 核心 SQL 范围

`sql/install.core.sql` 保留：

- 素材与上传：`la_album`、`la_album_cate`
- 字典：`la_dict_type`、`la_dict_data`
- 代码生成：`la_gen_table`、`la_gen_table_column`
- 权限与组织：`la_system_auth_admin`、`la_system_auth_role`、`la_system_auth_perm`、`la_system_auth_menu`、`la_system_auth_dept`、`la_system_auth_post`
- 系统配置与日志：`la_system_config`、`la_system_log_login`、`la_system_log_operate`

## 后续处理

P1 不继续给 `la_*` 打补丁，应按 `docs/P1_SCHEMA_PLAN.md` 新建 `ma_*` 自研模型，再做一次性种子迁移。
