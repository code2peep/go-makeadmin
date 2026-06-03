package gen

import (
	"os"
	"testing"

	"go-makeadmin/config"
	"go-makeadmin/core"
	"go-makeadmin/generator/schemas/req"
	"go-makeadmin/generator/schemas/resp"
)

const moduleManifestUninstallSmokeEnv = "MAKEADMIN_ALLOW_MODULE_UNINSTALL_SMOKE_WRITE"

func TestModuleManifestUninstallApplyLocalSmoke(t *testing.T) {
	if os.Getenv(moduleManifestUninstallSmokeEnv) != "1" {
		t.Skip(moduleManifestUninstallSmokeEnv + " is not set")
	}
	t.Setenv(moduleManifestInstallApplyEnv, "1")
	t.Setenv(moduleManifestUninstallApplyEnv, "1")
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
	start, err := moduleUninstallSnapshot(db, manifest)
	if err != nil {
		t.Fatalf("initial snapshot: %v", err)
	}
	if moduleInstallSnapshotTotal(start) != 0 {
		t.Fatalf("demo article rows already exist before uninstall smoke: %+v", start)
	}

	cleanup := func() {
		if cleanupErr := cleanupModuleInstallSmokeRows(db, manifest); cleanupErr != nil {
			t.Fatalf("cleanup module uninstall smoke rows: %v", cleanupErr)
		}
	}
	defer cleanup()

	srv := generateService{db: db}
	_, err = srv.ApplyModuleManifestInstall(req.ModuleManifestInstallApplyReq{
		ManifestPath:    "examples/demo/manifest.json",
		ConfirmModule:   "article",
		ConfirmTenantID: &tenantID,
		ConfirmRoleID:   &roleID,
		ConfirmInstall:  true,
	})
	if err != nil {
		t.Fatalf("seed install apply: %v", err)
	}
	installed, err := moduleUninstallSnapshot(db, manifest)
	if err != nil {
		t.Fatalf("installed snapshot: %v", err)
	}
	assertInstallSnapshot(t, installed, resp.ModuleManifestInstallSnapshotResp{
		Permissions:     5,
		Menus:           1,
		MenuPermissions: 1,
		RolePermissions: 5,
	})

	first, err := srv.ApplyModuleManifestUninstall(req.ModuleManifestUninstallApplyReq{
		ManifestPath:  "examples/demo/manifest.json",
		ConfirmModule: "article",
		ConfirmDelete: true,
	})
	if err != nil {
		t.Fatalf("first uninstall apply: %v", err)
	}
	assertInstallSnapshot(t, first.Before, installed)
	assertInstallSnapshot(t, first.After, resp.ModuleManifestInstallSnapshotResp{})

	second, err := srv.ApplyModuleManifestUninstall(req.ModuleManifestUninstallApplyReq{
		ManifestPath:  "examples/demo/manifest.json",
		ConfirmModule: "article",
		ConfirmDelete: true,
	})
	if err != nil {
		t.Fatalf("second uninstall apply: %v", err)
	}
	assertInstallSnapshot(t, second.Before, resp.ModuleManifestInstallSnapshotResp{})
	assertInstallSnapshot(t, second.After, resp.ModuleManifestInstallSnapshotResp{})
}
