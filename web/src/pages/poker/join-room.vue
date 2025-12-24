<template>
  <view class="game-bg" style="min-height: 100vh; display: flex; align-items: center; justify-content: center; padding: 32rpx;">
    <view class="animate-fade-in" style="width: 100%; max-width: 700rpx; position: relative; z-index: 10;">
      <!-- Header -->
      <view style="margin-bottom: 64rpx;">
        <view class="text-glow text-white" style="font-size: 56rpx; font-weight: bold; margin-bottom: 16rpx;">加入房间</view>
        <view class="text-white-70" style="font-size: 28rpx;">输入8位房间ID即可加入</view>
      </view>

      <!-- Form Card -->
      <view class="glass-card" style="padding: 48rpx;">
        <!-- Room ID Input -->
        <view style="margin-bottom: 48rpx;">
          <view class="text-white" style="font-size: 32rpx; font-weight: 600; margin-bottom: 24rpx;">房间ID</view>
          <u-input
            v-model="roomId"
            placeholder="请输入8位ID"
            :maxlength="8"
            :customStyle="{
              width: '100%',
              background: 'rgba(255, 255, 255, 0.1)',
              border: '2rpx solid rgba(255, 255, 255, 0.2)',
              borderRadius: '12rpx',
              color: 'white',
              textAlign: 'center',
              fontSize: '48rpx',
              fontWeight: 'bold',
              letterSpacing: '0.1em',
              padding: '32rpx'
            }"
            placeholderStyle="color: rgba(255, 255, 255, 0.5)"
          ></u-input>
        </view>

        <!-- Join Button -->
        <u-button
          type="success"
          text="加入房间"
          @click="handleJoin"
          :customStyle="{
            width: '100%',
            marginTop: '32rpx',
            fontSize: '32rpx',
            fontWeight: 'bold'
          }"
          :custom-class="'btn-success'"
        ></u-button>
      </view>
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
/* 使用全局样式 */
</style>
