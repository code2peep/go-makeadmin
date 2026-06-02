package routers

import (
	"os"

	"go-makeadmin/core"
	"go-makeadmin/modules/routers/demo"
)

const EnableDemoModuleEnv = "MAKEADMIN_ENABLE_DEMO_MODULE"

func InitRouters() []*core.GroupBase {
	if os.Getenv(EnableDemoModuleEnv) != "1" {
		return nil
	}
	return demo.InitRouters
}
