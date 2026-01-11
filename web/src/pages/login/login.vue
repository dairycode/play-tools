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
          欢迎使用
        </view>
        <view class="text-white-70" style="font-size: 28rpx;">游戏助手小程序</view>
      </view>

      <!-- Login Card -->
      <view class="glass-card animate-fade-in" style="padding: 48rpx;">
        <!-- Wechat Login Button (仅在微信小程序环境显示) -->
        <!-- #ifdef MP-WEIXIN -->
        <u-button
          type="success"
          text="微信一键登录"
          @click="handleWechatLogin"
          :customStyle="{
            width: '100%',
            marginBottom: '32rpx',
            fontSize: '32rpx',
            fontWeight: 'bold',
            background: 'linear-gradient(135deg, #07c160 0%, #05a854 100%)',
            border: 'none'
          }"
        ></u-button>
        <!-- #endif -->

        <view class="text-white-70" style="text-align: center; font-size: 24rpx; padding: 0 32rpx;">
          仅支持微信登录，点击按钮一键登录
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { wechatLogin } from '@/api/user'
import { setToken, setUserInfo } from '@/utils/auth'

const statusBarHeight = ref(0)
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

// 微信一键登录
const handleWechatLogin = async () => {
  try {
    uni.showLoading({
      title: '登录中...'
    })

    // 1. 获取微信登录 code
    const loginRes = await uni.login({
      provider: 'weixin'
    })

    if (!loginRes.code) {
      throw new Error('获取微信登录code失败')
    }

    // 2. 调用后端微信登录接口（使用默认昵称）
    const res = await wechatLogin({
      code: loginRes.code,
      nickname: '微信用户',
      avatar: ''
    })

    // 3. 保存 token 和用户信息
    setToken(res.data.token)
    setUserInfo(res.data.user)

    uni.hideLoading()

    // 4. 判断是否需要完善信息
    const needSetup = !res.data.user.nickname || res.data.user.nickname === '微信用户'

    if (needSetup) {
      // 跳转到个人信息填写页面
      uni.showToast({
        title: '登录成功，请完善信息',
        icon: 'success'
      })

      setTimeout(() => {
        uni.redirectTo({
          url: `/pages/login/profile-setup?redirect=${encodeURIComponent(redirectUrl.value)}`
        })
      }, 1000)
    } else {
      // 直接跳转到首页
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
    }
  } catch (error: any) {
    uni.hideLoading()

    // 处理用户拒绝授权的情况
    if (error.errMsg && error.errMsg.includes('auth deny')) {
      uni.showToast({
        title: '您拒绝了授权',
        icon: 'none'
      })
    } else {
      uni.showToast({
        title: error.message || '微信登录失败',
        icon: 'none'
      })
    }
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
