import {
    buildModuleRuntimeStatus,
    buildRegistryAcceptanceRows,
    buildRegistryManualChecklistRows,
    countRegistryFailures,
    hasBrokenRegistryFixture,
    isRegistryEmptyState,
    registryErrorDetailText,
    registryTableEmptyTextFromState,
    type RegistryAcceptanceRow,
    type RegistryManualChecklistRow,
    type ModuleRuntimeStatusView,
    type RegistryStateInput,
    type RegistryStateModule
} from './registry-state'

const defaultModules = [
    {
        module: 'article',
        registryStatusCode: 'passed',
        entry: '/demo/article',
        registryCheckCount: 7
    }
] satisfies RegistryStateModule[]

const brokenFixtureModules = [
    ...defaultModules,
    {
        module: 'broken_fixture',
        registryStatusCode: 'failed',
        entry: '/broken-fixture',
        registryCheckCount: 7
    }
] satisfies RegistryStateModule[]

const defaultRegistryState = {
    modules: defaultModules,
    registryLoaded: true,
    registryLoading: false,
    registryError: ''
} satisfies RegistryStateInput

const brokenRegistryState = {
    modules: brokenFixtureModules,
    registryLoaded: true,
    registryLoading: false,
    registryError: ''
} satisfies RegistryStateInput

const emptyRegistryState = {
    modules: [],
    registryLoaded: true,
    registryLoading: false,
    registryError: ''
} satisfies RegistryStateInput

const failedRegistryState = {
    modules: [],
    registryLoaded: true,
    registryLoading: false,
    registryError: 'network error'
} satisfies RegistryStateInput

const defaultRows: RegistryAcceptanceRow[] = buildRegistryAcceptanceRows(defaultModules)
const brokenRows: RegistryAcceptanceRow[] = buildRegistryAcceptanceRows(brokenFixtureModules)
const defaultChecklistRows: RegistryManualChecklistRow[] =
    buildRegistryManualChecklistRows(defaultRegistryState)
const brokenChecklistRows: RegistryManualChecklistRow[] =
    buildRegistryManualChecklistRows(brokenRegistryState)
const emptyChecklistRows: RegistryManualChecklistRow[] =
    buildRegistryManualChecklistRows(emptyRegistryState)
const defaultFailureCount: number = countRegistryFailures(defaultModules)
const brokenFailureCount: number = countRegistryFailures(brokenFixtureModules)
const brokenEnabled: boolean = hasBrokenRegistryFixture(brokenFixtureModules)
const emptyStateMatched: boolean = isRegistryEmptyState(emptyRegistryState)
const emptyText: string = registryTableEmptyTextFromState(emptyRegistryState)
const failedText: string = registryTableEmptyTextFromState(failedRegistryState)
const failedDetail: string = registryErrorDetailText(failedRegistryState.registryError)
const demoNoticeRuntimeStatus: ModuleRuntimeStatusView = buildModuleRuntimeStatus({
    runtimeRegistered: false,
    runtimeHint: 'No runtime env gate is defined for this module yet.'
})
const demoArticleRuntimeStatus: ModuleRuntimeStatusView = buildModuleRuntimeStatus({
    runtimeRegistered: true,
    runtimeEnv: 'MAKEADMIN_ENABLE_DEMO_MODULE',
    runtimeEnabled: true,
    runtimeHint: 'MAKEADMIN_ENABLE_DEMO_MODULE=1'
})

export const registryStateFixture = {
    defaultRows,
    brokenRows,
    defaultChecklistRows,
    brokenChecklistRows,
    emptyChecklistRows,
    defaultFailureCount,
    brokenFailureCount,
    brokenEnabled,
    emptyStateMatched,
    emptyText,
    failedText,
    failedDetail,
    demoNoticeRuntimeStatus,
    demoArticleRuntimeStatus
}
