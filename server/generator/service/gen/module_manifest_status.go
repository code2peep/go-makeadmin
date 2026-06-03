package gen

import (
	"os"
	"strings"

	"go-makeadmin/core/response"
	"go-makeadmin/generator/schemas/req"
	"go-makeadmin/generator/schemas/resp"
	"go-makeadmin/model/makeadmin"
	"gorm.io/gorm"
)

// ReadModuleManifestInstallStatus reads the current local install state for a module manifest.
func (genSrv generateService) ReadModuleManifestInstallStatus(previewReq req.ModuleManifestPreviewReq) (res resp.ModuleManifestInstallStatusResp, e error) {
	manifest, source, err := loadModuleManifest(previewReq)
	if err != nil {
		return res, err
	}
	if err = validateModuleManifest(manifest); err != nil {
		return res, err
	}

	tenantID := manifestTenantID(previewReq)
	roleID := manifestRoleID(previewReq)
	runtimeHint := moduleRuntimeHint(manifest)
	runtimeEnv := moduleRuntimeEnv(runtimeHint)
	res = resp.ModuleManifestInstallStatusResp{
		Source: source,
		Manifest: resp.ModuleManifestSummaryResp{
			Module:         manifest.Module,
			Entity:         manifest.Entity,
			Table:          manifest.Table,
			MenuName:       manifest.Menu.Name,
			RequiresSchema: manifest.RequiresSchema,
		},
		TenantID:          tenantID,
		RoleID:            roleID,
		RuntimeHint:       runtimeHint,
		RuntimeEnv:        runtimeEnv,
		RuntimeRegistered: manifest.RuntimeRegistered,
		RuntimeEnabled:    runtimeEnv == "" || os.Getenv(runtimeEnv) == "1",
		Summary:           buildModuleManifestApplySummary(manifest, "status"),
		Expected:          moduleInstallExpectedSnapshot(manifest),
		Checks: []resp.ModuleManifestInstallCheckResp{
			moduleInstallCheck("manifest", "passed", "module manifest loaded"),
		},
	}

	if err = validateModuleInstallDatabaseTarget(); err != nil {
		res.Status = "blocked"
		res.Message = err.Error()
		res.Checks = append(res.Checks, moduleInstallCheck("databaseTarget", "failed", err.Error()))
		return res, response.AssertArgumentError.Make(err.Error()).MakeData(res)
	}
	res.Checks = append(res.Checks, moduleInstallCheck("databaseTarget", "passed", "local go_makeadmin database confirmed"))

	if genSrv.db == nil {
		res.Status = "blocked"
		res.Message = "module install status requires configured database"
		res.Checks = append(res.Checks, moduleInstallCheck("database", "failed", res.Message))
		return res, response.AssertArgumentError.Make(res.Message).MakeData(res)
	}

	snapshot, err := moduleInstallSnapshot(genSrv.db, manifest, tenantID, roleID)
	if err != nil {
		res.Status = "failed"
		res.Message = "module install status snapshot failed"
		res.Checks = append(res.Checks, moduleInstallCheck("snapshot", "failed", err.Error()))
		return res, response.CheckErr(err, "ModuleManifestInstallStatus Snapshot err")
	}
	res.Snapshot = snapshot
	res.Missing = moduleInstallMissingSnapshot(res.Expected, snapshot)
	res.MenuVisible = moduleInstallMenuVisible(genSrv.db, manifest)
	res.Status = moduleInstallStatus(res.Expected, snapshot)
	res.Message = moduleInstallStatusMessage(res.Status)
	res.Checks = append(res.Checks, moduleInstallCheck("snapshot", "passed", "module install snapshot loaded"))
	return res, nil
}

func moduleRuntimeEnv(runtimeHint string) string {
	parts := strings.SplitN(strings.TrimSpace(runtimeHint), "=", 2)
	if len(parts) != 2 {
		return ""
	}
	name := strings.TrimSpace(parts[0])
	if !strings.HasPrefix(name, "MAKEADMIN_ENABLE_") {
		return ""
	}
	return name
}

func moduleInstallExpectedSnapshot(manifest moduleManifest) resp.ModuleManifestInstallSnapshotResp {
	expected := resp.ModuleManifestInstallSnapshotResp{
		Permissions:     int64(len(manifest.Permissions)),
		Menus:           1,
		RolePermissions: int64(len(manifest.Permissions)),
	}
	if strings.TrimSpace(manifest.Menu.Permission) != "" {
		expected.MenuPermissions = 1
	}
	return expected
}

func moduleInstallMissingSnapshot(expected, snapshot resp.ModuleManifestInstallSnapshotResp) resp.ModuleManifestInstallSnapshotResp {
	return resp.ModuleManifestInstallSnapshotResp{
		Permissions:     missingCount(expected.Permissions, snapshot.Permissions),
		Menus:           missingCount(expected.Menus, snapshot.Menus),
		MenuPermissions: missingCount(expected.MenuPermissions, snapshot.MenuPermissions),
		RolePermissions: missingCount(expected.RolePermissions, snapshot.RolePermissions),
	}
}

func missingCount(expected, actual int64) int64 {
	if actual >= expected {
		return 0
	}
	return expected - actual
}

func moduleInstallStatus(expected, snapshot resp.ModuleManifestInstallSnapshotResp) string {
	if moduleManifestSnapshotTotal(snapshot) == 0 {
		return "uninstalled"
	}
	missing := moduleInstallMissingSnapshot(expected, snapshot)
	if moduleManifestSnapshotTotal(missing) == 0 {
		return "installed"
	}
	return "partial"
}

func moduleManifestSnapshotTotal(snapshot resp.ModuleManifestInstallSnapshotResp) int64 {
	return snapshot.Permissions + snapshot.Menus + snapshot.MenuPermissions + snapshot.RolePermissions
}

func moduleInstallStatusMessage(status string) string {
	switch status {
	case "installed":
		return "module is installed"
	case "partial":
		return "module is partially installed"
	case "uninstalled":
		return "module is not installed"
	default:
		return "module status is unknown"
	}
}

func moduleInstallMenuVisible(db *gorm.DB, manifest moduleManifest) bool {
	var count int64
	err := db.Model(&makeadmin.Menu{}).
		Where("route_name = ? AND is_visible = 1 AND status = ? AND delete_time = 0", manifest.Menu.RouteName, makeadmin.StatusEnabled).
		Count(&count).Error
	return err == nil && count > 0
}
