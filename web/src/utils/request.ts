// HTTP请求工具类
// 根据当前访问地址自动确定后端地址
const getBaseUrl = () => {
  // @ts-ignore
  if (typeof window !== 'undefined' && window.location) {
    const hostname = window.location.hostname
    // 如果是局域网IP或localhost，使用相同的host
    return `http://${hostname}:8081/api`
  }
  // 默认使用localhost（用于非H5环境）
  return 'http://192.168.0.102:8081/api'
}

const BASE_URL = getBaseUrl()

interface RequestOptions {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  data?: any
  header?: any
}

interface Response<T = any> {
  code: number
  msg: string
  data: T
}

export const request = <T = any>(options: RequestOptions): Promise<Response<T>> => {
  return new Promise((resolve, reject) => {
    const token = uni.getStorageSync('token')

    uni.request({
      url: BASE_URL + options.url,
      method: options.method || 'GET',
      data: options.data,
      timeout: 60000, // 设置60秒超时
      withCredentials: true, // H5环境下携带cookie
      header: {
        'Content-Type': 'application/json',
        'Authorization': token ? `Bearer ${token}` : '',
        ...options.header
      },
      success: (res: any) => {
        console.log('API Response:', options.url, res)
        const data = res.data as Response<T>
        if (data.code === 200) {
          resolve(data)
        } else if (data.code === 401) {
          // 未登录或token过期
          uni.removeStorageSync('token')
          uni.removeStorageSync('userInfo')
          uni.reLaunch({
            url: '/pages/index/index'
          })
          reject(data)
        } else {
          uni.showToast({
            title: data.msg || '请求失败',
            icon: 'none'
          })
          reject(data)
        }
      },
      fail: (err) => {
        console.error('API Error:', options.url, err)
        uni.showToast({
          title: '网络请求失败',
          icon: 'none'
        })
        reject(err)
      }
    })
  })
}
