<template>
    <div class="module-manifest-preview">
        <popup
            ref="popupRef"
            :clickModalClose="false"
            title="Manifest 预览"
            width="960px"
            :async="true"
            confirmButtonText="生成预览"
            @confirm="handlePreview"
        >
            <template #trigger>
                <slot>
                    <el-button>
                        <template #icon>
                            <icon name="el-icon-Document" />
                        </template>
                        Manifest 预览
                    </el-button>
                </slot>
            </template>

            <el-form class="ls-form" :model="formData" label-width="90px">
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
                        :autosize="{ minRows: 10, maxRows: 16 }"
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

            <div v-if="preview" class="manifest-result">
                <el-descriptions :column="3" border>
                    <el-descriptions-item label="来源">
                        {{ preview.source }}
                    </el-descriptions-item>
                    <el-descriptions-item label="模块">
                        {{ preview.manifest.module }}
                    </el-descriptions-item>
                    <el-descriptions-item label="实体">
                        {{ preview.manifest.entity }}
                    </el-descriptions-item>
                    <el-descriptions-item label="表名">
                        {{ preview.detail.base.tableName }}
                    </el-descriptions-item>
                    <el-descriptions-item label="功能">
                        {{ preview.detail.gen.functionName }}
                    </el-descriptions-item>
                    <el-descriptions-item label="模板">
                        {{ preview.detail.gen.genTpl }}
                    </el-descriptions-item>
                    <el-descriptions-item label="租户">
                        {{ preview.plan.tenantId }}
                    </el-descriptions-item>
                    <el-descriptions-item label="角色">
                        {{ preview.plan.roleId }}
                    </el-descriptions-item>
                    <el-descriptions-item label="运行时">
                        {{ preview.plan.runtimeHint }}
                    </el-descriptions-item>
                </el-descriptions>

                <el-form class="install-gate-form" :model="confirmData" label-width="90px">
                    <el-form-item label="确认模块">
                        <el-input
                            class="w-[280px]"
                            v-model="confirmData.confirmModule"
                            clearable
                        />
                    </el-form-item>
                    <el-form-item label="写入确认">
                        <el-checkbox v-model="confirmData.confirmInstall">安装写入</el-checkbox>
                        <el-checkbox
                            v-model="confirmData.confirmSchemaRisk"
                            :disabled="!preview.manifest.requiresSchema"
                        >
                            Schema 风险
                        </el-checkbox>
                        <el-checkbox v-model="confirmData.confirmDelete">删除确认</el-checkbox>
                    </el-form-item>
                </el-form>

                <div v-if="installResult || uninstallResult" class="apply-result">
                    <el-tabs v-model="resultTab">
                        <el-tab-pane v-if="installResult" label="安装结果" name="install">
                            <el-alert
                                :title="resultTitle(installResult, '安装写入已阻断')"
                                :type="resultAlertType(installResult)"
                                show-icon
                                :closable="false"
                            />
                            <el-descriptions class="mt-3" :column="4" border>
                                <el-descriptions-item label="状态">
                                    {{ installResult.status || '-' }}
                                </el-descriptions-item>
                                <el-descriptions-item label="模块">
                                    {{ installResult.manifest?.module || '-' }}
                                </el-descriptions-item>
                                <el-descriptions-item label="来源">
                                    {{ installResult.source || '-' }}
                                </el-descriptions-item>
                                <el-descriptions-item label="环境变量">
                                    {{ installResult.requiredEnv || '-' }}
                                </el-descriptions-item>
                            </el-descriptions>
                            <el-table
                                v-if="hasSnapshot(installResult)"
                                class="mt-3"
                                :data="snapshotRows(installResult)"
                                size="large"
                            >
                                <el-table-column label="对象" prop="name" min-width="130" />
                                <el-table-column label="执行前" prop="before" min-width="100" />
                                <el-table-column label="执行后" prop="after" min-width="100" />
                            </el-table>
                            <el-table
                                v-if="installResult.checks?.length"
                                class="mt-3"
                                :data="installResult.checks"
                                size="large"
                            >
                                <el-table-column label="检查项" prop="name" min-width="140" />
                                <el-table-column label="状态" prop="status" min-width="120" />
                                <el-table-column label="说明" prop="message" min-width="280" />
                            </el-table>
                        </el-tab-pane>
                        <el-tab-pane v-if="uninstallResult" label="卸载结果" name="uninstall">
                            <el-alert
                                :title="resultTitle(uninstallResult, '卸载写入已阻断')"
                                :type="resultAlertType(uninstallResult)"
                                show-icon
                                :closable="false"
                            />
                            <el-descriptions class="mt-3" :column="4" border>
                                <el-descriptions-item label="状态">
                                    {{ uninstallResult.status || '-' }}
                                </el-descriptions-item>
                                <el-descriptions-item label="模块">
                                    {{ uninstallResult.manifest?.module || '-' }}
                                </el-descriptions-item>
                                <el-descriptions-item label="来源">
                                    {{ uninstallResult.source || '-' }}
                                </el-descriptions-item>
                                <el-descriptions-item label="环境变量">
                                    {{ uninstallResult.requiredEnv || '-' }}
                                </el-descriptions-item>
                            </el-descriptions>
                            <el-table
                                v-if="hasSnapshot(uninstallResult)"
                                class="mt-3"
                                :data="snapshotRows(uninstallResult)"
                                size="large"
                            >
                                <el-table-column label="对象" prop="name" min-width="130" />
                                <el-table-column label="执行前" prop="before" min-width="100" />
                                <el-table-column label="执行后" prop="after" min-width="100" />
                            </el-table>
                            <el-table
                                v-if="uninstallResult.checks?.length"
                                class="mt-3"
                                :data="uninstallResult.checks"
                                size="large"
                            >
                                <el-table-column label="检查项" prop="name" min-width="140" />
                                <el-table-column label="状态" prop="status" min-width="120" />
                                <el-table-column label="说明" prop="message" min-width="280" />
                            </el-table>
                        </el-tab-pane>
                    </el-tabs>
                </div>

                <el-table class="mt-4" :data="preview.detail.column" size="large" height="260">
                    <el-table-column label="字段" prop="columnName" min-width="130" />
                    <el-table-column label="Go 字段" prop="goField" min-width="120" />
                    <el-table-column label="Go 类型" prop="goType" min-width="100" />
                    <el-table-column label="表单" prop="htmlType" min-width="100" />
                    <el-table-column label="查询" prop="queryType" min-width="100" />
                    <el-table-column label="字典" prop="dictType" min-width="120" />
                </el-table>

                <div class="flex justify-end mt-4">
                    <el-button @click="handlePlanPreview">
                        <template #icon>
                            <icon name="el-icon-DocumentCopy" />
                        </template>
                        安装计划
                    </el-button>
                    <el-button type="warning" @click="handleInstallGate">
                        <template #icon>
                            <icon name="el-icon-Lock" />
                        </template>
                        安装执行
                    </el-button>
                    <el-button type="danger" @click="handleUninstallGate">
                        <template #icon>
                            <icon name="el-icon-Delete" />
                        </template>
                        卸载执行
                    </el-button>
                    <el-button type="primary" @click="handleCodePreview">
                        <template #icon>
                            <icon name="el-icon-View" />
                        </template>
                        代码预览
                    </el-button>
                </div>
            </div>
        </popup>

        <code-preview
            v-if="previewState.show"
            v-model="previewState.show"
            :code="previewState.code"
        />
    </div>
</template>

<script lang="ts" setup>
import Popup from '@/components/popup/index.vue'
import CodePreview from './code-preview.vue'
import {
    applyModuleManifestInstall,
    applyModuleManifestUninstall,
    previewModuleManifest
} from '@/api/tools/code'
import feedback from '@/utils/feedback'

const popupRef = shallowRef<InstanceType<typeof Popup>>()
const inputMode = ref<'path' | 'body'>('path')
const formData = reactive({
    manifestPath: 'examples/demo/manifest.json',
    manifestBody: '',
    authorName: 'codepeep',
    tenantId: 0,
    roleId: 1
})

const preview = ref<any>()
const installResult = ref<any>()
const uninstallResult = ref<any>()
const resultTab = ref<'install' | 'uninstall'>('install')
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

const manifestParams = () =>
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

const handlePreview = async () => {
    const params = manifestParams()
    preview.value = await previewModuleManifest(params)
    confirmData.confirmModule = preview.value?.manifest?.module || ''
    confirmData.confirmInstall = false
    confirmData.confirmSchemaRisk = false
    confirmData.confirmDelete = false
    installResult.value = undefined
    uninstallResult.value = undefined
    feedback.msgSuccess('预览生成成功')
}

const handleCodePreview = () => {
    previewState.code = preview.value?.code || {}
    previewState.show = true
}

const handlePlanPreview = () => {
    const plan = preview.value?.plan || {}
    previewState.code = {
        'registry.sql': plan.registrySql || '',
        'role_grant.sql': plan.roleGrantSql || '',
        'install.sql': plan.installSql || '',
        'uninstall.sql': plan.uninstallSql || ''
    }
    previewState.show = true
}

const handleInstallGate = async () => {
    const params = {
        ...manifestParams(),
        confirmModule: confirmData.confirmModule,
        confirmTenantId: formData.tenantId,
        confirmRoleId: formData.roleId,
        confirmInstall: confirmData.confirmInstall,
        confirmSchemaRisk: confirmData.confirmSchemaRisk
    }
    try {
        installResult.value = await applyModuleManifestInstall(params)
        feedback.msgSuccess('安装执行完成')
    } catch (error) {
        installResult.value = error || {
            message: '安装写入已阻断',
            checks: []
        }
    }
    resultTab.value = 'install'
}

const handleUninstallGate = async () => {
    const params = {
        ...manifestParams(),
        confirmModule: confirmData.confirmModule,
        confirmDelete: confirmData.confirmDelete
    }
    try {
        uninstallResult.value = await applyModuleManifestUninstall(params)
        feedback.msgSuccess('卸载执行完成')
    } catch (error) {
        uninstallResult.value = error || {
            message: '卸载写入已阻断',
            checks: []
        }
    }
    resultTab.value = 'uninstall'
}

const resultTitle = (result: any, fallback: string) => result?.message || fallback

const resultAlertType = (result: any) => (result?.status === 'applied' ? 'success' : 'warning')

const hasSnapshot = (result: any) =>
    result?.status === 'applied' ||
    snapshotRows(result).some((row) => Number(row.before) > 0 || Number(row.after) > 0)

const snapshotRows = (result: any) => {
    const before = result?.before || {}
    const after = result?.after || {}
    return [
        {
            name: '权限',
            before: before.permissions || 0,
            after: after.permissions || 0
        },
        {
            name: '菜单',
            before: before.menus || 0,
            after: after.menus || 0
        },
        {
            name: '菜单权限',
            before: before.menuPermissions || 0,
            after: after.menuPermissions || 0
        },
        {
            name: '角色授权',
            before: before.rolePermissions || 0,
            after: after.rolePermissions || 0
        }
    ]
}
</script>

<style scoped lang="scss">
.module-manifest-preview {
    display: inline-block;
}

.manifest-result {
    margin-top: 16px;
}

.install-gate-form,
.apply-result {
    margin-top: 16px;
}
</style>
