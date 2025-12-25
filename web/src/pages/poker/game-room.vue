<template>
  <view class="game-bg" style="min-height: 100vh; padding: 32rpx;">
    <view style="width: 100%; max-width: 1200rpx; margin: 0 auto; position: relative; z-index: 10;">
      <!-- Header -->
      <view class="animate-fade-in" style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 48rpx;">
        <view style="flex: 1;">
          <view class="text-glow text-white" style="font-size: 48rpx; font-weight: bold;">{{ roomDetail?.name }}</view>
          <view class="text-white" style="font-size: 28rpx; margin-top: 8rpx;">ID: {{ roomId }}</view>
        </view>
        <u-button
          v-if="roomDetail?.status === 'waiting'"
          text="邀请好友"
          @click="handleShare"
          :customStyle="{
            fontSize: '28rpx',
            padding: '16rpx 32rpx',
            height: 'auto',
            width: 'auto'
          }"
          :custom-class="'btn-secondary'"
        ></u-button>
        <u-button
          v-if="roomDetail?.isOwner && roomDetail?.status === 'playing'"
          text="结束游戏"
          @click="handleFinish"
          :customStyle="{
            fontSize: '28rpx',
            padding: '16rpx 32rpx',
            height: 'auto',
            width: 'auto'
          }"
          :custom-class="'btn-danger'"
        ></u-button>
      </view>

      <!-- Waiting Status -->
      <view v-if="roomDetail?.status === 'waiting'" class="animate-fade-in">
        <view class="text-white" style="text-align: center; font-size: 36rpx; font-weight: 600; margin-bottom: 32rpx;">
          房间成员 ({{ roomDetail?.members?.length || 0 }}/{{ roomDetail?.maxPlayers || 0 }}人)
        </view>

        <!-- Members Grid -->
        <view style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 24rpx; margin-bottom: 32rpx;">
          <view
            v-for="member in roomDetail?.members"
            :key="member.userId"
            class="glass-card"
            style="padding: 32rpx;"
          >
            <view style="display: flex; flex-direction: column; align-items: center;">
              <view class="btn-primary" style="width: 128rpx; height: 128rpx; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 56rpx; font-weight: bold; color: white; margin-bottom: 16rpx;">
                {{ member.nickname.substring(0, 1) }}
              </view>
              <view class="text-white" style="font-size: 28rpx; margin-bottom: 8rpx;">{{ member.nickname }}</view>
              <view :style="{ fontSize: '24rpx', color: member.isReady ? '#66bb6a' : 'rgba(255, 255, 255, 0.5)', fontWeight: member.isReady ? 'bold' : 'normal' }">
                {{ member.isReady ? '已准备' : '未准备' }}
              </view>
            </view>
          </view>
        </view>

        <!-- Action Buttons -->
        <view class="glass-card" style="padding: 48rpx; margin-top: 32rpx;">
          <u-button
            v-if="!roomDetail?.isOwner"
            :text="currentUserReady ? '取消准备' : '准备'"
            @click="handleToggleReady"
            :customStyle="{ width: '100%', fontSize: '32rpx', fontWeight: 'bold' }"
            :custom-class="currentUserReady ? 'btn-secondary' : 'btn-success'"
          ></u-button>
          <u-button
            v-if="roomDetail?.isOwner"
            :text="startButtonText"
            :disabled="!allReady"
            @click="handleStart"
            :customStyle="{ width: '100%', fontSize: '32rpx', fontWeight: 'bold' }"
            :custom-class="'btn-success'"
          ></u-button>
        </view>
      </view>

      <!-- Playing Status -->
      <view v-if="roomDetail?.status === 'playing'" class="animate-fade-in">
        <!-- Members Section -->
        <view class="glass-card" style="padding: 48rpx; margin-bottom: 32rpx;">
          <view class="text-white" style="font-size: 36rpx; font-weight: 600; margin-bottom: 32rpx;">
            房间成员 ({{ roomDetail?.members?.length || 0 }}人)
          </view>
          <view style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 24rpx;">
            <view
              v-for="member in roomDetail?.members"
              :key="member.userId"
              style="background: rgba(255, 255, 255, 0.05); backdrop-filter: blur(12px); border-radius: 12rpx; border: 2rpx solid rgba(255, 255, 255, 0.1); padding: 32rpx;"
            >
              <view style="display: flex; flex-direction: column; align-items: center;">
                <view class="btn-primary" style="width: 128rpx; height: 128rpx; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 56rpx; font-weight: bold; color: white; margin-bottom: 16rpx;">
                  {{ member.nickname.substring(0, 1) }}
                </view>
                <view class="text-white" style="font-size: 28rpx; margin-bottom: 8rpx;">{{ member.nickname }}</view>
                <view :style="{
                  fontSize: '32rpx',
                  fontWeight: 'bold',
                  color: member.balance > 0 ? '#f44336' : member.balance < 0 ? '#66bb6a' : 'rgba(255, 255, 255, 0.5)'
                }">
                  {{ member.balance >= 0 ? '+' : '' }}{{ member.balance }}
                </view>
              </view>
            </view>
          </view>
        </view>

        <!-- Transfer Section -->
        <view class="glass-card" style="padding: 48rpx; margin-bottom: 32rpx;">
          <view class="text-white" style="font-size: 36rpx; font-weight: 600; margin-bottom: 32rpx;">我要转账</view>

          <!-- Amount Input -->
          <u-input
            v-model.number="transferAmount"
            type="digit"
            placeholder="输入金额"
            color="#ffffff"
            :customStyle="{
              background: 'rgba(255, 255, 255, 0.1)',
              border: '2rpx solid rgba(255, 255, 255, 0.2)',
              borderRadius: '12rpx',
              textAlign: 'center',
              fontSize: '48rpx',
              fontWeight: 'bold',
              padding: '24rpx',
              marginBottom: '32rpx'
            }"
            placeholderStyle="color: rgba(255, 255, 255, 0.5)"
          ></u-input>

          <!-- Quick Amounts -->
          <view style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 16rpx; margin-bottom: 32rpx;">
            <view
              v-for="amount in quickAmounts"
              :key="amount"
              @click="setAmount(amount)"
              style="height: 96rpx; background: rgba(255, 255, 255, 0.1); border-radius: 12rpx; display: flex; align-items: center; justify-content: center; color: rgba(255, 255, 255, 0.7); font-size: 28rpx; cursor: pointer; transition: all 0.3s;"
            >
              {{ amount }}
            </view>
          </view>

          <!-- Target Members -->
          <view style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 24rpx;">
            <view
              v-for="member in otherMembers"
              :key="member.userId"
              @click="handleTransfer(member)"
              style="display: flex; flex-direction: column; align-items: center; cursor: pointer;"
            >
              <view class="btn-primary" style="width: 128rpx; height: 128rpx; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 48rpx; font-weight: bold; color: white; margin-bottom: 16rpx;">
                {{ member.nickname.substring(0, 1) }}
              </view>
              <view class="text-white-80" style="font-size: 24rpx; text-align: center;">{{ member.nickname }}</view>
            </view>
          </view>
        </view>

        <!-- Transaction History -->
        <view v-if="roomDetail?.transactions && roomDetail.transactions.length > 0" class="glass-card" style="padding: 48rpx;">
          <view class="text-white" style="font-size: 36rpx; font-weight: 600; margin-bottom: 32rpx;">转账记录</view>
          <view style="max-height: 768rpx; overflow-y: auto;">
            <view
              v-for="tx in roomDetail.transactions"
              :key="tx.id"
              style="display: flex; align-items: center; justify-content: space-between; padding: 24rpx; background: rgba(255, 255, 255, 0.05); border-radius: 12rpx; margin-bottom: 16rpx;"
            >
              <view style="display: flex; align-items: center; gap: 16rpx; font-size: 28rpx;">
                <text class="text-white-70">{{ tx.fromNickname }}</text>
                <text class="text-success">→</text>
                <text class="text-white" style="font-weight: bold;">{{ tx.toNickname }}</text>
              </view>
              <view class="text-danger" style="font-weight: bold;">¥{{ tx.amount }}</view>
            </view>
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
  const isFull = roomDetail.value.members.length === roomDetail.value.maxPlayers
  const allMembersReady = roomDetail.value.members.every(m => m.isReady)
  return isFull && allMembersReady
})

const startButtonText = computed(() => {
  if (!roomDetail.value?.members) return '等待玩家加入...'
  const currentCount = roomDetail.value.members.length
  const maxCount = roomDetail.value.maxPlayers
  const isFull = currentCount === maxCount
  const allMembersReady = roomDetail.value.members.every(m => m.isReady)

  if (!isFull) {
    return `等待玩家加入 (${currentCount}/${maxCount})`
  }
  if (!allMembersReady) {
    return '等待玩家准备...'
  }
  return '开始游戏'
})

const otherMembers = computed(() => {
  if (!roomDetail.value?.members) return []
  return roomDetail.value.members.filter(m => m.userId !== currentUserId.value)
})

onLoad(async (options: any) => {
  if (!isLogin()) {
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

    uni.navigateTo({
      url: '/pages/login/login?redirect=' + encodeURIComponent(currentPage)
    })
    return
  }

  if (options.roomId) {
    roomId.value = options.roomId

    if (options.join === '1') {
      try {
        await joinRoom({ roomId: roomId.value })
        uni.showToast({
          title: '成功加入房间',
          icon: 'success'
        })
      } catch (error: any) {
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
  refreshTimer = setInterval(() => {
    loadRoomInfo(true)
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
  let inviteUrl = ''

  // @ts-ignore
  if (typeof window !== 'undefined') {
    const origin = window.location.origin
    const pathname = window.location.pathname
    inviteUrl = `${origin}${pathname}#/pages/poker/game-room?roomId=${roomId.value}&join=1`
  } else {
    inviteUrl = `/pages/poker/game-room?roomId=${roomId.value}&join=1`
  }

  uni.setClipboardData({
    data: inviteUrl,
    success: () => {
      uni.showToast({
        title: '已复制邀请链接',
        icon: 'success'
      })
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
/* 使用全局样式 */
</style>
