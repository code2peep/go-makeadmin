package gen

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"go-makeadmin/config"
	"go-makeadmin/core/response"
	"go-makeadmin/generator/schemas/req"
	"go-makeadmin/generator/schemas/resp"
	"go-makeadmin/model/makeadmin"
	"gorm.io/gorm"
)

const moduleManifestInstallApplyEnv = "MAKEADMIN_ALLOW_MODULE_INSTALL_APPLY"

// ApplyModuleManifestInstall validates the write gate for a module install request.
// P3.10 deliberately stops before any database access or SQL execution.
func (genSrv generateService) ApplyModuleManifestInstall(applyReq req.ModuleManifestInstallApplyReq) (res resp.ModuleManifestInstallApplyResp, e error) {
	previewReq := moduleManifestPreviewReqFromInstallApply(applyReq)
	manifest, source, err := loadModuleManifest(previewReq)
	if err != nil {
		return res, err
	}
	if err = validateModuleManifest(manifest); err != nil {
		return res, err
	}

	tenantID := moduleManifestInstallTenantID(applyReq)
	roleID := moduleManifestInstallRoleID(applyReq)
	previewReq.TenantID = tenantID
	previewReq.RoleID = roleID
	res = resp.ModuleManifestInstallApplyResp{
		Source: source,
		Manifest: resp.ModuleManifestSummaryResp{
			Module:         manifest.Module,
			Entity:         manifest.Entity,
			Table:          manifest.Table,
			MenuName:       manifest.Menu.Name,
			RequiresSchema: manifest.RequiresSchema,
		},
		TenantID:    tenantID,
		RoleID:      roleID,
		Status:      "blocked",
		RequiredEnv: moduleManifestInstallApplyEnv,
		Plan:        buildModuleManifestPlan(manifest, previewReq),
		Summary:     buildModuleManifestApplySummary(manifest, "install"),
	}

	if os.Getenv(moduleManifestInstallApplyEnv) != "1" {
		res.Checks = append(res.Checks, moduleInstallCheck("environment", "failed", moduleManifestInstallApplyEnv+"=1 is required"))
		return res, moduleInstallGateError(res, moduleManifestInstallApplyEnv+"=1 is required; no database access was attempted")
	}
	res.Checks = append(res.Checks, moduleInstallCheck("environment", "passed", moduleManifestInstallApplyEnv+"=1 is present"))

	confirmModule := strings.TrimSpace(applyReq.ConfirmModule)
	if confirmModule != manifest.Module {
		res.Checks = append(res.Checks, moduleInstallCheck("confirmModule", "failed", fmt.Sprintf("confirmModule must be %q", manifest.Module)))
		return res, moduleInstallGateError(res, fmt.Sprintf("confirmModule must be %q; no database access was attempted", manifest.Module))
	}
	res.Checks = append(res.Checks, moduleInstallCheck("confirmModule", "passed", "module name confirmed"))

	if applyReq.ConfirmTenantID == nil || *applyReq.ConfirmTenantID != tenantID {
		res.Checks = append(res.Checks, moduleInstallCheck("confirmTenantId", "failed", fmt.Sprintf("confirmTenantId must be %d", tenantID)))
		return res, moduleInstallGateError(res, fmt.Sprintf("confirmTenantId must be %d; no database access was attempted", tenantID))
	}
	res.Checks = append(res.Checks, moduleInstallCheck("confirmTenantId", "passed", "tenant id confirmed"))

	if applyReq.ConfirmRoleID == nil || *applyReq.ConfirmRoleID != roleID {
		res.Checks = append(res.Checks, moduleInstallCheck("confirmRoleId", "failed", fmt.Sprintf("confirmRoleId must be %d", roleID)))
		return res, moduleInstallGateError(res, fmt.Sprintf("confirmRoleId must be %d; no database access was attempted", roleID))
	}
	res.Checks = append(res.Checks, moduleInstallCheck("confirmRoleId", "passed", "role id confirmed"))

	if !applyReq.ConfirmInstall {
		res.Checks = append(res.Checks, moduleInstallCheck("confirmInstall", "failed", "confirmInstall must be true"))
		return res, moduleInstallGateError(res, "confirmInstall must be true; no database access was attempted")
	}
	res.Checks = append(res.Checks, moduleInstallCheck("confirmInstall", "passed", "install intent confirmed"))

	if manifest.RequiresSchema && !applyReq.ConfirmSchemaRisk {
		res.Checks = append(res.Checks, moduleInstallCheck("confirmSchemaRisk", "failed", "confirmSchemaRisk must be true because manifest.requiresSchema is true"))
		return res, moduleInstallGateError(res, "confirmSchemaRisk must be true because manifest.requiresSchema is true; no database access was attempted")
	}
	if manifest.RequiresSchema {
		res.Checks = append(res.Checks, moduleInstallCheck("confirmSchemaRisk", "passed", "schema risk confirmed"))
	} else {
		res.Checks = append(res.Checks, moduleInstallCheck("confirmSchemaRisk", "skipped", "manifest.requiresSchema is false"))
	}

	if err = validateModuleInstallDatabaseTarget(); err != nil {
		res.Checks = append(res.Checks, moduleInstallCheck("databaseTarget", "failed", err.Error()))
		return res, moduleInstallGateError(res, err.Error()+"; no database access was attempted")
	}
	res.Checks = append(res.Checks, moduleInstallCheck("databaseTarget", "passed", "local go_makeadmin database confirmed"))

	if genSrv.db == nil {
		res.Checks = append(res.Checks, moduleInstallCheck("database", "failed", "module install apply requires configured database"))
		return res, moduleInstallGateError(res, "module install apply requires configured database; no database access was attempted")
	}

	if err = genSrv.db.Transaction(func(tx *gorm.DB) error {
		before, snapErr := moduleInstallSnapshot(tx, manifest, tenantID, roleID)
		if snapErr != nil {
			return snapErr
		}
		res.Before = before
		if applyErr := applyModuleInstall(tx, manifest, tenantID, roleID); applyErr != nil {
			return applyErr
		}
		after, snapErr := moduleInstallSnapshot(tx, manifest, tenantID, roleID)
		if snapErr != nil {
			return snapErr
		}
		res.After = after
		return nil
	}); err != nil {
		res.Checks = append(res.Checks, moduleInstallCheck("executor", "failed", "module install apply transaction failed"))
		return res, response.CheckErr(err, "ModuleManifestInstallApply Transaction err")
	}
	res.Checks = append(res.Checks, moduleInstallCheck("executor", "passed", "module install apply transaction completed"))
	res.Status = "applied"
	res.Message = "module install apply completed"
	return res, nil
}

func moduleManifestPreviewReqFromInstallApply(applyReq req.ModuleManifestInstallApplyReq) req.ModuleManifestPreviewReq {
	return req.ModuleManifestPreviewReq{
		ManifestPath: applyReq.ManifestPath,
		ManifestBody: applyReq.ManifestBody,
		TenantID:     moduleManifestInstallTenantID(applyReq),
		RoleID:       moduleManifestInstallRoleID(applyReq),
	}
}

func moduleManifestInstallTenantID(applyReq req.ModuleManifestInstallApplyReq) uint64 {
	if applyReq.TenantID > 0 {
		return applyReq.TenantID
	}
	return makeadmin.GlobalTenantID
}

func moduleManifestInstallRoleID(applyReq req.ModuleManifestInstallApplyReq) uint64 {
	if applyReq.RoleID > 0 {
		return applyReq.RoleID
	}
	return 1
}

func moduleInstallCheck(name string, status string, message string) resp.ModuleManifestInstallCheckResp {
	return resp.ModuleManifestInstallCheckResp{Name: name, Status: status, Message: message}
}

func moduleInstallGateError(res resp.ModuleManifestInstallApplyResp, message string) error {
	res.Message = message
	return response.AssertArgumentError.Make(message).MakeData(res)
}

func validateModuleInstallDatabaseTarget() error {
	dsn := config.Config.DatabaseUrl
	isLocalHost := strings.Contains(dsn, "@tcp(localhost:") || strings.Contains(dsn, "@tcp(127.0.0.1:")
	isLocalDatabase := strings.Contains(dsn, ")/go_makeadmin?") || strings.HasSuffix(dsn, ")/go_makeadmin")
	if !isLocalHost || !isLocalDatabase {
		return fmt.Errorf("module install apply requires local go_makeadmin database")
	}
	return nil
}

func applyModuleInstall(tx *gorm.DB, manifest moduleManifest, tenantID uint64, roleID uint64) error {
	now := time.Now().Unix()
	for index, permission := range manifest.Permissions {
		if err := ensureModulePermission(tx, permission, uint16(1000-index*10), now); err != nil {
			return err
		}
	}
	menuID, err := ensureModuleMenu(tx, manifest.Menu, now)
	if err != nil {
		return err
	}
	permissionID, err := findPermissionID(tx, manifest.Menu.Permission)
	if err != nil {
		return err
	}
	if menuID > 0 && permissionID > 0 {
		if err := ensureModuleMenuPermission(tx, menuID, permissionID, now); err != nil {
			return err
		}
	}
	return ensureModuleRolePermissions(tx, manifest, tenantID, roleID, now)
}

func ensureModulePermission(tx *gorm.DB, permission moduleManifestPermission, sort uint16, now int64) error {
	var count int64
	if err := tx.Model(&makeadmin.Permission{}).Where("code = ?", permission.Code).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return tx.Create(&makeadmin.Permission{
		Code:       permission.Code,
		Name:       permission.Name,
		Module:     permission.Module,
		Resource:   permission.Resource,
		Action:     permission.Action,
		Status:     makeadmin.StatusEnabled,
		Sort:       sort,
		CreateTime: now,
		UpdateTime: now,
	}).Error
}

func ensureModuleMenu(tx *gorm.DB, menu moduleManifestMenu, now int64) (uint64, error) {
	menuID, err := findMenuID(tx, menu.RouteName)
	if err != nil || menuID > 0 {
		return menuID, err
	}
	parentID, err := findMenuID(tx, menu.Parent)
	if err != nil {
		return 0, err
	}
	visible := uint8(0)
	if menu.Visible {
		visible = 1
	}
	next := makeadmin.Menu{
		ParentID:   parentID,
		MenuType:   menu.Type,
		Name:       menu.Name,
		RoutePath:  menu.RoutePath,
		RouteName:  menu.RouteName,
		Component:  menu.Component,
		Meta:       "{}",
		IsVisible:  visible,
		IsCache:    1,
		Status:     makeadmin.StatusEnabled,
		Sort:       uint16(menu.Sort),
		CreateTime: now,
		UpdateTime: now,
	}
	if err := tx.Create(&next).Error; err != nil {
		return 0, err
	}
	return next.ID, nil
}

func ensureModuleMenuPermission(tx *gorm.DB, menuID uint64, permissionID uint64, now int64) error {
	var count int64
	err := tx.Model(&makeadmin.MenuPermission{}).
		Where("menu_id = ? AND permission_id = ?", menuID, permissionID).
		Count(&count).Error
	if err != nil || count > 0 {
		return err
	}
	return tx.Create(&makeadmin.MenuPermission{
		MenuID:       menuID,
		PermissionID: permissionID,
		CreateTime:   now,
	}).Error
}

func ensureModuleRolePermissions(tx *gorm.DB, manifest moduleManifest, tenantID uint64, roleID uint64, now int64) error {
	roleExists, err := moduleRoleExists(tx, tenantID, roleID)
	if err != nil || !roleExists {
		return err
	}
	for _, permission := range manifest.Permissions {
		permissionID, findErr := findPermissionID(tx, permission.Code)
		if findErr != nil {
			return findErr
		}
		if permissionID == 0 {
			continue
		}
		var count int64
		err = tx.Model(&makeadmin.RolePermission{}).
			Where("tenant_id = ? AND role_id = ? AND permission_id = ?", tenantID, roleID, permissionID).
			Count(&count).Error
		if err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		err = tx.Create(&makeadmin.RolePermission{
			TenantID:     tenantID,
			RoleID:       roleID,
			PermissionID: permissionID,
			CreateTime:   now,
		}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func moduleRoleExists(tx *gorm.DB, tenantID uint64, roleID uint64) (bool, error) {
	var role makeadmin.Role
	err := tx.Select("id").
		Where("tenant_id = ? AND id = ? AND status = ? AND delete_time = 0", tenantID, roleID, makeadmin.StatusEnabled).
		Take(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return err == nil, err
}

func findMenuID(tx *gorm.DB, routeName string) (uint64, error) {
	var menu makeadmin.Menu
	err := tx.Select("id").Where("route_name = ? AND delete_time = 0", routeName).Take(&menu).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return menu.ID, err
}

func findPermissionID(tx *gorm.DB, code string) (uint64, error) {
	var permission makeadmin.Permission
	err := tx.Select("id").Where("code = ?", code).Take(&permission).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return permission.ID, err
}

func moduleInstallSnapshot(tx *gorm.DB, manifest moduleManifest, tenantID uint64, roleID uint64) (resp.ModuleManifestInstallSnapshotResp, error) {
	var snapshot resp.ModuleManifestInstallSnapshotResp
	codes := moduleManifestPermissionCodes(manifest)
	if err := tx.Model(&makeadmin.Permission{}).Where("code IN ?", codes).Count(&snapshot.Permissions).Error; err != nil {
		return snapshot, err
	}
	if err := tx.Model(&makeadmin.Menu{}).Where("route_name = ? AND delete_time = 0", manifest.Menu.RouteName).Count(&snapshot.Menus).Error; err != nil {
		return snapshot, err
	}
	err := tx.Table("ma_menu_permission AS mp").
		Joins("LEFT JOIN ma_menu AS m ON m.id = mp.menu_id").
		Joins("LEFT JOIN ma_permission AS p ON p.id = mp.permission_id").
		Where("m.route_name = ? OR p.code IN ?", manifest.Menu.RouteName, codes).
		Count(&snapshot.MenuPermissions).Error
	if err != nil {
		return snapshot, err
	}
	err = tx.Table("ma_role_permission AS rp").
		Joins("INNER JOIN ma_permission AS p ON p.id = rp.permission_id").
		Where("rp.tenant_id = ? AND rp.role_id = ? AND p.code IN ?", tenantID, roleID, codes).
		Count(&snapshot.RolePermissions).Error
	return snapshot, err
}

func moduleManifestPermissionCodes(manifest moduleManifest) []string {
	codes := make([]string, 0, len(manifest.Permissions))
	for _, permission := range manifest.Permissions {
		codes = append(codes, permission.Code)
	}
	return codes
}
