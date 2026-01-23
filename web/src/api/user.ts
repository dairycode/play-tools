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

// 上传头像
export const uploadAvatar = (filePath: string): Promise<{ code: number; msg: string; data: { url: string } }> => {
  return new Promise((resolve, reject) => {
    const token = uni.getStorageSync('token')
    const baseUrl = 'https://dairycode.top:8080/api'

    uni.uploadFile({
      url: baseUrl + '/user/upload-avatar',
      filePath: filePath,
      name: 'avatar',
      header: {
        'Authorization': token ? `Bearer ${token}` : ''
      },
      success: (res) => {
        console.log('上传头像响应:', res)
        try {
          const data = JSON.parse(res.data)
          if (data.code === 200) {
            resolve(data)
          } else {
            uni.showToast({
              title: data.msg || '上传失败',
              icon: 'none'
            })
            reject(data)
          }
        } catch (error) {
          console.error('解析上传响应失败:', error)
          uni.showToast({
            title: '上传失败',
            icon: 'none'
          })
          reject(error)
        }
      },
      fail: (err) => {
        console.error('上传头像失败:', err)
        uni.showToast({
          title: '上传失败',
          icon: 'none'
        })
        reject(err)
      }
    })
  })
}
