<template>
  <view class="min-h-screen game-bg p-4 md:p-6">
    <view class="w-full max-w-2xl mx-auto relative z-10 animate-fade-in">
      <!-- Header -->
      <view class="text-center mb-8">
        <view class="text-4xl font-bold text-white text-glow mb-2">游戏结算</view>
        <view class="text-sm text-white/70">根据记账记录自动计算</view>
      </view>

      <!-- Settlements List -->
      <view v-if="settlements.length > 0" class="space-y-4 mb-6">
        <view
          v-for="(item, index) in settlements"
          :key="index"
          class="glass-card p-6"
        >
          <view class="flex items-center justify-between">
            <!-- From Player -->
            <view class="flex flex-col items-center flex-1">
              <view class="w-20 h-20 rounded-full bg-gradient-to-br from-primary-from to-primary-to flex items-center justify-center mb-2">
                <text class="text-2xl font-bold text-white">{{ item.fromNickname.substring(0, 1) }}</text>
              </view>
              <view class="text-sm text-white text-center">{{ item.fromNickname }}</view>
            </view>

            <!-- Arrow and Amount -->
            <view class="flex flex-col items-center mx-6">
              <view class="text-4xl text-success-from mb-2">→</view>
              <view class="text-2xl font-bold text-danger-from">¥{{ item.amount }}</view>
            </view>

            <!-- To Player -->
            <view class="flex flex-col items-center flex-1">
              <view class="w-20 h-20 rounded-full bg-gradient-to-br from-primary-from to-primary-to flex items-center justify-center mb-2">
                <text class="text-2xl font-bold text-white">{{ item.toNickname.substring(0, 1) }}</text>
              </view>
              <view class="text-sm text-white text-center">{{ item.toNickname }}</view>
            </view>
          </view>
        </view>
      </view>

      <!-- No Settlement -->
      <view v-else class="glass-card p-16 text-center mb-6">
        <text class="text-2xl text-white/50">无需结算</text>
      </view>

      <!-- Summary -->
      <view class="glass-card p-6 mb-6">
        <view class="text-lg font-bold text-white mb-3">结算说明</view>
        <view class="text-sm text-white/70 leading-relaxed">
          以上是最优结算方案，通过最少的转账次数完成所有账务结清。
        </view>
      </view>

      <!-- Action Button -->
      <button class="w-full btn-primary text-base" @click="backToHome">
        返回首页
      </button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import type { Settlement } from '@/types'
import { getRoomSettlements } from '@/api/game'

const settlements = ref<Settlement[]>([])
const roomId = ref('')

onLoad(async (options: any) => {
  if (options.roomId) {
    roomId.value = options.roomId
    await loadSettlements()
  } else if (options.settlements) {
    // 兼容旧的通过URL参数传递的方式
    try {
      settlements.value = JSON.parse(decodeURIComponent(options.settlements))
    } catch (error) {
      console.error('解析结算数据失败', error)
    }
  }
})

const loadSettlements = async () => {
  try {
    const res = await getRoomSettlements(roomId.value)
    settlements.value = res.data
  } catch (error: any) {
    uni.showToast({
      title: error.msg || '加载结算数据失败',
      icon: 'none'
    })
  }
}

const backToHome = () => {
  uni.reLaunch({
    url: '/pages/index/index'
  })
}
</script>

<style scoped>
button::after {
  border: none;
}
</style>
