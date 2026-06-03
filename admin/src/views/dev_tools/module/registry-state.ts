export type ElementTagType = 'primary' | 'success' | 'warning' | 'danger' | 'info'

export type RegistryStateModule = {
    module: string
    registryStatusCode: string
    entry?: string
    registryCheckCount?: number
}

export type RegistryFilterModule = RegistryStateModule & {
    installStatusCode: string
}

export type RegistryStateInput = {
    modules: ReadonlyArray<RegistryStateModule>
    registryLoaded: boolean
    registryLoading: boolean
    registryError: string
}

export type RegistryAcceptanceRow = {
    key: string
    label: string
    value: string
    type: ElementTagType
}

export type RegistryManualChecklistRow = {
    key: string
    label: string
    status: string
    statusType: ElementTagType
    detail: string
}

export type ModuleRuntimeStatusInput = {
    runtimeRegistered?: boolean
    runtimeEnv?: string
    runtimeEnabled?: boolean
    runtimeHint?: string
}

export type ModuleRuntimeStatusView = {
    label: string
    type: ElementTagType
    detail: string
}

export type ModuleStatusSummaryRow = {
    key: string
    label: string
    value: number
}

export const registrySmokeCommand = 'scripts/check-module-registry-smoke.sh'

export const countRegistryFailures = (modules: ReadonlyArray<RegistryStateModule>) =>
    modules.filter((item) => item.registryStatusCode === 'failed').length

export const hasBrokenRegistryFixture = (modules: ReadonlyArray<RegistryStateModule>) =>
    modules.some((item) => item.module === 'broken_fixture')

export const isRegistryEmptyState = (state: RegistryStateInput) =>
    state.registryLoaded && !state.registryLoading && !state.registryError && state.modules.length === 0

export const registryErrorDetailText = (registryError: string) =>
    registryError ? `${registryError} · ${registrySmokeCommand}` : registrySmokeCommand

export const registryTableEmptyTextFromState = (state: RegistryStateInput) => {
    if (state.registryError) {
        return 'registry 读取失败'
    }
    if (isRegistryEmptyState(state)) {
        return 'registry 暂无模块'
    }
    return '暂无匹配模块'
}

export const buildRegistryAcceptanceRows = (
    modules: ReadonlyArray<RegistryStateModule>
): RegistryAcceptanceRow[] => {
    const failedCount = countRegistryFailures(modules)
    const brokenFixtureEnabled = hasBrokenRegistryFixture(modules)
    return [
        {
            key: 'source',
            label: '来源',
            value: '/api/gen/moduleRegistry',
            type: 'info'
        },
        {
            key: 'total',
            label: '模块',
            value: `${modules.length}`,
            type: modules.length ? 'primary' : 'info'
        },
        {
            key: 'failed',
            label: '校验异常',
            value: `${failedCount}`,
            type: failedCount ? 'danger' : 'success'
        },
        {
            key: 'fixture',
            label: 'Broken Fixture',
            value: brokenFixtureEnabled ? '已开启' : '未开启',
            type: brokenFixtureEnabled ? 'warning' : 'info'
        },
        {
            key: 'smoke',
            label: 'Smoke',
            value: registrySmokeCommand,
            type: 'success'
        },
        {
            key: 'manual',
            label: '人工入口',
            value: '/module',
            type: 'warning'
        }
    ]
}

export const buildRegistryManualChecklistRows = (
    state: RegistryStateInput
): RegistryManualChecklistRow[] => {
    const failedCount = countRegistryFailures(state.modules)
    const brokenFixtureEnabled = hasBrokenRegistryFixture(state.modules)
    const articleModule = state.modules.find((item) => item.module === 'article')
    const hasCheckDetail = state.modules.some((item) => (item.registryCheckCount || 0) > 0)
    const registryReady = state.registryLoaded && !state.registryLoading && !state.registryError

    const registryStatus = () => {
        if (state.registryError) {
            return { label: '读取失败', type: 'danger' as ElementTagType }
        }
        if (state.registryLoading || !state.registryLoaded) {
            return { label: '读取中', type: 'info' as ElementTagType }
        }
        if (!state.modules.length) {
            return { label: '空清单', type: 'warning' as ElementTagType }
        }
        return { label: '已就绪', type: 'success' as ElementTagType }
    }

    const registry = registryStatus()
    return [
        {
            key: 'registry',
            label: '默认 Registry',
            status: registry.label,
            statusType: registry.type,
            detail: '/api/gen/moduleRegistry'
        },
        {
            key: 'fixture',
            label: 'Broken Fixture',
            status: brokenFixtureEnabled ? '已开启' : '未开启',
            statusType: brokenFixtureEnabled ? 'warning' : 'info',
            detail: 'MAKEADMIN_ENABLE_BROKEN_MODULE_REGISTRY_FIXTURE=1'
        },
        {
            key: 'failed_filter',
            label: '异常筛选',
            status: failedCount ? '可验收' : '无异常',
            statusType: failedCount ? 'warning' : 'success',
            detail: `failed=${failedCount}`
        },
        {
            key: 'check_detail',
            label: '校验明细',
            status: hasCheckDetail ? '可打开' : '待返回',
            statusType: hasCheckDetail ? 'success' : 'info',
            detail: 'manifestChecks'
        },
        {
            key: 'demo_entry',
            label: 'Demo 入口',
            status: articleModule?.entry ? '可打开' : registryReady ? '未配置' : '待加载',
            statusType: articleModule?.entry ? 'success' : registryReady ? 'warning' : 'info',
            detail: articleModule?.entry || '/demo/article'
        }
    ]
}

export const buildModuleRuntimeStatus = (
    status: ModuleRuntimeStatusInput
): ModuleRuntimeStatusView => {
    if (!status.runtimeRegistered) {
        return {
            label: '未注册',
            type: 'warning',
            detail: status.runtimeHint || '-'
        }
    }
    if (status.runtimeEnv && !status.runtimeEnabled) {
        return {
            label: '未开启',
            type: 'warning',
            detail: `${status.runtimeEnv}=1`
        }
    }
    return {
        label: '已开启',
        type: 'success',
        detail: status.runtimeEnv ? `${status.runtimeEnv}=1` : status.runtimeHint || '-'
    }
}

export const isRegistryModuleFailed = (module: RegistryFilterModule) =>
    module.registryStatusCode === 'failed' || ['blocked', 'failed'].includes(module.installStatusCode)

export const filterRegistryModules = <T extends RegistryFilterModule>(
    modules: ReadonlyArray<T>,
    filter: string
): T[] => {
    if (filter === 'all') {
        return [...modules]
    }
    if (filter === 'failed') {
        return modules.filter(isRegistryModuleFailed)
    }
    return modules.filter((item) => item.installStatusCode === filter)
}

export const buildModuleStatusSummary = (
    modules: ReadonlyArray<RegistryFilterModule>
): ModuleStatusSummaryRow[] => {
    const countBy = (codes: string[]) =>
        modules.filter((item) => codes.includes(item.installStatusCode)).length
    return [
        { key: 'total', label: '总数', value: modules.length },
        { key: 'installed', label: '已安装', value: countBy(['installed']) },
        { key: 'partial', label: '部分', value: countBy(['partial']) },
        { key: 'uninstalled', label: '未安装', value: countBy(['uninstalled']) },
        { key: 'failed', label: '异常', value: modules.filter(isRegistryModuleFailed).length }
    ]
}
