<template>
  <view class="game-bg" style="min-height: 100vh; display: flex; align-items: flex-start; justify-content: center; padding: 32rpx; padding-top: 120rpx;">
    <view class="animate-fade-in" style="width: 100%; max-width: 700rpx; position: relative; z-index: 10;">
      <!-- Header -->
      <view style="margin-bottom: 64rpx;">
        <view class="text-glow text-white" style="font-size: 56rpx; font-weight: bold; margin-bottom: 16rpx;">创建房间</view>
        <view class="text-white-70" style="font-size: 28rpx;">输入房间名称，创建后分享给好友加入</view>
      </view>

      <!-- Form Card -->
      <view class="glass-card" style="padding: 48rpx;">
        <!-- Room Name -->
        <view style="margin-bottom: 48rpx;">
          <view class="text-white" style="font-size: 32rpx; font-weight: 600; margin-bottom: 24rpx;">房间名称</view>
          <u-input
            v-model="roomName"
            placeholder="例如：周末德州局"
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

        <!-- Player Count -->
        <view style="margin-bottom: 48rpx;">
          <view class="text-white" style="font-size: 32rpx; font-weight: 600; margin-bottom: 24rpx;">房间人数</view>
          <view style="display: grid; grid-template-columns: repeat(5, 1fr); gap: 16rpx;">
            <view
              v-for="count in [2, 3, 4, 5, 6, 7, 8, 9, 10]"
              :key="count"
              :style="{
                height: '112rpx',
                borderRadius: '12rpx',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                fontSize: '28rpx',
                cursor: 'pointer',
                transition: 'all 0.3s',
                background: maxPlayers === count
                  ? 'linear-gradient(135deg, #4fc3f7 0%, #2196f3 100%)'
                  : 'rgba(255, 255, 255, 0.1)',
                color: maxPlayers === count ? 'white' : 'rgba(255, 255, 255, 0.6)',
                fontWeight: maxPlayers === count ? 'bold' : 'normal',
                boxShadow: maxPlayers === count ? '0 8rpx 24rpx rgba(33, 150, 243, 0.3)' : 'none'
              }"
              @click="maxPlayers = count"
            >
              {{ count }}人
            </view>
          </view>
        </view>

        <!-- Create Button -->
        <u-button
          type="primary"
          text="创建房间"
          @click="handleCreate"
          :customStyle="{
            width: '100%',
            marginTop: '32rpx',
            fontSize: '32rpx',
            fontWeight: 'bold'
          }"
          :custom-class="'btn-primary'"
        ></u-button>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { createRoom } from '@/api/game'
import { isLogin } from '@/utils/auth'

const roomName = ref('')
const maxPlayers = ref(5) // 默认5人

onLoad(() => {
  // 检查登录状态
  if (!isLogin()) {
    uni.navigateTo({
      url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/poker/create-room')
    })
  }
})

const handleCreate = async () => {
  if (!roomName.value.trim()) {
    uni.showToast({
      title: '请输入房间名称',
      icon: 'none'
    })
    return
  }

  try {
    uni.showLoading({ title: '创建中...' })

    const res = await createRoom({
      name: roomName.value.trim(),
      maxPlayers: maxPlayers.value
    })

    uni.hideLoading()

    uni.showToast({
      title: '创建成功',
      icon: 'success'
    })

    // 直接跳转到游戏房间
    setTimeout(() => {
      uni.redirectTo({
        url: `/pages/poker/game-room?roomId=${res.data.id}`
      })
    }, 500)
  } catch (error) {
    uni.hideLoading()
  }
}
</script>

<style scoped>
/* 使用全局样式 */
</style>
