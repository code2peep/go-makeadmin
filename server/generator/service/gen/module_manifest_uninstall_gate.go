package gen

import (
	"fmt"
	"os"
	"strings"

	"go-makeadmin/core/response"
	"go-makeadmin/generator/schemas/req"
	"go-makeadmin/generator/schemas/resp"
	"go-makeadmin/model/makeadmin"
	"gorm.io/gorm"
)

const moduleManifestUninstallApplyEnv = "MAKEADMIN_ALLOW_MODULE_UNINSTALL_APPLY"

// ApplyModuleManifestUninstall validates and applies a local module uninstall request.
func (genSrv generateService) ApplyModuleManifestUninstall(applyReq req.ModuleManifestUninstallApplyReq) (res resp.ModuleManifestUninstallApplyResp, e error) {
	previewReq := moduleManifestPreviewReqFromUninstallApply(applyReq)
	manifest, source, err := loadModuleManifest(previewReq)
	if err != nil {
		return res, err
	}
	if err = validateModuleManifest(manifest); err != nil {
		return res, err
	}

	res = resp.ModuleManifestUninstallApplyResp{
		Source: source,
		Manifest: resp.ModuleManifestSummaryResp{
			Module:         manifest.Module,
			Entity:         manifest.Entity,
			Table:          manifest.Table,
			MenuName:       manifest.Menu.Name,
			RequiresSchema: manifest.RequiresSchema,
		},
		Status:      "blocked",
		RequiredEnv: moduleManifestUninstallApplyEnv,
		Plan:        buildModuleManifestPlan(manifest, previewReq),
		Summary:     buildModuleManifestApplySummary(manifest, "uninstall"),
	}

	if os.Getenv(moduleManifestUninstallApplyEnv) != "1" {
		res.Checks = append(res.Checks, moduleInstallCheck("environment", "failed", moduleManifestUninstallApplyEnv+"=1 is required"))
		return res, moduleUninstallGateError(res, moduleManifestUninstallApplyEnv+"=1 is required; no database access was attempted")
	}
	res.Checks = append(res.Checks, moduleInstallCheck("environment", "passed", moduleManifestUninstallApplyEnv+"=1 is present"))

	confirmModule := strings.TrimSpace(applyReq.ConfirmModule)
	if confirmModule != manifest.Module {
		res.Checks = append(res.Checks, moduleInstallCheck("confirmModule", "failed", fmt.Sprintf("confirmModule must be %q", manifest.Module)))
		return res, moduleUninstallGateError(res, fmt.Sprintf("confirmModule must be %q; no database access was attempted", manifest.Module))
	}
	res.Checks = append(res.Checks, moduleInstallCheck("confirmModule", "passed", "module name confirmed"))

	if !applyReq.ConfirmDelete {
		res.Checks = append(res.Checks, moduleInstallCheck("confirmDelete", "failed", "confirmDelete must be true"))
		return res, moduleUninstallGateError(res, "confirmDelete must be true; no database access was attempted")
	}
	res.Checks = append(res.Checks, moduleInstallCheck("confirmDelete", "passed", "delete intent confirmed"))

	if err = validateModuleInstallDatabaseTarget(); err != nil {
		res.Checks = append(res.Checks, moduleInstallCheck("databaseTarget", "failed", err.Error()))
		return res, moduleUninstallGateError(res, err.Error()+"; no database access was attempted")
	}
	res.Checks = append(res.Checks, moduleInstallCheck("databaseTarget", "passed", "local go_makeadmin database confirmed"))

	if genSrv.db == nil {
		res.Checks = append(res.Checks, moduleInstallCheck("database", "failed", "module uninstall apply requires configured database"))
		return res, moduleUninstallGateError(res, "module uninstall apply requires configured database; no database access was attempted")
	}

	if err = genSrv.db.Transaction(func(tx *gorm.DB) error {
		before, snapErr := moduleUninstallSnapshot(tx, manifest)
		if snapErr != nil {
			return snapErr
		}
		res.Before = before
		if applyErr := applyModuleUninstall(tx, manifest); applyErr != nil {
			return applyErr
		}
		after, snapErr := moduleUninstallSnapshot(tx, manifest)
		if snapErr != nil {
			return snapErr
		}
		res.After = after
		return nil
	}); err != nil {
		res.Checks = append(res.Checks, moduleInstallCheck("executor", "failed", "module uninstall apply transaction failed"))
		return res, response.CheckErr(err, "ModuleManifestUninstallApply Transaction err")
	}
	res.Checks = append(res.Checks, moduleInstallCheck("executor", "passed", "module uninstall apply transaction completed"))
	res.Status = "applied"
	res.Message = "module uninstall apply completed"
	return res, nil
}

func moduleManifestPreviewReqFromUninstallApply(applyReq req.ModuleManifestUninstallApplyReq) req.ModuleManifestPreviewReq {
	return req.ModuleManifestPreviewReq{
		ManifestPath: applyReq.ManifestPath,
		ManifestBody: applyReq.ManifestBody,
	}
}

func moduleUninstallGateError(res resp.ModuleManifestUninstallApplyResp, message string) error {
	res.Message = message
	return response.AssertArgumentError.Make(message).MakeData(res)
}

func applyModuleUninstall(tx *gorm.DB, manifest moduleManifest) error {
	codes := moduleManifestPermissionCodes(manifest)
	if err := tx.Exec(
		"DELETE rp FROM ma_role_permission AS rp INNER JOIN ma_permission AS p ON p.id = rp.permission_id WHERE p.code IN ?",
		codes,
	).Error; err != nil {
		return err
	}
	if err := tx.Exec(
		"DELETE mp FROM ma_menu_permission AS mp LEFT JOIN ma_menu AS m ON m.id = mp.menu_id LEFT JOIN ma_permission AS p ON p.id = mp.permission_id WHERE m.route_name = ? OR p.code IN ?",
		manifest.Menu.RouteName,
		codes,
	).Error; err != nil {
		return err
	}
	if err := tx.Exec("DELETE FROM ma_menu WHERE route_name = ?", manifest.Menu.RouteName).Error; err != nil {
		return err
	}
	return tx.Exec("DELETE FROM ma_permission WHERE code IN ?", codes).Error
}

func moduleUninstallSnapshot(tx *gorm.DB, manifest moduleManifest) (resp.ModuleManifestInstallSnapshotResp, error) {
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
		Where("p.code IN ?", codes).
		Count(&snapshot.RolePermissions).Error
	return snapshot, err
}
