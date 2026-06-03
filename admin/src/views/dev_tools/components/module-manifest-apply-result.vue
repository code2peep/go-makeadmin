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
                <div class="permission-tags">
                    <el-tag v-for="code in permissionCodes" :key="code" size="small">
                        {{ code }}
                    </el-tag>
                </div>
            </el-descriptions-item>
        </el-descriptions>
        <el-table v-if="hasSnapshot" class="mt-3" :data="snapshotRows" size="large">
            <el-table-column label="对象" prop="name" min-width="130" />
            <el-table-column label="执行前" prop="before" min-width="100" />
            <el-table-column label="执行后" prop="after" min-width="100" />
        </el-table>
        <el-table v-if="result.checks?.length" class="mt-3" :data="result.checks" size="large">
            <el-table-column label="检查项" prop="name" min-width="140" />
            <el-table-column label="状态" prop="status" min-width="120" />
            <el-table-column label="说明" prop="message" min-width="280" />
        </el-table>
    </div>
</template>

<script lang="ts" setup>
import type { ModuleManifestApplyResult } from '@/api/tools/code'

interface SnapshotRow {
    name: string
    before: number
    after: number
}

const props = defineProps<{
    result: ModuleManifestApplyResult
    fallbackTitle: string
}>()

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
</script>

<style scoped lang="scss">
.permission-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
}
</style>
