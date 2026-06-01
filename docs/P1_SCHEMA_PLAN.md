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
- P1 支持全新初始化；如需从 `la_*` 迁移，只提供一次性迁移脚本，不做长期双写。

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

## P1.1 当前落地

- 新增 `docs/P1_AUTH_PERMISSION_MODEL.md`，定义认证、权限、组织、租户预留和日志的 `ma_*` SQL 草案。
- 新增 `server/model/makeadmin/auth.go`，定义同名 Go model 草案。
- 本阶段不执行数据库迁移，不接入当前运行链路，不改变 `la_*` 表。

## P1.1 已定事项

- 权限编码先沿用 `module:resource:action`，降低当前菜单和接口权限迁移成本。
- 租户能力先预留但默认关闭，`tenant_id=0` 表示全局默认上下文。
- 数据范围先采用组织树相关策略：`all`、`self`、`org`、`org_tree`、`custom_org`。

## P1.2 当前落地

- 新增 `docs/P1_MINIMAL_SEED.md`，定义 `ma_*` 最小种子清单。
- 最小种子范围包括超级管理员、超级管理员角色、根组织、管理员岗位、全部数据范围、核心菜单、核心权限、网站设置、存储设置、框架级字典和素材基础分类。
- 本阶段不生成真实密码、不写入数据库、不导入初始化数据。

## P1.2 已定事项

- 超级管理员密码必须在初始化时生成，不在仓库中保存默认密码或 hash。
- 核心菜单只覆盖 P0 当前可用后台能力，不包含业务演示模块。
- P1.2 暂保留少量历史双段权限编码，后续切服务时再决定是否统一为三段式并提供兼容映射。

## P1.3 当前落地

- 新增 `server/model/makeadmin/system.go`，补齐 `ma_setting`、`ma_dict_type`、`ma_dict_item`、`ma_file`、`ma_file_category`、`ma_codegen_table`、`ma_codegen_column` Go model 草案。
- 新增 `sql/p1.schema.sql`，集中定义 25 张 `ma_*` P1 表。
- 新增 `sql/p1.seed.sql`，集中定义最小初始化种子。
- 新增 `scripts/check-p1-seed.sh`，用于只读检查 P1 表和种子完整性。
- 已用独立验证库导入 `sql/p1.schema.sql` 和 `sql/p1.seed.sql`，并通过 `./scripts/check-p1-seed.sh`。

## P1.3 已定事项

- go-makeadmin 使用独立数据库，和 zyai 业务库分离；框架完成后再考虑业务迁移。
- P1 SQL 可以在独立验证库中建表和导入，当前不修改 `server/.env`，不切换运行链路。
- `ma_setting` 使用 `setting_group`、`setting_key`、`setting_value` 字段，避免使用 MySQL 保留字。

## P1.4 当前落地

- 新增 `server/makeadmin/repository/auth.go`，定义并实现基于 `ma_*` 的只读认证权限 repository 骨架。
- 新增 `server/makeadmin/service/auth.go`，定义并实现 `BuildIdentityByUsername` 和 `ListRouteMenus` 服务骨架。
- 新增 `server/makeadmin/service/auth_test.go`，覆盖超级管理员权限通配和普通角色菜单父级补全。
- P1.4 不接入当前 `system` 登录服务，不写 Redis，不签发 token，不替换 middleware。

## P1.4 已定事项

- P1 采用并行 repository/service 结构，先在 `server/makeadmin` 内实现，再逐步切换现有 `admin/service/system`。
- 第一阶段只做读取链路：管理员身份、角色权限、菜单路由。
- 密码校验、登录日志写入、token 缓存和最后登录信息更新留到切换登录链路时处理。

## P1.5 当前落地

- 新增 `docs/P1_PASSWORD_STRATEGY.md`，明确 P1 新账号使用 bcrypt，旧 MD5+salt 只做迁移期兼容校验。
- 新增 `server/makeadmin/security/password.go`，实现 bcrypt 生成/校验、旧 MD5+salt 校验、安装占位符拒绝。
- `server/makeadmin/service/auth.go` 新增 `AuthenticateByUsername`，只做账号密码校验和身份构建，不签发 token、不写 Redis、不更新登录信息。
- `sql/p1.seed.sql` 的超级管理员密码改为安装时 bcrypt 占位符，`password_salt` 对新账号保持空字符串。

## P1.5 已定事项

- P1 新写入账号只使用 bcrypt。
- `ma_admin.password_salt` 对新账号为空，只为旧 `la_*` MD5+salt 迁移兼容保留。
- 登录成功后的旧 hash 自动升级需要写库，留到 P1 登录写路径确认后再做。

## P1.6 当前落地

- 新增 `docs/P1_LOGIN_SWITCH_PLAN.md`，定义从旧 `la_*` 登录链路切到 `ma_*` 的顺序。
- `server/makeadmin/repository/auth.go` 新增 `UpdateAdminLoginInfo` 和 `CreateLoginLog` 写接口。
- 新增 `server/makeadmin/service/session.go`，提供 token 生成抽象和 Redis session store。
- `server/makeadmin/service/auth.go` 新增 `Login` 和 `Logout`，完成新框架内部登录闭环：校验密码、签 token、写 Redis session、更新最后登录信息、写 `ma_login_log`。
- `server/makeadmin/service/auth_test.go` 覆盖登录成功写 session/审计、密码失败写失败日志。

## P1.6 已定事项

- P1 登录应用服务先在 `server/makeadmin` 内闭环，不直接替换当前 `/api/system/login`。
- 新 Redis session key 使用 `makeadmin:token:*` 和 `makeadmin:token:set:*`，避免和旧 `backstage:*` 缓存结构耦合。
- 当前不实现旧 MD5+salt 登录成功后的自动升级写回，避免在只读迁移判断尚未完成前改数据。

## P1.7 当前落地

- 新增 `server/makeadmin/adapter/system.go`，把 `makeadmin` 新服务适配到现有后台 API 响应形状。
- `server/admin/routers/system/login.go` 接入新登录适配：检测到可用 `ma_admin` 时走 `ma_*`，否则旧 `la_*` 兜底。
- `server/middleware/auth.go` 支持识别 `makeadmin:token:*` session，并基于 `ma_permission` 校验权限。
- `server/admin/routers/system/admin.go` 的 `/system/admin/self` 支持新 token 分流到 `ma_*`。
- `server/admin/routers/system/menu.go` 的 `/system/menu/route` 支持新 token 分流到 `ma_*`。

## P1.7 已定事项

- P1 切换采用并行分流，不在当前阶段删除旧 `la_*` 运行路径。
- `ma_*` token 与旧 `backstage:*` token 隔离，避免缓存结构互相污染。
- 只有检测到至少一个非占位密码的 `ma_admin` 账号时，登录接口才启用新链路。

## P1.8 当前落地

- 新增 `server/cmd/makeadmin-password`，根据 `MAKEADMIN_PASSWORD` 生成 bcrypt hash。
- 新增 `scripts/init-p1-db.sh`，用于初始化独立 P1 数据库：创建数据库、生成 admin bcrypt hash、导入 `sql/p1.schema.sql` 和替换后的 `sql/p1.seed.sql`、运行 `check-p1-seed`。
- 初始化脚本默认不覆盖已有 `ma_*` 表；需要重建时必须显式设置 `INIT_P1_DROP=1`。

## P1.8 已定事项

- P1 初始化不修改 `.env`。
- P1 初始化不在仓库写入真实密码或 hash。
- 默认 P1 数据库名沿用 `go_makeadmin`，可用 `MYSQL_DATABASE` 覆盖。
- P1 只支持一次性初始化或迁移，不做长期 `la_*`/`ma_*` 双写。

## P1.9 当前落地

- 新增 `server/makeadmin/repository/setting.go`，实现 `ma_setting` 按分组读取和 upsert 写入。
- 新增 `server/makeadmin/service/setting.go`，实现网站基础设置的读取和保存。
- 新增 `server/makeadmin/adapter/website.go`，把 `ma_setting` 适配到现有 `/setting/website/detail` 和 `/setting/website/save` 响应形状。
- `server/admin/routers/setting/website.go` 接入新适配：检测到可用 `ma_setting` 时走 `ma_*`，否则旧 `la_system_config` 兜底。

## P1.9 已定事项

- 网站设置先迁 `name`、`logo`、`favicon`、`backdrop` 四个基础字段。
- `ma_setting` 写入采用 `tenant_id + setting_group + setting_key` upsert，不依赖固定自增 ID。
- 不在本阶段迁移存储密钥、协议内容和备案字段，避免把配置面扩大。

## P1.10 当前落地

- `server/makeadmin/service/setting.go` 新增备案和政策协议读写服务。
- 新增 `server/makeadmin/adapter/website_policy.go`，把 `ma_setting` 适配到现有 `/setting/copyright/*` 和 `/setting/protocol/*` 响应形状。
- `server/admin/routers/setting/copyright.go` 和 `server/admin/routers/setting/protocol.go` 接入新适配：检测到对应 `ma_setting` 时走 `ma_*`，否则旧 `la_system_config` 兜底。
- `sql/p1.seed.sql` 将 `protocol.service` 和 `protocol.privacy` 种子改为 JSON 对象，匹配现有前端协议页数据形状。

## P1.10 已定事项

- 备案继续存储在 `ma_setting.website.copyright`，值为 JSON 数组。
- 政策协议继续存储在 `ma_setting.protocol.service` 和 `ma_setting.protocol.privacy`，值为 JSON 对象。
- 存储设置包含密钥字段，继续留到后续单独处理。

## P1.11 当前落地

- `server/makeadmin/service/setting.go` 新增 `ma_setting.storage` 读写服务，覆盖 `local`、`qiniu`、`aliyun`、`qcloud` 四种别名。
- 新增 `server/makeadmin/adapter/storage.go`，把 `ma_setting.storage` 适配到现有 `/setting/storage/*` 响应形状。
- `server/admin/routers/setting/storage.go` 接入新适配：检测到完整 `ma_setting.storage` 种子时走 `ma_*`，否则旧 `la_system_config` 兜底。
- `server/makeadmin/service/setting_test.go` 新增存储列表、详情、保存和默认存储切换测试。

## P1.11 已定事项

- `storage.default` 只保存当前启用别名，不保存真实密钥。
- 云存储配置只按请求写入 `bucket`、`accessKey`、`secretKey`、`domain`、`region` 字段；仓库种子保持空值。
- 关闭非当前默认存储时不改动 `storage.default`，避免误关其它已启用存储。

## P1.12 当前落地

- 新增 `server/makeadmin/repository/dict.go` 和 `server/makeadmin/service/dict.go`，实现 `ma_dict_type`、`ma_dict_item` 的列表、详情、新增、编辑和软删除。
- 新增 `server/makeadmin/adapter/dict.go`，把 `ma_dict_type`、`ma_dict_item` 适配到现有 `/setting/dict/type/*` 和 `/setting/dict/data/*` 响应形状。
- `server/admin/routers/setting/dict_type.go` 和 `server/admin/routers/setting/dict_data.go` 接入新适配：检测到 `ma_dict_type` 和 `ma_dict_item` 表时走 `ma_*`，否则旧 `la_dict_*` 兜底。
- `server/makeadmin/service/dict_test.go` 覆盖字典类型唯一性、字典数据按类型编码读取和字典数据值唯一性。

## P1.12 已定事项

- 对外 API 继续使用旧字段名：`dictName`、`dictType`、`typeId`、`name`、`value`，避免前端同步改造。
- `ma_dict_type.code` 对应旧 `dict_type`，`ma_dict_type.name` 对应旧 `dict_name`。
- `ma_dict_item.item_label` 对应旧 `name`，`ma_dict_item.item_value` 对应旧 `value`。
- 字典数据唯一性按 `type_id + item_value` 约束，不继续沿用旧链路的全局 `name` 唯一判断。

## P1.13 当前落地

- 新增 `server/makeadmin/repository/file.go` 和 `server/makeadmin/service/file.go`，实现 `ma_file`、`ma_file_category` 的素材列表、分类、移动、重命名、软删除和上传后元数据写入。
- 新增 `server/makeadmin/adapter/file.go`，把 `ma_file`、`ma_file_category` 适配到现有 `/common/album/*` 和 `/common/upload/*` 响应形状。
- `server/admin/routers/common/album.go` 和 `server/admin/routers/common/upload.go` 接入新适配：检测到 `ma_file_category` 与 `ma_file` 表时走 `ma_*`，否则旧 `la_album*` 兜底。
- `server/makeadmin/service/file_test.go` 覆盖文件新增、移动、分类删除保护和分类编码生成。

## P1.13 已定事项

- 对外 API 继续使用旧素材字段：`cid`、`type`、`name`、`uri`、`path`，避免前端素材组件同步改造。
- 旧 `type=10/20/30` 分别映射为 `image`、`video`、`file`，当前种子只内置图片和视频根分类。
- 上传文件的物理存储仍沿用当前 `plugin.StorageDriver`；P1.13 只切换上传成功后的元数据落表。
- 分类删除会额外阻止删除存在子分类的分类，避免产生孤儿分类。

## P1.14 当前落地

- `ma_audit_log` 补充 `ip` 字段，匹配现有操作日志查询和返回字段。
- 新增 `server/makeadmin/repository/log.go`、`server/makeadmin/service/log.go` 和 `server/makeadmin/adapter/log.go`，实现 `ma_login_log`、`ma_audit_log` 的后台日志查询适配。
- `server/admin/routers/system/log.go` 接入新适配：检测到 `ma_login_log` 与 `ma_audit_log` 表时走 `ma_*`，否则旧 `la_system_log_*` 兜底。
- `server/middleware/log.go` 在存在 `ma_audit_log` 时写入新操作审计表，否则继续写旧操作日志表。
- `server/makeadmin/service/log_test.go` 覆盖登录日志和操作审计日志查询过滤参数传递。

## P1.14 已定事项

- `ma_login_log.status` 与 `ma_audit_log.status` 内部保持 `1=成功,0=失败`，对外响应继续映射为旧接口的 `1=成功,2=失败`。
- `ma_audit_log.action` 对应旧操作日志 `title`，`ma_audit_log.path` 对应旧 `url`，`ma_audit_log.request_body` 对应旧 `args`。
- P1.14 不迁移旧 `la_system_log_*` 历史数据，只保证新 P1 链路产生和查询新日志。

## P1.15 当前落地

- 新增 `server/makeadmin/repository/role.go`、`server/makeadmin/service/role.go` 和 `server/makeadmin/adapter/role.go`，实现 `ma_role` 的全部、列表、详情、新增、编辑和软删除。
- `server/admin/routers/system/role.go` 接入新适配：检测到 `ma_role`、`ma_role_permission` 和 `ma_menu_permission` 表时走 `ma_*`，否则旧 `la_system_auth_role` 兜底。
- 角色授权继续接收旧接口 `menuIds`，在 P1 内部通过 `ma_menu_permission` 展开为 `ma_role_permission`。
- `server/makeadmin/service/role_test.go` 覆盖角色新增、重名校验、详情菜单/成员数、系统角色保护、使用中角色保护和 `menuIds` 解析。

## P1.15 已定事项

- 对外接口继续使用旧字段：`name`、`remark`、`sort`、`isDisable`、`menus`、`menuIds`。
- `ma_role.status=1/0` 对外映射为 `isDisable=0/1`。
- 旧接口没有角色编码字段，P1 新增角色自动生成内部 `code=role_<timestamp>`。
- 系统角色 `is_system=1` 不允许通过后台角色删除接口删除。

## 已定事项

- `la_* -> ma_*` 只支持一次性迁移。
- P1 独立库最终命名沿用 `go_makeadmin`。
