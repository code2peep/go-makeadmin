package gen

import (
	"os"
	"testing"

	"go-makeadmin/config"
	"go-makeadmin/core"
	"go-makeadmin/generator/schemas/req"
	"go-makeadmin/generator/schemas/resp"
	"gorm.io/gorm"
)

const moduleManifestInstallSmokeEnv = "MAKEADMIN_ALLOW_MODULE_INSTALL_SMOKE_WRITE"

func TestModuleManifestInstallApplyLocalSmoke(t *testing.T) {
	if os.Getenv(moduleManifestInstallSmokeEnv) != "1" {
		t.Skip(moduleManifestInstallSmokeEnv + " is not set")
	}
	t.Setenv(moduleManifestInstallApplyEnv, "1")
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		config.Config.DatabaseUrl = databaseURL
	}

	db := core.GetDB()
	manifest, _, err := loadModuleManifest(req.ModuleManifestPreviewReq{ManifestPath: "examples/demo/manifest.json"})
	if err != nil {
		t.Fatalf("load manifest: %v", err)
	}
	if err = validateModuleManifest(manifest); err != nil {
		t.Fatalf("validate manifest: %v", err)
	}
	tenantID := uint64(0)
	roleID := uint64(1)
	start, err := moduleInstallSnapshot(db, manifest, tenantID, roleID)
	if err != nil {
		t.Fatalf("initial snapshot: %v", err)
	}
	if moduleInstallSnapshotTotal(start) != 0 {
		t.Fatalf("demo article rows already exist before install smoke: %+v", start)
	}

	cleanup := func() {
		if cleanupErr := cleanupModuleInstallSmokeRows(db, manifest); cleanupErr != nil {
			t.Fatalf("cleanup module install smoke rows: %v", cleanupErr)
		}
	}
	defer cleanup()

	srv := generateService{db: db}
	applyReq := req.ModuleManifestInstallApplyReq{
		ManifestPath:    "examples/demo/manifest.json",
		ConfirmModule:   "article",
		ConfirmTenantID: &tenantID,
		ConfirmRoleID:   &roleID,
		ConfirmInstall:  true,
	}
	first, err := srv.ApplyModuleManifestInstall(applyReq)
	if err != nil {
		t.Fatalf("first install apply: %v", err)
	}
	assertInstallSnapshot(t, first.Before, resp.ModuleManifestInstallSnapshotResp{})
	assertInstallSnapshot(t, first.After, resp.ModuleManifestInstallSnapshotResp{
		Permissions:     5,
		Menus:           1,
		MenuPermissions: 1,
		RolePermissions: 5,
	})

	second, err := srv.ApplyModuleManifestInstall(applyReq)
	if err != nil {
		t.Fatalf("second install apply: %v", err)
	}
	assertInstallSnapshot(t, second.Before, first.After)
	assertInstallSnapshot(t, second.After, first.After)

	cleanup()
	finalSnapshot, err := moduleInstallSnapshot(db, manifest, tenantID, roleID)
	if err != nil {
		t.Fatalf("final snapshot: %v", err)
	}
	assertInstallSnapshot(t, finalSnapshot, resp.ModuleManifestInstallSnapshotResp{})
}

func cleanupModuleInstallSmokeRows(db *gorm.DB, manifest moduleManifest) error {
	codes := moduleManifestPermissionCodes(manifest)
	return db.Transaction(func(tx *gorm.DB) error {
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
	})
}

func moduleInstallSnapshotTotal(snapshot resp.ModuleManifestInstallSnapshotResp) int64 {
	return snapshot.Permissions + snapshot.Menus + snapshot.MenuPermissions + snapshot.RolePermissions
}

func assertInstallSnapshot(t *testing.T, got resp.ModuleManifestInstallSnapshotResp, want resp.ModuleManifestInstallSnapshotResp) {
	t.Helper()
	if got != want {
		t.Fatalf("install snapshot = %+v, want %+v", got, want)
	}
}
