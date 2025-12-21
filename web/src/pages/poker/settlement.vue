<template>
  <view class="container">
    <view class="header">
      <view class="title">游戏结算</view>
      <view class="subtitle">根据记账记录自动计算</view>
    </view>

    <view v-if="settlements.length > 0" class="settlements">
      <view
        v-for="(item, index) in settlements"
        :key="index"
        class="settlement-item"
      >
        <view class="from-player">
          <view class="player-avatar">
            {{ item.fromNickname.substring(0, 1) }}
          </view>
          <view class="player-name">{{ item.fromNickname }}</view>
        </view>

        <view class="arrow-wrapper">
          <view class="arrow">→</view>
          <view class="amount">¥{{ item.amount }}</view>
        </view>

        <view class="to-player">
          <view class="player-avatar">
            {{ item.toNickname.substring(0, 1) }}
          </view>
          <view class="player-name">{{ item.toNickname }}</view>
        </view>
      </view>
    </view>

    <view v-else class="no-settlement">
      <text class="no-settlement-text">无需结算</text>
    </view>

    <view class="summary">
      <view class="summary-title">结算说明</view>
      <view class="summary-text">
        以上是最优结算方案，通过最少的转账次数完成所有账务结清。
      </view>
    </view>

    <view class="button-group">
      <button class="back-btn" @click="backToHome">返回首页</button>
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
.container {
  min-height: 100vh;
  background: linear-gradient(135deg, #1AAD19 0%, #0e7d0e 100%);
  padding: 40rpx;
}

.header {
  margin-bottom: 40rpx;
  text-align: center;
}

.title {
  font-size: 48rpx;
  font-weight: bold;
  color: #fff;
  margin-bottom: 10rpx;
}

.subtitle {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.8);
}

.settlements {
  margin-bottom: 40rpx;
}

.settlement-item {
  background-color: rgba(255, 255, 255, 0.95);
  border-radius: 20rpx;
  padding: 40rpx;
  margin-bottom: 20rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.from-player,
.to-player {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
}

.player-avatar {
  width: 80rpx;
  height: 80rpx;
  border-radius: 40rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32rpx;
  font-weight: bold;
  color: #fff;
  margin-bottom: 10rpx;
}

.player-name {
  font-size: 26rpx;
  color: #333;
  text-align: center;
  word-break: break-all;
}

.arrow-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin: 0 30rpx;
}

.arrow {
  font-size: 48rpx;
  color: #1AAD19;
  margin-bottom: 10rpx;
}

.amount {
  font-size: 32rpx;
  font-weight: bold;
  color: #ff4444;
}

.no-settlement {
  background-color: rgba(255, 255, 255, 0.95);
  border-radius: 20rpx;
  padding: 100rpx 40rpx;
  text-align: center;
  margin-bottom: 40rpx;
}

.no-settlement-text {
  font-size: 32rpx;
  color: #999;
}

.summary {
  background-color: rgba(255, 255, 255, 0.95);
  border-radius: 20rpx;
  padding: 30rpx;
  margin-bottom: 40rpx;
}

.summary-title {
  font-size: 30rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 15rpx;
}

.summary-text {
  font-size: 26rpx;
  color: #666;
  line-height: 1.6;
}

.button-group {
  margin-top: 40rpx;
}

.back-btn {
  width: 100%;
  height: 90rpx;
  background-color: #fff;
  border-radius: 45rpx;
  font-size: 32rpx;
  font-weight: bold;
  color: #1AAD19;
  border: none;
}

.back-btn::after {
  border: none;
}
</style>
