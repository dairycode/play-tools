<template>
  <view class="container">
    <view class="header">
      <view class="title">创建房间</view>
      <view class="subtitle">输入房间名称，创建后分享给好友加入</view>
    </view>

    <view class="form">
      <view class="form-item">
        <view class="label">房间名称</view>
        <input
          class="input"
          v-model="roomName"
          placeholder="例如：周末德州局"
          maxlength="20"
        />
      </view>

      <view class="form-item">
        <view class="label">房间人数</view>
        <view class="player-count-grid">
          <view
            v-for="count in [2, 3, 4, 5, 6, 7, 8, 9, 10]"
            :key="count"
            :class="['count-option', { active: maxPlayers === count }]"
            @click="maxPlayers = count"
          >
            {{ count }}人
          </view>
        </view>
      </view>

      <button class="create-btn" @click="handleCreate">创建房间</button>
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

    // 显示房间ID，方便分享
    uni.showModal({
      title: '创建成功',
      content: `房间ID: ${res.data.id}\n点击确定进入房间`,
      showCancel: false,
      success: () => {
        // 跳转到游戏房间
        uni.redirectTo({
          url: `/pages/poker/game-room?roomId=${res.data.id}`
        })
      }
    })
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
  font-size: 28rpx;
  box-sizing: border-box;
}

.player-count-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 15rpx;
}

.count-option {
  height: 70rpx;
  background-color: #f5f5f5;
  border-radius: 10rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 26rpx;
  color: #666;
  transition: all 0.3s;
}

.count-option.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  font-weight: bold;
}

.create-btn {
  width: 100%;
  height: 90rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 45rpx;
  font-size: 32rpx;
  font-weight: bold;
  color: #fff;
  border: none;
  margin-top: 40rpx;
}

.create-btn::after {
  border: none;
}
</style>
