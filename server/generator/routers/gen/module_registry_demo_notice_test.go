package gen

import (
	"testing"

	gensvc "go-makeadmin/generator/service/gen"
)

func TestListModuleRegistryRouteDemoNoticeResponse(t *testing.T) {
	t.Setenv(gensvc.EnableDemoNoticeModuleRegistryEnv, "1")

	items := moduleRegistryRouteResponse(t)
	notice := findModuleRegistryRouteItem(items, "demo_notice")
	if notice == nil {
		t.Fatalf("demo notice route item not found: %+v", items)
	}
	if notice.Manifest != "examples/demo_notice/manifest.json" ||
		notice.Entry != "/demo/notice" ||
		notice.ManifestStatus != "passed" {
		t.Fatalf("unexpected demo notice route item: %+v", notice)
	}
}
