<template>
  <view class="game-bg" style="min-height: 100vh; padding: 32rpx;">
    <view class="animate-fade-in" style="width: 100%; max-width: 1000rpx; margin: 0 auto; position: relative; z-index: 10;">
      <!-- Header -->
      <view style="text-align: center; margin-bottom: 64rpx;">
        <view class="text-glow text-white" style="font-size: 72rpx; font-weight: bold; margin-bottom: 16rpx;">游戏结算</view>
        <view class="text-white-70" style="font-size: 28rpx;">根据记账记录自动计算</view>
      </view>

      <!-- Settlements List -->
      <view v-if="settlements.length > 0" style="margin-bottom: 48rpx;">
        <view
          v-for="(item, index) in settlements"
          :key="index"
          class="glass-card"
          style="padding: 48rpx; margin-bottom: 32rpx;"
        >
          <view style="display: flex; align-items: center; justify-content: space-between;">
            <!-- From Player -->
            <view style="display: flex; flex-direction: column; align-items: center; flex: 1;">
              <view class="btn-primary" style="width: 160rpx; height: 160rpx; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin-bottom: 16rpx;">
                <text style="font-size: 64rpx; font-weight: bold; color: white;">{{ item.fromNickname.substring(0, 1) }}</text>
              </view>
              <view class="text-white" style="font-size: 28rpx; text-align: center;">{{ item.fromNickname }}</view>
            </view>

            <!-- Arrow and Amount -->
            <view style="display: flex; flex-direction: column; align-items: center; margin: 0 48rpx;">
              <view class="text-success" style="font-size: 80rpx; margin-bottom: 16rpx;">→</view>
              <view class="text-danger" style="font-size: 56rpx; font-weight: bold;">¥{{ item.amount }}</view>
            </view>

            <!-- To Player -->
            <view style="display: flex; flex-direction: column; align-items: center; flex: 1;">
              <view class="btn-primary" style="width: 160rpx; height: 160rpx; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin-bottom: 16rpx;">
                <text style="font-size: 64rpx; font-weight: bold; color: white;">{{ item.toNickname.substring(0, 1) }}</text>
              </view>
              <view class="text-white" style="font-size: 28rpx; text-align: center;">{{ item.toNickname }}</view>
            </view>
          </view>
        </view>
      </view>

      <!-- No Settlement -->
      <view v-else class="glass-card" style="padding: 128rpx; text-align: center; margin-bottom: 48rpx;">
        <text class="text-white-50" style="font-size: 56rpx;">无需结算</text>
      </view>

      <!-- Summary -->
      <view class="glass-card" style="padding: 48rpx; margin-bottom: 48rpx;">
        <view class="text-white" style="font-size: 36rpx; font-weight: bold; margin-bottom: 24rpx;">结算说明</view>
        <view class="text-white-70" style="font-size: 28rpx; line-height: 1.6;">
          以上是最优结算方案，通过最少的转账次数完成所有账务结清。
        </view>
      </view>

      <!-- Action Button -->
      <u-button
        type="primary"
        text="返回首页"
        @click="backToHome"
        :customStyle="{
          width: '100%',
          fontSize: '32rpx',
          fontWeight: 'bold'
        }"
        :custom-class="'btn-primary'"
      ></u-button>
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
/* 使用全局样式 */
</style>
