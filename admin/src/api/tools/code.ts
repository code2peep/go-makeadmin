import request from '@/utils/request'

export interface ModuleManifestPreviewParams {
    manifestPath?: string
    manifestBody?: string
    authorName?: string
    tenantId?: number
    roleId?: number
}

export interface ModuleManifestInstallApplyParams extends ModuleManifestPreviewParams {
    confirmModule: string
    confirmTenantId: number
    confirmRoleId: number
    confirmInstall: boolean
    confirmSchemaRisk: boolean
}

export interface ModuleManifestUninstallApplyParams extends ModuleManifestPreviewParams {
    confirmModule: string
    confirmDelete: boolean
}

export interface GenTableBaseResult {
    id?: number
    tableName: string
    tableComment?: string
    entityName?: string
    authorName?: string
    remarks?: string
    createTime?: string
    updateTime?: string
}

export interface GenTableGenResult {
    genTpl: string
    genType?: number
    genPath?: string
    moduleName?: string
    functionName: string
    treePrimary?: string
    treeParent?: string
    treeName?: string
    subTableName?: string
    subTableFk?: string
}

export interface GenColumnResult {
    id?: number
    columnName: string
    columnComment?: string
    columnLength?: number
    columnType?: string
    goType: string
    goField: string
    isRequired?: number
    isInsert?: number
    isEdit?: number
    isList?: number
    isQuery?: number
    queryType: string
    htmlType: string
    dictType: string
    createTime?: string
    updateTime?: string
}

export interface GenTableDetailResult {
    base: GenTableBaseResult
    gen: GenTableGenResult
    column: GenColumnResult[]
}

export interface ModuleManifestSummaryResult {
    module: string
    entity: string
    table: string
    menuName: string
    requiresSchema: boolean
}

export interface ModuleManifestPlanResult {
    tenantId: number
    roleId: number
    registrySql: string
    roleGrantSql: string
    installSql: string
    uninstallSql: string
    runtimeHint: string
}

export interface ModuleManifestPreviewResult {
    source: string
    warning: string
    manifest: ModuleManifestSummaryResult
    detail: GenTableDetailResult
    code: Record<string, string>
    plan: ModuleManifestPlanResult
}

export interface ModuleManifestApplySummaryResult {
    operation?: 'install' | 'uninstall' | string
    module?: string
    entity?: string
    table?: string
    routeName?: string
    permissionCodes?: string[]
    requiresSchema?: boolean
    databaseScope?: string
    runtimeHint?: string
}

export interface ModuleManifestApplyCheckResult {
    name: string
    status: string
    message: string
}

export interface ModuleManifestApplySnapshotResult {
    permissions?: number
    menus?: number
    menuPermissions?: number
    rolePermissions?: number
}

export interface ModuleManifestApplyResult {
    source?: string
    manifest?: ModuleManifestSummaryResult
    tenantId?: number
    roleId?: number
    status?: string
    message?: string
    requiredEnv?: string
    plan?: ModuleManifestPlanResult
    summary?: ModuleManifestApplySummaryResult
    checks?: ModuleManifestApplyCheckResult[]
    before?: ModuleManifestApplySnapshotResult
    after?: ModuleManifestApplySnapshotResult
}

export interface ModuleManifestApplyAuditActorResult {
    id?: number
    name?: string
    type?: string
}

export interface ModuleManifestApplyAuditScopeResult {
    tenantId?: number
    roleId?: number
    databaseScope?: string
    requiresSchema?: boolean
}

export interface ModuleManifestApplyAuditEventResult {
    eventId?: string
    operation?: 'install' | 'uninstall' | string
    source?: string
    manifest?: ModuleManifestSummaryResult
    summary?: ModuleManifestApplySummaryResult
    scope?: ModuleManifestApplyAuditScopeResult
    status?: string
    message?: string
    requiredEnv?: string
    checks?: ModuleManifestApplyCheckResult[]
    before?: ModuleManifestApplySnapshotResult
    after?: ModuleManifestApplySnapshotResult
    actor?: ModuleManifestApplyAuditActorResult
    requestedAt?: string
    completedAt?: string
}

export interface ModuleManifestApplyAuditPreviewOptions {
    eventId?: string
    actor?: ModuleManifestApplyAuditActorResult
    requestedAt?: string
    completedAt?: string
}

export interface ModuleManifestApplyAuditPreviewSummaryResult {
    operation?: string
    module?: string
    status?: string
    routeName?: string
    permissionCount: number
    checkCount: number
    beforeTotal: number
    afterTotal: number
    databaseScope?: string
    actorType?: string
}

export function buildModuleManifestApplyAuditPreview(
    result: ModuleManifestApplyResult,
    options: ModuleManifestApplyAuditPreviewOptions = {}
): ModuleManifestApplyAuditEventResult {
    const previewAt = new Date().toISOString()
    return {
        eventId: options.eventId || 'preview',
        operation: result.summary?.operation || '',
        source: result.source,
        manifest: result.manifest,
        summary: result.summary,
        scope: {
            tenantId: result.tenantId ?? result.plan?.tenantId,
            roleId: result.roleId ?? result.plan?.roleId,
            databaseScope: result.summary?.databaseScope,
            requiresSchema: result.manifest?.requiresSchema ?? result.summary?.requiresSchema
        },
        status: result.status,
        message: result.message,
        requiredEnv: result.requiredEnv,
        checks: result.checks || [],
        before: result.before || {},
        after: result.after || {},
        actor: options.actor || {
            type: 'frontend-preview'
        },
        requestedAt: options.requestedAt || previewAt,
        completedAt: options.completedAt || previewAt
    }
}

export function buildModuleManifestApplyAuditPreviewSummary(
    event: ModuleManifestApplyAuditEventResult
): ModuleManifestApplyAuditPreviewSummaryResult {
    return {
        operation: event.operation,
        module: event.manifest?.module || event.summary?.module,
        status: event.status,
        routeName: event.summary?.routeName,
        permissionCount: event.summary?.permissionCodes?.length || 0,
        checkCount: event.checks?.length || 0,
        beforeTotal: moduleManifestApplySnapshotTotal(event.before),
        afterTotal: moduleManifestApplySnapshotTotal(event.after),
        databaseScope: event.scope?.databaseScope,
        actorType: event.actor?.type
    }
}

const moduleManifestApplySnapshotTotal = (snapshot?: ModuleManifestApplySnapshotResult) =>
    (snapshot?.permissions || 0) +
    (snapshot?.menus || 0) +
    (snapshot?.menuPermissions || 0) +
    (snapshot?.rolePermissions || 0)

const isRecord = (value: unknown): value is Record<string, unknown> =>
    typeof value === 'object' && value !== null

const isModuleManifestApplyResult = (value: unknown): value is ModuleManifestApplyResult => {
    if (!isRecord(value)) {
        return false
    }
    return Boolean(
        value.source ||
            value.manifest ||
            value.status ||
            value.message ||
            value.requiredEnv ||
            value.plan ||
            value.summary ||
            value.checks ||
            value.before ||
            value.after
    )
}

export function normalizeModuleManifestApplyError(
    error: unknown,
    fallbackMessage: string
): ModuleManifestApplyResult {
    if (isModuleManifestApplyResult(error)) {
        return {
            ...error,
            message: error.message || fallbackMessage,
            checks: error.checks || []
        }
    }

    if (error instanceof Error && error.message) {
        return {
            message: error.message,
            checks: []
        }
    }

    if (typeof error === 'string' && error) {
        return {
            message: error,
            checks: []
        }
    }

    return {
        message: fallbackMessage,
        checks: []
    }
}

// 代码生成已选数据表列表接口
export function generateTable(params: any) {
    return request.get({ url: '/gen/list', params })
}

// 数据表列表接口
export function dataTable(params: any) {
    return request.get({ url: '/gen/db', params })
}

//选择要生成代码的数据表
export function selectTable(params: any) {
    return request.post(
        { url: '/gen/importTable', params },
        {
            isParamsToData: false
        }
    )
}

// 已选择的数据表详情
export function tableDetail(params: any) {
    return request.get({ url: '/gen/detail', params })
}

//同步字段
export function syncColumn(params: any) {
    return request.post(
        { url: '/gen/syncTable', params },
        {
            isParamsToData: false
        }
    )
}

//删除已选择的数据表
export function generateDelete(params: any) {
    return request.post({ url: '/gen/delTable', params })
}

//编辑已选表字段
export function generateEdit(params: any) {
    return request.post({ url: '/gen/editTable', params })
}

//预览代码
export function generatePreview(params: any) {
    return request.get({ url: '/gen/previewCode', params })
}

// 模块 manifest 预览代码
export function previewModuleManifest(
    params: ModuleManifestPreviewParams
): Promise<ModuleManifestPreviewResult> {
    return request.post({ url: '/gen/previewCode', params })
}

// 模块 manifest 安装写入门禁
export function applyModuleManifestInstall(
    params: ModuleManifestInstallApplyParams
): Promise<ModuleManifestApplyResult> {
    return request.request({ url: '/gen/previewCode', method: 'PUT', data: params })
}

// 模块 manifest 卸载写入门禁
export function applyModuleManifestUninstall(
    params: ModuleManifestUninstallApplyParams
): Promise<ModuleManifestApplyResult> {
    return request.request({ url: '/gen/previewCode', method: 'DELETE', data: params })
}

//生成代码
export function generateCode(params: any) {
    return request.get({ url: '/gen/genCode', params })
}

//下载代码
export function downloadCode(params: any) {
    return request.get(
        { responseType: 'blob', url: '/gen/downloadCode', params },
        {
            isTransformResponse: false
        }
    )
}
