export type ElementTagType = 'primary' | 'success' | 'warning' | 'danger' | 'info'

export type RegistryStateModule = {
    module: string
    registryStatusCode: string
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
