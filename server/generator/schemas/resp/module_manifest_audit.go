package resp

// ModuleManifestApplyAuditActorResp 模块安装/卸载审计操作人草图
type ModuleManifestApplyAuditActorResp struct {
	ID   uint64 `json:"id" structs:"id"`     // 操作人ID
	Name string `json:"name" structs:"name"` // 操作人名称
	Type string `json:"type" structs:"type"` // 操作人类型
}

// ModuleManifestApplyAuditScopeResp 模块安装/卸载审计范围草图
type ModuleManifestApplyAuditScopeResp struct {
	TenantID       uint64 `json:"tenantId" structs:"tenantId"`             // 租户ID
	RoleID         uint64 `json:"roleId" structs:"roleId"`                 // 角色ID
	DatabaseScope  string `json:"databaseScope" structs:"databaseScope"`   // 数据库范围
	RequiresSchema bool   `json:"requiresSchema" structs:"requiresSchema"` // 是否需要业务表
}

// ModuleManifestApplyAuditEventResp 模块安装/卸载审计事件草图
type ModuleManifestApplyAuditEventResp struct {
	EventID     string                            `json:"eventId" structs:"eventId"`         // 审计事件ID
	Operation   string                            `json:"operation" structs:"operation"`     // 操作类型
	Source      string                            `json:"source" structs:"source"`           // manifest 来源
	Manifest    ModuleManifestSummaryResp         `json:"manifest" structs:"manifest"`       // manifest 摘要
	Summary     ModuleManifestApplySummaryResp    `json:"summary" structs:"summary"`         // 操作摘要
	Scope       ModuleManifestApplyAuditScopeResp `json:"scope" structs:"scope"`             // 执行范围
	Status      string                            `json:"status" structs:"status"`           // 执行状态
	Message     string                            `json:"message" structs:"message"`         // 执行说明
	RequiredEnv string                            `json:"requiredEnv" structs:"requiredEnv"` // 写入环境变量
	Checks      []ModuleManifestInstallCheckResp  `json:"checks" structs:"checks"`           // 检查结果
	Before      ModuleManifestInstallSnapshotResp `json:"before" structs:"before"`           // 执行前快照
	After       ModuleManifestInstallSnapshotResp `json:"after" structs:"after"`             // 执行后快照
	Actor       ModuleManifestApplyAuditActorResp `json:"actor" structs:"actor"`             // 操作人
	RequestedAt string                            `json:"requestedAt" structs:"requestedAt"` // 请求时间
	CompletedAt string                            `json:"completedAt" structs:"completedAt"` // 完成时间
}
