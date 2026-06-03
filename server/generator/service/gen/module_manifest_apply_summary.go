package gen

import "go-makeadmin/generator/schemas/resp"

const moduleManifestLocalDatabaseScope = "local go_makeadmin only"

func buildModuleManifestApplySummary(manifest moduleManifest, operation string) resp.ModuleManifestApplySummaryResp {
	return resp.ModuleManifestApplySummaryResp{
		Operation:       operation,
		Module:          manifest.Module,
		Entity:          manifest.Entity,
		Table:           manifest.Table,
		RouteName:       manifest.Menu.RouteName,
		PermissionCodes: moduleManifestPermissionCodes(manifest),
		RequiresSchema:  manifest.RequiresSchema,
		DatabaseScope:   moduleManifestLocalDatabaseScope,
		RuntimeHint:     moduleRuntimeHint(manifest),
	}
}
