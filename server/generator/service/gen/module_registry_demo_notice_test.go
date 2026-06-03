package gen

import (
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
