<template>
    <div class="article-page">
        <el-card class="!border-none" shadow="never">
            <div class="module-heading">
                <div>
                    <div class="module-title">Demo Article</div>
                    <div class="module-subtitle">article / ma_demo_article</div>
                </div>
                <div class="module-tags">
                    <el-tag type="success" effect="plain">P5.1</el-tag>
                    <el-tag effect="plain">只读示例</el-tag>
                </div>
            </div>
            <el-form class="mt-4 mb-[-16px]" :model="formData" inline>
                <el-form-item label="标题">
                    <el-input
                        class="w-[260px]"
                        v-model="formData.title"
                        clearable
                        @keyup.enter="resetPage"
                    />
                </el-form-item>
                <el-form-item label="状态">
                    <el-select class="w-[180px]" v-model="formData.status" clearable>
                        <el-option label="全部" value="" />
                        <el-option label="启用" value="1" />
                        <el-option label="停用" value="0" />
                    </el-select>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="resetPage">查询</el-button>
                    <el-button @click="resetParams">重置</el-button>
                </el-form-item>
            </el-form>
        </el-card>

        <el-card class="!border-none mt-4" shadow="never" v-loading="pager.loading">
            <div class="mb-4">
                <el-button type="primary" @click="handleReadonly">
                    <template #icon>
                        <icon name="el-icon-Plus" />
                    </template>
                    新增
                </el-button>
                <el-button @click="handleRuntimeDetail">运行时详情</el-button>
            </div>
            <el-table :data="pager.lists" size="large" empty-text="暂无数据">
                <el-table-column label="ID" prop="id" min-width="80" />
                <el-table-column label="标题" prop="title" min-width="180" />
                <el-table-column label="状态" prop="status" min-width="100" />
                <el-table-column label="创建时间" prop="createTime" min-width="180" />
                <el-table-column label="操作" width="160" fixed="right">
                    <template #default>
                        <el-button link type="primary" @click="handleRuntimeDetail">详情</el-button>
                        <el-button link type="primary" @click="handleReadonly">编辑</el-button>
                        <el-button link type="danger" @click="handleReadonly">删除</el-button>
                    </template>
                </el-table-column>
            </el-table>
            <div class="flex justify-end mt-4">
                <pagination v-model="pager" @change="getLists" />
            </div>
        </el-card>

        <el-dialog v-model="detailVisible" title="运行时详情" width="520px">
            <el-descriptions :column="1" border>
                <el-descriptions-item label="module">
                    {{ runtimeDetail.module || 'article' }}
                </el-descriptions-item>
                <el-descriptions-item label="runtimeRegistered">
                    {{ runtimeDetail.runtimeRegistered === true ? 'true' : 'false' }}
                </el-descriptions-item>
            </el-descriptions>
        </el-dialog>
    </div>
</template>

<script lang="ts" setup name="article">
import { articleDetail, articleLists } from '@/api/article'
import { usePaging } from '@/hooks/usePaging'
import feedback from '@/utils/feedback'

const formData = reactive({
    title: '',
    status: ''
})

const runtimeDetail = ref<Record<string, any>>({})
const detailVisible = ref(false)

const { pager, getLists, resetParams, resetPage } = usePaging({
    fetchFun: articleLists,
    params: formData
})

const handleReadonly = () => {
    feedback.msgWarning('demo module is read-only')
}

const handleRuntimeDetail = async () => {
    runtimeDetail.value = await articleDetail()
    detailVisible.value = true
}

getLists()
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
