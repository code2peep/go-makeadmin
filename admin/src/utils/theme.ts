type MixMethod = 'shade' | 'tint'
type MixConfig = Record<string, [MixMethod, number]>

const lightConfig = {
    'dark-2': ['shade', 20],
    'light-3': ['tint', 30],
    'light-5': ['tint', 50],
    'light-7': ['tint', 70],
    'light-8': ['tint', 80],
    'light-9': ['tint', 90]
} satisfies MixConfig

const darkConfig = {
    'light-3': ['shade', 20],
    'light-5': ['shade', 30],
    'light-7': ['shade', 50],
    'light-8': ['shade', 60],
    'light-9': ['shade', 70],
    'dark-2': ['tint', 20]
} satisfies MixConfig

const themeId = 'theme-vars'

const normalizeHexColor = (color: string) => {
    const hex = color.trim().replace(/^#/, '')
    if (/^[0-9a-fA-F]{3}$/.test(hex)) {
        return hex
            .split('')
            .map((char) => char + char)
            .join('')
    }
    if (/^[0-9a-fA-F]{6}$/.test(hex)) {
        return hex
    }
    return ''
}

const mixHexColor = (color: string, method: MixMethod, percent: number) => {
    const hex = normalizeHexColor(color)
    if (!hex) {
        return color
    }
    const target = method === 'tint' ? 255 : 0
    const ratio = Math.min(Math.max(percent, 0), 100) / 100
    const rgb = [0, 2, 4].map((index) => {
        const value = parseInt(hex.slice(index, index + 2), 16)
        const mixed = Math.round(value + (target - value) * ratio)
        return mixed.toString(16).padStart(2, '0')
    })
    return `#${rgb.join('')}`
}

/**
 * @author Jason
 * @description 用于生成elementui主题的行为变量
 * 可选值有primary、success、warning、danger、error、info
 */

export const generateVars = (color: string, type = 'primary', isDark = false) => {
    const colors = {
        [`--el-color-${type}`]: color
    }
    const config: MixConfig = isDark ? darkConfig : lightConfig
    for (const key in config) {
        const [method, percent] = config[key]
        colors[`--el-color-${type}-${key}`] = mixHexColor(color, method, percent)
    }
    return colors
}

/**
 * @author Jason
 * @description 用于设置css变量
 * @param key css变量key 如 --color-primary
 * @param value css变量值 如 #f40
 * @param dom dom元素
 */
export const setCssVar = (key: string, value: string, dom = document.documentElement) => {
    dom.style.setProperty(key, value)
}

/**
 * @author Jason
 * @description 设置主题
 */
export const setTheme = (options: Record<string, string>, isDark = false) => {
    const varsMap: Record<string, string> = Object.keys(options).reduce((prev, key) => {
        return Object.assign(prev, generateVars(options[key], key, isDark))
    }, {})

    let theme = Object.keys(varsMap).reduce((prev, key) => `${prev}${key}:${varsMap[key]};`, '')
    theme = `:root{${theme}}`
    let style = document.getElementById(themeId)
    if (style) {
        style.innerHTML = theme
        return
    }
    style = document.createElement('style')
    style.setAttribute('type', 'text/css')
    style.setAttribute('id', themeId)
    style.innerHTML = theme
    document.head.append(style)
}
