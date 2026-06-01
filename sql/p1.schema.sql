-- go-makeadmin P1 ma_* schema draft.
-- Intended for a dedicated go-makeadmin database or disposable validation database.

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `ma_codegen_column`;
DROP TABLE IF EXISTS `ma_codegen_table`;
DROP TABLE IF EXISTS `ma_file`;
DROP TABLE IF EXISTS `ma_file_category`;
DROP TABLE IF EXISTS `ma_dict_item`;
DROP TABLE IF EXISTS `ma_dict_type`;
DROP TABLE IF EXISTS `ma_setting`;
DROP TABLE IF EXISTS `ma_audit_log`;
DROP TABLE IF EXISTS `ma_login_log`;
DROP TABLE IF EXISTS `ma_role_data_scope`;
DROP TABLE IF EXISTS `ma_data_scope`;
DROP TABLE IF EXISTS `ma_admin_org`;
DROP TABLE IF EXISTS `ma_position`;
DROP TABLE IF EXISTS `ma_org_unit`;
DROP TABLE IF EXISTS `ma_menu_permission`;
DROP TABLE IF EXISTS `ma_menu`;
DROP TABLE IF EXISTS `ma_role_permission`;
DROP TABLE IF EXISTS `ma_permission`;
DROP TABLE IF EXISTS `ma_admin_role`;
DROP TABLE IF EXISTS `ma_role`;
DROP TABLE IF EXISTS `ma_tenant_setting`;
DROP TABLE IF EXISTS `ma_tenant_member`;
DROP TABLE IF EXISTS `ma_tenant`;
DROP TABLE IF EXISTS `ma_admin_profile`;
DROP TABLE IF EXISTS `ma_admin`;

CREATE TABLE `ma_admin` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(64) NOT NULL DEFAULT '',
  `password_hash` VARCHAR(255) NOT NULL DEFAULT '',
  `password_salt` VARCHAR(64) NOT NULL DEFAULT '',
  `is_super` TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `last_login_ip` VARCHAR(64) NOT NULL DEFAULT '',
  `last_login_time` BIGINT NOT NULL DEFAULT 0,
  `create_time` BIGINT NOT NULL DEFAULT 0,
  `update_time` BIGINT NOT NULL DEFAULT 0,
  `delete_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ma_admin_username_live` (`username`, `delete_time`),
  KEY `idx_ma_admin_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='admin account';

CREATE TABLE `ma_admin_profile` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `admin_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `nickname` VARCHAR(64) NOT NULL DEFAULT '',
  `avatar` VARCHAR(255) NOT NULL DEFAULT '',
  `email` VARCHAR(128) NOT NULL DEFAULT '',
  `mobile` VARCHAR(32) NOT NULL DEFAULT '',
  `remark` VARCHAR(255) NOT NULL DEFAULT '',
  `create_time` BIGINT NOT NULL DEFAULT 0,
  `update_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ma_admin_profile_admin` (`admin_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='admin profile';

CREATE TABLE `ma_tenant` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(64) NOT NULL DEFAULT '',
  `name` VARCHAR(128) NOT NULL DEFAULT '',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `create_time` BIGINT NOT NULL DEFAULT 0,
  `update_time` BIGINT NOT NULL DEFAULT 0,
  `delete_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ma_tenant_code_live` (`code`, `delete_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='tenant';

CREATE TABLE `ma_tenant_member` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `admin_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `member_type` VARCHAR(32) NOT NULL DEFAULT 'member',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `create_time` BIGINT NOT NULL DEFAULT 0,
  `update_time` BIGINT NOT NULL DEFAULT 0,
  `delete_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ma_tenant_member_live` (`tenant_id`, `admin_id`, `delete_time`),
  KEY `idx_ma_tenant_member_admin` (`admin_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='tenant member';

CREATE TABLE `ma_tenant_setting` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `setting_key` VARCHAR(128) NOT NULL DEFAULT '',
  `setting_value` TEXT NOT NULL,
  `create_time` BIGINT NOT NULL DEFAULT 0,
  `update_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ma_tenant_setting_key` (`tenant_id`, `setting_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='tenant setting';

CREATE TABLE `ma_role` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `code` VARCHAR(64) NOT NULL DEFAULT '',
  `name` VARCHAR(64) NOT NULL DEFAULT '',
  `remark` VARCHAR(255) NOT NULL DEFAULT '',
  `is_system` TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `sort` SMALLINT UNSIGNED NOT NULL DEFAULT 0,
  `create_time` BIGINT NOT NULL DEFAULT 0,
  `update_time` BIGINT NOT NULL DEFAULT 0,
  `delete_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ma_role_code_live` (`tenant_id`, `code`, `delete_time`),
  KEY `idx_ma_role_tenant_status` (`tenant_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='role';

CREATE TABLE `ma_admin_role` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `admin_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `role_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `create_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ma_admin_role` (`tenant_id`, `admin_id`, `role_id`),
  KEY `idx_ma_admin_role_role` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='admin role';

CREATE TABLE `ma_permission` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(128) NOT NULL DEFAULT '',
  `name` VARCHAR(64) NOT NULL DEFAULT '',
  `module` VARCHAR(64) NOT NULL DEFAULT '',
  `resource` VARCHAR(64) NOT NULL DEFAULT '',
  `action` VARCHAR(64) NOT NULL DEFAULT '',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `sort` SMALLINT UNSIGNED NOT NULL DEFAULT 0,
  `create_time` BIGINT NOT NULL DEFAULT 0,
  `update_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ma_permission_code` (`code`),
  KEY `idx_ma_permission_module` (`module`, `resource`, `action`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='permission';

CREATE TABLE `ma_role_permission` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `role_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `permission_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `create_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ma_role_permission` (`tenant_id`, `role_id`, `permission_id`),
  KEY `idx_ma_role_permission_permission` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='role permission';

CREATE TABLE `ma_menu` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `menu_type` VARCHAR(16) NOT NULL DEFAULT 'page',
  `name` VARCHAR(64) NOT NULL DEFAULT '',
  `icon` VARCHAR(64) NOT NULL DEFAULT '',
  `route_path` VARCHAR(128) NOT NULL DEFAULT '',
  `route_name` VARCHAR(128) NOT NULL DEFAULT '',
  `component` VARCHAR(255) NOT NULL DEFAULT '',
  `redirect` VARCHAR(255) NOT NULL DEFAULT '',
  `active_path` VARCHAR(128) NOT NULL DEFAULT '',
  `meta` TEXT NOT NULL,
  `is_visible` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `is_cache` TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `sort` SMALLINT UNSIGNED NOT NULL DEFAULT 0,
  `create_time` BIGINT NOT NULL DEFAULT 0,
  `update_time` BIGINT NOT NULL DEFAULT 0,
  `delete_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_ma_menu_parent_sort` (`parent_id`, `sort`),
  KEY `idx_ma_menu_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='menu';

CREATE TABLE `ma_menu_permission` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `menu_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `permission_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `create_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ma_menu_permission` (`menu_id`, `permission_id`),
  KEY `idx_ma_menu_permission_permission` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='menu permission';

CREATE TABLE `ma_org_unit` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `code` VARCHAR(64) NOT NULL DEFAULT '',
  `name` VARCHAR(128) NOT NULL DEFAULT '',
  `leader_admin_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `sort` SMALLINT UNSIGNED NOT NULL DEFAULT 0,
  `create_time` BIGINT NOT NULL DEFAULT 0,
  `update_time` BIGINT NOT NULL DEFAULT 0,
  `delete_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ma_org_unit_code_live` (`tenant_id`, `code`, `delete_time`),
  KEY `idx_ma_org_unit_parent_sort` (`tenant_id`, `parent_id`, `sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='organization unit';

CREATE TABLE `ma_position` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `code` VARCHAR(64) NOT NULL DEFAULT '',
  `name` VARCHAR(64) NOT NULL DEFAULT '',
  `remark` VARCHAR(255) NOT NULL DEFAULT '',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `sort` SMALLINT UNSIGNED NOT NULL DEFAULT 0,
  `create_time` BIGINT NOT NULL DEFAULT 0,
  `update_time` BIGINT NOT NULL DEFAULT 0,
  `delete_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ma_position_code_live` (`tenant_id`, `code`, `delete_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='position';

CREATE TABLE `ma_admin_org` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `admin_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `org_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `position_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `is_primary` TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `create_time` BIGINT NOT NULL DEFAULT 0,
  `update_time` BIGINT NOT NULL DEFAULT 0,
  `delete_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ma_admin_org_live` (`tenant_id`, `admin_id`, `org_id`, `position_id`, `delete_time`),
  KEY `idx_ma_admin_org_org` (`tenant_id`, `org_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='admin organization';

CREATE TABLE `ma_data_scope` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `code` VARCHAR(64) NOT NULL DEFAULT '',
  `name` VARCHAR(64) NOT NULL DEFAULT '',
  `scope_type` VARCHAR(32) NOT NULL DEFAULT 'self',
  `scope_value` TEXT NOT NULL,
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `create_time` BIGINT NOT NULL DEFAULT 0,
  `update_time` BIGINT NOT NULL DEFAULT 0,
  `delete_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ma_data_scope_code_live` (`tenant_id`, `code`, `delete_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='data scope';

CREATE TABLE `ma_role_data_scope` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `role_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `data_scope_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `create_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ma_role_data_scope` (`tenant_id`, `role_id`, `data_scope_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='role data scope';

CREATE TABLE `ma_login_log` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `admin_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `username` VARCHAR(64) NOT NULL DEFAULT '',
  `ip` VARCHAR(64) NOT NULL DEFAULT '',
  `os` VARCHAR(64) NOT NULL DEFAULT '',
  `browser` VARCHAR(64) NOT NULL DEFAULT '',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `message` VARCHAR(255) NOT NULL DEFAULT '',
  `create_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_ma_login_log_admin_time` (`admin_id`, `create_time`),
  KEY `idx_ma_login_log_tenant_time` (`tenant_id`, `create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='login log';

CREATE TABLE `ma_audit_log` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `admin_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `trace_id` VARCHAR(64) NOT NULL DEFAULT '',
  `action` VARCHAR(128) NOT NULL DEFAULT '',
  `method` VARCHAR(16) NOT NULL DEFAULT '',
  `path` VARCHAR(255) NOT NULL DEFAULT '',
  `request_body` TEXT NOT NULL,
  `response_code` INT NOT NULL DEFAULT 0,
  `error` TEXT NOT NULL,
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `start_time` BIGINT NOT NULL DEFAULT 0,
  `end_time` BIGINT NOT NULL DEFAULT 0,
  `duration_ms` BIGINT NOT NULL DEFAULT 0,
  `create_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_ma_audit_log_admin_time` (`admin_id`, `create_time`),
  KEY `idx_ma_audit_log_tenant_time` (`tenant_id`, `create_time`),
  KEY `idx_ma_audit_log_trace` (`trace_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='audit log';

CREATE TABLE `ma_setting` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `setting_group` VARCHAR(64) NOT NULL DEFAULT '',
  `setting_key` VARCHAR(128) NOT NULL DEFAULT '',
  `setting_value` TEXT NOT NULL,
  `value_type` VARCHAR(32) NOT NULL DEFAULT 'string',
  `is_public` TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `remark` VARCHAR(255) NOT NULL DEFAULT '',
  `create_time` BIGINT NOT NULL DEFAULT 0,
  `update_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ma_setting_key` (`tenant_id`, `setting_group`, `setting_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='setting';

CREATE TABLE `ma_dict_type` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(64) NOT NULL DEFAULT '',
  `name` VARCHAR(64) NOT NULL DEFAULT '',
  `remark` VARCHAR(255) NOT NULL DEFAULT '',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `sort` SMALLINT UNSIGNED NOT NULL DEFAULT 0,
  `create_time` BIGINT NOT NULL DEFAULT 0,
  `update_time` BIGINT NOT NULL DEFAULT 0,
  `delete_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ma_dict_type_code_live` (`code`, `delete_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='dict type';

CREATE TABLE `ma_dict_item` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `type_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `item_label` VARCHAR(64) NOT NULL DEFAULT '',
  `item_value` VARCHAR(128) NOT NULL DEFAULT '',
  `remark` VARCHAR(255) NOT NULL DEFAULT '',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `sort` SMALLINT UNSIGNED NOT NULL DEFAULT 0,
  `create_time` BIGINT NOT NULL DEFAULT 0,
  `update_time` BIGINT NOT NULL DEFAULT 0,
  `delete_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ma_dict_item_value_live` (`type_id`, `item_value`, `delete_time`),
  KEY `idx_ma_dict_item_type` (`type_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='dict item';

CREATE TABLE `ma_file_category` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `code` VARCHAR(64) NOT NULL DEFAULT '',
  `name` VARCHAR(64) NOT NULL DEFAULT '',
  `file_type` VARCHAR(32) NOT NULL DEFAULT 'image',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `sort` SMALLINT UNSIGNED NOT NULL DEFAULT 0,
  `create_time` BIGINT NOT NULL DEFAULT 0,
  `update_time` BIGINT NOT NULL DEFAULT 0,
  `delete_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ma_file_category_code_live` (`tenant_id`, `code`, `delete_time`),
  KEY `idx_ma_file_category_parent_sort` (`tenant_id`, `parent_id`, `sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='file category';

CREATE TABLE `ma_file` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `category_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `owner_admin_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `file_type` VARCHAR(32) NOT NULL DEFAULT 'image',
  `storage_driver` VARCHAR(32) NOT NULL DEFAULT 'local',
  `original_name` VARCHAR(255) NOT NULL DEFAULT '',
  `file_name` VARCHAR(255) NOT NULL DEFAULT '',
  `uri` VARCHAR(512) NOT NULL DEFAULT '',
  `url` VARCHAR(512) NOT NULL DEFAULT '',
  `mime_type` VARCHAR(128) NOT NULL DEFAULT '',
  `ext` VARCHAR(32) NOT NULL DEFAULT '',
  `size` BIGINT NOT NULL DEFAULT 0,
  `checksum` VARCHAR(128) NOT NULL DEFAULT '',
  `status` TINYINT UNSIGNED NOT NULL DEFAULT 1,
  `create_time` BIGINT NOT NULL DEFAULT 0,
  `update_time` BIGINT NOT NULL DEFAULT 0,
  `delete_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_ma_file_tenant_category` (`tenant_id`, `category_id`),
  KEY `idx_ma_file_owner` (`owner_admin_id`),
  KEY `idx_ma_file_checksum` (`checksum`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='file';

CREATE TABLE `ma_codegen_table` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `tenant_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `table_name` VARCHAR(128) NOT NULL DEFAULT '',
  `table_comment` VARCHAR(255) NOT NULL DEFAULT '',
  `module_name` VARCHAR(64) NOT NULL DEFAULT '',
  `package_name` VARCHAR(128) NOT NULL DEFAULT '',
  `business_name` VARCHAR(64) NOT NULL DEFAULT '',
  `entity_name` VARCHAR(64) NOT NULL DEFAULT '',
  `function_name` VARCHAR(64) NOT NULL DEFAULT '',
  `author_name` VARCHAR(64) NOT NULL DEFAULT '',
  `template_type` VARCHAR(32) NOT NULL DEFAULT 'crud',
  `generate_type` VARCHAR(32) NOT NULL DEFAULT 'zip',
  `generate_path` VARCHAR(255) NOT NULL DEFAULT '',
  `options` TEXT NOT NULL,
  `remark` VARCHAR(255) NOT NULL DEFAULT '',
  `create_time` BIGINT NOT NULL DEFAULT 0,
  `update_time` BIGINT NOT NULL DEFAULT 0,
  `delete_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ma_codegen_table_live` (`tenant_id`, `table_name`, `delete_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='codegen table';

CREATE TABLE `ma_codegen_column` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `table_id` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `column_name` VARCHAR(128) NOT NULL DEFAULT '',
  `column_comment` VARCHAR(255) NOT NULL DEFAULT '',
  `column_type` VARCHAR(64) NOT NULL DEFAULT '',
  `column_length` INT NOT NULL DEFAULT 0,
  `go_type` VARCHAR(64) NOT NULL DEFAULT '',
  `go_field` VARCHAR(64) NOT NULL DEFAULT '',
  `json_field` VARCHAR(64) NOT NULL DEFAULT '',
  `is_pk` TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `is_increment` TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `is_required` TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `is_insert` TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `is_edit` TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `is_list` TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `is_query` TINYINT UNSIGNED NOT NULL DEFAULT 0,
  `query_type` VARCHAR(32) NOT NULL DEFAULT '=',
  `html_type` VARCHAR(32) NOT NULL DEFAULT '',
  `dict_type` VARCHAR(64) NOT NULL DEFAULT '',
  `sort` SMALLINT UNSIGNED NOT NULL DEFAULT 0,
  `create_time` BIGINT NOT NULL DEFAULT 0,
  `update_time` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_ma_codegen_column` (`table_id`, `column_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='codegen column';

SET FOREIGN_KEY_CHECKS = 1;
