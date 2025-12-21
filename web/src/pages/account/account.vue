<template>
  <view class="min-h-screen game-bg">
    <!-- Header Profile Card -->
    <view class="bg-gradient-to-br from-primary-from to-primary-to pt-16 pb-12 px-4 md:px-6">
      <view class="w-full max-w-md mx-auto flex items-center">
        <view class="w-24 h-24 rounded-full bg-white/30 backdrop-blur-md flex items-center justify-center mr-6">
          <text class="text-5xl font-bold text-white">{{ avatarText }}</text>
        </view>
        <view class="flex-1">
          <view class="text-3xl font-bold text-white mb-2">{{ userInfo?.nickname || '未设置昵称' }}</view>
          <view class="text-sm text-white/80">@{{ userInfo?.username }}</view>
        </view>
      </view>
    </view>

    <!-- Menu List -->
    <view class="w-full max-w-md mx-auto px-4 md:px-6 -mt-8 relative z-10">
      <view class="glass-card p-2">
        <view
          class="flex items-center justify-between p-4 cursor-pointer transition-all duration-300 hover:bg-white/5 rounded-xl"
          @click="showEditNickname"
        >
          <view class="flex items-center">
            <text class="text-3xl mr-4">✏️</text>
            <text class="text-base text-white">修改昵称</text>
          </view>
          <text class="text-4xl text-white/30">›</text>
        </view>

        <view class="h-px bg-white/10 my-1"></view>

        <view
          class="flex items-center justify-between p-4 cursor-pointer transition-all duration-300 hover:bg-white/5 rounded-xl"
          @click="handleLogout"
        >
          <view class="flex items-center">
            <text class="text-3xl mr-4">🚪</text>
            <text class="text-base text-white">退出登录</text>
          </view>
          <text class="text-4xl text-white/30">›</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { getUserInfo as getStorageUserInfo, logout, isLogin } from '@/utils/auth'
import { getUserInfo as fetchUserInfo, updateUserInfo } from '@/api/user'
import { setUserInfo } from '@/utils/auth'
import type { User } from '@/types'

const userInfo = ref<User | null>(null)

const avatarText = computed(() => {
  return userInfo.value?.nickname?.substring(0, 1) || '?'
})

onLoad(() => {
  // 检查登录状态
  if (!isLogin()) {
    uni.navigateTo({
      url: '/pages/login/login?redirect=' + encodeURIComponent('/pages/account/account')
    })
    return
  }
})

onMounted(async () => {
  // 先从本地获取
  userInfo.value = getStorageUserInfo()

  // 再从服务器获取最新数据
  try {
    const res = await fetchUserInfo()
    userInfo.value = res.data
    setUserInfo(res.data)
  } catch (error) {
    console.error('获取用户信息失败', error)
  }
})

const showEditNickname = () => {
  uni.showModal({
    title: '修改昵称',
    editable: true,
    placeholderText: '请输入新昵称',
    content: userInfo.value?.nickname || '',
    success: async (res) => {
      if (res.confirm && res.content) {
        try {
          uni.showLoading({ title: '修改中...' })
          await updateUserInfo({ nickname: res.content })

          // 更新本地信息
          if (userInfo.value) {
            userInfo.value.nickname = res.content
            setUserInfo(userInfo.value)
          }

          uni.hideLoading()
          uni.showToast({
            title: '修改成功',
            icon: 'success'
          })
        } catch (error) {
          uni.hideLoading()
        }
      }
    }
  })
}

const handleLogout = () => {
  uni.showModal({
    title: '提示',
    content: '确定要退出登录吗？',
    success: (res) => {
      if (res.confirm) {
        logout()
      }
    }
  })
}
</script>

<style scoped>
/* No custom styles needed - all using Tailwind */
</style>
