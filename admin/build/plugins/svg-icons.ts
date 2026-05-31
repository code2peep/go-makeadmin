import fs from 'node:fs'
import path from 'node:path'
import type { Plugin } from 'vite'

const REGISTER_ID = 'virtual:svg-icons-register'
const NAMES_ID = 'virtual:svg-icons-names'
const RESOLVED_REGISTER_ID = `\0${REGISTER_ID}`
const RESOLVED_NAMES_ID = `\0${NAMES_ID}`

interface SvgIconOptions {
    iconDirs: string[]
    symbolId: string
}

interface SvgIconData {
    file: string
    id: string
    symbol: string
}

function collectSvgFiles(dir: string) {
    if (!fs.existsSync(dir)) {
        return []
    }
    const files: string[] = []
    for (const name of fs.readdirSync(dir)) {
        const filepath = path.join(dir, name)
        const stat = fs.statSync(filepath)
        if (stat.isDirectory()) {
            files.push(...collectSvgFiles(filepath))
            continue
        }
        if (stat.isFile() && path.extname(filepath).toLowerCase() === '.svg') {
            files.push(filepath)
        }
    }
    return files.sort()
}

function getSymbolId(pattern: string, iconDir: string, file: string) {
    const relative = path.relative(iconDir, file)
    const parsed = path.parse(relative)
    const dir = parsed.dir.split(path.sep).filter(Boolean).join('-')
    return pattern
        .replace('[dir]', dir)
        .replace('[name]', parsed.name)
        .replace(/--+/g, '-')
        .replace(/-$/g, '')
}

function getAttribute(source: string, name: string) {
    return source.match(new RegExp(`\\s${name}=(["'])(.*?)\\1`, 'i'))?.[2]
}

function createSymbol(file: string, id: string) {
    const content = fs
        .readFileSync(file, 'utf-8')
        .replace(/<\?xml[\s\S]*?\?>/g, '')
        .replace(/<!DOCTYPE[\s\S]*?>/gi, '')
    const svg = content.match(/<svg([^>]*)>([\s\S]*?)<\/svg>/i)
    if (!svg) {
        return ''
    }
    const [, attributes, body] = svg
    const className = getAttribute(attributes, 'class')
    const viewBox = getAttribute(attributes, 'viewBox')
    const classAttr = className ? ` class="${className}"` : ''
    const viewBoxAttr = viewBox ? ` viewBox="${viewBox}"` : ''
    return `<symbol${classAttr}${viewBoxAttr} id="${id}">${body}</symbol>`
}

function loadIcons(options: SvgIconOptions) {
    const icons: SvgIconData[] = []
    for (const iconDir of options.iconDirs) {
        for (const file of collectSvgFiles(iconDir)) {
            const id = getSymbolId(options.symbolId, iconDir, file)
            icons.push({
                file,
                id,
                symbol: createSymbol(file, id)
            })
        }
    }
    return icons
}

function createRegisterModule(sprite: string) {
    return `
const sprite = ${JSON.stringify(sprite)}
const id = '__makeadmin_svg_icons__'

function injectSvgIcons() {
    let container = document.getElementById(id)
    if (!container) {
        container = document.createElement('div')
        container.id = id
        container.style.position = 'absolute'
        container.style.width = '0'
        container.style.height = '0'
        container.style.overflow = 'hidden'
        container.setAttribute('aria-hidden', 'true')
        document.body.insertBefore(container, document.body.firstChild)
    }
    container.innerHTML = sprite
}

if (typeof document !== 'undefined') {
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', injectSvgIcons)
    } else {
        injectSvgIcons()
    }
}
`
}

export function createLocalSvgIconsPlugin(options: SvgIconOptions): Plugin {
    return {
        name: 'makeadmin-svg-icons',
        resolveId(id) {
            if (id === REGISTER_ID) {
                return RESOLVED_REGISTER_ID
            }
            if (id === NAMES_ID) {
                return RESOLVED_NAMES_ID
            }
        },
        load(id) {
            if (id !== RESOLVED_REGISTER_ID && id !== RESOLVED_NAMES_ID) {
                return
            }
            const icons = loadIcons(options)
            for (const icon of icons) {
                this.addWatchFile(icon.file)
            }
            if (id === RESOLVED_NAMES_ID) {
                return `export default ${JSON.stringify(icons.map((icon) => icon.id))}`
            }
            const sprite = `<svg xmlns="http://www.w3.org/2000/svg" style="display:none">${icons
                .map((icon) => icon.symbol)
                .join('')}</svg>`
            return createRegisterModule(sprite)
        }
    }
}
