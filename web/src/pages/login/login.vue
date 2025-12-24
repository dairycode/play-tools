<template>
  <view class="game-bg">
    <view class="status-bar" :style="{ height: statusBarHeight + 'px' }"></view>

    <view style="padding: 32rpx; padding-top: 160rpx; position: relative; z-index: 10;">
      <!-- Logo/Icon -->
      <view style="display: flex; justify-content: center; margin-bottom: 32rpx;">
        <view class="glass-card" style="width: 160rpx; height: 160rpx; display: flex; align-items: center; justify-content: center;">
          <text style="font-size: 80rpx;">🎮</text>
        </view>
      </view>

      <!-- Title -->
      <view style="text-align: center; margin-bottom: 64rpx;">
        <view class="text-glow" style="font-size: 56rpx; font-weight: bold; color: white; margin-bottom: 16rpx;">
          {{ isLoginMode ? '欢迎登录' : '注册账号' }}
        </view>
        <view class="text-white-70" style="font-size: 28rpx;">游戏助手小程序</view>
      </view>

      <!-- Form Card -->
      <view class="glass-card animate-fade-in" style="padding: 48rpx;">
        <!-- Username -->
        <view style="margin-bottom: 32rpx;">
          <u-input
            v-model="username"
            placeholder="请输入用户名"
            :maxlength="20"
            color="#ffffff"
            :customStyle="{
              background: 'rgba(255, 255, 255, 0.1)',
              border: '2rpx solid rgba(255, 255, 255, 0.2)',
              borderRadius: '12rpx',
              padding: '24rpx 32rpx',
              fontSize: '32rpx'
            }"
            placeholderStyle="color: rgba(255, 255, 255, 0.5)"
          ></u-input>
        </view>

        <!-- Password -->
        <view style="margin-bottom: 32rpx;">
          <u-input
            v-model="password"
            type="password"
            placeholder="请输入密码"
            :password="true"
            :maxlength="20"
            color="#ffffff"
            :customStyle="{
              background: 'rgba(255, 255, 255, 0.1)',
              border: '2rpx solid rgba(255, 255, 255, 0.2)',
              borderRadius: '12rpx',
              padding: '24rpx 32rpx',
              fontSize: '32rpx'
            }"
            placeholderStyle="color: rgba(255, 255, 255, 0.5)"
          ></u-input>
        </view>

        <!-- Nickname (Register mode only) -->
        <view v-if="!isLoginMode" style="margin-bottom: 32rpx;">
          <u-input
            v-model="nickname"
            placeholder="请输入昵称"
            :maxlength="20"
            color="#ffffff"
            :customStyle="{
              background: 'rgba(255, 255, 255, 0.1)',
              border: '2rpx solid rgba(255, 255, 255, 0.2)',
              borderRadius: '12rpx',
              padding: '24rpx 32rpx',
              fontSize: '32rpx'
            }"
            placeholderStyle="color: rgba(255, 255, 255, 0.5)"
          ></u-input>
        </view>

        <!-- Submit Button -->
        <u-button
          type="primary"
          :text="isLoginMode ? '登录' : '注册'"
          @click="handleSubmit"
          :customStyle="{
            width: '100%',
            marginBottom: '32rpx',
            fontSize: '32rpx',
            fontWeight: 'bold'
          }"
          :custom-class="'btn-primary'"
        ></u-button>

        <!-- Switch Mode -->
        <view class="text-white-90" style="text-align: center; font-size: 28rpx; cursor: pointer;" @click="switchMode">
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
const redirectUrl = ref('/pages/index/index')

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
        uni.switchTab({
          url: redirectUrl.value,
          fail: () => {
            uni.redirectTo({
              url: redirectUrl.value,
              fail: () => {
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
</style>
