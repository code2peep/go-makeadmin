# P1.1 认证权限模型草案

更新时间：2026-05-31

## 范围

本文件定义 P1 第一小步的 `ma_*` 认证、权限、组织、租户预留和日志模型草案。

本阶段只做设计和可编译 Go model 草案：

- 不连接或迁移任何真实业务库。
- 不执行 `AutoMigrate`。
- 不修改当前 `la_*` 登录、菜单和权限运行链路。
- 不导入初始化数据。

## 已定边界

- 权限编码沿用 `module:resource:action`，例如 `system:admin:list`。改成点号编码收益不高，会增加 P1 迁移成本。
- 菜单和权限分离：`ma_menu` 只描述前端导航，`ma_permission` 描述后端能力，二者通过 `ma_menu_permission` 关联。
- 管理员和角色多对多：`ma_admin_role` 替代当前单字段角色绑定。
- 组织和岗位不直接塞进管理员主表：`ma_admin_org` 表达管理员在组织内的岗位和主组织关系。
- 租户默认关闭，`tenant_id=0` 表示全局默认上下文；后续启用租户时再加 middleware 和租户上下文校验。
- `ma_permission` 和 `ma_menu` 暂按全局能力定义；`ma_role`、组织、数据范围和角色授权按 `tenant_id` 隔离。
- 软删除使用 `delete_time=0` 表示未删除，便于唯一索引允许删除后复用编码。

## 表分组

认证与账号：

```text
ma_admin
ma_admin_profile
ma_admin_role
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
ma_admin_org
ma_data_scope
ma_role_data_scope
```

租户预留：

```text
ma_tenant
ma_tenant_member
ma_tenant_setting
```

审计：

```text
ma_audit_log
```

## SQL 草案

以下 SQL 是 P1.1 认证权限模型草案。P1.3 已将完整可导入版本集中到 `sql/p1.schema.sql`，最小种子集中到 `sql/p1.seed.sql`。

P1.5 已确定密码策略：新账号 `password_hash` 保存 bcrypt hash，`password_salt` 为空；旧 `la_*` 的 MD5+salt 只在迁移期兼容校验。

```sql
CREATE TABLE ma_admin (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  username VARCHAR(64) NOT NULL DEFAULT '',
  password_hash VARCHAR(255) NOT NULL DEFAULT '',
  password_salt VARCHAR(64) NOT NULL DEFAULT '',
  is_super TINYINT UNSIGNED NOT NULL DEFAULT 0,
  status TINYINT UNSIGNED NOT NULL DEFAULT 1,
  last_login_ip VARCHAR(64) NOT NULL DEFAULT '',
  last_login_time BIGINT NOT NULL DEFAULT 0,
  create_time BIGINT NOT NULL DEFAULT 0,
  update_time BIGINT NOT NULL DEFAULT 0,
  delete_time BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE KEY uk_ma_admin_username_live (username, delete_time),
  KEY idx_ma_admin_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='admin account';

CREATE TABLE ma_admin_profile (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  admin_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  nickname VARCHAR(64) NOT NULL DEFAULT '',
  avatar VARCHAR(255) NOT NULL DEFAULT '',
  email VARCHAR(128) NOT NULL DEFAULT '',
  mobile VARCHAR(32) NOT NULL DEFAULT '',
  remark VARCHAR(255) NOT NULL DEFAULT '',
  create_time BIGINT NOT NULL DEFAULT 0,
  update_time BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE KEY uk_ma_admin_profile_admin (admin_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='admin profile';

CREATE TABLE ma_tenant (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  code VARCHAR(64) NOT NULL DEFAULT '',
  name VARCHAR(128) NOT NULL DEFAULT '',
  status TINYINT UNSIGNED NOT NULL DEFAULT 1,
  create_time BIGINT NOT NULL DEFAULT 0,
  update_time BIGINT NOT NULL DEFAULT 0,
  delete_time BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE KEY uk_ma_tenant_code_live (code, delete_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='tenant';

CREATE TABLE ma_tenant_member (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  tenant_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  admin_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  member_type VARCHAR(32) NOT NULL DEFAULT 'member',
  status TINYINT UNSIGNED NOT NULL DEFAULT 1,
  create_time BIGINT NOT NULL DEFAULT 0,
  update_time BIGINT NOT NULL DEFAULT 0,
  delete_time BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE KEY uk_ma_tenant_member_live (tenant_id, admin_id, delete_time),
  KEY idx_ma_tenant_member_admin (admin_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='tenant member';

CREATE TABLE ma_tenant_setting (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  tenant_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  setting_key VARCHAR(128) NOT NULL DEFAULT '',
  setting_value TEXT NOT NULL,
  create_time BIGINT NOT NULL DEFAULT 0,
  update_time BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE KEY uk_ma_tenant_setting_key (tenant_id, setting_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='tenant setting';

CREATE TABLE ma_role (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  tenant_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  code VARCHAR(64) NOT NULL DEFAULT '',
  name VARCHAR(64) NOT NULL DEFAULT '',
  remark VARCHAR(255) NOT NULL DEFAULT '',
  is_system TINYINT UNSIGNED NOT NULL DEFAULT 0,
  status TINYINT UNSIGNED NOT NULL DEFAULT 1,
  sort SMALLINT UNSIGNED NOT NULL DEFAULT 0,
  create_time BIGINT NOT NULL DEFAULT 0,
  update_time BIGINT NOT NULL DEFAULT 0,
  delete_time BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE KEY uk_ma_role_code_live (tenant_id, code, delete_time),
  KEY idx_ma_role_tenant_status (tenant_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='role';

CREATE TABLE ma_admin_role (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  tenant_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  admin_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  role_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  create_time BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE KEY uk_ma_admin_role (tenant_id, admin_id, role_id),
  KEY idx_ma_admin_role_role (role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='admin role';

CREATE TABLE ma_permission (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  code VARCHAR(128) NOT NULL DEFAULT '',
  name VARCHAR(64) NOT NULL DEFAULT '',
  module VARCHAR(64) NOT NULL DEFAULT '',
  resource VARCHAR(64) NOT NULL DEFAULT '',
  action VARCHAR(64) NOT NULL DEFAULT '',
  status TINYINT UNSIGNED NOT NULL DEFAULT 1,
  sort SMALLINT UNSIGNED NOT NULL DEFAULT 0,
  create_time BIGINT NOT NULL DEFAULT 0,
  update_time BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE KEY uk_ma_permission_code (code),
  KEY idx_ma_permission_module (module, resource, action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='permission';

CREATE TABLE ma_role_permission (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  tenant_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  role_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  permission_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  create_time BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE KEY uk_ma_role_permission (tenant_id, role_id, permission_id),
  KEY idx_ma_role_permission_permission (permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='role permission';

CREATE TABLE ma_menu (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  parent_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  menu_type VARCHAR(16) NOT NULL DEFAULT 'page',
  name VARCHAR(64) NOT NULL DEFAULT '',
  icon VARCHAR(64) NOT NULL DEFAULT '',
  route_path VARCHAR(128) NOT NULL DEFAULT '',
  route_name VARCHAR(128) NOT NULL DEFAULT '',
  component VARCHAR(255) NOT NULL DEFAULT '',
  redirect VARCHAR(255) NOT NULL DEFAULT '',
  active_path VARCHAR(128) NOT NULL DEFAULT '',
  meta TEXT NOT NULL,
  is_visible TINYINT UNSIGNED NOT NULL DEFAULT 1,
  is_cache TINYINT UNSIGNED NOT NULL DEFAULT 0,
  status TINYINT UNSIGNED NOT NULL DEFAULT 1,
  sort SMALLINT UNSIGNED NOT NULL DEFAULT 0,
  create_time BIGINT NOT NULL DEFAULT 0,
  update_time BIGINT NOT NULL DEFAULT 0,
  delete_time BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  KEY idx_ma_menu_parent_sort (parent_id, sort),
  KEY idx_ma_menu_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='menu';

CREATE TABLE ma_menu_permission (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  menu_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  permission_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  create_time BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE KEY uk_ma_menu_permission (menu_id, permission_id),
  KEY idx_ma_menu_permission_permission (permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='menu permission';

CREATE TABLE ma_org_unit (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  tenant_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  parent_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  code VARCHAR(64) NOT NULL DEFAULT '',
  name VARCHAR(128) NOT NULL DEFAULT '',
  leader_admin_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  status TINYINT UNSIGNED NOT NULL DEFAULT 1,
  sort SMALLINT UNSIGNED NOT NULL DEFAULT 0,
  create_time BIGINT NOT NULL DEFAULT 0,
  update_time BIGINT NOT NULL DEFAULT 0,
  delete_time BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE KEY uk_ma_org_unit_code_live (tenant_id, code, delete_time),
  KEY idx_ma_org_unit_parent_sort (tenant_id, parent_id, sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='organization unit';

CREATE TABLE ma_position (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  tenant_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  code VARCHAR(64) NOT NULL DEFAULT '',
  name VARCHAR(64) NOT NULL DEFAULT '',
  remark VARCHAR(255) NOT NULL DEFAULT '',
  status TINYINT UNSIGNED NOT NULL DEFAULT 1,
  sort SMALLINT UNSIGNED NOT NULL DEFAULT 0,
  create_time BIGINT NOT NULL DEFAULT 0,
  update_time BIGINT NOT NULL DEFAULT 0,
  delete_time BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE KEY uk_ma_position_code_live (tenant_id, code, delete_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='position';

CREATE TABLE ma_admin_org (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  tenant_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  admin_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  org_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  position_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  is_primary TINYINT UNSIGNED NOT NULL DEFAULT 0,
  status TINYINT UNSIGNED NOT NULL DEFAULT 1,
  create_time BIGINT NOT NULL DEFAULT 0,
  update_time BIGINT NOT NULL DEFAULT 0,
  delete_time BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE KEY uk_ma_admin_org_live (tenant_id, admin_id, org_id, position_id, delete_time),
  KEY idx_ma_admin_org_org (tenant_id, org_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='admin organization';

CREATE TABLE ma_data_scope (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  tenant_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  code VARCHAR(64) NOT NULL DEFAULT '',
  name VARCHAR(64) NOT NULL DEFAULT '',
  scope_type VARCHAR(32) NOT NULL DEFAULT 'self',
  scope_value TEXT NOT NULL,
  status TINYINT UNSIGNED NOT NULL DEFAULT 1,
  create_time BIGINT NOT NULL DEFAULT 0,
  update_time BIGINT NOT NULL DEFAULT 0,
  delete_time BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE KEY uk_ma_data_scope_code_live (tenant_id, code, delete_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='data scope';

CREATE TABLE ma_role_data_scope (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  tenant_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  role_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  data_scope_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  create_time BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE KEY uk_ma_role_data_scope (tenant_id, role_id, data_scope_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='role data scope';

CREATE TABLE ma_login_log (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  tenant_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  admin_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  username VARCHAR(64) NOT NULL DEFAULT '',
  ip VARCHAR(64) NOT NULL DEFAULT '',
  os VARCHAR(64) NOT NULL DEFAULT '',
  browser VARCHAR(64) NOT NULL DEFAULT '',
  status TINYINT UNSIGNED NOT NULL DEFAULT 0,
  message VARCHAR(255) NOT NULL DEFAULT '',
  create_time BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  KEY idx_ma_login_log_admin_time (admin_id, create_time),
  KEY idx_ma_login_log_tenant_time (tenant_id, create_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='login log';

CREATE TABLE ma_audit_log (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  tenant_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  admin_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
  trace_id VARCHAR(64) NOT NULL DEFAULT '',
  action VARCHAR(128) NOT NULL DEFAULT '',
  method VARCHAR(16) NOT NULL DEFAULT '',
  path VARCHAR(255) NOT NULL DEFAULT '',
  request_body TEXT NOT NULL,
  response_code INT NOT NULL DEFAULT 0,
  error TEXT NOT NULL,
  status TINYINT UNSIGNED NOT NULL DEFAULT 0,
  start_time BIGINT NOT NULL DEFAULT 0,
  end_time BIGINT NOT NULL DEFAULT 0,
  duration_ms BIGINT NOT NULL DEFAULT 0,
  create_time BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  KEY idx_ma_audit_log_admin_time (admin_id, create_time),
  KEY idx_ma_audit_log_tenant_time (tenant_id, create_time),
  KEY idx_ma_audit_log_trace (trace_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='audit log';
```

## 后续切换顺序

1. 新增 `ma_*` repository/service，不替换当前接口。
2. 登录链路切到 `ma_admin`、`ma_admin_profile`、`ma_admin_role`。
3. 菜单和权限链路切到 `ma_menu`、`ma_permission`、`ma_role_permission`。
4. 组织、数据范围、日志按模块切换。
5. 冻结 `la_*` 运行路径，只保留迁移文档和一次性迁移脚本。
