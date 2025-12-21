<template>
  <view class="min-h-screen game-bg p-4 md:p-6">
    <view class="w-full max-w-4xl mx-auto relative z-10">
      <!-- Header -->
      <view class="flex items-center justify-between mb-6 animate-fade-in">
        <view class="flex-1">
          <view class="text-2xl font-bold text-white text-glow">{{ roomDetail?.name }}</view>
          <view class="text-sm text-white/80 mt-1">ID: {{ roomId }}</view>
        </view>
        <button
          v-if="roomDetail?.status === 'waiting'"
          class="btn-secondary text-sm"
          @click="handleShare"
        >
          邀请好友
        </button>
        <button
          v-if="roomDetail?.isOwner && roomDetail?.status === 'playing'"
          class="btn-danger text-sm"
          @click="handleFinish"
        >
          结束游戏
        </button>
      </view>

      <!-- Waiting Status -->
      <view v-if="roomDetail?.status === 'waiting'" class="animate-fade-in space-y-4">
        <view class="text-center text-lg font-semibold text-white mb-4">
          房间成员 ({{ roomDetail?.members?.length || 0 }}/{{ roomDetail?.maxPlayers || 0 }}人)
        </view>

        <!-- Members Grid -->
        <view class="grid grid-cols-2 md:grid-cols-3 gap-3">
          <view
            v-for="member in roomDetail?.members"
            :key="member.userId"
            class="glass-card p-4"
          >
            <view class="flex flex-col items-center">
              <view class="w-16 h-16 rounded-full bg-gradient-to-br from-primary-from to-primary-to flex items-center justify-center text-white font-bold text-xl mb-2">
                {{ member.nickname.substring(0, 1) }}
              </view>
              <view class="text-white text-sm mb-1">{{ member.nickname }}</view>
              <view :class="['text-xs', member.isReady ? 'text-success-from font-bold' : 'text-white/50']">
                {{ member.isReady ? '已准备' : '未准备' }}
              </view>
            </view>
          </view>
        </view>

        <!-- Action Buttons -->
        <view class="glass-card p-6 mt-4">
          <button
            v-if="!roomDetail?.isOwner"
            :class="['w-full', currentUserReady ? 'btn-secondary' : 'btn-success']"
            @click="handleToggleReady"
          >
            {{ currentUserReady ? '取消准备' : '准备' }}
          </button>
          <button
            v-if="roomDetail?.isOwner"
            class="w-full btn-success"
            :disabled="!allReady"
            @click="handleStart"
          >
            {{ allReady ? '开始游戏' : '等待玩家准备...' }}
          </button>
        </view>
      </view>

      <!-- Playing Status -->
      <view v-if="roomDetail?.status === 'playing'" class="animate-fade-in space-y-4">
        <!-- Members Section -->
        <view class="glass-card p-6">
          <view class="text-lg font-semibold text-white mb-4">
            房间成员 ({{ roomDetail?.members?.length || 0 }}人)
          </view>
          <view class="grid grid-cols-2 md:grid-cols-3 gap-3">
            <view
              v-for="member in roomDetail?.members"
              :key="member.userId"
              class="bg-white/5 backdrop-blur-md rounded-xl border-2 border-white/10 p-4"
            >
              <view class="flex flex-col items-center">
                <view class="w-16 h-16 rounded-full bg-gradient-to-br from-primary-from to-primary-to flex items-center justify-center text-white font-bold text-xl mb-2">
                  {{ member.nickname.substring(0, 1) }}
                </view>
                <view class="text-white text-sm mb-1">{{ member.nickname }}</view>
                <view :class="[
                  'text-base font-bold',
                  member.balance > 0 ? 'text-danger-from' : member.balance < 0 ? 'text-success-from' : 'text-white/50'
                ]">
                  {{ member.balance >= 0 ? '+' : '' }}{{ member.balance }}
                </view>
              </view>
            </view>
          </view>
        </view>

        <!-- Transfer Section -->
        <view class="glass-card p-6">
          <view class="text-lg font-semibold text-white mb-4">我要转账</view>

          <!-- Amount Input -->
          <input
            class="w-full px-4 py-3 bg-white/10 rounded-xl text-white text-center text-xl font-bold mb-4 border-2 border-white/20"
            v-model.number="transferAmount"
            type="digit"
            placeholder="输入金额"
          />

          <!-- Quick Amounts -->
          <view class="grid grid-cols-3 gap-2 mb-4">
            <view
              v-for="amount in quickAmounts"
              :key="amount"
              class="h-12 bg-white/10 rounded-xl flex items-center justify-center text-white/70 text-sm cursor-pointer hover:bg-white/15 transition-all"
              @click="setAmount(amount)"
            >
              {{ amount }}
            </view>
          </view>

          <!-- Target Members -->
          <view class="grid grid-cols-4 gap-3">
            <view
              v-for="member in otherMembers"
              :key="member.userId"
              class="flex flex-col items-center cursor-pointer"
              @click="handleTransfer(member)"
            >
              <view class="w-16 h-16 rounded-full bg-gradient-to-br from-primary-from to-primary-to flex items-center justify-center text-white font-bold text-lg mb-2">
                {{ member.nickname.substring(0, 1) }}
              </view>
              <view class="text-white/80 text-xs text-center">{{ member.nickname }}</view>
            </view>
          </view>
        </view>

        <!-- Transaction History -->
        <view v-if="roomDetail?.transactions && roomDetail.transactions.length > 0" class="glass-card p-6">
          <view class="text-lg font-semibold text-white mb-4">转账记录</view>
          <view class="space-y-2 max-h-96 overflow-y-auto">
            <view
              v-for="tx in roomDetail.transactions"
              :key="tx.id"
              class="flex items-center justify-between p-3 bg-white/5 rounded-xl"
            >
              <view class="flex items-center gap-2 text-sm">
                <text class="text-white/70">{{ tx.fromNickname }}</text>
                <text class="text-success-from">→</text>
                <text class="text-white font-bold">{{ tx.toNickname }}</text>
              </view>
              <view class="text-danger-from font-bold">¥{{ tx.amount }}</view>
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
button::after {
  border: none;
}

/* Custom scrollbar for transaction history */
.max-h-96::-webkit-scrollbar {
  width: 4px;
}

.max-h-96::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 4px;
}

.max-h-96::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 4px;
}

.max-h-96::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}
</style>
