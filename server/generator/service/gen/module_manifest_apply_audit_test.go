package gen

import (
	"testing"

	"go-makeadmin/generator/schemas/resp"
)

func TestBuildModuleManifestInstallAuditEvent(t *testing.T) {
	result := resp.ModuleManifestInstallApplyResp{
		Source: "examples/demo/manifest.json",
		Manifest: resp.ModuleManifestSummaryResp{
			Module:         "article",
			Entity:         "DemoArticle",
			Table:          "ma_demo_article",
			MenuName:       "Demo Article",
			RequiresSchema: true,
		},
		TenantID:    7,
		RoleID:      3,
		Status:      "blocked",
		Message:     "install blocked",
		RequiredEnv: moduleManifestInstallApplyEnv,
		Plan: resp.ModuleManifestPlanResp{
			TenantID:    7,
			RoleID:      3,
			RuntimeHint: "MAKEADMIN_ENABLE_DEMO_MODULE=1",
		},
		Summary: resp.ModuleManifestApplySummaryResp{
			Operation:       "install",
			Module:          "article",
			Entity:          "DemoArticle",
			Table:           "ma_demo_article",
			RouteName:       "demo.article",
			PermissionCodes: []string{"article:list"},
			RequiresSchema:  true,
			DatabaseScope:   moduleManifestLocalDatabaseScope,
			RuntimeHint:     "MAKEADMIN_ENABLE_DEMO_MODULE=1",
		},
		Checks: []resp.ModuleManifestInstallCheckResp{
			moduleInstallCheck("environment", "failed", "env missing"),
		},
		Before: resp.ModuleManifestInstallSnapshotResp{Permissions: 1},
		After:  resp.ModuleManifestInstallSnapshotResp{Permissions: 2},
	}
	actor := resp.ModuleManifestApplyAuditActorResp{ID: 10, Name: "admin", Type: "admin"}

	event := buildModuleManifestInstallAuditEvent("evt-install", result, actor, "2026-06-03T10:00:00+08:00", "2026-06-03T10:00:01+08:00")

	if event.EventID != "evt-install" || event.Operation != "install" || event.Source != result.Source {
		t.Fatalf("unexpected event identity: %+v", event)
	}
	if event.Scope.TenantID != 7 || event.Scope.RoleID != 3 || event.Scope.DatabaseScope != moduleManifestLocalDatabaseScope || !event.Scope.RequiresSchema {
		t.Fatalf("unexpected event scope: %+v", event.Scope)
	}
	if event.Actor != actor {
		t.Fatalf("unexpected actor: %+v", event.Actor)
	}
	if event.Before.Permissions != 1 || event.After.Permissions != 2 {
		t.Fatalf("unexpected snapshots: before=%+v after=%+v", event.Before, event.After)
	}
	if len(event.Checks) != 1 || event.Checks[0].Name != "environment" {
		t.Fatalf("unexpected checks: %+v", event.Checks)
	}
}

func TestBuildModuleManifestUninstallAuditEvent(t *testing.T) {
	result := resp.ModuleManifestUninstallApplyResp{
		Source: "examples/demo/manifest.json",
		Manifest: resp.ModuleManifestSummaryResp{
			Module:         "article",
			Entity:         "DemoArticle",
			Table:          "ma_demo_article",
			MenuName:       "Demo Article",
			RequiresSchema: false,
		},
		Status:      "applied",
		Message:     "uninstall completed",
		RequiredEnv: moduleManifestUninstallApplyEnv,
		Plan: resp.ModuleManifestPlanResp{
			TenantID: 0,
			RoleID:   1,
		},
		Summary: resp.ModuleManifestApplySummaryResp{
			Operation:       "uninstall",
			Module:          "article",
			Entity:          "DemoArticle",
			Table:           "ma_demo_article",
			RouteName:       "demo.article",
			PermissionCodes: []string{"article:list"},
			DatabaseScope:   moduleManifestLocalDatabaseScope,
		},
		Before: resp.ModuleManifestInstallSnapshotResp{Menus: 1},
		After:  resp.ModuleManifestInstallSnapshotResp{},
	}
	actor := resp.ModuleManifestApplyAuditActorResp{ID: 11, Name: "operator", Type: "admin"}

	event := buildModuleManifestUninstallAuditEvent("evt-uninstall", result, actor, "2026-06-03T11:00:00+08:00", "2026-06-03T11:00:01+08:00")

	if event.EventID != "evt-uninstall" || event.Operation != "uninstall" || event.Status != "applied" {
		t.Fatalf("unexpected event identity: %+v", event)
	}
	if event.Scope.TenantID != 0 || event.Scope.RoleID != 1 || event.Scope.RequiresSchema {
		t.Fatalf("unexpected event scope: %+v", event.Scope)
	}
	if event.RequiredEnv != moduleManifestUninstallApplyEnv {
		t.Fatalf("required env = %q", event.RequiredEnv)
	}
	if event.Before.Menus != 1 || event.After.Menus != 0 {
		t.Fatalf("unexpected snapshots: before=%+v after=%+v", event.Before, event.After)
	}
}
