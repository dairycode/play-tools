<template>
  <view class="min-h-screen game-bg flex items-center justify-center p-4 md-p-6">
    <view class="w-full max-w-md relative z-10 animate-fade-in">
      <!-- Header -->
      <view class="mb-8">
        <view class="text-3xl font-bold text-white text-glow mb-2">创建房间</view>
        <view class="text-sm text-white-70">输入房间名称，创建后分享给好友加入</view>
      </view>

      <!-- Form Card -->
      <view class="glass-card p-6 md-p-8">
        <!-- Room Name -->
        <view class="mb-6">
          <view class="text-base font-semibold text-white mb-3">房间名称</view>
          <input
            class="input-glass text-base"
            v-model="roomName"
            placeholder="例如：周末德州局"
            maxlength="20"
          />
        </view>

        <!-- Player Count -->
        <view class="mb-6">
          <view class="text-base font-semibold text-white mb-3">房间人数</view>
          <view class="grid grid-cols-5 gap-2">
            <view
              v-for="count in [2, 3, 4, 5, 6, 7, 8, 9, 10]"
              :key="count"
              :class="[
                'h-14 rounded-xl flex items-center justify-center text-sm cursor-pointer transition-all duration-300',
                maxPlayers === count
                  ? 'bg-gradient-to-br from-primary-from to-primary-to text-white font-bold shadow-glow-blue'
                  : 'bg-white-10 text-white-60 hover-bg-white-15'
              ]"
              @click="maxPlayers = count"
            >
              {{ count }}人
            </view>
          </view>
        </view>

        <!-- Create Button -->
        <button class="w-full btn-primary text-base mt-4" @click="handleCreate">
          创建房间
        </button>
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
button::after {
  border: none;
}
</style>
