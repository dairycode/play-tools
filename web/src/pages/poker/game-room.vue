<template>
  <view class="container">
    <view class="header">
      <view class="room-title">
        <view class="room-name">{{ roomDetail?.name }}</view>
        <view class="room-id">ID: {{ roomId }}</view>
      </view>
      <button
        v-if="roomDetail?.status === 'waiting'"
        class="share-btn"
        size="mini"
        @click="handleShare"
      >
        邀请好友
      </button>
      <button
        v-if="roomDetail?.isOwner && roomDetail?.status === 'playing'"
        class="finish-btn"
        size="mini"
        @click="handleFinish"
      >
        结束游戏
      </button>
    </view>

    <!-- 等待中状态 -->
    <view v-if="roomDetail?.status === 'waiting'" class="waiting-section">
      <view class="section-title">
        房间成员 ({{ roomDetail?.members?.length || 0 }}/{{ roomDetail?.maxPlayers || 0 }}人)
      </view>
      <view class="members-grid">
        <view
          v-for="member in roomDetail?.members"
          :key="member.userId"
          class="member-card"
        >
          <view class="member-avatar">
            {{ member.nickname.substring(0, 1) }}
          </view>
          <view class="member-name">{{ member.nickname }}</view>
          <view :class="['ready-status', { ready: member.isReady }]">
            {{ member.isReady ? '已准备' : '未准备' }}
          </view>
        </view>
      </view>

      <!-- 准备/开始按钮 -->
      <view class="action-btns">
        <button
          v-if="!roomDetail?.isOwner"
          :class="['ready-btn', { ready: currentUserReady }]"
          @click="handleToggleReady"
        >
          {{ currentUserReady ? '取消准备' : '准备' }}
        </button>
        <button
          v-if="roomDetail?.isOwner"
          class="start-btn"
          :disabled="!allReady"
          @click="handleStart"
        >
          {{ allReady ? '开始游戏' : '等待玩家准备...' }}
        </button>
      </view>
    </view>

    <!-- 游戏中状态 -->
    <view v-if="roomDetail?.status === 'playing'" class="playing-section">
      <view class="members-section">
        <view class="section-title">房间成员 ({{ roomDetail?.members?.length || 0 }}人)</view>
        <view class="members-grid">
          <view
            v-for="member in roomDetail?.members"
            :key="member.userId"
            class="member-card"
          >
            <view class="member-avatar">
              {{ member.nickname.substring(0, 1) }}
            </view>
            <view class="member-name">{{ member.nickname }}</view>
            <view :class="['member-balance', getBalanceClass(member.balance)]">
              {{ member.balance >= 0 ? '+' : '' }}{{ member.balance }}
            </view>
          </view>
        </view>
      </view>

      <view class="transfer-section">
        <view class="section-title">我要转账</view>
        <view class="transfer-form">
          <view class="amount-input-wrapper">
            <input
              class="amount-input"
              v-model.number="transferAmount"
              type="digit"
              placeholder="输入金额"
            />
          </view>

          <view class="quick-amounts">
            <view
              v-for="amount in quickAmounts"
              :key="amount"
              class="quick-amount"
              @click="setAmount(amount)"
            >
              {{ amount }}
            </view>
          </view>

          <view class="target-members">
            <view
              v-for="member in otherMembers"
              :key="member.userId"
              class="target-member"
              @click="handleTransfer(member)"
            >
              <view class="target-avatar">
                {{ member.nickname.substring(0, 1) }}
              </view>
              <view class="target-name">{{ member.nickname }}</view>
            </view>
          </view>
        </view>
      </view>

      <view v-if="roomDetail?.transactions && roomDetail.transactions.length > 0" class="history-section">
        <view class="section-title">转账记录</view>
        <view class="history-list">
          <view
            v-for="tx in roomDetail.transactions"
            :key="tx.id"
            class="history-item"
          >
            <view class="history-content">
              <text class="from-user">{{ tx.fromNickname }}</text>
              <text class="arrow">→</text>
              <text class="to-user">{{ tx.toNickname }}</text>
            </view>
            <view class="history-amount">¥{{ tx.amount }}</view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onUnmounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { getRoomInfo, addTransaction, finishGame, toggleReady, startGame, joinRoom } from '@/api/game'
import { getUserInfo, isLogin } from '@/utils/auth'
import type { RoomDetail, Member } from '@/types'

const roomId = ref('')
const roomDetail = ref<RoomDetail | null>(null)
const transferAmount = ref<number | string>('')
let refreshTimer: number | null = null

const quickAmounts = [10, 20, 50, 100, 200, 500]

const currentUserId = computed(() => {
  const userInfo = getUserInfo()
  return userInfo?.id
})

const currentUserReady = computed(() => {
  if (!roomDetail.value?.members) return false
  const currentMember = roomDetail.value.members.find(m => m.userId === currentUserId.value)
  return currentMember?.isReady || false
})

const allReady = computed(() => {
  if (!roomDetail.value?.members) return false
  return roomDetail.value.members.every(m => m.isReady)
})

const otherMembers = computed(() => {
  if (!roomDetail.value?.members) return []
  return roomDetail.value.members.filter(m => m.userId !== currentUserId.value)
})

onLoad(async (options: any) => {
  // 检查登录状态
  if (!isLogin()) {
    // 构建当前页面完整路径（包含所有参数）
    let currentPage = '/pages/poker/game-room'
    const params: string[] = []

    if (options.roomId) {
      params.push(`roomId=${options.roomId}`)
    }
    if (options.join) {
      params.push(`join=${options.join}`)
    }

    if (params.length > 0) {
      currentPage += '?' + params.join('&')
    }

    // 跳转到登录页，携带返回路径
    uni.navigateTo({
      url: '/pages/login/login?redirect=' + encodeURIComponent(currentPage)
    })
    return
  }

  if (options.roomId) {
    roomId.value = options.roomId

    // 如果是通过链接进入，先尝试加入房间
    if (options.join === '1') {
      try {
        await joinRoom({ roomId: roomId.value })
        uni.showToast({
          title: '成功加入房间',
          icon: 'success'
        })
      } catch (error: any) {
        // 如果已经加入过，忽略错误
        if (!error.msg?.includes('已经加入')) {
          uni.showToast({
            title: error.msg || '加入失败',
            icon: 'none'
          })
          return
        }
      }
    }

    loadRoomInfo()
    startAutoRefresh()
  } else {
    uni.showToast({
      title: '房间ID不存在',
      icon: 'none'
    })
  }
})

onUnmounted(() => {
  stopAutoRefresh()
})

const startAutoRefresh = () => {
  // 每3秒刷新一次房间信息
  refreshTimer = setInterval(() => {
    loadRoomInfo(true) // 静默刷新
  }, 3000) as unknown as number
}

const stopAutoRefresh = () => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

const loadRoomInfo = async (silent = false) => {
  try {
    const res = await getRoomInfo(roomId.value)
    roomDetail.value = res.data

    // 如果房间已结束，跳转到结算页面
    if (res.data.status === 'finished') {
      stopAutoRefresh()
      uni.redirectTo({
        url: `/pages/poker/settlement?roomId=${roomId.value}`
      })
    }
  } catch (error) {
    if (!silent) {
      console.error('加载房间信息失败', error)
    }
  }
}

const handleShare = () => {
  // 构建完整的邀请链接
  let inviteUrl = ''

  // H5环境下获取当前完整URL
  // @ts-ignore
  if (typeof window !== 'undefined') {
    const origin = window.location.origin
    const pathname = window.location.pathname
    inviteUrl = `${origin}${pathname}#/pages/poker/game-room?roomId=${roomId.value}&join=1`
  } else {
    // 非H5环境（小程序等），使用相对路径
    inviteUrl = `/pages/poker/game-room?roomId=${roomId.value}&join=1`
  }

  uni.showModal({
    title: '邀请好友',
    content: `房间ID: ${roomId.value}\n点击"复制链接"后发送给好友`,
    confirmText: '复制链接',
    cancelText: '复制ID',
    success: (res) => {
      if (res.confirm) {
        // 复制完整链接
        uni.setClipboardData({
          data: inviteUrl,
          success: () => {
            uni.showToast({
              title: '已复制邀请链接',
              icon: 'success'
            })
          }
        })
      } else if (res.cancel) {
        // 只复制房间ID
        uni.setClipboardData({
          data: roomId.value,
          success: () => {
            uni.showToast({
              title: '已复制房间ID',
              icon: 'success'
            })
          }
        })
      }
    }
  })
}

const handleToggleReady = async () => {
  try {
    await toggleReady(roomId.value)
    loadRoomInfo()
  } catch (error) {
    console.error('切换准备状态失败', error)
  }
}

const handleStart = async () => {
  try {
    uni.showLoading({ title: '开始游戏...' })
    await startGame(roomId.value)
    uni.hideLoading()
    loadRoomInfo()
  } catch (error) {
    uni.hideLoading()
  }
}

const getBalanceClass = (balance: number) => {
  if (balance > 0) return 'positive'
  if (balance < 0) return 'negative'
  return 'zero'
}

const setAmount = (amount: number) => {
  transferAmount.value = amount
}

const handleTransfer = async (member: Member) => {
  const amount = Number(transferAmount.value)
  if (!amount || amount <= 0) {
    uni.showToast({
      title: '请输入有效金额',
      icon: 'none'
    })
    return
  }

  try {
    await addTransaction({
      roomId: roomId.value,
      toUserId: member.userId,
      amount: amount
    })

    uni.showToast({
      title: `已转账给 ${member.nickname} ¥${amount}`,
      icon: 'success'
    })

    transferAmount.value = ''
    loadRoomInfo()
  } catch (error) {
    console.error('转账失败', error)
  }
}

const handleFinish = () => {
  uni.showModal({
    title: '确认结束',
    content: '确定要结束游戏并进行结算吗？',
    success: async (res) => {
      if (res.confirm) {
        try {
          uni.showLoading({ title: '结算中...' })

          const result = await finishGame(roomId.value)

          uni.hideLoading()

          uni.redirectTo({
            url: `/pages/poker/settlement?settlements=${encodeURIComponent(JSON.stringify(result.data))}`
          })
        } catch (error) {
          uni.hideLoading()
        }
      }
    }
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
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 30rpx;
}

.room-title {
  flex: 1;
}

.room-name {
  font-size: 36rpx;
  font-weight: bold;
  color: #fff;
  margin-bottom: 5rpx;
}

.room-id {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.8);
}

.share-btn,
.finish-btn {
  background-color: #fff;
  color: #1AAD19;
  font-weight: bold;
  border: none;
}

.waiting-section,
.playing-section {
  margin-top: 20rpx;
}

.waiting-section > .section-title {
  font-size: 30rpx;
  font-weight: bold;
  color: #fff;
  margin-bottom: 20rpx;
  text-align: center;
}

.members-section,
.transfer-section,
.history-section {
  background-color: rgba(255, 255, 255, 0.95);
  border-radius: 20rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
}

.section-title {
  font-size: 30rpx;
  font-weight: bold;
  color: #333;
  margin-bottom: 20rpx;
}

.members-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20rpx;
}

.member-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20rpx;
  background-color: #f5f5f5;
  border-radius: 12rpx;
}

.member-avatar {
  width: 60rpx;
  height: 60rpx;
  border-radius: 30rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24rpx;
  font-weight: bold;
  color: #fff;
  margin-bottom: 8rpx;
}

.member-name {
  font-size: 24rpx;
  color: #333;
  margin-bottom: 5rpx;
  text-align: center;
  word-break: break-all;
}

.ready-status {
  font-size: 22rpx;
  color: #999;
}

.ready-status.ready {
  color: #1AAD19;
  font-weight: bold;
}

.member-balance {
  font-size: 26rpx;
  font-weight: bold;
}

.member-balance.positive {
  color: #ff4444;
}

.member-balance.negative {
  color: #1AAD19;
}

.member-balance.zero {
  color: #999;
}

.action-btns {
  margin-top: 30rpx;
  background-color: rgba(255, 255, 255, 0.95);
  border-radius: 20rpx;
  padding: 30rpx;
}

.ready-btn,
.start-btn {
  width: 100%;
  height: 90rpx;
  border-radius: 45rpx;
  font-size: 32rpx;
  font-weight: bold;
  border: none;
}

.ready-btn {
  background-color: #1AAD19;
  color: #fff;
}

.ready-btn.ready {
  background-color: #999;
}

.ready-btn::after,
.start-btn::after {
  border: none;
}

.start-btn {
  background-color: #1AAD19;
  color: #fff;
}

.start-btn:disabled {
  background-color: #ccc;
  color: #999;
}

.amount-input-wrapper {
  margin-bottom: 20rpx;
}

.amount-input {
  width: 100%;
  height: 80rpx;
  background-color: #f5f5f5;
  border-radius: 10rpx;
  padding: 0 30rpx;
  font-size: 36rpx;
  font-weight: bold;
  text-align: center;
  box-sizing: border-box;
}

.quick-amounts {
  display: flex;
  flex-wrap: wrap;
  gap: 15rpx;
  margin-bottom: 30rpx;
}

.quick-amount {
  flex: 0 0 calc(33.333% - 10rpx);
  height: 60rpx;
  background-color: #f5f5f5;
  border-radius: 10rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 26rpx;
  color: #666;
}

.target-members {
  display: flex;
  flex-wrap: wrap;
  gap: 20rpx;
}

.target-member {
  flex: 0 0 calc(25% - 15rpx);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.target-avatar {
  width: 70rpx;
  height: 70rpx;
  border-radius: 35rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28rpx;
  font-weight: bold;
  color: #fff;
  margin-bottom: 8rpx;
}

.target-name {
  font-size: 22rpx;
  color: #666;
  text-align: center;
  word-break: break-all;
}

.history-list {
  max-height: 400rpx;
  overflow-y: auto;
}

.history-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20rpx;
  background-color: #f5f5f5;
  border-radius: 10rpx;
  margin-bottom: 10rpx;
}

.history-item:last-child {
  margin-bottom: 0;
}

.history-content {
  flex: 1;
  font-size: 26rpx;
}

.from-user {
  color: #666;
}

.arrow {
  margin: 0 10rpx;
  color: #1AAD19;
}

.to-user {
  color: #333;
  font-weight: bold;
}

.history-amount {
  font-size: 28rpx;
  font-weight: bold;
  color: #ff4444;
}
</style>
