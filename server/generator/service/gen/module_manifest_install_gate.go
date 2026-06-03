package gen

import (
	"fmt"
	"os"
	"strings"

	"go-makeadmin/core/response"
	"go-makeadmin/generator/schemas/req"
	"go-makeadmin/generator/schemas/resp"
	"go-makeadmin/model/makeadmin"
)

const moduleManifestInstallApplyEnv = "MAKEADMIN_ALLOW_MODULE_INSTALL_APPLY"

// ApplyModuleManifestInstall validates the write gate for a module install request.
// P3.10 deliberately stops before any database access or SQL execution.
func (genSrv generateService) ApplyModuleManifestInstall(applyReq req.ModuleManifestInstallApplyReq) (res resp.ModuleManifestInstallApplyResp, e error) {
	previewReq := moduleManifestPreviewReqFromInstallApply(applyReq)
	manifest, source, err := loadModuleManifest(previewReq)
	if err != nil {
		return res, err
	}
	if err = validateModuleManifest(manifest); err != nil {
		return res, err
	}

	tenantID := moduleManifestInstallTenantID(applyReq)
	roleID := moduleManifestInstallRoleID(applyReq)
	previewReq.TenantID = tenantID
	previewReq.RoleID = roleID
	res = resp.ModuleManifestInstallApplyResp{
		Source: source,
		Manifest: resp.ModuleManifestSummaryResp{
			Module:         manifest.Module,
			Entity:         manifest.Entity,
			Table:          manifest.Table,
			MenuName:       manifest.Menu.Name,
			RequiresSchema: manifest.RequiresSchema,
		},
		TenantID:    tenantID,
		RoleID:      roleID,
		Status:      "blocked",
		RequiredEnv: moduleManifestInstallApplyEnv,
		Plan:        buildModuleManifestPlan(manifest, previewReq),
	}

	if os.Getenv(moduleManifestInstallApplyEnv) != "1" {
		res.Checks = append(res.Checks, moduleInstallCheck("environment", "failed", moduleManifestInstallApplyEnv+"=1 is required"))
		return res, moduleInstallGateError(res, moduleManifestInstallApplyEnv+"=1 is required; no database access was attempted")
	}
	res.Checks = append(res.Checks, moduleInstallCheck("environment", "passed", moduleManifestInstallApplyEnv+"=1 is present"))

	confirmModule := strings.TrimSpace(applyReq.ConfirmModule)
	if confirmModule != manifest.Module {
		res.Checks = append(res.Checks, moduleInstallCheck("confirmModule", "failed", fmt.Sprintf("confirmModule must be %q", manifest.Module)))
		return res, moduleInstallGateError(res, fmt.Sprintf("confirmModule must be %q; no database access was attempted", manifest.Module))
	}
	res.Checks = append(res.Checks, moduleInstallCheck("confirmModule", "passed", "module name confirmed"))

	if applyReq.ConfirmTenantID == nil || *applyReq.ConfirmTenantID != tenantID {
		res.Checks = append(res.Checks, moduleInstallCheck("confirmTenantId", "failed", fmt.Sprintf("confirmTenantId must be %d", tenantID)))
		return res, moduleInstallGateError(res, fmt.Sprintf("confirmTenantId must be %d; no database access was attempted", tenantID))
	}
	res.Checks = append(res.Checks, moduleInstallCheck("confirmTenantId", "passed", "tenant id confirmed"))

	if applyReq.ConfirmRoleID == nil || *applyReq.ConfirmRoleID != roleID {
		res.Checks = append(res.Checks, moduleInstallCheck("confirmRoleId", "failed", fmt.Sprintf("confirmRoleId must be %d", roleID)))
		return res, moduleInstallGateError(res, fmt.Sprintf("confirmRoleId must be %d; no database access was attempted", roleID))
	}
	res.Checks = append(res.Checks, moduleInstallCheck("confirmRoleId", "passed", "role id confirmed"))

	if !applyReq.ConfirmInstall {
		res.Checks = append(res.Checks, moduleInstallCheck("confirmInstall", "failed", "confirmInstall must be true"))
		return res, moduleInstallGateError(res, "confirmInstall must be true; no database access was attempted")
	}
	res.Checks = append(res.Checks, moduleInstallCheck("confirmInstall", "passed", "install intent confirmed"))

	if manifest.RequiresSchema && !applyReq.ConfirmSchemaRisk {
		res.Checks = append(res.Checks, moduleInstallCheck("confirmSchemaRisk", "failed", "confirmSchemaRisk must be true because manifest.requiresSchema is true"))
		return res, moduleInstallGateError(res, "confirmSchemaRisk must be true because manifest.requiresSchema is true; no database access was attempted")
	}
	if manifest.RequiresSchema {
		res.Checks = append(res.Checks, moduleInstallCheck("confirmSchemaRisk", "passed", "schema risk confirmed"))
	} else {
		res.Checks = append(res.Checks, moduleInstallCheck("confirmSchemaRisk", "skipped", "manifest.requiresSchema is false"))
	}

	res.Checks = append(res.Checks, moduleInstallCheck("executor", "blocked", "module install apply executor is not open in P3.10"))
	return res, moduleInstallGateError(res, "module install apply executor is not open in P3.10; no database access was attempted")
}

func moduleManifestPreviewReqFromInstallApply(applyReq req.ModuleManifestInstallApplyReq) req.ModuleManifestPreviewReq {
	return req.ModuleManifestPreviewReq{
		ManifestPath: applyReq.ManifestPath,
		ManifestBody: applyReq.ManifestBody,
		TenantID:     moduleManifestInstallTenantID(applyReq),
		RoleID:       moduleManifestInstallRoleID(applyReq),
	}
}

func moduleManifestInstallTenantID(applyReq req.ModuleManifestInstallApplyReq) uint64 {
	if applyReq.TenantID > 0 {
		return applyReq.TenantID
	}
	return makeadmin.GlobalTenantID
}

func moduleManifestInstallRoleID(applyReq req.ModuleManifestInstallApplyReq) uint64 {
	if applyReq.RoleID > 0 {
		return applyReq.RoleID
	}
	return 1
}

func moduleInstallCheck(name string, status string, message string) resp.ModuleManifestInstallCheckResp {
	return resp.ModuleManifestInstallCheckResp{Name: name, Status: status, Message: message}
}

func moduleInstallGateError(res resp.ModuleManifestInstallApplyResp, message string) error {
	res.Message = message
	return response.AssertArgumentError.Make(message).MakeData(res)
}
