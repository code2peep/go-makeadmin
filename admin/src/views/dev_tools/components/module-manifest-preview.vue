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
                </el-descriptions>

                <el-table class="mt-4" :data="preview.detail.column" size="large" height="260">
                    <el-table-column label="字段" prop="columnName" min-width="130" />
                    <el-table-column label="Go 字段" prop="goField" min-width="120" />
                    <el-table-column label="Go 类型" prop="goType" min-width="100" />
                    <el-table-column label="表单" prop="htmlType" min-width="100" />
                    <el-table-column label="查询" prop="queryType" min-width="100" />
                    <el-table-column label="字典" prop="dictType" min-width="120" />
                </el-table>

                <div class="flex justify-end mt-4">
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
import { previewModuleManifest } from '@/api/tools/code'
import feedback from '@/utils/feedback'

const popupRef = shallowRef<InstanceType<typeof Popup>>()
const inputMode = ref<'path' | 'body'>('path')
const formData = reactive({
    manifestPath: 'examples/demo/manifest.json',
    manifestBody: '',
    authorName: 'codepeep'
})

const preview = ref<any>()
const previewState = reactive({
    show: false,
    code: {} as Record<string, string>
})

const handlePreview = async () => {
    const params =
        inputMode.value === 'path'
            ? {
                  manifestPath: formData.manifestPath,
                  authorName: formData.authorName
              }
            : {
                  manifestBody: formData.manifestBody,
                  authorName: formData.authorName
              }
    preview.value = await previewModuleManifest(params)
    feedback.msgSuccess('预览生成成功')
}

const handleCodePreview = () => {
    previewState.code = preview.value?.code || {}
    previewState.show = true
}
</script>

<style scoped lang="scss">
.module-manifest-preview {
    display: inline-block;
}

.manifest-result {
    margin-top: 16px;
}
</style>
