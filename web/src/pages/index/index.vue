<template>
  <view class="min-h-screen game-bg p-4 md-p-6">
    <view class="w-full max-w-4xl mx-auto relative z-10">
      <!-- Header -->
      <view class="flex items-center justify-between mb-8 animate-fade-in">
        <view class="text-3xl md-text-4xl font-bold text-white text-glow">游戏助手</view>
        <view
          class="w-11 h-11 glass-card flex items-center justify-center cursor-pointer transition-all duration-300 hover-bg-glass-lighter"
          @click="goToAccount"
        >
          <text class="text-2xl">👤</text>
        </view>
      </view>

      <!-- Tools Grid -->
      <view class="grid grid-cols-1 md-grid-cols-2 gap-4 animate-fade-in">
        <!-- Poker Tool Card -->
        <view
          class="glass-card p-8 cursor-pointer transition-all duration-300 hover-bg-glass-lighter hover-translate-y-1"
          @click="goToCreateRoom"
        >
          <view class="flex flex-col items-center text-center">
            <view class="text-6xl mb-4">🃏</view>
            <view class="text-2xl font-bold text-white mb-2">打牌记账</view>
            <view class="text-sm text-white-60">创建房间开始游戏</view>
          </view>
        </view>

        <!-- Coming Soon Card -->
        <view class="glass-card p-8 opacity-60">
          <view class="flex flex-col items-center text-center">
            <view class="text-6xl mb-4">🎮</view>
            <view class="text-2xl font-bold text-white mb-2">敬请期待</view>
            <view class="text-sm text-white-60">更多工具开发中</view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { onLoad } from '@dcloudio/uni-app'
import { isLogin } from '@/utils/auth'

onLoad(() => {
  // 检查登录状态
  if (!isLogin()) {
    // 跳转到登录页，并携带返回路径
    uni.navigateTo({
      url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/index/index')
    })
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
/* No custom styles needed - all using Tailwind */
</style>
