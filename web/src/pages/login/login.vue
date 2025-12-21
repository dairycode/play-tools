<template>
  <view class="login-container">
    <view class="status-bar" :style="{ height: statusBarHeight + 'px' }"></view>

    <view class="login-content">
      <view class="title">{{ isLoginMode ? '欢迎登录' : '注册账号' }}</view>
      <view class="subtitle">游戏助手小程序</view>

      <view class="form">
        <view class="form-item">
          <input
            class="input"
            v-model="username"
            placeholder="请输入用户名"
            maxlength="20"
          />
        </view>

        <view class="form-item">
          <input
            class="input"
            v-model="password"
            type="password"
            placeholder="请输入密码"
            maxlength="20"
          />
        </view>

        <view class="form-item" v-if="!isLoginMode">
          <input
            class="input"
            v-model="nickname"
            placeholder="请输入昵称"
            maxlength="20"
          />
        </view>

        <button class="submit-btn" @click="handleSubmit">
          {{ isLoginMode ? '登录' : '注册' }}
        </button>

        <view class="switch-mode" @click="switchMode">
          {{ isLoginMode ? '没有账号？去注册' : '已有账号？去登录' }}
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { login, register } from '@/api/user'
import { setToken, setUserInfo } from '@/utils/auth'

const statusBarHeight = ref(0)
const isLoginMode = ref(true)
const username = ref('')
const password = ref('')
const nickname = ref('')
const redirectUrl = ref('/pages/index/index') // 默认跳转首页

// 获取状态栏高度
uni.getSystemInfo({
  success: (res) => {
    statusBarHeight.value = res.statusBarHeight || 0
  }
})

// 获取重定向路径
onLoad((options: any) => {
  if (options.redirect) {
    redirectUrl.value = decodeURIComponent(options.redirect)
  }
})

const switchMode = () => {
  isLoginMode.value = !isLoginMode.value
  username.value = ''
  password.value = ''
  nickname.value = ''
}

const handleSubmit = async () => {
  if (!username.value || !password.value) {
    uni.showToast({
      title: '请输入用户名和密码',
      icon: 'none'
    })
    return
  }

  if (!isLoginMode.value && !nickname.value) {
    uni.showToast({
      title: '请输入昵称',
      icon: 'none'
    })
    return
  }

  try {
    uni.showLoading({
      title: isLoginMode.value ? '登录中...' : '注册中...'
    })

    if (isLoginMode.value) {
      // 登录
      const res = await login({
        username: username.value,
        password: password.value
      })

      setToken(res.data.token)
      setUserInfo(res.data.user)

      uni.hideLoading()
      uni.showToast({
        title: '登录成功',
        icon: 'success'
      })

      setTimeout(() => {
        // 登录成功后跳转到之前的页面
        uni.switchTab({
          url: redirectUrl.value,
          fail: () => {
            // 如果不是 tabBar 页面，使用 redirectTo
            uni.redirectTo({
              url: redirectUrl.value,
              fail: () => {
                // 如果 redirectTo 也失败，使用 reLaunch
                uni.reLaunch({
                  url: redirectUrl.value
                })
              }
            })
          }
        })
      }, 1000)
    } else {
      // 注册
      await register({
        username: username.value,
        password: password.value,
        nickname: nickname.value
      })

      uni.hideLoading()
      uni.showToast({
        title: '注册成功，请登录',
        icon: 'success'
      })

      // 切换到登录模式
      setTimeout(() => {
        isLoginMode.value = true
        password.value = ''
        nickname.value = ''
      }, 1000)
    }
  } catch (error) {
    uni.hideLoading()
  }
}
</script>

<style scoped>
.login-container {
  width: 100vw;
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  flex-direction: column;
}

.status-bar {
  width: 100%;
}

.login-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 80rpx 60rpx;
}

.title {
  font-size: 48rpx;
  font-weight: bold;
  color: #fff;
  margin-bottom: 20rpx;
}

.subtitle {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.8);
  margin-bottom: 100rpx;
}

.form {
  width: 100%;
}

.form-item {
  margin-bottom: 40rpx;
}

.input {
  width: 100%;
  height: 90rpx;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 45rpx;
  padding: 0 40rpx;
  font-size: 28rpx;
  box-sizing: border-box;
}

.submit-btn {
  width: 100%;
  height: 90rpx;
  background: #fff;
  border-radius: 45rpx;
  font-size: 32rpx;
  font-weight: bold;
  color: #667eea;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 40rpx;
  border: none;
}

.submit-btn::after {
  border: none;
}

.switch-mode {
  text-align: center;
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.9);
  margin-top: 40rpx;
}
</style>
