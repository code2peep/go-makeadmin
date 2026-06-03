-- go-makeadmin P1 minimal ma_* seed draft.
-- This file intentionally uses an install-time bcrypt placeholder.
-- Replace it at install time; do not use this placeholder in a real login flow.

SET NAMES utf8mb4;
SET @seed_time = 1700000000;

INSERT INTO `ma_admin` (`id`, `username`, `password_hash`, `password_salt`, `is_super`, `status`, `last_login_ip`, `last_login_time`, `create_time`, `update_time`, `delete_time`) VALUES
(1, 'admin', 'INSTALL_TIME_PASSWORD_BCRYPT_REPLACE_ME', '', 1, 1, '', 0, @seed_time, @seed_time, 0);

INSERT INTO `ma_admin_profile` (`id`, `admin_id`, `nickname`, `avatar`, `email`, `mobile`, `remark`, `create_time`, `update_time`) VALUES
(1, 1, 'admin', '/api/static/backend_avatar.png', '', '', 'Initial super administrator profile.', @seed_time, @seed_time);

INSERT INTO `ma_role` (`id`, `tenant_id`, `code`, `name`, `remark`, `is_system`, `status`, `sort`, `create_time`, `update_time`, `delete_time`) VALUES
(1, 0, 'super_admin', '超级管理员', 'Built-in full access role.', 1, 1, 1000, @seed_time, @seed_time, 0);

INSERT INTO `ma_admin_role` (`id`, `tenant_id`, `admin_id`, `role_id`, `create_time`) VALUES
(1, 0, 1, 1, @seed_time);

INSERT INTO `ma_org_unit` (`id`, `tenant_id`, `parent_id`, `code`, `name`, `leader_admin_id`, `status`, `sort`, `create_time`, `update_time`, `delete_time`) VALUES
(1, 0, 0, 'root', '总部', 1, 1, 1000, @seed_time, @seed_time, 0);

INSERT INTO `ma_position` (`id`, `tenant_id`, `code`, `name`, `remark`, `status`, `sort`, `create_time`, `update_time`, `delete_time`) VALUES
(1, 0, 'admin', '管理员', 'Default administrator position.', 1, 1000, @seed_time, @seed_time, 0);

INSERT INTO `ma_admin_org` (`id`, `tenant_id`, `admin_id`, `org_id`, `position_id`, `is_primary`, `status`, `create_time`, `update_time`, `delete_time`) VALUES
(1, 0, 1, 1, 1, 1, 1, @seed_time, @seed_time, 0);

INSERT INTO `ma_data_scope` (`id`, `tenant_id`, `code`, `name`, `scope_type`, `scope_value`, `status`, `create_time`, `update_time`, `delete_time`) VALUES
(1, 0, 'all', '全部数据', 'all', '{}', 1, @seed_time, @seed_time, 0);

INSERT INTO `ma_role_data_scope` (`id`, `tenant_id`, `role_id`, `data_scope_id`, `create_time`) VALUES
(1, 0, 1, 1, @seed_time);

INSERT INTO `ma_menu` (`id`, `parent_id`, `menu_type`, `name`, `icon`, `route_path`, `route_name`, `component`, `redirect`, `active_path`, `meta`, `is_visible`, `is_cache`, `status`, `sort`, `create_time`, `update_time`, `delete_time`) VALUES
(1, 0, 'page', '工作台', 'el-icon-Monitor', '/workbench', 'dashboard.workbench', 'workbench/index', '', '', '{}', 1, 1, 1, 1000, @seed_time, @seed_time, 0),
(100, 0, 'catalog', '权限管理', 'el-icon-Lock', '/permission', 'permission', '', '', '', '{}', 1, 0, 1, 900, @seed_time, @seed_time, 0),
(101, 100, 'page', '管理员', 'local-icon-wode', '/permission/admin', 'permission.admin', 'permission/admin/index', '', '', '{}', 1, 1, 1, 300, @seed_time, @seed_time, 0),
(110, 100, 'page', '角色管理', 'el-icon-Female', '/permission/role', 'permission.role', 'permission/role/index', '', '', '{}', 1, 1, 1, 200, @seed_time, @seed_time, 0),
(120, 100, 'page', '菜单管理', 'el-icon-Operation', '/permission/menu', 'permission.menu', 'permission/menu/index', '', '', '{}', 1, 1, 1, 100, @seed_time, @seed_time, 0),
(130, 0, 'catalog', '组织管理', 'el-icon-OfficeBuilding', '/organization', 'organization', '', '', '', '{}', 1, 0, 1, 800, @seed_time, @seed_time, 0),
(131, 130, 'page', '部门管理', 'el-icon-Coordinate', '/organization/department', 'organization.department', 'organization/department/index', '', '', '{}', 1, 1, 1, 200, @seed_time, @seed_time, 0),
(140, 130, 'page', '岗位管理', 'el-icon-PriceTag', '/organization/post', 'organization.post', 'organization/post/index', '', '', '{}', 1, 1, 1, 100, @seed_time, @seed_time, 0),
(700, 0, 'page', '素材管理', 'el-icon-Picture', '/material/index', 'material', 'material/index', '', '', '{}', 1, 1, 1, 700, @seed_time, @seed_time, 0),
(500, 0, 'catalog', '系统设置', 'el-icon-Setting', '/setting', 'setting', '', '', '', '{}', 1, 0, 1, 600, @seed_time, @seed_time, 0),
(501, 500, 'catalog', '网站设置', 'el-icon-Basketball', '/setting/website', 'setting.website', '', '', '', '{}', 1, 0, 1, 300, @seed_time, @seed_time, 0),
(502, 501, 'page', '网站信息', '', '/setting/website/information', 'setting.website.information', 'setting/website/information', '', '', '{}', 1, 0, 1, 300, @seed_time, @seed_time, 0),
(505, 501, 'page', '网站备案', '', '/setting/website/filing', 'setting.website.filing', 'setting/website/filing', '', '', '{}', 1, 0, 1, 200, @seed_time, @seed_time, 0),
(510, 501, 'page', '政策协议', '', '/setting/website/protocol', 'setting.website.protocol', 'setting/website/protocol', '', '', '{}', 1, 0, 1, 100, @seed_time, @seed_time, 0),
(555, 500, 'page', '存储设置', 'el-icon-FolderOpened', '/setting/storage', 'setting.storage', 'setting/storage/index', '', '', '{}', 1, 0, 1, 250, @seed_time, @seed_time, 0),
(550, 500, 'catalog', '系统维护', 'el-icon-SetUp', '/setting/system', 'setting.system', '', '', '', '{}', 1, 0, 1, 200, @seed_time, @seed_time, 0),
(551, 550, 'page', '系统环境', '', '/setting/system/environment', 'setting.system.environment', 'setting/system/environment', '', '', '{}', 1, 0, 1, 300, @seed_time, @seed_time, 0),
(552, 550, 'page', '系统缓存', '', '/setting/system/cache', 'setting.system.cache', 'setting/system/cache', '', '', '{}', 1, 0, 1, 200, @seed_time, @seed_time, 0),
(553, 550, 'page', '系统日志', '', '/setting/system/journal', 'setting.system.journal', 'setting/system/journal', '', '', '{}', 1, 0, 1, 100, @seed_time, @seed_time, 0),
(600, 0, 'catalog', '开发工具', 'el-icon-EditPen', '/dev_tools', 'dev_tools', '', '', '', '{}', 1, 0, 1, 500, @seed_time, @seed_time, 0),
(515, 600, 'page', '字典管理', 'el-icon-Box', '/dev_tools/dict', 'dev_tools.dict', 'setting/dict/type/index', '', '', '{}', 1, 0, 1, 300, @seed_time, @seed_time, 0),
(620, 600, 'page', '模块中心', 'el-icon-Box', '/dev_tools/module', 'dev_tools.module', 'dev_tools/module/index', '', '', '{}', 1, 0, 1, 200, @seed_time, @seed_time, 0),
(610, 600, 'page', '代码生成器', 'el-icon-DocumentAdd', '/dev_tools/code', 'dev_tools.code', 'dev_tools/code/index', '', '', '{}', 1, 0, 1, 100, @seed_time, @seed_time, 0);

INSERT INTO `ma_permission` (`id`, `code`, `name`, `module`, `resource`, `action`, `status`, `sort`, `create_time`, `update_time`) VALUES
(1, 'dashboard:workbench:view', '工作台查看', 'dashboard', 'workbench', 'view', 1, 1000, @seed_time, @seed_time),
(2, 'system:admin:self', '管理员个人信息', 'system', 'admin', 'self', 1, 1000, @seed_time, @seed_time),
(3, 'system:admin:list', '管理员列表', 'system', 'admin', 'list', 1, 990, @seed_time, @seed_time),
(4, 'system:admin:detail', '管理员详情', 'system', 'admin', 'detail', 1, 980, @seed_time, @seed_time),
(5, 'system:admin:add', '管理员新增', 'system', 'admin', 'add', 1, 970, @seed_time, @seed_time),
(6, 'system:admin:edit', '管理员编辑', 'system', 'admin', 'edit', 1, 960, @seed_time, @seed_time),
(7, 'system:admin:upInfo', '管理员更新个人信息', 'system', 'admin', 'upInfo', 1, 950, @seed_time, @seed_time),
(8, 'system:admin:del', '管理员删除', 'system', 'admin', 'del', 1, 940, @seed_time, @seed_time),
(9, 'system:admin:disable', '管理员状态', 'system', 'admin', 'disable', 1, 930, @seed_time, @seed_time),
(10, 'system:role:all', '角色全集', 'system', 'role', 'all', 1, 1000, @seed_time, @seed_time),
(11, 'system:role:list', '角色列表', 'system', 'role', 'list', 1, 990, @seed_time, @seed_time),
(12, 'system:role:detail', '角色详情', 'system', 'role', 'detail', 1, 980, @seed_time, @seed_time),
(13, 'system:role:add', '角色新增', 'system', 'role', 'add', 1, 970, @seed_time, @seed_time),
(14, 'system:role:edit', '角色编辑', 'system', 'role', 'edit', 1, 960, @seed_time, @seed_time),
(15, 'system:role:del', '角色删除', 'system', 'role', 'del', 1, 950, @seed_time, @seed_time),
(16, 'system:menu:route', '菜单路由', 'system', 'menu', 'route', 1, 1000, @seed_time, @seed_time),
(17, 'system:menu:list', '菜单列表', 'system', 'menu', 'list', 1, 990, @seed_time, @seed_time),
(18, 'system:menu:detail', '菜单详情', 'system', 'menu', 'detail', 1, 980, @seed_time, @seed_time),
(19, 'system:menu:add', '菜单新增', 'system', 'menu', 'add', 1, 970, @seed_time, @seed_time),
(20, 'system:menu:edit', '菜单编辑', 'system', 'menu', 'edit', 1, 960, @seed_time, @seed_time),
(21, 'system:menu:del', '菜单删除', 'system', 'menu', 'del', 1, 950, @seed_time, @seed_time),
(22, 'system:dept:all', '部门全集', 'system', 'dept', 'all', 1, 1000, @seed_time, @seed_time),
(23, 'system:dept:list', '部门列表', 'system', 'dept', 'list', 1, 990, @seed_time, @seed_time),
(24, 'system:dept:detail', '部门详情', 'system', 'dept', 'detail', 1, 980, @seed_time, @seed_time),
(25, 'system:dept:add', '部门新增', 'system', 'dept', 'add', 1, 970, @seed_time, @seed_time),
(26, 'system:dept:edit', '部门编辑', 'system', 'dept', 'edit', 1, 960, @seed_time, @seed_time),
(27, 'system:dept:del', '部门删除', 'system', 'dept', 'del', 1, 950, @seed_time, @seed_time),
(28, 'system:post:all', '岗位全集', 'system', 'post', 'all', 1, 1000, @seed_time, @seed_time),
(29, 'system:post:list', '岗位列表', 'system', 'post', 'list', 1, 990, @seed_time, @seed_time),
(30, 'system:post:detail', '岗位详情', 'system', 'post', 'detail', 1, 980, @seed_time, @seed_time),
(31, 'system:post:add', '岗位新增', 'system', 'post', 'add', 1, 970, @seed_time, @seed_time),
(32, 'system:post:edit', '岗位编辑', 'system', 'post', 'edit', 1, 960, @seed_time, @seed_time),
(33, 'system:post:del', '岗位删除', 'system', 'post', 'del', 1, 950, @seed_time, @seed_time),
(34, 'material:file:list', '素材文件列表', 'material', 'file', 'list', 1, 1000, @seed_time, @seed_time),
(35, 'material:file:rename', '素材文件重命名', 'material', 'file', 'rename', 1, 990, @seed_time, @seed_time),
(36, 'material:file:move', '素材文件移动', 'material', 'file', 'move', 1, 980, @seed_time, @seed_time),
(37, 'material:file:delete', '素材文件删除', 'material', 'file', 'delete', 1, 970, @seed_time, @seed_time),
(38, 'material:category:list', '素材分类列表', 'material', 'category', 'list', 1, 960, @seed_time, @seed_time),
(39, 'material:category:add', '素材分类新增', 'material', 'category', 'add', 1, 950, @seed_time, @seed_time),
(40, 'material:category:rename', '素材分类重命名', 'material', 'category', 'rename', 1, 940, @seed_time, @seed_time),
(41, 'material:category:delete', '素材分类删除', 'material', 'category', 'delete', 1, 930, @seed_time, @seed_time),
(42, 'upload:image', '上传图片', 'upload', 'image', '', 1, 1000, @seed_time, @seed_time),
(43, 'upload:video', '上传视频', 'upload', 'video', '', 1, 990, @seed_time, @seed_time),
(44, 'setting:website:detail', '网站信息详情', 'setting', 'website', 'detail', 1, 1000, @seed_time, @seed_time),
(45, 'setting:website:save', '网站信息保存', 'setting', 'website', 'save', 1, 990, @seed_time, @seed_time),
(46, 'setting:copyright:detail', '备案详情', 'setting', 'copyright', 'detail', 1, 980, @seed_time, @seed_time),
(47, 'setting:copyright:save', '备案保存', 'setting', 'copyright', 'save', 1, 970, @seed_time, @seed_time),
(48, 'setting:protocol:detail', '协议详情', 'setting', 'protocol', 'detail', 1, 960, @seed_time, @seed_time),
(49, 'setting:protocol:save', '协议保存', 'setting', 'protocol', 'save', 1, 950, @seed_time, @seed_time),
(50, 'setting:storage:list', '存储列表', 'setting', 'storage', 'list', 1, 940, @seed_time, @seed_time),
(51, 'setting:storage:detail', '存储详情', 'setting', 'storage', 'detail', 1, 930, @seed_time, @seed_time),
(52, 'setting:storage:edit', '存储编辑', 'setting', 'storage', 'edit', 1, 920, @seed_time, @seed_time),
(53, 'setting:storage:change', '存储切换', 'setting', 'storage', 'change', 1, 910, @seed_time, @seed_time),
(54, 'setting:dict:type:all', '字典类型全集', 'setting', 'dict:type', 'all', 1, 1000, @seed_time, @seed_time),
(55, 'setting:dict:type:list', '字典类型列表', 'setting', 'dict:type', 'list', 1, 990, @seed_time, @seed_time),
(56, 'setting:dict:type:detail', '字典类型详情', 'setting', 'dict:type', 'detail', 1, 980, @seed_time, @seed_time),
(57, 'setting:dict:type:add', '字典类型新增', 'setting', 'dict:type', 'add', 1, 970, @seed_time, @seed_time),
(58, 'setting:dict:type:edit', '字典类型编辑', 'setting', 'dict:type', 'edit', 1, 960, @seed_time, @seed_time),
(59, 'setting:dict:type:del', '字典类型删除', 'setting', 'dict:type', 'del', 1, 950, @seed_time, @seed_time),
(60, 'setting:dict:data:all', '字典数据全集', 'setting', 'dict:data', 'all', 1, 940, @seed_time, @seed_time),
(61, 'setting:dict:data:list', '字典数据列表', 'setting', 'dict:data', 'list', 1, 930, @seed_time, @seed_time),
(62, 'setting:dict:data:detail', '字典数据详情', 'setting', 'dict:data', 'detail', 1, 920, @seed_time, @seed_time),
(63, 'setting:dict:data:add', '字典数据新增', 'setting', 'dict:data', 'add', 1, 910, @seed_time, @seed_time),
(64, 'setting:dict:data:edit', '字典数据编辑', 'setting', 'dict:data', 'edit', 1, 900, @seed_time, @seed_time),
(65, 'setting:dict:data:del', '字典数据删除', 'setting', 'dict:data', 'del', 1, 890, @seed_time, @seed_time),
(66, 'monitor:server', '系统环境', 'monitor', 'server', '', 1, 1000, @seed_time, @seed_time),
(67, 'monitor:cache', '系统缓存', 'monitor', 'cache', '', 1, 990, @seed_time, @seed_time),
(68, 'system:log:operate', '操作日志', 'system', 'log', 'operate', 1, 980, @seed_time, @seed_time),
(69, 'system:log:login', '登录日志', 'system', 'log', 'login', 1, 970, @seed_time, @seed_time),
(70, 'gen:db', '代码生成数据库表', 'gen', 'db', '', 1, 1000, @seed_time, @seed_time),
(71, 'gen:list', '代码生成列表', 'gen', 'list', '', 1, 990, @seed_time, @seed_time),
(72, 'gen:detail', '代码生成详情', 'gen', 'detail', '', 1, 980, @seed_time, @seed_time),
(73, 'gen:importTable', '导入数据表', 'gen', 'importTable', '', 1, 970, @seed_time, @seed_time),
(74, 'gen:syncTable', '同步表结构', 'gen', 'syncTable', '', 1, 960, @seed_time, @seed_time),
(75, 'gen:editTable', '编辑数据表', 'gen', 'editTable', '', 1, 950, @seed_time, @seed_time),
(76, 'gen:delTable', '删除数据表', 'gen', 'delTable', '', 1, 940, @seed_time, @seed_time),
(77, 'gen:previewCode', '预览代码', 'gen', 'previewCode', '', 1, 930, @seed_time, @seed_time),
(78, 'gen:genCode', '生成代码', 'gen', 'genCode', '', 1, 920, @seed_time, @seed_time),
(79, 'gen:downloadCode', '下载代码', 'gen', 'downloadCode', '', 1, 910, @seed_time, @seed_time),
(80, 'module:center:view', '模块中心查看', 'module', 'center', 'view', 1, 1000, @seed_time, @seed_time);

INSERT INTO `ma_menu_permission` (`menu_id`, `permission_id`, `create_time`) VALUES
(1, 1, @seed_time),
(101, 3, @seed_time),
(110, 11, @seed_time),
(120, 17, @seed_time),
(131, 23, @seed_time),
(140, 29, @seed_time),
(700, 34, @seed_time),
(502, 44, @seed_time),
(505, 46, @seed_time),
(510, 48, @seed_time),
(555, 50, @seed_time),
(551, 66, @seed_time),
(552, 67, @seed_time),
(553, 68, @seed_time),
(515, 55, @seed_time),
(620, 80, @seed_time),
(610, 71, @seed_time);

INSERT INTO `ma_role_permission` (`tenant_id`, `role_id`, `permission_id`, `create_time`)
SELECT 0, 1, `id`, @seed_time FROM `ma_permission`;

INSERT INTO `ma_setting` (`id`, `tenant_id`, `setting_group`, `setting_key`, `setting_value`, `value_type`, `is_public`, `remark`, `create_time`, `update_time`) VALUES
(1, 0, 'website', 'name', 'go-makeadmin', 'string', 1, 'Website name.', @seed_time, @seed_time),
(2, 0, 'website', 'logo', '/api/static/backend_logo.png', 'string', 1, 'Website logo.', @seed_time, @seed_time),
(3, 0, 'website', 'favicon', '/api/static/backend_favicon.ico', 'string', 1, 'Website favicon.', @seed_time, @seed_time),
(4, 0, 'website', 'backdrop', '/api/static/backend_backdrop.png', 'string', 1, 'Login backdrop.', @seed_time, @seed_time),
(5, 0, 'website', 'copyright', '[{"name":"go-makeadmin","link":"https://www.go-makeadmin.cn"}]', 'json', 1, 'Copyright links.', @seed_time, @seed_time),
(6, 0, 'storage', 'default', 'local', 'string', 0, 'Default storage driver.', @seed_time, @seed_time),
(7, 0, 'storage', 'local', '{"name":"本地存储"}', 'json', 0, 'Local storage.', @seed_time, @seed_time),
(8, 0, 'storage', 'qiniu', '{"name":"七牛云存储","bucket":"","secretKey":"","accessKey":"","domain":""}', 'json', 0, 'Qiniu storage.', @seed_time, @seed_time),
(9, 0, 'storage', 'aliyun', '{"name":"阿里云存储","bucket":"","secretKey":"","accessKey":"","domain":""}', 'json', 0, 'Aliyun OSS storage.', @seed_time, @seed_time),
(10, 0, 'storage', 'qcloud', '{"name":"腾讯云存储","bucket":"","secretKey":"","accessKey":"","domain":"","region":""}', 'json', 0, 'Tencent COS storage.', @seed_time, @seed_time),
(11, 0, 'protocol', 'service', '{"name":"","content":""}', 'json', 1, 'Service protocol.', @seed_time, @seed_time),
(12, 0, 'protocol', 'privacy', '{"name":"","content":""}', 'json', 1, 'Privacy protocol.', @seed_time, @seed_time);

INSERT INTO `ma_dict_type` (`id`, `code`, `name`, `remark`, `status`, `sort`, `create_time`, `update_time`, `delete_time`) VALUES
(1, 'system_status', '系统状态', 'Framework enabled/disabled status.', 1, 1000, @seed_time, @seed_time, 0),
(2, 'menu_type', '菜单类型', 'Framework menu type.', 1, 990, @seed_time, @seed_time, 0),
(3, 'storage_type', '存储类型', 'Storage driver type.', 1, 980, @seed_time, @seed_time, 0),
(4, 'data_scope_type', '数据范围类型', 'Role data scope type.', 1, 970, @seed_time, @seed_time, 0),
(5, 'common_status', '通用状态', 'Numeric enabled/disabled status for generated CRUD modules.', 1, 960, @seed_time, @seed_time, 0);

INSERT INTO `ma_dict_item` (`id`, `type_id`, `item_label`, `item_value`, `remark`, `status`, `sort`, `create_time`, `update_time`, `delete_time`) VALUES
(1, 1, '启用', 'enabled', '', 1, 1000, @seed_time, @seed_time, 0),
(2, 1, '禁用', 'disabled', '', 1, 990, @seed_time, @seed_time, 0),
(3, 2, '目录', 'catalog', '', 1, 1000, @seed_time, @seed_time, 0),
(4, 2, '页面', 'page', '', 1, 990, @seed_time, @seed_time, 0),
(5, 2, '操作', 'action', '', 1, 980, @seed_time, @seed_time, 0),
(6, 3, '本地', 'local', '', 1, 1000, @seed_time, @seed_time, 0),
(7, 3, '七牛云', 'qiniu', '', 1, 990, @seed_time, @seed_time, 0),
(8, 3, '阿里云', 'aliyun', '', 1, 980, @seed_time, @seed_time, 0),
(9, 3, '腾讯云', 'qcloud', '', 1, 970, @seed_time, @seed_time, 0),
(10, 4, '全部数据', 'all', '', 1, 1000, @seed_time, @seed_time, 0),
(11, 4, '本人数据', 'self', '', 1, 990, @seed_time, @seed_time, 0),
(12, 4, '本组织', 'org', '', 1, 980, @seed_time, @seed_time, 0),
(13, 4, '本组织及下级', 'org_tree', '', 1, 970, @seed_time, @seed_time, 0),
(14, 4, '自定义组织', 'custom_org', '', 1, 960, @seed_time, @seed_time, 0),
(15, 5, '启用', '1', '', 1, 1000, @seed_time, @seed_time, 0),
(16, 5, '禁用', '0', '', 1, 990, @seed_time, @seed_time, 0);

INSERT INTO `ma_file_category` (`id`, `tenant_id`, `parent_id`, `code`, `name`, `file_type`, `status`, `sort`, `create_time`, `update_time`, `delete_time`) VALUES
(1, 0, 0, 'image', '图片', 'image', 1, 1000, @seed_time, @seed_time, 0),
(2, 0, 0, 'video', '视频', 'video', 1, 990, @seed_time, @seed_time, 0);
