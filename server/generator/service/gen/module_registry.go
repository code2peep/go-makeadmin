package gen

import (
	"os"
	"strings"

	"go-makeadmin/generator/schemas/req"
	"go-makeadmin/generator/schemas/resp"
)

const EnableBrokenModuleRegistryFixtureEnv = "MAKEADMIN_ENABLE_BROKEN_MODULE_REGISTRY_FIXTURE"
const EnableDemoNoticeModuleRegistryEnv = "MAKEADMIN_ENABLE_DEMO_NOTICE_MODULE"

// ListModuleRegistry returns built-in modules known to the local framework.
func (genSrv generateService) ListModuleRegistry() []resp.ModuleRegistryItemResp {
	items := []resp.ModuleRegistryItemResp{
		{
			Name:       "Demo Article",
			Module:     "article",
			Manifest:   "examples/demo/manifest.json",
			Table:      "ma_demo_article",
			Runtime:    "MAKEADMIN_ENABLE_DEMO_MODULE=1",
			Entry:      "/demo/article",
			Status:     "可安装",
			StatusType: "success",
		},
	}
	if strings.TrimSpace(os.Getenv(EnableDemoNoticeModuleRegistryEnv)) == "1" {
		items = append(items, resp.ModuleRegistryItemResp{
			Name:       "Demo Notice",
			Module:     "demo_notice",
			Manifest:   "examples/demo_notice/manifest.json",
			Table:      "ma_demo_notice",
			Runtime:    moduleRuntimeNoGate,
			Entry:      "/demo/notice",
			Status:     "前端只读示例",
			StatusType: "warning",
		})
	}
	if strings.TrimSpace(os.Getenv(EnableBrokenModuleRegistryFixtureEnv)) == "1" {
		items = append(items, resp.ModuleRegistryItemResp{
			Name:       "Broken Manifest Fixture",
			Module:     "broken_fixture",
			Manifest:   "examples/demo/missing/manifest.json",
			Table:      "ma_broken_fixture",
			Runtime:    EnableBrokenModuleRegistryFixtureEnv + "=1",
			Entry:      "/demo/broken-fixture",
			Status:     "本地异常示例",
			StatusType: "warning",
		})
	}
	for index := range items {
		items[index] = genSrv.validateModuleRegistryItem(items[index])
	}
	return items
}

func (genSrv generateService) validateModuleRegistryItem(item resp.ModuleRegistryItemResp) resp.ModuleRegistryItemResp {
	item.ManifestStatus = "passed"
	item.ManifestMessage = "manifest registry check passed"

	manifest, _, err := loadModuleManifest(req.ModuleManifestPreviewReq{ManifestPath: item.Manifest})
	if err != nil {
		return moduleRegistryItemFailed(item, "manifest", err.Error())
	}
	item.ManifestChecks = append(item.ManifestChecks, moduleInstallCheck("manifest", "passed", "manifest loaded"))

	if err = validateModuleManifest(manifest); err != nil {
		return moduleRegistryItemFailed(item, "manifestValidation", err.Error())
	}
	item.ManifestChecks = append(item.ManifestChecks, moduleInstallCheck("manifestValidation", "passed", "manifest shape is valid"))

	if manifest.Module != item.Module {
		return moduleRegistryItemFailed(item, "module", "registry module does not match manifest module")
	}
	item.ManifestChecks = append(item.ManifestChecks, moduleInstallCheck("module", "passed", "module matches manifest"))

	if manifest.Table != item.Table {
		return moduleRegistryItemFailed(item, "table", "registry table does not match manifest table")
	}
	item.ManifestChecks = append(item.ManifestChecks, moduleInstallCheck("table", "passed", "table matches manifest"))

	if moduleRuntimeHint(manifest) != item.Runtime {
		return moduleRegistryItemFailed(item, "runtime", "registry runtime does not match manifest runtime hint")
	}
	item.ManifestChecks = append(item.ManifestChecks, moduleInstallCheck("runtime", "passed", "runtime hint matches manifest"))

	if strings.TrimSpace(item.Entry) == "" || !strings.HasPrefix(item.Entry, "/") {
		return moduleRegistryItemFailed(item, "entry", "registry entry must be an absolute admin route")
	}
	item.ManifestChecks = append(item.ManifestChecks, moduleInstallCheck("entry", "passed", "admin entry is configured"))

	if strings.TrimSpace(manifest.Menu.RouteName) == "" ||
		strings.TrimSpace(manifest.Menu.RoutePath) == "" ||
		strings.TrimSpace(manifest.Menu.Component) == "" {
		return moduleRegistryItemFailed(item, "menu", "manifest menu routeName, routePath and component are required")
	}
	item.ManifestChecks = append(item.ManifestChecks, moduleInstallCheck("menu", "passed", "manifest menu route is complete"))
	return item
}

func moduleRegistryItemFailed(item resp.ModuleRegistryItemResp, name string, message string) resp.ModuleRegistryItemResp {
	item.ManifestStatus = "failed"
	item.ManifestMessage = message
	item.ManifestChecks = append(item.ManifestChecks, moduleInstallCheck(name, "failed", message))
	return item
}
