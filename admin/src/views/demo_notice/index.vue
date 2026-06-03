<template>
    <div class="demo-notice-page">
        <el-card class="!border-none" shadow="never">
            <div class="module-heading">
                <div>
                    <div class="module-title">Demo Notice</div>
                    <div class="module-subtitle">demo_notice / ma_demo_notice</div>
                </div>
                <div class="module-tags">
                    <el-tag type="success" effect="plain">P5.20</el-tag>
                    <el-tag type="warning" effect="plain">前端只读示例</el-tag>
                </div>
            </div>
        </el-card>

        <el-card class="!border-none mt-4" shadow="never">
            <el-table :data="rows" size="large" empty-text="暂无数据">
                <el-table-column label="场景" prop="name" min-width="160" />
                <el-table-column label="状态" min-width="120">
                    <template #default="{ row }">
                        <el-tag :type="row.type" size="small">{{ row.status }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column label="说明" prop="detail" min-width="280" />
            </el-table>
        </el-card>
    </div>
</template>

<script lang="ts" setup name="demoNotice">
const rows = [
    {
        name: 'Registry',
        status: 'env gated',
        type: 'warning',
        detail: 'MAKEADMIN_ENABLE_DEMO_NOTICE_MODULE=1'
    },
    {
        name: 'Manifest',
        status: 'passed',
        type: 'success',
        detail: 'examples/demo_notice/manifest.json'
    },
    {
        name: 'Runtime',
        status: 'unregistered',
        type: 'info',
        detail: 'front-only read-only example'
    }
]
</script>

<style lang="scss" scoped>
.module-heading {
    align-items: flex-start;
    display: flex;
    gap: 16px;
    justify-content: space-between;
}

.module-title {
    color: #111827;
    font-size: 18px;
    font-weight: 600;
    line-height: 28px;
}

.module-subtitle {
    color: #667085;
    font-size: 13px;
    line-height: 20px;
    margin-top: 2px;
}

.module-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    justify-content: flex-end;
}

@media (max-width: 640px) {
    .module-heading {
        flex-direction: column;
    }

    .module-tags {
        justify-content: flex-start;
    }
}
</style>
