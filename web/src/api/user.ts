// 用户相关API
import { request } from '@/utils/request'
import type { User } from '@/types'

// 获取用户信息
export const getUserInfo = () => {
  return request<User>({
    url: '/user/info',
    method: 'GET'
  })
}

// 更新用户信息
export const updateUserInfo = (data: { nickname?: string; avatar?: string }) => {
  return request({
    url: '/user/update',
    method: 'PUT',
    data
  })
}

// 微信小程序登录
export const wechatLogin = (data: { code: string; nickname: string; avatar?: string }) => {
  return request<{ token: string; user: User }>({
    url: '/user/wechat-login',
    method: 'POST',
    data
  })
}
