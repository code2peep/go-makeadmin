<template>
    <main class="main-wrap h-full bg-page">
        <el-scrollbar>
            <div class="p-4">
                <router-view v-if="isRouteShow" v-slot="{ Component, route }">
                    <keep-alive
                        v-if="shouldKeepAlive(Component, route)"
                        :include="includeList"
                        :max="20"
                    >
                        <component :is="Component" :key="route.fullPath" />
                    </keep-alive>
                    <component v-else :is="Component" :key="route.fullPath" />
                </router-view>
            </div>
        </el-scrollbar>
    </main>
</template>

<script setup lang="ts">
import useAppStore from '@/stores/modules/app'
import useTabsStore from '@/stores/modules/multipleTabs'
import useSettingStore from '@/stores/modules/setting'
import type { Component } from 'vue'
import type { RouteLocationNormalizedLoaded } from 'vue-router'
const appStore = useAppStore()
const tabsStore = useTabsStore()
const settingStore = useSettingStore()
const isRouteShow = computed(() => appStore.isRouteShow)
const includeList = computed(() => (settingStore.openMultipleTabs ? tabsStore.getCacheTabList : []))
const shouldKeepAlive = (component: Component, route: RouteLocationNormalizedLoaded) => {
    return !!route.meta.keepAlive && getComponentName(component) !== 'RouterView'
}

const getComponentName = (component: Component) => {
    return typeof component === 'object' && component !== null && 'name' in component
        ? component.name
        : ''
}

</script>

<style></style>
