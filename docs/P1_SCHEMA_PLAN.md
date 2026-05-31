# P1 Schema 设计基线

更新时间：2026-05-31

## 目标

P1 的目标不是把 `la_*` 改名前缀，而是建立 `go-makeadmin` 自己的后台基础模型：认证、权限、组织、设置、文件、日志、代码生成都应服务于可复用后台框架，而不是某个商城/C 端业务。

## 设计原则

- 新表统一使用 `ma_*` 前缀。
- 权限能力和菜单展示分离：权限描述“能做什么”，菜单描述“在哪里展示”。
- 管理员、角色、组织支持多对多，避免继续沿用单角色字段。
- 租户能力作为基础设施预留，不绑定具体 SaaS 业务。
- 初始化数据只包含后台框架必需项，不包含文章、商城、用户端、公众号、装修等业务演示数据。
- P1 先支持全新初始化，再决定是否提供 `la_* -> ma_*` 一次性迁移脚本。

## 核心表草案

认证与账号：

```text
ma_admin
ma_admin_profile
ma_admin_role
ma_admin_org
ma_login_log
```

权限与菜单：

```text
ma_role
ma_permission
ma_role_permission
ma_menu
ma_menu_permission
```

组织与数据范围：

```text
ma_org_unit
ma_position
ma_data_scope
ma_role_data_scope
```

租户预留：

```text
ma_tenant
ma_tenant_member
ma_tenant_setting
```

系统基础：

```text
ma_setting
ma_dict_type
ma_dict_item
ma_file
ma_file_category
ma_audit_log
```

开发工具：

```text
ma_codegen_table
ma_codegen_column
```

## 与当前 `la_*` 的映射方向

- `la_system_auth_admin` -> `ma_admin`、`ma_admin_profile`
- `la_system_auth_role` -> `ma_role`
- `la_system_auth_perm` -> `ma_role_permission`
- `la_system_auth_menu` -> `ma_menu` 和 `ma_permission`
- `la_system_auth_dept` -> `ma_org_unit`
- `la_system_auth_post` -> `ma_position`
- `la_system_config` -> `ma_setting`
- `la_album`、`la_album_cate` -> `ma_file`、`ma_file_category`
- `la_system_log_login` -> `ma_login_log`
- `la_system_log_operate` -> `ma_audit_log`
- `la_gen_table`、`la_gen_table_column` -> `ma_codegen_table`、`ma_codegen_column`

业务演示表不进入核心映射：

```text
la_article*
la_user*
la_decorate*
la_hot_search
la_notice_setting
la_official_reply
la_system_log_sms
```

## P1 实施顺序

1. 定义 `ma_*` SQL 和 Go model，不修改当前 `la_*` 链路。
2. 建立 `ma_*` 最小种子：超级管理员、基础角色、基础权限、核心菜单、网站设置、存储设置。
3. 让登录、菜单、角色权限先切到 `ma_*`。
4. 再切组织、字典、素材、日志、代码生成。
5. 删除或冻结 `la_*` 运行路径，只保留迁移脚本和来源文档。

## 未定事项

- 是否在 P1 就启用多租户，还是先设计表但关闭入口。
- 数据范围采用部门树、组织树，还是策略表达式。
- 权限编码是否保持 `module:resource:action`，或改为更严格的 `domain.resource.action`。
