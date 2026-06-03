<template>
    <div class="module-center">
        <el-card class="!border-none mb-4" shadow="never">
            <div class="module-hero">
                <div>
                    <div class="module-title">模块中心</div>
                    <div class="module-subtitle">
                        manifest、codegen、安装计划、安装执行、卸载执行和审计预览的统一入口
                    </div>
                </div>
                <el-button type="primary" :loading="previewLoading" @click="handlePreview">
                    <template #icon>
                        <icon name="el-icon-Document" />
                    </template>
                    生成预览
                </el-button>
            </div>
        </el-card>

        <el-card class="!border-none mb-4" shadow="never">
            <el-form class="module-form" :model="formData" label-width="90px">
                <el-form-item label="来源">
                    <el-radio-group v-model="inputMode">
                        <el-radio-button label="path">仓库路径</el-radio-button>
                        <el-radio-button label="body">JSON</el-radio-button>
                    </el-radio-group>
                </el-form-item>
                <el-form-item v-if="inputMode === 'path'" label="路径">
                    <el-input v-model="formData.manifestPath" clearable />
                </el-form-item>
                <el-form-item v-else label="JSON">
                    <el-input
                        v-model="formData.manifestBody"
                        type="textarea"
                        :autosize="{ minRows: 8, maxRows: 14 }"
                        clearable
                    />
                </el-form-item>
                <el-form-item label="作者">
                    <el-input class="w-[280px]" v-model="formData.authorName" clearable />
                </el-form-item>
                <el-form-item label="租户/角色">
                    <div class="flex gap-3">
                        <el-input-number
                            class="w-[160px]"
                            v-model="formData.tenantId"
                            :min="0"
                            :controls="false"
                        />
                        <el-input-number
                            class="w-[160px]"
                            v-model="formData.roleId"
                            :min="1"
                            :controls="false"
                        />
                    </div>
                </el-form-item>
            </el-form>
        </el-card>

        <div class="module-grid mb-4">
            <el-card v-for="item in capabilityCards" :key="item.name" class="!border-none" shadow="never">
                <div class="module-card">
                    <div class="module-icon">
                        <icon :name="item.icon" :size="24" />
                    </div>
                    <div>
                        <div class="module-card-title">{{ item.name }}</div>
                        <div class="module-card-desc">{{ item.desc }}</div>
                    </div>
                </div>
            </el-card>
        </div>

        <el-card v-if="preview" class="!border-none mb-4" shadow="never">
            <template #header>
                <div class="section-header">
                    <span class="card-title">预览结果</span>
                    <el-tag type="success" size="small">{{ preview.manifest.module }}</el-tag>
                </div>
            </template>
            <el-descriptions :column="3" border>
                <el-descriptions-item label="来源">{{ preview.source }}</el-descriptions-item>
                <el-descriptions-item label="实体">{{ preview.manifest.entity }}</el-descriptions-item>
                <el-descriptions-item label="表名">{{ preview.detail.base.tableName }}</el-descriptions-item>
                <el-descriptions-item label="功能">
                    {{ preview.detail.gen.functionName }}
                </el-descriptions-item>
                <el-descriptions-item label="模板">{{ preview.detail.gen.genTpl }}</el-descriptions-item>
                <el-descriptions-item label="运行时">
                    <span class="wrap-text">{{ preview.plan.runtimeHint }}</span>
                </el-descriptions-item>
            </el-descriptions>

            <div class="section-label">模块状态</div>
            <el-descriptions :column="2" border>
                <el-descriptions-item label="预览">
                    <el-tag :type="previewStatusType" size="small">{{ previewStatus }}</el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="安装">
                    <el-tag :type="installStatusType" size="small">{{ installStatus }}</el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="卸载">
                    <el-tag :type="uninstallStatusType" size="small">{{ uninstallStatus }}</el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="写入门禁" :span="2">
                    <div class="gate-tags">
                        <el-tag v-if="!writeGateEnvList.length" :type="writeGateStatusType" size="small">
                            {{ writeGateStatus }}
                        </el-tag>
                        <template v-else>
                            <el-tag
                                v-for="env in writeGateEnvList"
                                :key="env"
                                class="gate-tag"
                                type="warning"
                                size="small"
                            >
                                {{ env }}
                            </el-tag>
                        </template>
                    </div>
                </el-descriptions-item>
                <el-descriptions-item label="运行时" :span="2">
                    <span class="wrap-text">{{ preview.plan.runtimeHint }}</span>
                </el-descriptions-item>
            </el-descriptions>

            <el-form class="apply-form" :model="confirmData" label-width="90px">
                <el-form-item label="确认模块">
                    <el-input
                        class="w-[280px]"
                        v-model="confirmData.confirmModule"
                        :disabled="isApplyLoading"
                        clearable
                    />
                </el-form-item>
                <el-form-item label="写入确认">
                    <el-checkbox v-model="confirmData.confirmInstall" :disabled="isApplyLoading">
                        安装写入
                    </el-checkbox>
                    <el-checkbox
                        v-model="confirmData.confirmSchemaRisk"
                        :disabled="!requiresSchemaConfirm || isApplyLoading"
                    >
                        Schema 风险
                    </el-checkbox>
                    <el-checkbox v-model="confirmData.confirmDelete" :disabled="isApplyLoading">
                        删除确认
                    </el-checkbox>
                </el-form-item>
            </el-form>

            <div class="preview-actions">
                <el-button :disabled="!hasCurrentPreview || isApplyLoading" @click="handlePlanPreview">
                    <template #icon>
                        <icon name="el-icon-DocumentCopy" />
                    </template>
                    安装计划
                </el-button>
                <el-button
                    v-perms="['gen:previewCode']"
                    type="warning"
                    :disabled="!canInstallApply"
                    :loading="installApplyLoading"
                    @click="handleInstallApply"
                >
                    <template #icon>
                        <icon name="el-icon-Lock" />
                    </template>
                    安装执行
                </el-button>
                <el-button
                    v-perms="['gen:previewCode']"
                    type="danger"
                    :disabled="!canUninstallApply"
                    :loading="uninstallApplyLoading"
                    @click="handleUninstallApply"
                >
                    <template #icon>
                        <icon name="el-icon-Delete" />
                    </template>
                    卸载执行
                </el-button>
                <el-button
                    type="primary"
                    :disabled="!hasCurrentPreview || isApplyLoading"
                    @click="handleCodePreview"
                >
                    <template #icon>
                        <icon name="el-icon-View" />
                    </template>
                    代码预览
                </el-button>
            </div>

            <div v-if="installResult || uninstallResult" class="apply-result">
                <el-tabs v-model="resultTab">
                    <el-tab-pane v-if="installResult" label="安装结果" name="install">
                        <module-manifest-apply-result-view
                            :result="installResult"
                            fallback-title="安装写入已阻断"
                        />
                    </el-tab-pane>
                    <el-tab-pane v-if="uninstallResult" label="卸载结果" name="uninstall">
                        <module-manifest-apply-result-view
                            :result="uninstallResult"
                            fallback-title="卸载写入已阻断"
                        />
                    </el-tab-pane>
                </el-tabs>
            </div>

            <div class="section-label">人工测试清单</div>
            <el-table :data="testChecklistRows" size="large">
                <el-table-column label="项目" prop="name" min-width="130" />
                <el-table-column label="状态" min-width="120">
                    <template #default="{ row }">
                        <el-tag :type="row.statusType" size="small">{{ row.status }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column label="结果" prop="detail" min-width="280" />
            </el-table>

            <el-table class="mt-4" :data="preview.detail.column" size="large">
                <el-table-column label="字段" prop="columnName" min-width="130" />
                <el-table-column label="Go 字段" prop="goField" min-width="120" />
                <el-table-column label="Go 类型" prop="goType" min-width="100" />
                <el-table-column label="表单" prop="htmlType" min-width="100" />
                <el-table-column label="查询" prop="queryType" min-width="100" />
                <el-table-column label="字典" prop="dictType" min-width="120" />
            </el-table>
        </el-card>

        <el-card class="!border-none" shadow="never">
            <template #header>
                <div class="section-header">
                    <span class="card-title">内置模块清单</span>
                    <div class="section-actions">
                        <el-tag type="success" size="small">P5.9</el-tag>
                        <el-button
                            type="primary"
                            link
                            :loading="registryLoading || statusLoading"
                            @click="refreshModuleCenter"
                        >
                            刷新状态
                        </el-button>
                    </div>
                </div>
            </template>
            <div class="module-list-toolbar">
                <el-radio-group v-model="moduleStatusFilter" size="small">
                    <el-radio-button
                        v-for="item in moduleStatusFilterOptions"
                        :key="item.value"
                        :label="item.value"
                    >
                        {{ item.label }}
                    </el-radio-button>
                </el-radio-group>
                <div class="module-summary">
                    <div v-for="item in moduleStatusSummary" :key="item.key" class="module-summary-item">
                        <span>{{ item.label }}</span>
                        <strong>{{ item.value }}</strong>
                    </div>
                </div>
            </div>
            <el-table
                v-loading="registryLoading || statusLoading"
                :data="filteredModules"
                size="large"
                empty-text="暂无匹配模块"
            >
                <el-table-column label="模块" prop="name" min-width="130" />
                <el-table-column label="Manifest" prop="manifest" min-width="240" />
                <el-table-column label="校验" min-width="180">
                    <template #default="{ row }">
                        <div class="status-stack">
                            <el-tag :type="row.registryStatusType" size="small">
                                {{ row.registryStatus }}
                            </el-tag>
                            <el-button
                                type="primary"
                                link
                                :disabled="!row.registryChecks.length"
                                @click="handleRegistryCheckDetail(row)"
                            >
                                <template #icon>
                                    <icon name="el-icon-DocumentChecked" />
                                </template>
                                明细
                            </el-button>
                            <span class="status-detail">{{ row.registryDetail }}</span>
                        </div>
                    </template>
                </el-table-column>
                <el-table-column label="表名" prop="table" min-width="160" />
                <el-table-column label="运行时" prop="runtime" min-width="260" />
                <el-table-column label="页面" prop="entry" min-width="160" />
                <el-table-column label="安装" min-width="180">
                    <template #default="{ row }">
                        <div class="status-stack">
                            <el-tag :type="row.installStatusType" size="small">
                                {{ row.installStatus }}
                            </el-tag>
                            <span class="status-detail">{{ row.statusDetail }}</span>
                        </div>
                    </template>
                </el-table-column>
                <el-table-column label="快照" min-width="220">
                    <template #default="{ row }">
                        <span class="wrap-text">{{ row.snapshotText }}</span>
                    </template>
                </el-table-column>
                <el-table-column label="运行时状态" min-width="220">
                    <template #default="{ row }">
                        <div class="status-stack">
                            <el-tag :type="row.runtimeStatusType" size="small">
                                {{ row.runtimeStatus }}
                            </el-tag>
                            <span class="status-detail">{{ row.runtimeDetail }}</span>
                        </div>
                    </template>
                </el-table-column>
                <el-table-column label="状态" width="120">
                    <template #default="{ row }">
                        <el-tag :type="row.statusType" size="small">{{ row.status }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column label="入口" width="220" fixed="right">
                    <template #default="{ row }">
                        <el-button type="primary" link @click="handleModulePreview(row.manifest)">
                            <template #icon>
                                <icon name="el-icon-View" />
                            </template>
                            预览
                        </el-button>
                        <el-button type="primary" link @click="goTo(row.entry)">打开</el-button>
                    </template>
                </el-table-column>
            </el-table>
        </el-card>

        <el-dialog
            v-model="registryCheckDialog.show"
            width="760px"
            title="Manifest 校验明细"
        >
            <el-descriptions :column="2" border>
                <el-descriptions-item label="模块">
                    {{ registryCheckDialog.row?.name || '-' }}
                </el-descriptions-item>
                <el-descriptions-item label="状态">
                    <el-tag
                        :type="registryDialogStatusType"
                        size="small"
                    >
                        {{ registryCheckDialog.row?.registryStatus || '-' }}
                    </el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="Manifest">
                    <span class="wrap-text">{{ registryCheckDialog.row?.manifest || '-' }}</span>
                </el-descriptions-item>
                <el-descriptions-item label="说明">
                    <span class="wrap-text">{{ registryCheckDialog.row?.registryDetail || '-' }}</span>
                </el-descriptions-item>
            </el-descriptions>
            <el-table
                class="mt-3"
                :data="registryCheckRows"
                size="large"
                empty-text="暂无检查项"
            >
                <el-table-column label="检查项" prop="name" min-width="140" />
                <el-table-column label="状态" min-width="120">
                    <template #default="{ row }">
                        <el-tag :type="registryCheckStatusType(row.status)" size="small">
                            {{ registryCheckStatusLabel(row.status) }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column label="说明" prop="message" min-width="320" />
            </el-table>
        </el-dialog>

        <code-preview
            v-if="previewState.show"
            v-model="previewState.show"
            :code="previewState.code"
        />
    </div>
</template>

<script lang="ts" setup name="moduleCenter">
import {
    applyModuleManifestInstall,
    applyModuleManifestUninstall,
    listModuleRegistry,
    normalizeModuleManifestApplyError,
    previewModuleManifest,
    readModuleManifestInstallStatus,
    type ModuleManifestApplyResult,
    type ModuleManifestApplySnapshotResult,
    type ModuleManifestApplyCheckResult,
    type ModuleManifestInstallApplyParams,
    type ModuleManifestInstallStatusResult,
    type ModuleManifestPreviewParams,
    type ModuleManifestPreviewResult,
    type ModuleManifestUninstallApplyParams,
    type ModuleRegistryItemResult
} from '@/api/tools/code'
import CodePreview from '../components/code-preview.vue'
import ModuleManifestApplyResultView from '../components/module-manifest-apply-result.vue'
import feedback from '@/utils/feedback'

const router = useRouter()
const inputMode = ref<'path' | 'body'>('path')
const formData = reactive({
    manifestPath: 'examples/demo/manifest.json',
    manifestBody: '',
    authorName: 'codepeep',
    tenantId: 0,
    roleId: 1
})

const preview = ref<ModuleManifestPreviewResult>()
const previewSnapshotKey = ref('')
const installResult = ref<ModuleManifestApplyResult>()
const uninstallResult = ref<ModuleManifestApplyResult>()
const resultTab = ref<'install' | 'uninstall'>('install')
const previewLoading = ref(false)
const registryLoading = ref(false)
const statusLoading = ref(false)
const installApplyLoading = ref(false)
const uninstallApplyLoading = ref(false)
const planPreviewOpened = ref(false)
const codePreviewOpened = ref(false)
const confirmData = reactive({
    confirmModule: '',
    confirmInstall: false,
    confirmSchemaRisk: false,
    confirmDelete: false
})
const previewState = reactive({
    show: false,
    code: {} as Record<string, string>
})
const moduleStatusFilter = ref('all')

const capabilityCards = [
    {
        name: 'Manifest',
        desc: '从仓库路径或 JSON 读取模块清单并生成结构化预览。',
        icon: 'el-icon-Document'
    },
    {
        name: 'Codegen',
        desc: '把模块清单转换为 ma_codegen_table 和 ma_codegen_column 预览。',
        icon: 'el-icon-Tools'
    },
    {
        name: 'Install',
        desc: '查看安装计划，并在本地门禁满足时执行安装 apply。',
        icon: 'el-icon-CirclePlus'
    },
    {
        name: 'Audit',
        desc: '在 apply 结果中查看摘要和审计 dry-run 预览。',
        icon: 'el-icon-DocumentChecked'
    }
]

type ModuleCenterModule = {
    name: string
    manifest: string
    table: string
    runtime: string
    entry: string
    status: string
    statusType: string
    registryStatusCode: string
    registryStatus: string
    registryStatusType: string
    registryDetail: string
    registryChecks: ModuleManifestApplyCheckResult[]
    installStatusCode: string
    installStatus: string
    installStatusType: string
    statusDetail: string
    snapshotText: string
    runtimeStatus: string
    runtimeStatusType: string
    runtimeDetail: string
}

type ElementTagType = 'primary' | 'success' | 'warning' | 'danger' | 'info'

const modules = reactive<ModuleCenterModule[]>([])
const registryCheckDialog = reactive<{
    show: boolean
    row?: ModuleCenterModule
}>({
    show: false
})

const moduleStatusFilterOptions = [
    { label: '全部', value: 'all' },
    { label: '已安装', value: 'installed' },
    { label: '部分安装', value: 'partial' },
    { label: '未安装', value: 'uninstalled' },
    { label: '异常', value: 'failed' }
]

const manifestParams = (): ModuleManifestPreviewParams =>
    inputMode.value === 'path'
        ? {
              manifestPath: formData.manifestPath,
              authorName: formData.authorName,
              tenantId: formData.tenantId,
              roleId: formData.roleId
          }
        : {
              manifestBody: formData.manifestBody,
              authorName: formData.authorName,
              tenantId: formData.tenantId,
              roleId: formData.roleId
          }

const manifestInputKey = computed(() =>
    JSON.stringify({
        inputMode: inputMode.value,
        ...manifestParams()
    })
)

const hasCurrentPreview = computed(
    () => Boolean(preview.value) && previewSnapshotKey.value === manifestInputKey.value
)

const requiresSchemaConfirm = computed(() => Boolean(preview.value?.manifest?.requiresSchema))

const isApplyLoading = computed(
    () => previewLoading.value || installApplyLoading.value || uninstallApplyLoading.value
)

const expectedModule = computed(() => preview.value?.manifest?.module || '')

const isConfirmModuleMatched = computed(
    () => Boolean(expectedModule.value) && confirmData.confirmModule === expectedModule.value
)

const canInstallApply = computed(
    () =>
        hasCurrentPreview.value &&
        isConfirmModuleMatched.value &&
        confirmData.confirmInstall &&
        (!requiresSchemaConfirm.value || confirmData.confirmSchemaRisk) &&
        !isApplyLoading.value
)

const canUninstallApply = computed(
    () =>
        hasCurrentPreview.value &&
        isConfirmModuleMatched.value &&
        confirmData.confirmDelete &&
        !isApplyLoading.value
)

const applyStatusText = (result?: ModuleManifestApplyResult) => result?.status || '未执行'

const applyStatusType = (status?: string) => {
    if (status === 'applied') {
        return 'success'
    }
    if (status === 'blocked') {
        return 'warning'
    }
    return 'info'
}

const previewStatus = computed(() => (hasCurrentPreview.value ? '已生成' : '未生成'))
const previewStatusType = computed(() => (hasCurrentPreview.value ? 'success' : 'info'))
const installStatus = computed(() => applyStatusText(installResult.value))
const installStatusType = computed(() => applyStatusType(installResult.value?.status))
const uninstallStatus = computed(() => applyStatusText(uninstallResult.value))
const uninstallStatusType = computed(() => applyStatusType(uninstallResult.value?.status))

const writeGateEnvList = computed(() =>
    Array.from(
        new Set([installResult.value?.requiredEnv, uninstallResult.value?.requiredEnv].filter(Boolean))
    )
)

const writeGateStatus = computed(() => {
    if (writeGateEnvList.value.length) {
        return '门禁阻断'
    }
    if (installResult.value?.status === 'applied' || uninstallResult.value?.status === 'applied') {
        return '已执行本地写入'
    }
    return '未执行'
})

const writeGateStatusType = computed(() => {
    if (installResult.value?.requiredEnv || uninstallResult.value?.requiredEnv) {
        return 'warning'
    }
    if (installResult.value?.status === 'applied' || uninstallResult.value?.status === 'applied') {
        return 'success'
    }
    return 'info'
})

const filteredModules = computed(() => {
    if (moduleStatusFilter.value === 'all') {
        return modules
    }
    if (moduleStatusFilter.value === 'failed') {
        return modules.filter(
            (item) =>
                item.registryStatusCode === 'failed' ||
                ['blocked', 'failed'].includes(item.installStatusCode)
        )
    }
    return modules.filter((item) => item.installStatusCode === moduleStatusFilter.value)
})

const moduleStatusSummary = computed(() => {
    const countBy = (codes: string[]) =>
        modules.filter((item) => codes.includes(item.installStatusCode)).length
    const failedCount = modules.filter(
        (item) =>
            item.registryStatusCode === 'failed' ||
            ['blocked', 'failed'].includes(item.installStatusCode)
    ).length
    return [
        { key: 'total', label: '总数', value: modules.length },
        { key: 'installed', label: '已安装', value: countBy(['installed']) },
        { key: 'partial', label: '部分', value: countBy(['partial']) },
        { key: 'uninstalled', label: '未安装', value: countBy(['uninstalled']) },
        { key: 'failed', label: '异常', value: failedCount }
    ]
})

const registryCheckRows = computed(() => registryCheckDialog.row?.registryChecks || [])
const registryDialogStatusType = computed<ElementTagType>(() =>
    moduleRegistryStatusType(registryCheckDialog.row?.registryStatusCode)
)

const testChecklistRows = computed(() => [
    {
        name: 'Manifest 预览',
        status: hasCurrentPreview.value ? '通过' : '待执行',
        statusType: hasCurrentPreview.value ? 'success' : 'info',
        detail: preview.value?.source || '-'
    },
    {
        name: '安装计划',
        status: planPreviewOpened.value ? '已打开' : '待打开',
        statusType: planPreviewOpened.value ? 'success' : 'info',
        detail: preview.value?.plan?.runtimeHint || '-'
    },
    {
        name: '代码预览',
        status: codePreviewOpened.value ? '已打开' : '待打开',
        statusType: codePreviewOpened.value ? 'success' : 'info',
        detail: preview.value ? `${Object.keys(preview.value.code || {}).length} files` : '-'
    },
    {
        name: '安装执行',
        status: installStatus.value,
        statusType: installStatusType.value,
        detail: installResult.value?.message || installResult.value?.requiredEnv || '-'
    },
    {
        name: '卸载执行',
        status: uninstallStatus.value,
        statusType: uninstallStatusType.value,
        detail: uninstallResult.value?.message || uninstallResult.value?.requiredEnv || '-'
    },
    {
        name: '审计预览',
        status: installResult.value || uninstallResult.value ? '可展开' : '待结果',
        statusType: installResult.value || uninstallResult.value ? 'success' : 'info',
        detail: installResult.value || uninstallResult.value ? 'apply result' : '-'
    }
])

const resetConfirmState = (module = '') => {
    confirmData.confirmModule = module
    confirmData.confirmInstall = false
    confirmData.confirmSchemaRisk = false
    confirmData.confirmDelete = false
}

const clearApplyResults = () => {
    installResult.value = undefined
    uninstallResult.value = undefined
    resultTab.value = 'install'
}

const clearTestState = () => {
    planPreviewOpened.value = false
    codePreviewOpened.value = false
    clearApplyResults()
}

const clearPreviewState = () => {
    preview.value = undefined
    previewSnapshotKey.value = ''
    resetConfirmState()
    clearTestState()
}

watch(manifestInputKey, () => {
    if (preview.value) {
        clearPreviewState()
    }
})

const handlePreview = async () => {
    if (previewLoading.value) {
        return
    }
    const params = manifestParams()
    const snapshotKey = manifestInputKey.value
    previewLoading.value = true
    try {
        const data = await previewModuleManifest(params)
        if (snapshotKey !== manifestInputKey.value) {
            return
        }
        preview.value = data
        previewSnapshotKey.value = snapshotKey
        resetConfirmState(preview.value?.manifest?.module || '')
        clearTestState()
        feedback.msgSuccess('预览生成成功')
    } finally {
        previewLoading.value = false
    }
}

const handleModulePreview = async (manifestPath: string) => {
    inputMode.value = 'path'
    formData.manifestPath = manifestPath
    await handlePreview()
}

const goTo = (url: string) => {
    router.push(url)
}

const toModuleCenterModule = (item: ModuleRegistryItemResult): ModuleCenterModule => ({
    name: item.name,
    manifest: item.manifest,
    table: item.table,
    runtime: item.runtime,
    entry: item.entry,
    status: item.status,
    statusType: item.statusType,
    registryStatusCode: item.manifestStatus || 'failed',
    registryStatus: moduleRegistryStatusLabel(item.manifestStatus),
    registryStatusType: moduleRegistryStatusType(item.manifestStatus),
    registryDetail: item.manifestMessage || '-',
    registryChecks: item.manifestChecks || [],
    installStatusCode: 'loading',
    installStatus: '读取中',
    installStatusType: 'info',
    statusDetail: '-',
    snapshotText: '-',
    runtimeStatus: '读取中',
    runtimeStatusType: 'info',
    runtimeDetail: '-'
})

const handleRegistryCheckDetail = (row: Record<string, unknown>) => {
    registryCheckDialog.row = row as ModuleCenterModule
    registryCheckDialog.show = true
}

const moduleRegistryStatusLabel = (status?: string) => {
    if (status === 'passed') {
        return '已通过'
    }
    if (status === 'failed') {
        return '异常'
    }
    return '未知'
}

const moduleRegistryStatusType = (status?: string): ElementTagType => {
    if (status === 'passed') {
        return 'success'
    }
    if (status === 'failed') {
        return 'danger'
    }
    return 'info'
}

const registryCheckStatusLabel = (status?: string) => {
    if (status === 'passed') {
        return '通过'
    }
    if (status === 'failed') {
        return '异常'
    }
    if (status === 'blocked') {
        return '阻断'
    }
    return status || '未知'
}

const registryCheckStatusType = (status?: string): ElementTagType => {
    if (status === 'passed') {
        return 'success'
    }
    if (status === 'failed') {
        return 'danger'
    }
    if (status === 'blocked') {
        return 'warning'
    }
    return 'info'
}

const statusTypeMap: Record<string, string> = {
    installed: 'success',
    partial: 'warning',
    uninstalled: 'info',
    blocked: 'danger',
    failed: 'danger'
}

const statusLabelMap: Record<string, string> = {
    installed: '已安装',
    partial: '部分安装',
    uninstalled: '未安装',
    blocked: '已阻断',
    failed: '读取失败'
}

const moduleSnapshotValue = (
    snapshot: ModuleManifestApplySnapshotResult | undefined,
    key: keyof ModuleManifestApplySnapshotResult
) => snapshot?.[key] || 0

const moduleSnapshotText = (status: ModuleManifestInstallStatusResult) => {
    const snapshot = status.snapshot || {}
    const expected = status.expected || {}
    return [
        `权限 ${moduleSnapshotValue(snapshot, 'permissions')}/${moduleSnapshotValue(expected, 'permissions')}`,
        `菜单 ${moduleSnapshotValue(snapshot, 'menus')}/${moduleSnapshotValue(expected, 'menus')}`,
        `菜单权限 ${moduleSnapshotValue(snapshot, 'menuPermissions')}/${moduleSnapshotValue(expected, 'menuPermissions')}`,
        `角色授权 ${moduleSnapshotValue(snapshot, 'rolePermissions')}/${moduleSnapshotValue(expected, 'rolePermissions')}`
    ].join(' · ')
}

const runtimeStatusFrom = (status: ModuleManifestInstallStatusResult) => {
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

const applyModuleStatusToRow = (row: ModuleCenterModule, status: ModuleManifestInstallStatusResult) => {
    const rawStatus = status.status || 'failed'
    const runtimeStatus = runtimeStatusFrom(status)
    row.installStatusCode = rawStatus
    row.installStatus = statusLabelMap[rawStatus] || rawStatus
    row.installStatusType = statusTypeMap[rawStatus] || 'info'
    row.statusDetail = status.message || '-'
    row.snapshotText = moduleSnapshotText(status)
    row.runtimeStatus = runtimeStatus.label
    row.runtimeStatusType = runtimeStatus.type
    row.runtimeDetail = runtimeStatus.detail
}

const moduleStatusErrorMessage = (error: unknown) => {
    if (error instanceof Error && error.message) {
        return error.message
    }
    if (typeof error === 'string' && error) {
        return error
    }
    return 'status api failed'
}

const loadModuleStatuses = async () => {
    if (statusLoading.value) {
        return
    }
    statusLoading.value = true
    try {
        await Promise.all(
            modules.map(async (item) => {
                try {
                    const status = await readModuleManifestInstallStatus({
                        manifestPath: item.manifest,
                        tenantId: formData.tenantId,
                        roleId: formData.roleId
                    })
                    applyModuleStatusToRow(item, status)
                } catch (error) {
                    item.installStatusCode = 'failed'
                    item.installStatus = '读取失败'
                    item.installStatusType = 'danger'
                    item.statusDetail = moduleStatusErrorMessage(error)
                    item.snapshotText = '-'
                    item.runtimeStatus = '未知'
                    item.runtimeStatusType = 'info'
                    item.runtimeDetail = '-'
                }
            })
        )
    } finally {
        statusLoading.value = false
    }
}

const loadModuleRegistry = async () => {
    if (registryLoading.value) {
        return
    }
    registryLoading.value = true
    try {
        const items = await listModuleRegistry()
        modules.splice(0, modules.length, ...items.map(toModuleCenterModule))
    } finally {
        registryLoading.value = false
    }
}

const refreshModuleCenter = async () => {
    await loadModuleRegistry()
    await loadModuleStatuses()
}

const handlePlanPreview = () => {
    const currentPreview = preview.value
    if (!hasCurrentPreview.value || !currentPreview) {
        return
    }
    const plan = currentPreview.plan
    previewState.code = {
        'registry.sql': plan.registrySql || '',
        'role_grant.sql': plan.roleGrantSql || '',
        'install.sql': plan.installSql || '',
        'uninstall.sql': plan.uninstallSql || ''
    }
    planPreviewOpened.value = true
    previewState.show = true
}

const handleCodePreview = () => {
    if (!hasCurrentPreview.value || !preview.value) {
        return
    }
    previewState.code = preview.value.code
    codePreviewOpened.value = true
    previewState.show = true
}

const handleInstallApply = async () => {
    if (!canInstallApply.value) {
        return
    }
    const params: ModuleManifestInstallApplyParams = {
        ...manifestParams(),
        confirmModule: confirmData.confirmModule,
        confirmTenantId: formData.tenantId,
        confirmRoleId: formData.roleId,
        confirmInstall: confirmData.confirmInstall,
        confirmSchemaRisk: confirmData.confirmSchemaRisk
    }
    installApplyLoading.value = true
    installResult.value = undefined
    try {
        installResult.value = await applyModuleManifestInstall(params)
        feedback.msgSuccess('安装执行完成')
    } catch (error) {
        installResult.value = normalizeModuleManifestApplyError(error, '安装写入已阻断')
    } finally {
        installApplyLoading.value = false
    }
    await loadModuleStatuses()
    resultTab.value = 'install'
}

const handleUninstallApply = async () => {
    if (!canUninstallApply.value) {
        return
    }
    const params: ModuleManifestUninstallApplyParams = {
        ...manifestParams(),
        confirmModule: confirmData.confirmModule,
        confirmDelete: confirmData.confirmDelete
    }
    uninstallApplyLoading.value = true
    uninstallResult.value = undefined
    try {
        uninstallResult.value = await applyModuleManifestUninstall(params)
        feedback.msgSuccess('卸载执行完成')
    } catch (error) {
        uninstallResult.value = normalizeModuleManifestApplyError(error, '卸载写入已阻断')
    } finally {
        uninstallApplyLoading.value = false
    }
    await loadModuleStatuses()
    resultTab.value = 'uninstall'
}

onMounted(() => {
    refreshModuleCenter()
})
</script>

<style lang="scss" scoped>
.module-center {
    .el-card {
        border-radius: 8px;
    }
}

.module-hero {
    align-items: center;
    display: flex;
    gap: 20px;
    justify-content: space-between;
}

.module-title {
    color: #111827;
    font-size: 24px;
    font-weight: 600;
    line-height: 32px;
}

.module-subtitle {
    color: #667085;
    font-size: 14px;
    line-height: 24px;
    margin-top: 6px;
}

.module-form {
    max-width: 720px;
}

.module-grid {
    display: grid;
    gap: 16px;
    grid-template-columns: repeat(4, minmax(0, 1fr));
}

.module-card {
    display: flex;
    gap: 14px;
    min-height: 90px;
}

.module-icon {
    align-items: center;
    background: #eef4ff;
    border-radius: 8px;
    color: #3b5bdb;
    display: flex;
    flex: 0 0 44px;
    height: 44px;
    justify-content: center;
}

.module-card-title {
    color: #111827;
    font-size: 16px;
    font-weight: 600;
    line-height: 24px;
}

.module-card-desc {
    color: #667085;
    font-size: 13px;
    line-height: 20px;
}

.section-header {
    align-items: center;
    display: flex;
    justify-content: space-between;
}

.section-actions {
    align-items: center;
    display: flex;
    gap: 10px;
}

.module-list-toolbar {
    align-items: center;
    display: flex;
    flex-wrap: wrap;
    gap: 14px;
    justify-content: space-between;
    margin-bottom: 14px;
}

.module-summary {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    justify-content: flex-end;
}

.module-summary-item {
    align-items: center;
    border: 1px solid #eaecf0;
    border-radius: 8px;
    color: #667085;
    display: inline-flex;
    font-size: 12px;
    gap: 8px;
    line-height: 18px;
    min-height: 30px;
    padding: 4px 10px;

    strong {
        color: #111827;
        font-size: 14px;
        font-weight: 600;
    }
}

.section-label {
    color: #111827;
    font-size: 15px;
    font-weight: 600;
    line-height: 22px;
    margin: 18px 0 10px;
}

.gate-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    min-width: 0;
}

.gate-tag {
    height: auto;
    max-width: 100%;
    padding: 4px 8px;

    :deep(.el-tag__content) {
        line-height: 18px;
        overflow-wrap: anywhere;
        white-space: normal;
    }
}

.wrap-text {
    overflow-wrap: anywhere;
    white-space: normal;
}

.status-stack {
    display: flex;
    flex-direction: column;
    gap: 6px;
    min-width: 0;
}

.status-detail {
    color: #667085;
    font-size: 12px;
    line-height: 18px;
    overflow-wrap: anywhere;
}

.module-center :deep(.el-table .cell) {
    overflow-wrap: anywhere;
}

.module-center :deep(.el-descriptions__label) {
    min-width: 76px;
    white-space: nowrap;
}

.preview-actions {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    justify-content: flex-end;
    margin-top: 16px;
}

.apply-form,
.apply-result {
    margin-top: 16px;
}

@media (max-width: 1024px) {
    .module-grid {
        grid-template-columns: repeat(2, minmax(0, 1fr));
    }
}

@media (max-width: 640px) {
    .module-hero {
        align-items: flex-start;
        flex-direction: column;
    }

    .module-grid {
        grid-template-columns: 1fr;
    }

    .preview-actions {
        justify-content: flex-start;
    }
}
</style>
