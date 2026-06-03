import request from '@/utils/request'

export function articleLists(params: any) {
    return request.get({ url: '/article/list', params })
}

export function articleDetail(params?: any) {
    return request.get({ url: '/article/detail', params })
}

export function articleAdd(params: any) {
    return request.post({ url: '/article/add', params })
}

export function articleEdit(params: any) {
    return request.post({ url: '/article/edit', params })
}

export function articleDelete(params: any) {
    return request.post({ url: '/article/del', params })
}
