package routers

import (
	"go-makeadmin/core"
	"go-makeadmin/generator/routers/gen"
)

var InitRouters = []*core.GroupBase{
	// gen
	gen.GenGroup,
}
