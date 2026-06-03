<template>
    <el-card class="!border-none" shadow="never">
        <div class="edit-placeholder">
            <div class="edit-title">Demo Article Edit</div>
            <el-tag effect="plain">只读示例</el-tag>
        </div>
    </el-card>
</template>

<script lang="ts" setup name="articleEdit"></script>

<style lang="scss" scoped>
.edit-placeholder {
    align-items: center;
    display: flex;
    gap: 12px;
}

.edit-title {
    color: #111827;
    font-size: 16px;
    font-weight: 600;
    line-height: 24px;
}
</style>
