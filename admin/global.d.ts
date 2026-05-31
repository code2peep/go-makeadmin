/// <reference types="vite/client" />

declare module 'virtual:svg-icons-register'

declare module 'virtual:svg-icons-names' {
    const icons: string[]
    export default icons
}

declare module '@wangeditor/editor-for-vue' {
    export const Editor: any
    export const Toolbar: any
}
