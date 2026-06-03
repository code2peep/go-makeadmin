import request from '@/utils/request'

export function demoNoticeLists(params?: any) {
    return request.get({ url: '/demo_notice/list', params })
}

export function demoNoticeDetail(params?: any) {
    return request.get({ url: '/demo_notice/detail', params })
}
