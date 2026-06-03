<template>
    <div class="module-manifest-apply-result">
        <el-alert
            :title="resultTitle"
            :type="resultAlertType"
            show-icon
            :closable="false"
        />
        <el-descriptions class="mt-3" :column="4" border>
            <el-descriptions-item label="状态">
                {{ result.status || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="模块">
                {{ result.manifest?.module || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="来源">
                {{ result.source || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="环境变量">
                {{ result.requiredEnv || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="操作">
                {{ result.summary?.operation || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="路由">
                {{ result.summary?.routeName || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="权限" :span="2">
                <div v-if="permissionCodes.length" class="permission-tags">
                    <el-tag v-for="code in permissionCodes" :key="code" size="small">
                        {{ code }}
                    </el-tag>
                </div>
                <span v-else class="apply-empty-text">无权限编码</span>
            </el-descriptions-item>
        </el-descriptions>
        <el-table v-if="hasSnapshot" class="mt-3" :data="snapshotRows" size="large">
            <el-table-column label="对象" prop="name" min-width="130" />
            <el-table-column label="执行前" prop="before" min-width="100" />
            <el-table-column label="执行后" prop="after" min-width="100" />
        </el-table>
        <div v-else class="apply-empty-text mt-3">无快照</div>
        <el-table v-if="result.checks?.length" class="mt-3" :data="result.checks" size="large">
            <el-table-column label="检查项" prop="name" min-width="140" />
            <el-table-column label="状态" prop="status" min-width="120" />
            <el-table-column label="说明" prop="message" min-width="280" />
        </el-table>
        <div v-else class="apply-empty-text mt-3">无检查项</div>
        <div class="audit-preview-toolbar">
            <el-button type="primary" link @click="toggleAuditPreview">
                <template #icon>
                    <icon name="el-icon-DocumentCopy" />
                </template>
                审计预览
            </el-button>
            <el-button v-if="auditPreviewVisible" type="primary" link @click="toggleAuditPreviewCode">
                JSON
            </el-button>
        </div>
        <el-descriptions v-if="auditPreviewVisible" class="mt-2" :column="4" border>
            <el-descriptions-item label="操作">
                {{ auditPreviewSummary.operation || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="模块">
                {{ auditPreviewSummary.module || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="状态">
                {{ auditPreviewSummary.status || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="路由">
                {{ auditPreviewSummary.routeName || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="权限">
                {{ auditPreviewSummary.permissionCount }}
            </el-descriptions-item>
            <el-descriptions-item label="检查">
                {{ auditPreviewSummary.checkCount }}
            </el-descriptions-item>
            <el-descriptions-item label="执行前">
                {{ auditPreviewSummary.beforeTotal }}
            </el-descriptions-item>
            <el-descriptions-item label="执行后">
                {{ auditPreviewSummary.afterTotal }}
            </el-descriptions-item>
        </el-descriptions>
        <pre v-if="auditPreviewVisible && auditPreviewCodeVisible" class="audit-preview-code">{{ auditPreviewCode }}</pre>
    </div>
</template>

<script lang="ts" setup>
import {
    buildModuleManifestApplyAuditPreview,
    buildModuleManifestApplyAuditPreviewSummary,
    type ModuleManifestApplyResult
} from '@/api/tools/code'

interface SnapshotRow {
    name: string
    before: number
    after: number
}

const props = defineProps<{
    result: ModuleManifestApplyResult
    fallbackTitle: string
}>()

const auditPreviewVisible = ref(false)
const auditPreviewCodeVisible = ref(false)

const resultTitle = computed(() => props.result.message || props.fallbackTitle)

const resultAlertType = computed(() => (props.result.status === 'applied' ? 'success' : 'warning'))

const permissionCodes = computed(() => props.result.summary?.permissionCodes || [])

const snapshotRows = computed<SnapshotRow[]>(() => {
    const before = props.result.before || {}
    const after = props.result.after || {}
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
})

const hasSnapshot = computed(
    () =>
        props.result.status === 'applied' ||
        snapshotRows.value.some((row) => Number(row.before) > 0 || Number(row.after) > 0)
)

const auditPreview = computed(() => buildModuleManifestApplyAuditPreview(props.result))

const auditPreviewSummary = computed(() =>
    buildModuleManifestApplyAuditPreviewSummary(auditPreview.value)
)

const auditPreviewCode = computed(() => JSON.stringify(auditPreview.value, null, 2))

const toggleAuditPreview = () => {
    auditPreviewVisible.value = !auditPreviewVisible.value
    if (!auditPreviewVisible.value) {
        auditPreviewCodeVisible.value = false
    }
}

const toggleAuditPreviewCode = () => {
    auditPreviewCodeVisible.value = !auditPreviewCodeVisible.value
}
</script>

<style scoped lang="scss">
.permission-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
}

.apply-empty-text {
    color: #909399;
    font-size: 13px;
}

.audit-preview-toolbar {
    display: flex;
    justify-content: flex-end;
    margin-top: 12px;
}

.audit-preview-code {
    max-height: 320px;
    margin-top: 8px;
    overflow: auto;
    padding: 12px;
    border-radius: 6px;
    background: #f6f8fa;
    color: #1f2937;
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
    font-size: 12px;
    line-height: 1.6;
    white-space: pre-wrap;
}
</style>
