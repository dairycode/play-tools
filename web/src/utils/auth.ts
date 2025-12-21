// 认证工具类
export const getToken = (): string => {
  return uni.getStorageSync('token') || ''
}

export const setToken = (token: string): void => {
  uni.setStorageSync('token', token)
}

export const removeToken = (): void => {
  uni.removeStorageSync('token')
}

export const isLogin = (): boolean => {
  return !!getToken()
}

export const getUserInfo = (): any => {
  return uni.getStorageSync('userInfo') || null
}

export const setUserInfo = (userInfo: any): void => {
  uni.setStorageSync('userInfo', userInfo)
}

export const removeUserInfo = (): void => {
  uni.removeStorageSync('userInfo')
}

export const logout = (): void => {
  removeToken()
  removeUserInfo()
  uni.reLaunch({
    url: '/pages/index/index'
  })
}
