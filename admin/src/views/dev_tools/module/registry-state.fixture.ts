import {
    buildModuleRuntimeStatus,
    buildModuleInstallWizardRows,
    buildModuleMarketRows,
    buildModuleStatusSummary,
    filterRegistryModules,
    buildRegistryAcceptanceRows,
    buildRegistryManualChecklistRows,
    countRegistryFailures,
    hasBrokenRegistryFixture,
    isRegistryEmptyState,
    registryErrorDetailText,
    registryTableEmptyTextFromState,
    type RegistryAcceptanceRow,
    type RegistryManualChecklistRow,
    type ModuleInstallWizardRow,
    type ModuleMarketRow,
    type ModuleRuntimeStatusView,
    type ModuleStatusSummaryRow,
    type RegistryFilterModule,
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

const multiRegistryModules = [
    ...defaultModules,
    {
        module: 'demo_notice',
        registryStatusCode: 'passed',
        entry: '/demo/notice',
        registryCheckCount: 7
    }
] satisfies RegistryStateModule[]

const multiStatusModules = [
    {
        module: 'article',
        registryStatusCode: 'passed',
        entry: '/demo/article',
        registryCheckCount: 7,
        installStatusCode: 'installed'
    },
    {
        module: 'demo_notice',
        registryStatusCode: 'passed',
        entry: '/demo/notice',
        registryCheckCount: 7,
        installStatusCode: 'uninstalled'
    }
] satisfies RegistryFilterModule[]

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

const multiRegistryState = {
    modules: multiRegistryModules,
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
const multiStatusSummaryRows: ModuleStatusSummaryRow[] =
    buildModuleStatusSummary(multiStatusModules)
const marketRows: ModuleMarketRow[] = buildModuleMarketRows(multiStatusModules)
const selectedWizardRows: ModuleInstallWizardRow[] =
    buildModuleInstallWizardRows(multiStatusModules[0])
const emptyWizardRows: ModuleInstallWizardRow[] = buildModuleInstallWizardRows()
const multiAllModules: RegistryFilterModule[] = filterRegistryModules(multiStatusModules, 'all')
const multiUninstalledModules: RegistryFilterModule[] = filterRegistryModules(
    multiStatusModules,
    'uninstalled'
)
const multiFailedModules: RegistryFilterModule[] = filterRegistryModules(multiStatusModules, 'failed')
const defaultChecklistRows: RegistryManualChecklistRow[] =
    buildRegistryManualChecklistRows(defaultRegistryState)
const brokenChecklistRows: RegistryManualChecklistRow[] =
    buildRegistryManualChecklistRows(brokenRegistryState)
const multiChecklistRows: RegistryManualChecklistRow[] =
    buildRegistryManualChecklistRows(multiRegistryState)
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
    multiStatusSummaryRows,
    marketRows,
    selectedWizardRows,
    emptyWizardRows,
    multiAllModules,
    multiUninstalledModules,
    multiFailedModules,
    defaultChecklistRows,
    brokenChecklistRows,
    multiChecklistRows,
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
