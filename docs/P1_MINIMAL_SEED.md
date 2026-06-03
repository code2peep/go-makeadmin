# P1.2 最小种子清单

更新时间：2026-05-31

## 范围

本文件定义 `ma_*` 初始化种子的最小清单。它是审阅用设计文档，不是可执行 SQL。

本阶段只确定种子边界：

- 不写入数据库。
- 不在文档或 SQL 中硬编码真实密码 hash；本地初始化脚本会默认生成 `admin / 123456` 的 bcrypt hash。
- 不执行导入脚本。
- 不修改当前 `la_*` 运行链路。
- 不包含文章、商城、用户端、公众号、装修等业务演示数据。

## 种子原则

- 用稳定 `code` 作为业务识别，不让业务逻辑依赖自增 ID。
- `id` 可以在 SQL 草案中固定，运行时代码仍应按 `code` 查询。
- `tenant_id=0` 表示全局默认上下文；P1 不启用租户入口。
- 超级管理员密码由初始化命令生成 bcrypt hash；本地默认密码为 `123456`，可通过 `ADMIN_PASSWORD` 覆盖，不能在仓库中硬编码最终 hash。
- P1 新账号统一使用 bcrypt；`password_salt` 只为旧 MD5+salt 迁移兼容保留。
- 超级管理员可以通过 `is_super=1` 拥有全权限；仍种下 `super_admin` 角色，方便页面展示和后续审计。
- 菜单只描述页面展示；接口能力统一进入 `ma_permission`。

## 初始化占位

```text
tenant_id: 0
create_time: {{unix_now}}
update_time: {{unix_now}}
admin_password_hash: {{bcrypt_generated_at_install_time}}
admin_password_salt: ""
```

## 账号与角色

`ma_admin`：

```text
code: admin
username: admin
is_super: 1
status: 1
password_hash: {{bcrypt_generated_at_install_time}}
password_salt: ""
```

`ma_admin_profile`：

```text
admin: admin
nickname: admin
avatar: /api/static/backend_avatar.png
```

`ma_role`：

```text
tenant_id: 0
code: super_admin
name: 超级管理员
is_system: 1
status: 1
sort: 1000
```

`ma_admin_role`：

```text
tenant_id: 0
admin: admin
role: super_admin
```

## 组织与数据范围

组织、岗位和数据范围要有最小数据，避免后续角色页、管理员页出现空选项。

`ma_org_unit`：

```text
tenant_id: 0
code: root
name: 总部
parent: none
status: 1
sort: 1000
```

`ma_position`：

```text
tenant_id: 0
code: admin
name: 管理员
status: 1
sort: 1000
```

`ma_admin_org`：

```text
tenant_id: 0
admin: admin
org: root
position: admin
is_primary: 1
status: 1
```

`ma_data_scope`：

```text
tenant_id: 0
code: all
name: 全部数据
scope_type: all
scope_value: {}
status: 1
```

`ma_role_data_scope`：

```text
tenant_id: 0
role: super_admin
data_scope: all
```

## 菜单种子

菜单树只保留核心后台页面。`code` 是文档层稳定标识，后续可映射到 `ma_menu.route_name` 或独立字段。

| code | parent | type | name | route_path | component | permission |
| --- | --- | --- | --- | --- | --- | --- |
| dashboard.workbench | root | page | 工作台 | /workbench | workbench/index | dashboard:workbench:view |
| permission | root | catalog | 权限管理 | /permission |  |  |
| permission.admin | permission | page | 管理员 | /permission/admin | permission/admin/index | system:admin:list |
| permission.role | permission | page | 角色管理 | /permission/role | permission/role/index | system:role:list |
| permission.menu | permission | page | 菜单管理 | /permission/menu | permission/menu/index | system:menu:list |
| organization | root | catalog | 组织管理 | /organization |  |  |
| organization.department | organization | page | 部门管理 | /organization/department | organization/department/index | system:dept:list |
| organization.post | organization | page | 岗位管理 | /organization/post | organization/post/index | system:post:list |
| material | root | page | 素材管理 | /material/index | material/index | material:file:list |
| setting | root | catalog | 系统设置 | /setting |  |  |
| setting.website | setting | catalog | 网站设置 | /setting/website |  |  |
| setting.website.information | setting.website | page | 网站信息 | /setting/website/information | setting/website/information | setting:website:detail |
| setting.website.filing | setting.website | page | 网站备案 | /setting/website/filing | setting/website/filing | setting:copyright:detail |
| setting.website.protocol | setting.website | page | 政策协议 | /setting/website/protocol | setting/website/protocol | setting:protocol:detail |
| setting.storage | setting | page | 存储设置 | /setting/storage | setting/storage/index | setting:storage:list |
| setting.system | setting | catalog | 系统维护 | /setting/system |  |  |
| setting.system.environment | setting.system | page | 系统环境 | /setting/system/environment | setting/system/environment | monitor:server |
| setting.system.cache | setting.system | page | 系统缓存 | /setting/system/cache | setting/system/cache | monitor:cache |
| setting.system.journal | setting.system | page | 系统日志 | /setting/system/journal | setting/system/journal | system:log:operate |
| dev_tools | root | catalog | 开发工具 | /dev_tools |  |  |
| dev_tools.dict | dev_tools | page | 字典管理 | /dev_tools/dict | setting/dict/type/index | setting:dict:type:list |
| dev_tools.code | dev_tools | page | 代码生成器 | /dev_tools/code | dev_tools/code/index | gen:list |

## 权限种子

权限种子以当前 P0 核心能力为边界。P1.2 暂保留少量历史双段权限编码，后续切服务时再决定是否统一为三段式并提供兼容映射。

工作台：

```text
dashboard:workbench:view
```

管理员：

```text
system:admin:self
system:admin:list
system:admin:detail
system:admin:add
system:admin:edit
system:admin:upInfo
system:admin:del
system:admin:disable
```

角色：

```text
system:role:all
system:role:list
system:role:detail
system:role:add
system:role:edit
system:role:del
```

菜单：

```text
system:menu:route
system:menu:list
system:menu:detail
system:menu:add
system:menu:edit
system:menu:del
```

组织：

```text
system:dept:all
system:dept:list
system:dept:detail
system:dept:add
system:dept:edit
system:dept:del
system:post:all
system:post:list
system:post:detail
system:post:add
system:post:edit
system:post:del
```

素材和上传：

```text
material:file:list
material:file:rename
material:file:move
material:file:delete
material:category:list
material:category:add
material:category:rename
material:category:delete
upload:image
upload:video
```

网站和存储：

```text
setting:website:detail
setting:website:save
setting:copyright:detail
setting:copyright:save
setting:protocol:detail
setting:protocol:save
setting:storage:list
setting:storage:detail
setting:storage:edit
setting:storage:change
```

字典：

```text
setting:dict:type:all
setting:dict:type:list
setting:dict:type:detail
setting:dict:type:add
setting:dict:type:edit
setting:dict:type:del
setting:dict:data:all
setting:dict:data:list
setting:dict:data:detail
setting:dict:data:add
setting:dict:data:edit
setting:dict:data:del
```

系统监控和日志：

```text
monitor:server
monitor:cache
system:log:operate
system:log:login
```

代码生成：

```text
gen:db
gen:list
gen:detail
gen:importTable
gen:syncTable
gen:editTable
gen:delTable
gen:previewCode
gen:genCode
gen:downloadCode
```

角色授权：

```text
role: super_admin
grant: all permissions above
```

## 系统设置种子

`ma_setting` 使用 `setting_group`、`setting_key`、`setting_value`，需要支持按 `tenant_id + setting_group + setting_key` 查询。

网站：

```text
website.name = go-makeadmin
website.logo = /api/static/backend_logo.png
website.favicon = /api/static/backend_favicon.ico
website.backdrop = /api/static/backend_backdrop.png
website.copyright = [{"name":"go-makeadmin","link":"https://www.go-makeadmin.cn"}]
```

存储：

```text
storage.default = local
storage.local = {"name":"本地存储"}
storage.qiniu = {"name":"七牛云存储","bucket":"","secretKey":"","accessKey":"","domain":""}
storage.aliyun = {"name":"阿里云存储","bucket":"","secretKey":"","accessKey":"","domain":""}
storage.qcloud = {"name":"腾讯云存储","bucket":"","secretKey":"","accessKey":"","domain":"","region":""}
```

协议：

```text
protocol.service = {"name":"","content":""}
protocol.privacy = {"name":"","content":""}
```

## 字典种子

只种框架级枚举，不放业务枚举。

```text
system_status: enabled=启用, disabled=禁用
menu_type: catalog=目录, page=页面, action=操作
storage_type: local=本地, qiniu=七牛云, aliyun=阿里云, qcloud=腾讯云
data_scope_type: all=全部数据, self=本人数据, org=本组织, org_tree=本组织及下级, custom_org=自定义组织
common_status: 1=启用, 0=禁用
```

## 文件分类种子

```text
image: 图片
video: 视频
```

## P1.3 落地

P1.3 已完成以下事情：

1. 补齐 `ma_setting`、`ma_dict_type`、`ma_dict_item`、`ma_file`、`ma_file_category`、`ma_codegen_*` 的 SQL 和 Go model 草案。
2. 把本文件转成 `sql/p1.seed.sql`。
3. 新增 `sql/p1.schema.sql`。
4. 用独立库 `go_makeadmin` 验证 schema 和 seed 可导入。

验证命令：

```bash
cd /Users/fengrongxin/AI/01-projects/go-makeadmin
./scripts/check-p1-seed.sh
```

验证结果：通过。25 张 `ma_*` 表存在，`admin`、`super_admin`、79 条权限、22 个菜单、12 条设置、5 类字典、16 个字典项和 2 个素材分类均已写入独立库。
