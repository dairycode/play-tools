<template>
  <view class="game-bg" style="padding: 32rpx; min-height: 100vh;">
    <view style="max-width: 1200rpx; margin: 0 auto; position: relative; z-index: 10;">
      <!-- Header -->
      <view class="animate-fade-in" style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 64rpx;">
        <view class="text-glow text-white" style="font-size: 60rpx; font-weight: bold;">游戏助手</view>
        <view
          class="glass-card"
          style="width: 88rpx; height: 88rpx; display: flex; align-items: center; justify-content: center; cursor: pointer;"
          @click="goToAccount"
        >
          <text style="font-size: 48rpx;">👤</text>
        </view>
      </view>

      <!-- Tools Grid -->
      <view class="animate-fade-in" style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 32rpx;">
        <!-- Poker Tool Card -->
        <view
          class="glass-card"
          style="padding: 64rpx; cursor: pointer; transition: all 0.3s;"
          @click="goToCreateRoom"
        >
          <view style="display: flex; flex-direction: column; align-items: center; text-align: center;">
            <view style="font-size: 120rpx; margin-bottom: 32rpx;">🃏</view>
            <view class="text-white" style="font-size: 48rpx; font-weight: bold; margin-bottom: 16rpx;">打牌记账</view>
            <view class="text-white-60" style="font-size: 28rpx;">创建房间开始游戏</view>
          </view>
        </view>

        <!-- Coming Soon Card -->
        <view class="glass-card" style="padding: 64rpx; opacity: 0.6;">
          <view style="display: flex; flex-direction: column; align-items: center; text-align: center;">
            <view style="font-size: 120rpx; margin-bottom: 32rpx;">🎮</view>
            <view class="text-white" style="font-size: 48rpx; font-weight: bold; margin-bottom: 16rpx;">敬请期待</view>
            <view class="text-white-60" style="font-size: 28rpx;">更多工具开发中</view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { isLogin } from '@/utils/auth'

onLoad(async () => {
  // 检查登录状态
  if (!isLogin()) {
    // 跳转到登录页，并携带返回路径
    uni.navigateTo({
      url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/index/index')
    })
    return
  }
})

const goToAccount = () => {
  uni.navigateTo({
    url: '/pages/account/account'
  })
}

const goToCreateRoom = () => {
  uni.navigateTo({
    url: '/pages/poker/create-room'
  })
}
</script>

<style scoped>
/* 使用全局样式 */
</style>
