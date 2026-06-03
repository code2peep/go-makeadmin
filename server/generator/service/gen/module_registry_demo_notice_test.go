package gen

import (
	"strings"
	"testing"

	"go-makeadmin/generator/schemas/req"
)

func TestListModuleRegistryIncludesDemoNoticeWhenEnabled(t *testing.T) {
	t.Setenv(EnableDemoNoticeModuleRegistryEnv, "1")

	srv := generateService{}
	items := srv.ListModuleRegistry()
	notice := findModuleRegistryItem(items, "demo_notice")
	if notice == nil {
		t.Fatalf("demo notice registry item not found: %+v", items)
	}
	if notice.Manifest != "examples/demo_notice/manifest.json" ||
		notice.Table != "ma_demo_notice" ||
		notice.Entry != "/demo/notice" {
		t.Fatalf("unexpected demo notice registry item: %+v", notice)
	}
	if notice.Runtime != moduleRuntimeNoGate {
		t.Fatalf("unexpected demo notice runtime: %s", notice.Runtime)
	}
	if notice.ManifestStatus != "passed" || notice.ManifestMessage == "" {
		t.Fatalf("unexpected demo notice manifest status: %+v", notice)
	}
}

func TestDemoNoticeManifestUsesNoRuntimeGate(t *testing.T) {
	manifest, _, err := loadModuleManifest(req.ModuleManifestPreviewReq{
		ManifestPath: "examples/demo_notice/manifest.json",
	})
	if err != nil {
		t.Fatalf("load demo notice manifest: %v", err)
	}
	if manifest.Module != "demo_notice" || manifest.RuntimeRegistered {
		t.Fatalf("unexpected demo notice manifest: %+v", manifest)
	}
	if moduleRuntimeHint(manifest) != moduleRuntimeNoGate {
		t.Fatalf("unexpected demo notice runtime hint: %s", moduleRuntimeHint(manifest))
	}
}

func TestPreviewDemoNoticeManifestIncludesInstallPlan(t *testing.T) {
	srv := generateService{}
	res, err := srv.PreviewModuleManifest(req.ModuleManifestPreviewReq{
		ManifestPath: "examples/demo_notice/manifest.json",
		TenantID:     0,
		RoleID:       2,
	})
	if err != nil {
		t.Fatalf("preview demo notice manifest: %v", err)
	}
	if res.Manifest.Module != "demo_notice" ||
		res.Manifest.Table != "ma_demo_notice" ||
		res.Manifest.MenuName != "Demo Notice" {
		t.Fatalf("unexpected demo notice preview summary: %+v", res.Manifest)
	}
	if res.Plan.RuntimeHint != moduleRuntimeNoGate {
		t.Fatalf("unexpected demo notice runtime hint: %s", res.Plan.RuntimeHint)
	}
	if res.Plan.TenantID != 0 || res.Plan.RoleID != 2 {
		t.Fatalf("unexpected demo notice plan ids: %+v", res.Plan)
	}
	for _, needle := range []string{
		"SET @parent_route_name = 'dev_tools';",
		"`route_name` = 'demo.notice'",
		"demo_notice:list",
		"demo_notice:detail",
		"DELETE FROM `ma_menu`",
	} {
		if !strings.Contains(res.Plan.InstallSQL+res.Plan.RegistrySQL+res.Plan.UninstallSQL, needle) {
			t.Fatalf("demo notice plan missing %q: %+v", needle, res.Plan)
		}
	}
}
