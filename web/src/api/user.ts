// 用户相关API
import { request } from '@/utils/request'
import type { User } from '@/types'

// 注册
export const register = (data: { username: string; password: string; nickname: string }) => {
  return request({
    url: '/user/register',
    method: 'POST',
    data
  })
}

// 登录
export const login = (data: { username: string; password: string }) => {
  return request<{ token: string; user: User }>({
    url: '/user/login',
    method: 'POST',
    data
  })
}

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

// 修改密码
export const changePassword = (data: { old_password: string; new_password: string }) => {
  return request({
    url: '/user/password',
    method: 'PUT',
    data
  })
}
