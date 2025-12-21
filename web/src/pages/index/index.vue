<template>
  <view class="container">
    <view class="header">
      <view class="header-title">游戏助手</view>
      <view class="account-btn" @click="goToAccount">
        <text class="account-icon">👤</text>
      </view>
    </view>

    <view class="tools-grid">
      <view class="tool-card" @click="goToCreateRoom">
        <view class="tool-icon">🃏</view>
        <view class="tool-name">打牌记账</view>
        <view class="tool-desc">创建房间开始游戏</view>
      </view>

      <view class="tool-card disabled">
        <view class="tool-icon">🎮</view>
        <view class="tool-name">敬请期待</view>
        <view class="tool-desc">更多工具开发中</view>
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
.container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 40rpx;
  box-sizing: border-box;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 60rpx;
}

.header-title {
  font-size: 48rpx;
  font-weight: bold;
  color: #fff;
}

.account-btn {
  width: 80rpx;
  height: 80rpx;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.account-icon {
  font-size: 40rpx;
}

.tools-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 30rpx;
}

.tool-card {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20rpx;
  padding: 50rpx 30rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 10rpx 30rpx rgba(0, 0, 0, 0.1);
}

.tool-card.disabled {
  opacity: 0.6;
}

.tool-icon {
  font-size: 80rpx;
  margin-bottom: 20rpx;
}

.tool-name {
  font-size: 32rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 10rpx;
}

.tool-desc {
  font-size: 24rpx;
  color: #999;
  text-align: center;
}
</style>
