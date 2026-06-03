import {
    buildRegistryAcceptanceRows,
    countRegistryFailures,
    hasBrokenRegistryFixture,
    isRegistryEmptyState,
    registryErrorDetailText,
    registryTableEmptyTextFromState,
    type RegistryAcceptanceRow,
    type RegistryStateInput,
    type RegistryStateModule
} from './registry-state'

const defaultModules = [
    {
        module: 'article',
        registryStatusCode: 'passed'
    }
] satisfies RegistryStateModule[]

const brokenFixtureModules = [
    ...defaultModules,
    {
        module: 'broken_fixture',
        registryStatusCode: 'failed'
    }
] satisfies RegistryStateModule[]

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
const defaultFailureCount: number = countRegistryFailures(defaultModules)
const brokenFailureCount: number = countRegistryFailures(brokenFixtureModules)
const brokenEnabled: boolean = hasBrokenRegistryFixture(brokenFixtureModules)
const emptyStateMatched: boolean = isRegistryEmptyState(emptyRegistryState)
const emptyText: string = registryTableEmptyTextFromState(emptyRegistryState)
const failedText: string = registryTableEmptyTextFromState(failedRegistryState)
const failedDetail: string = registryErrorDetailText(failedRegistryState.registryError)

export const registryStateFixture = {
    defaultRows,
    brokenRows,
    defaultFailureCount,
    brokenFailureCount,
    brokenEnabled,
    emptyStateMatched,
    emptyText,
    failedText,
    failedDetail
}
