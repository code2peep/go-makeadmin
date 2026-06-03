package gen

import "go-makeadmin/generator/schemas/resp"

// ListModuleRegistry returns built-in modules known to the local framework.
func (genSrv generateService) ListModuleRegistry() []resp.ModuleRegistryItemResp {
	return []resp.ModuleRegistryItemResp{
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
}
