<template>
  <view class="game-bg" style="min-height: 100vh;">
    <!-- Header Profile Card -->
    <view class="btn-primary" style="padding: 128rpx 32rpx 96rpx;">
      <view style="max-width: 700rpx; margin: 0 auto; display: flex; align-items: center;">
        <view style="width: 192rpx; height: 192rpx; border-radius: 50%; background: rgba(255, 255, 255, 0.3); backdrop-filter: blur(12px); display: flex; align-items: center; justify-content: center; margin-right: 48rpx;">
          <text style="font-size: 96rpx; font-weight: bold; color: white;">{{ avatarText }}</text>
        </view>
        <view style="flex: 1;">
          <view class="text-white" style="font-size: 56rpx; font-weight: bold; margin-bottom: 16rpx;">
            {{ userInfo?.nickname || '未设置昵称' }}
          </view>
          <view class="text-white-80" style="font-size: 28rpx;">@{{ userInfo?.username }}</view>
        </view>
      </view>
    </view>

    <!-- Menu List -->
    <view style="max-width: 700rpx; margin: 0 auto; padding: 0 32rpx; margin-top: -64rpx; position: relative; z-index: 10;">
      <view class="glass-card" style="padding: 16rpx;">
        <u-cell-group :border="false" :customStyle="{ background: 'transparent' }">
          <u-cell
            title="修改昵称"
            :isLink="true"
            @click="showEditNickname"
            :customStyle="{
              background: 'transparent',
              color: 'white'
            }"
          >
            <template #icon>
              <text style="font-size: 56rpx; margin-right: 32rpx;">✏️</text>
            </template>
          </u-cell>

          <u-line :color="'rgba(255, 255, 255, 0.1)'" :margin="'16rpx 0'"></u-line>

          <u-cell
            title="退出登录"
            :isLink="true"
            @click="handleLogout"
            :customStyle="{
              background: 'transparent',
              color: 'white'
            }"
          >
            <template #icon>
              <text style="font-size: 56rpx; margin-right: 32rpx;">🚪</text>
            </template>
          </u-cell>
        </u-cell-group>
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
/* 使用全局样式 */
</style>
