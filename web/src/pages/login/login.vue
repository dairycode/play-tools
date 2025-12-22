<template>
  <view class="min-h-screen game-bg flex items-center justify-center p-4 md-p-6">
    <view class="status-bar" :style="{ height: statusBarHeight + 'px' }"></view>

    <view class="w-full max-w-md relative z-10 animate-fade-in">
      <!-- Logo/Icon -->
      <view class="flex justify-center mb-4">
        <view class="glass-card w-20 h-20 flex items-center justify-center">
          <text class="text-5xl">🎮</text>
        </view>
      </view>

      <!-- Title -->
      <view class="text-center mb-8">
        <view class="text-3xl font-bold text-white text-glow mb-2">
          {{ isLoginMode ? '欢迎登录' : '注册账号' }}
        </view>
        <view class="text-sm text-white-70">游戏助手小程序</view>
      </view>

      <!-- Form Card -->
      <view class="glass-card p-6 md-p-8 space-y-4">
        <!-- Username -->
        <view>
          <input
            class="input-glass text-base"
            v-model="username"
            placeholder="请输入用户名"
            maxlength="20"
          />
        </view>

        <!-- Password -->
        <view>
          <input
            class="input-glass text-base"
            v-model="password"
            type="password"
            placeholder="请输入密码"
            maxlength="20"
          />
        </view>

        <!-- Nickname (Register mode only) -->
        <view v-if="!isLoginMode">
          <input
            class="input-glass text-base"
            v-model="nickname"
            placeholder="请输入昵称"
            maxlength="20"
          />
        </view>

        <!-- Submit Button -->
        <button class="w-full btn-primary text-base" @click="handleSubmit">
          {{ isLoginMode ? '登录' : '注册' }}
        </button>

        <!-- Switch Mode -->
        <view class="text-center text-sm text-white-90 mt-4 cursor-pointer" @click="switchMode">
          {{ isLoginMode ? '没有账号？去注册' : '已有账号?去登录' }}
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
.status-bar {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  z-index: 999;
}

/* Remove button default styles for uniapp */
button::after {
  border: none;
}
</style>
