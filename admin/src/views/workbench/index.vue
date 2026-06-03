<template>
    <div class="workbench">
        <el-card class="!border-none mb-4" shadow="never">
            <div class="workbench-hero">
                <div>
                    <div class="hero-title">go-makeadmin 基础框架</div>
                    <div class="hero-subtitle">
                        {{ workbenchData.version.based || 'Go、Gin、Gorm、Vue3、Element Plus、MySQL、Redis' }}
                    </div>
                </div>
                <div class="hero-tags">
                    <el-tag type="success" effect="plain">P3 已冻结</el-tag>
                    <el-tag type="primary" effect="plain">{{ workbenchData.framework.stage }}</el-tag>
                    <el-tag effect="plain">{{ workbenchData.framework.database }}</el-tag>
                </div>
            </div>
        </el-card>

        <div class="status-grid mb-4">
            <el-card v-for="item in capabilityCards" :key="item.name" class="!border-none" shadow="never">
                <div class="metric-card">
                    <div class="metric-icon">
                        <icon :name="item.icon" :size="24" />
                    </div>
                    <div>
                        <div class="metric-label">{{ item.name }}</div>
                        <div class="metric-value">{{ item.value }}</div>
                        <div class="metric-desc">{{ item.desc }}</div>
                    </div>
                </div>
            </el-card>
        </div>

        <div class="md:flex md:gap-4 mb-4">
            <el-card class="!border-none mb-4 md:mb-0 md:flex-1" shadow="never">
                <template #header>
                    <div class="section-header">
                        <span class="card-title">阶段状态</span>
                        <el-tag size="small" type="success">framework</el-tag>
                    </div>
                </template>
                <div class="timeline-list">
                    <div v-for="item in workbenchData.milestones" :key="item.name" class="timeline-item">
                        <div class="timeline-dot" :class="{ active: item.status === '进行中' }"></div>
                        <div class="timeline-content">
                            <div class="timeline-title">
                                <span>{{ item.name }}</span>
                                <el-tag size="small" :type="item.status === '进行中' ? 'primary' : 'success'">
                                    {{ item.status }}
                                </el-tag>
                            </div>
                            <div class="timeline-summary">{{ item.summary }}</div>
                        </div>
                    </div>
                </div>
            </el-card>

            <el-card class="!border-none md:w-[420px]" shadow="never">
                <template #header>
                    <div class="section-header">
                        <span class="card-title">验收状态</span>
                        <el-tag size="small" type="info">local</el-tag>
                    </div>
                </template>
                <div class="validation-list">
                    <div v-for="item in workbenchData.validation" :key="item.name" class="validation-item">
                        <div class="validation-main">
                            <span>{{ item.name }}</span>
                            <el-tag size="small" :type="item.status === '通过' ? 'success' : 'primary'">
                                {{ item.status }}
                            </el-tag>
                        </div>
                        <div class="validation-scope">{{ item.scope }}</div>
                    </div>
                </div>
            </el-card>
        </div>

        <el-card class="!border-none mb-4" shadow="never">
            <template #header>
                <div class="section-header">
                    <span class="card-title">核心页面验收</span>
                    <el-tag size="small" type="primary">P4.6</el-tag>
                </div>
            </template>
            <el-table :data="workbenchData.corePages" size="large">
                <el-table-column label="页面" prop="name" min-width="120" />
                <el-table-column label="状态" min-width="120">
                    <template #default="{ row }">
                        <el-tag size="small" type="success">{{ row.status }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column label="范围" prop="scope" min-width="220" />
                <el-table-column label="路由" prop="route" min-width="140" />
                <el-table-column label="入口" width="120" fixed="right">
                    <template #default="{ row }">
                        <el-button type="primary" link @click="goTo(row.route)">打开</el-button>
                    </template>
                </el-table-column>
            </el-table>
        </el-card>

        <el-card class="!border-none" shadow="never">
            <template #header>
                <span class="card-title">人工测试入口</span>
            </template>
            <div class="quick-grid">
                <div v-for="item in quickActions" :key="item.name" class="quick-link">
                    <el-button class="quick-button" @click="goTo(item.url)">
                        <template #icon>
                            <icon :name="item.icon" />
                        </template>
                        {{ item.name }}
                    </el-button>
                    <div class="quick-scope">{{ item.scope }}</div>
                </div>
            </div>
        </el-card>
    </div>
</template>

<script lang="ts" setup name="workbench">
import { getWorkbench } from '@/api/app'

const router = useRouter()

type ConsoleVersion = {
    version: string
    based: string
}

type ConsoleFramework = {
    stage: string
    database: string
    tables: string
    auth: string
    moduleLifecycle: string
}

type ConsoleMilestone = {
    name: string
    status: string
    summary: string
}

type ConsoleValidation = {
    name: string
    status: string
    scope: string
}

type ConsoleCorePage = {
    name: string
    route: string
    status: string
    scope: string
}

const defaultWorkbench = {
    version: {
        version: 'v0.1.0',
        based: 'Go、Gin、Gorm、Vue3、Element Plus、MySQL、Redis'
    } as ConsoleVersion,
    framework: {
        stage: 'P4.6 核心管理页可见验收',
        database: 'go_makeadmin',
        tables: 'ma_*',
        auth: 'JWT + Redis session',
        moduleLifecycle: 'manifest + codegen + install/uninstall apply'
    } as ConsoleFramework,
    milestones: [
        {
            name: 'P1 核心后台',
            status: '已冻结',
            summary: '登录、菜单、权限、设置、字典、文件、日志和代码生成器切到 ma_*。'
        },
        {
            name: 'P2 权限租户',
            status: '已冻结',
            summary: 'JWT、Redis session、租户上下文、数据权限和模块生命周期命令完成。'
        },
        {
            name: 'P3 模块产品化',
            status: '已冻结',
            summary: '脚手架、codegen、manifest、安装卸载和 apply 结果闭环完成。'
        },
        {
            name: 'P4 可见后台',
            status: '进行中',
            summary: '把底座能力沉到后台页面，进入人工测试和产品体验验收。'
        }
    ] as ConsoleMilestone[],
    validation: [
        {
            name: '无库验证',
            status: '通过',
            scope: 'runtime residue、Go test、type-check、build、npm audit'
        },
        {
            name: '模块工具链',
            status: '通过',
            scope: 'manifest、脚手架、codegen、安装卸载计划、写入门禁'
        },
        {
            name: '模块中心',
            status: '通过',
            scope: '内嵌预览、apply 结果、状态清单'
        },
        {
            name: '核心页面入口',
            status: '就绪',
            scope: '菜单、角色、管理员、部门、网站信息'
        },
        {
            name: '本地 API',
            status: '可用',
            scope: 'http://127.0.0.1:18000/api'
        },
        {
            name: '管理端',
            status: '可用',
            scope: 'http://127.0.0.1:5173'
        }
    ] as ConsoleValidation[],
    corePages: [
        {
            name: '菜单权限',
            route: '/menu',
            status: '入口就绪',
            scope: '菜单树、权限字符、路由显隐'
        },
        {
            name: '角色管理',
            route: '/role',
            status: '入口就绪',
            scope: '角色列表、授权入口、数据权限'
        },
        {
            name: '管理员',
            route: '/admin',
            status: '入口就绪',
            scope: '账号列表、组织岗位、启停'
        },
        {
            name: '组织部门',
            route: '/department',
            status: '入口就绪',
            scope: '部门树、负责人、状态'
        },
        {
            name: '网站信息',
            route: '/information',
            status: '入口就绪',
            scope: '站点名称、Logo、备案基础信息'
        },
        {
            name: '系统缓存',
            route: '/cache',
            status: '入口就绪',
            scope: '缓存清理、本地运行状态'
        },
        {
            name: '系统日志',
            route: '/journal',
            status: '入口就绪',
            scope: '管理员操作日志、登录日志'
        }
    ] as ConsoleCorePage[]
}

const workbenchData = reactive({ ...defaultWorkbench })

const capabilityCards = computed(() => [
    {
        name: '核心后台',
        value: workbenchData.framework.tables,
        desc: '登录、菜单、权限、设置、字典、文件、日志',
        icon: 'el-icon-Monitor'
    },
    {
        name: '认证权限',
        value: workbenchData.framework.auth,
        desc: '管理员、角色、菜单权限、租户上下文',
        icon: 'el-icon-Lock'
    },
    {
        name: '模块闭环',
        value: workbenchData.framework.moduleLifecycle,
        desc: '脚手架、代码生成、安装、卸载、回读',
        icon: 'el-icon-Tools'
    },
    {
        name: '当前版本',
        value: workbenchData.version.version,
        desc: 'P4 入口：可见后台与人工测试',
        icon: 'el-icon-Flag'
    }
])

const quickActions = [
    {
        name: '模块中心',
        url: '/module',
        icon: 'el-icon-Box',
        scope: 'manifest、安装计划、安装卸载、审计预览'
    },
    {
        name: '代码生成器',
        url: '/code',
        icon: 'el-icon-Tools',
        scope: 'manifest 预览、生成配置、安装卸载 apply'
    },
    {
        name: '菜单权限',
        url: '/menu',
        icon: 'el-icon-Menu',
        scope: '路由、权限字符、菜单显隐'
    },
    {
        name: '角色管理',
        url: '/role',
        icon: 'el-icon-Key',
        scope: '角色、授权、数据权限入口'
    },
    {
        name: '管理员',
        url: '/admin',
        icon: 'el-icon-User',
        scope: '账号、组织岗位、启停'
    },
    {
        name: '组织部门',
        url: '/department',
        icon: 'el-icon-OfficeBuilding',
        scope: '部门树、负责人、状态'
    },
    {
        name: '网站信息',
        url: '/information',
        icon: 'el-icon-Setting',
        scope: '站点名称、Logo、备案基础信息'
    }
]

const getData = async () => {
    const res = await getWorkbench()
    workbenchData.version = res.version || defaultWorkbench.version
    workbenchData.framework = res.framework || defaultWorkbench.framework
    workbenchData.milestones = res.milestones || defaultWorkbench.milestones
    workbenchData.validation = res.validation || defaultWorkbench.validation
    workbenchData.corePages = res.corePages || defaultWorkbench.corePages
}

const goTo = (url: string) => {
    router.push(url)
}

getData()
</script>

<style lang="scss" scoped>
.workbench {
    .el-card {
        border-radius: 8px;
    }
}

.workbench-hero {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 20px;
}

.hero-title {
    color: #1f2937;
    font-size: 24px;
    font-weight: 600;
    line-height: 32px;
}

.hero-subtitle {
    color: #667085;
    font-size: 14px;
    line-height: 24px;
    margin-top: 6px;
}

.hero-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    justify-content: flex-end;
}

.status-grid {
    display: grid;
    gap: 16px;
    grid-template-columns: repeat(4, minmax(0, 1fr));
}

.metric-card {
    display: flex;
    gap: 14px;
    min-height: 92px;
}

.metric-icon {
    align-items: center;
    background: #eef4ff;
    border-radius: 8px;
    color: #3b5bdb;
    display: flex;
    flex: 0 0 44px;
    height: 44px;
    justify-content: center;
}

.metric-label {
    color: #667085;
    font-size: 13px;
    line-height: 20px;
}

.metric-value {
    color: #111827;
    font-size: 17px;
    font-weight: 600;
    line-height: 26px;
}

.metric-desc,
.timeline-summary,
.validation-scope,
.quick-scope {
    color: #667085;
    font-size: 13px;
    line-height: 20px;
}

.section-header {
    align-items: center;
    display: flex;
    justify-content: space-between;
}

.timeline-list,
.validation-list {
    display: flex;
    flex-direction: column;
    gap: 14px;
}

.timeline-item {
    display: flex;
    gap: 12px;
}

.timeline-dot {
    background: #16a34a;
    border-radius: 50%;
    flex: 0 0 10px;
    height: 10px;
    margin-top: 8px;
}

.timeline-dot.active {
    background: #2563eb;
}

.timeline-content {
    flex: 1;
    min-width: 0;
}

.timeline-title,
.validation-main {
    align-items: center;
    color: #111827;
    display: flex;
    font-weight: 500;
    gap: 8px;
    justify-content: space-between;
    line-height: 24px;
}

.validation-item {
    border-bottom: 1px solid #eef0f3;
    padding-bottom: 12px;
}

.validation-item:last-child {
    border-bottom: 0;
    padding-bottom: 0;
}

.workbench :deep(.el-table .cell) {
    overflow-wrap: anywhere;
}

.quick-grid {
    display: grid;
    gap: 14px;
    grid-template-columns: repeat(3, minmax(0, 1fr));
}

.quick-link {
    border: 1px solid #eef0f3;
    border-radius: 8px;
    color: inherit;
    display: block;
    min-height: 92px;
    padding: 14px;
}

.quick-button {
    justify-content: flex-start;
    width: 100%;
}

@media (max-width: 1024px) {
    .status-grid,
    .quick-grid {
        grid-template-columns: repeat(2, minmax(0, 1fr));
    }
}

@media (max-width: 640px) {
    .workbench-hero {
        flex-direction: column;
    }

    .hero-tags {
        justify-content: flex-start;
    }

    .status-grid,
    .quick-grid {
        grid-template-columns: 1fr;
    }
}
</style>
