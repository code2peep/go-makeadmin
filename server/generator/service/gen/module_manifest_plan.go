package gen

import (
	"fmt"
	"strings"

	"go-makeadmin/generator/schemas/req"
	"go-makeadmin/generator/schemas/resp"
)

func buildModuleManifestPlan(manifest moduleManifest, previewReq req.ModuleManifestPreviewReq) resp.ModuleManifestPlanResp {
	tenantID := manifestTenantID(previewReq)
	roleID := manifestRoleID(previewReq)
	registrySQL := buildModuleRegistrySQL(manifest)
	roleGrantSQL := buildModuleRoleGrantSQL(manifest, tenantID, roleID)
	return resp.ModuleManifestPlanResp{
		TenantID:     tenantID,
		RoleID:       roleID,
		RegistrySQL:  registrySQL,
		RoleGrantSQL: roleGrantSQL,
		InstallSQL:   strings.Join([]string{registrySQL, roleGrantSQL}, "\n\n"),
		UninstallSQL: buildModuleUninstallSQL(manifest),
		RuntimeHint:  moduleRuntimeHint(manifest),
	}
}

func buildModuleRegistrySQL(manifest moduleManifest) string {
	statements := []string{
		"SET @now = UNIX_TIMESTAMP();",
		fmt.Sprintf("SET @parent_route_name = %s;", sqlQuote(manifest.Menu.Parent)),
		strings.Join([]string{
			"SET @parent_menu_id = COALESCE((",
			"    SELECT id FROM `ma_menu`",
			"    WHERE route_name = @parent_route_name AND delete_time = 0",
			"    LIMIT 1",
			"), 0);",
		}, "\n"),
	}
	for index, permission := range manifest.Permissions {
		statements = append(statements, modulePermissionInsertSQL(permission, 1000-index*10))
	}
	statements = append(
		statements,
		moduleMenuInsertSQL(manifest.Menu),
		strings.Join([]string{
			"SET @module_menu_id = COALESCE((",
			"    SELECT id FROM `ma_menu`",
			fmt.Sprintf("    WHERE route_name = %s AND delete_time = 0", sqlQuote(manifest.Menu.RouteName)),
			"    LIMIT 1",
			"), 0);",
		}, "\n"),
		strings.Join([]string{
			"SET @module_permission_id = COALESCE((",
			"    SELECT id FROM `ma_permission`",
			fmt.Sprintf("    WHERE code = %s", sqlQuote(manifest.Menu.Permission)),
			"    LIMIT 1",
			"), 0);",
		}, "\n"),
		moduleMenuPermissionInsertSQL(),
	)
	return strings.Join(statements, "\n\n")
}

func modulePermissionInsertSQL(permission moduleManifestPermission, sort int) string {
	return strings.Join([]string{
		"INSERT INTO `ma_permission`",
		"(`code`, `name`, `module`, `resource`, `action`, `status`, `sort`, `create_time`, `update_time`)",
		"SELECT",
		fmt.Sprintf("%s, %s, %s, %s, %s, 1, %d, @now, @now",
			sqlQuote(permission.Code),
			sqlQuote(permission.Name),
			sqlQuote(permission.Module),
			sqlQuote(permission.Resource),
			sqlQuote(permission.Action),
			sort,
		),
		"WHERE NOT EXISTS (",
		fmt.Sprintf("    SELECT 1 FROM `ma_permission` WHERE `code` = %s", sqlQuote(permission.Code)),
		");",
	}, "\n")
}

func moduleMenuInsertSQL(menu moduleManifestMenu) string {
	visible := 0
	if menu.Visible {
		visible = 1
	}
	return strings.Join([]string{
		"INSERT INTO `ma_menu`",
		"(`parent_id`, `menu_type`, `name`, `icon`, `route_path`, `route_name`, `component`, `redirect`,",
		"`active_path`, `meta`, `is_visible`, `is_cache`, `status`, `sort`, `create_time`, `update_time`, `delete_time`)",
		"SELECT",
		fmt.Sprintf("@parent_menu_id, %s, %s, '', %s, %s, %s, '', '', '{}', %d, 1, 1, %d, @now, @now, 0",
			sqlQuote(menu.Type),
			sqlQuote(menu.Name),
			sqlQuote(menu.RoutePath),
			sqlQuote(menu.RouteName),
			sqlQuote(menu.Component),
			visible,
			menu.Sort,
		),
		"WHERE NOT EXISTS (",
		fmt.Sprintf("    SELECT 1 FROM `ma_menu` WHERE `route_name` = %s AND `delete_time` = 0", sqlQuote(menu.RouteName)),
		");",
	}, "\n")
}

func moduleMenuPermissionInsertSQL() string {
	return strings.Join([]string{
		"INSERT INTO `ma_menu_permission`",
		"(`menu_id`, `permission_id`, `create_time`)",
		"SELECT @module_menu_id, @module_permission_id, @now",
		"WHERE @module_menu_id > 0",
		"  AND @module_permission_id > 0",
		"  AND NOT EXISTS (",
		"      SELECT 1 FROM `ma_menu_permission`",
		"      WHERE `menu_id` = @module_menu_id AND `permission_id` = @module_permission_id",
		"  );",
	}, "\n")
}

func buildModuleRoleGrantSQL(manifest moduleManifest, tenantID uint64, roleID uint64) string {
	codeList := modulePermissionCodeList(manifest)
	return strings.Join([]string{
		"SET @now = UNIX_TIMESTAMP();",
		fmt.Sprintf("SET @tenant_id = %d;", tenantID),
		fmt.Sprintf("SET @role_id = %d;", roleID),
		strings.Join([]string{
			"INSERT INTO `ma_role_permission`",
			"(`tenant_id`, `role_id`, `permission_id`, `create_time`)",
			"SELECT @tenant_id, @role_id, p.id, @now",
			"FROM `ma_permission` AS p",
			fmt.Sprintf("WHERE p.code IN (%s)", codeList),
			"  AND p.status = 1",
			"  AND EXISTS (",
			"      SELECT 1 FROM `ma_role` AS r",
			"      WHERE r.tenant_id = @tenant_id",
			"        AND r.id = @role_id",
			"        AND r.status = 1",
			"        AND r.delete_time = 0",
			"  )",
			"  AND NOT EXISTS (",
			"      SELECT 1 FROM `ma_role_permission` AS rp",
			"      WHERE rp.tenant_id = @tenant_id",
			"        AND rp.role_id = @role_id",
			"        AND rp.permission_id = p.id",
			"  );",
		}, "\n"),
	}, "\n\n")
}

func buildModuleUninstallSQL(manifest moduleManifest) string {
	codeList := modulePermissionCodeList(manifest)
	routeName := sqlQuote(manifest.Menu.RouteName)
	return strings.Join([]string{
		"SET @module_route_name = " + routeName + ";",
		strings.Join([]string{
			"DELETE rp FROM `ma_role_permission` AS rp",
			"INNER JOIN `ma_permission` AS p ON p.id = rp.permission_id",
			fmt.Sprintf("WHERE p.code IN (%s);", codeList),
		}, "\n"),
		strings.Join([]string{
			"DELETE mp FROM `ma_menu_permission` AS mp",
			"LEFT JOIN `ma_menu` AS m ON m.id = mp.menu_id",
			"LEFT JOIN `ma_permission` AS p ON p.id = mp.permission_id",
			fmt.Sprintf("WHERE m.route_name = @module_route_name OR p.code IN (%s);", codeList),
		}, "\n"),
		strings.Join([]string{
			"DELETE FROM `ma_menu`",
			"WHERE route_name = @module_route_name;",
		}, "\n"),
		strings.Join([]string{
			"DELETE FROM `ma_permission`",
			fmt.Sprintf("WHERE code IN (%s);", codeList),
		}, "\n"),
	}, "\n\n")
}

func modulePermissionCodeList(manifest moduleManifest) string {
	codes := make([]string, 0, len(manifest.Permissions))
	for _, permission := range manifest.Permissions {
		codes = append(codes, sqlQuote(permission.Code))
	}
	return strings.Join(codes, ", ")
}

func moduleRuntimeHint(manifest moduleManifest) string {
	if manifest.Module == "article" {
		return "MAKEADMIN_ENABLE_DEMO_MODULE=1"
	}
	return "No runtime env gate is defined for this module yet."
}

func sqlQuote(value any) string {
	text := fmt.Sprint(value)
	return "'" + strings.ReplaceAll(strings.ReplaceAll(text, `\`, `\\`), "'", "''") + "'"
}
