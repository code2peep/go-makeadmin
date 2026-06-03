package gen

import "go-makeadmin/generator/schemas/resp"

func buildModuleManifestInstallAuditEvent(
	eventID string,
	result resp.ModuleManifestInstallApplyResp,
	actor resp.ModuleManifestApplyAuditActorResp,
	requestedAt string,
	completedAt string,
) resp.ModuleManifestApplyAuditEventResp {
	return buildModuleManifestApplyAuditEvent(
		eventID,
		"install",
		result.Source,
		result.Manifest,
		result.Plan,
		result.Summary,
		result.Status,
		result.Message,
		result.RequiredEnv,
		result.Checks,
		result.Before,
		result.After,
		actor,
		requestedAt,
		completedAt,
	)
}

func buildModuleManifestUninstallAuditEvent(
	eventID string,
	result resp.ModuleManifestUninstallApplyResp,
	actor resp.ModuleManifestApplyAuditActorResp,
	requestedAt string,
	completedAt string,
) resp.ModuleManifestApplyAuditEventResp {
	return buildModuleManifestApplyAuditEvent(
		eventID,
		"uninstall",
		result.Source,
		result.Manifest,
		result.Plan,
		result.Summary,
		result.Status,
		result.Message,
		result.RequiredEnv,
		result.Checks,
		result.Before,
		result.After,
		actor,
		requestedAt,
		completedAt,
	)
}

func buildModuleManifestApplyAuditEvent(
	eventID string,
	operation string,
	source string,
	manifest resp.ModuleManifestSummaryResp,
	plan resp.ModuleManifestPlanResp,
	summary resp.ModuleManifestApplySummaryResp,
	status string,
	message string,
	requiredEnv string,
	checks []resp.ModuleManifestInstallCheckResp,
	before resp.ModuleManifestInstallSnapshotResp,
	after resp.ModuleManifestInstallSnapshotResp,
	actor resp.ModuleManifestApplyAuditActorResp,
	requestedAt string,
	completedAt string,
) resp.ModuleManifestApplyAuditEventResp {
	if summary.Operation != "" {
		operation = summary.Operation
	}
	return resp.ModuleManifestApplyAuditEventResp{
		EventID:   eventID,
		Operation: operation,
		Source:    source,
		Manifest:  manifest,
		Summary:   summary,
		Scope: resp.ModuleManifestApplyAuditScopeResp{
			TenantID:       plan.TenantID,
			RoleID:         plan.RoleID,
			DatabaseScope:  summary.DatabaseScope,
			RequiresSchema: manifest.RequiresSchema,
		},
		Status:      status,
		Message:     message,
		RequiredEnv: requiredEnv,
		Checks:      checks,
		Before:      before,
		After:       after,
		Actor:       actor,
		RequestedAt: requestedAt,
		CompletedAt: completedAt,
	}
}
