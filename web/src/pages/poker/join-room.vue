<template>
  <view class="min-h-screen game-bg flex items-center justify-center p-4 md-p-6">
    <view class="w-full max-w-md relative z-10 animate-fade-in">
      <!-- Header -->
      <view class="mb-8">
        <view class="text-3xl font-bold text-white text-glow mb-2">加入房间</view>
        <view class="text-sm text-white-70">输入8位房间ID即可加入</view>
      </view>

      <!-- Form Card -->
      <view class="glass-card p-6 md-p-8">
        <!-- Room ID Input -->
        <view class="mb-6">
          <view class="text-base font-semibold text-white mb-3">房间ID</view>
          <input
            class="w-full px-4 py-4 bg-white-10 rounded-xl text-white text-center text-2xl font-bold tracking-wider border-2 border-white-20 focus-border-primary-from-60 transition-all"
            v-model="roomId"
            placeholder="请输入8位ID"
            maxlength="8"
          />
        </view>

        <!-- Join Button -->
        <button class="w-full btn-success text-base mt-4" @click="handleJoin">
          加入房间
        </button>
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
button::after {
  border: none;
}
</style>
