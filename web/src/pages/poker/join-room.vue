<template>
  <view class="container">
    <view class="header">
      <view class="title">加入房间</view>
      <view class="subtitle">输入8位房间ID即可加入</view>
    </view>

    <view class="form">
      <view class="form-item">
        <view class="label">房间ID</view>
        <input
          class="input"
          v-model="roomId"
          placeholder="请输入8位房间ID"
          maxlength="8"
        />
      </view>

      <button class="join-btn" @click="handleJoin">加入房间</button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { joinRoom } from '@/api/game'
import { isLogin } from '@/utils/auth'

const roomId = ref('')

onLoad(() => {
  // 检查登录状态
  if (!isLogin()) {
    uni.navigateTo({
      url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/poker/join-room')
    })
  }
})

const handleJoin = async () => {
  if (!roomId.value.trim()) {
    uni.showToast({
      title: '请输入房间ID',
      icon: 'none'
    })
    return
  }

  if (roomId.value.trim().length !== 8) {
    uni.showToast({
      title: '房间ID必须是8位',
      icon: 'none'
    })
    return
  }

  try {
    uni.showLoading({ title: '加入中...' })

    await joinRoom({
      roomId: roomId.value.trim()
    })

    uni.hideLoading()

    uni.showToast({
      title: '加入成功',
      icon: 'success'
    })

    // 跳转到游戏房间
    setTimeout(() => {
      uni.redirectTo({
        url: `/pages/poker/game-room?roomId=${roomId.value.trim()}`
      })
    }, 1000)
  } catch (error) {
    uni.hideLoading()
  }
}
</script>

<style scoped>
.container {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 40rpx;
}

.header {
  margin-bottom: 60rpx;
}

.title {
  font-size: 44rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 10rpx;
}

.subtitle {
  font-size: 26rpx;
  color: #999;
}

.form {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 40rpx;
}

.form-item {
  margin-bottom: 40rpx;
}

.label {
  font-size: 28rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 20rpx;
}

.input {
  width: 100%;
  height: 80rpx;
  background-color: #f5f5f5;
  border-radius: 10rpx;
  padding: 0 30rpx;
  font-size: 32rpx;
  font-weight: bold;
  text-align: center;
  letter-spacing: 4rpx;
  box-sizing: border-box;
}

.join-btn {
  width: 100%;
  height: 90rpx;
  background: linear-gradient(135deg, #1AAD19 0%, #0e7d0e 100%);
  border-radius: 45rpx;
  font-size: 32rpx;
  font-weight: bold;
  color: #fff;
  border: none;
  margin-top: 40rpx;
}

.join-btn::after {
  border: none;
}
</style>
