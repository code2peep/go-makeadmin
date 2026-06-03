package gen

import (
	"fmt"
	"os"
	"strings"

	"go-makeadmin/core/response"
	"go-makeadmin/generator/schemas/req"
	"go-makeadmin/generator/schemas/resp"
)

const moduleManifestUninstallApplyEnv = "MAKEADMIN_ALLOW_MODULE_UNINSTALL_APPLY"

// ApplyModuleManifestUninstall validates the write gate for a module uninstall request.
// P3.12 deliberately stops before any database access or SQL execution.
func (genSrv generateService) ApplyModuleManifestUninstall(applyReq req.ModuleManifestUninstallApplyReq) (res resp.ModuleManifestUninstallApplyResp, e error) {
	previewReq := moduleManifestPreviewReqFromUninstallApply(applyReq)
	manifest, source, err := loadModuleManifest(previewReq)
	if err != nil {
		return res, err
	}
	if err = validateModuleManifest(manifest); err != nil {
		return res, err
	}

	res = resp.ModuleManifestUninstallApplyResp{
		Source: source,
		Manifest: resp.ModuleManifestSummaryResp{
			Module:         manifest.Module,
			Entity:         manifest.Entity,
			Table:          manifest.Table,
			MenuName:       manifest.Menu.Name,
			RequiresSchema: manifest.RequiresSchema,
		},
		Status:      "blocked",
		RequiredEnv: moduleManifestUninstallApplyEnv,
		Plan:        buildModuleManifestPlan(manifest, previewReq),
	}

	if os.Getenv(moduleManifestUninstallApplyEnv) != "1" {
		res.Checks = append(res.Checks, moduleInstallCheck("environment", "failed", moduleManifestUninstallApplyEnv+"=1 is required"))
		return res, moduleUninstallGateError(res, moduleManifestUninstallApplyEnv+"=1 is required; no database access was attempted")
	}
	res.Checks = append(res.Checks, moduleInstallCheck("environment", "passed", moduleManifestUninstallApplyEnv+"=1 is present"))

	confirmModule := strings.TrimSpace(applyReq.ConfirmModule)
	if confirmModule != manifest.Module {
		res.Checks = append(res.Checks, moduleInstallCheck("confirmModule", "failed", fmt.Sprintf("confirmModule must be %q", manifest.Module)))
		return res, moduleUninstallGateError(res, fmt.Sprintf("confirmModule must be %q; no database access was attempted", manifest.Module))
	}
	res.Checks = append(res.Checks, moduleInstallCheck("confirmModule", "passed", "module name confirmed"))

	if !applyReq.ConfirmDelete {
		res.Checks = append(res.Checks, moduleInstallCheck("confirmDelete", "failed", "confirmDelete must be true"))
		return res, moduleUninstallGateError(res, "confirmDelete must be true; no database access was attempted")
	}
	res.Checks = append(res.Checks, moduleInstallCheck("confirmDelete", "passed", "delete intent confirmed"))

	res.Checks = append(res.Checks, moduleInstallCheck("executor", "blocked", "module uninstall apply executor is not open in P3.12"))
	return res, moduleUninstallGateError(res, "module uninstall apply executor is not open in P3.12; no database access was attempted")
}

func moduleManifestPreviewReqFromUninstallApply(applyReq req.ModuleManifestUninstallApplyReq) req.ModuleManifestPreviewReq {
	return req.ModuleManifestPreviewReq{
		ManifestPath: applyReq.ManifestPath,
		ManifestBody: applyReq.ManifestBody,
	}
}

func moduleUninstallGateError(res resp.ModuleManifestUninstallApplyResp, message string) error {
	res.Message = message
	return response.AssertArgumentError.Make(message).MakeData(res)
}
