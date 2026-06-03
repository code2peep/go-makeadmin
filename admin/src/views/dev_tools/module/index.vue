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
                <module-manifest-preview>
                    <el-button type="primary">
                        <template #icon>
                            <icon name="el-icon-Document" />
                        </template>
                        Manifest 预览
                    </el-button>
                </module-manifest-preview>
            </div>
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

        <el-card class="!border-none" shadow="never">
            <template #header>
                <div class="section-header">
                    <span class="card-title">内置模块清单</span>
                    <el-tag type="success" size="small">P4.2</el-tag>
                </div>
            </template>
            <el-table :data="modules" size="large">
                <el-table-column label="模块" prop="name" min-width="130" />
                <el-table-column label="Manifest" prop="manifest" min-width="240" />
                <el-table-column label="表名" prop="table" min-width="160" />
                <el-table-column label="运行时" prop="runtime" min-width="260" />
                <el-table-column label="状态" width="120">
                    <template #default="{ row }">
                        <el-tag :type="row.statusType" size="small">{{ row.status }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column label="入口" width="160" fixed="right">
                    <template #default>
                        <module-manifest-preview>
                            <el-button type="primary" link>
                                <template #icon>
                                    <icon name="el-icon-View" />
                                </template>
                                预览
                            </el-button>
                        </module-manifest-preview>
                    </template>
                </el-table-column>
            </el-table>
        </el-card>
    </div>
</template>

<script lang="ts" setup name="moduleCenter">
import ModuleManifestPreview from '../components/module-manifest-preview.vue'

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

const modules = [
    {
        name: 'Demo Article',
        manifest: 'examples/demo/manifest.json',
        table: 'ma_demo_article',
        runtime: 'MAKEADMIN_ENABLE_DEMO_MODULE=1',
        status: '可预览',
        statusType: 'success'
    }
]
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
}
</style>
