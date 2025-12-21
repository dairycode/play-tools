<template>
  <view class="account-container">
    <view class="header">
      <view class="avatar">
        <text class="avatar-text">{{ avatarText }}</text>
      </view>
      <view class="user-info">
        <view class="nickname">{{ userInfo?.nickname || '未设置昵称' }}</view>
        <view class="username">@{{ userInfo?.username }}</view>
      </view>
    </view>

    <view class="menu-list">
      <view class="menu-item" @click="showEditNickname">
        <view class="menu-left">
          <text class="menu-icon">✏️</text>
          <text class="menu-title">修改昵称</text>
        </view>
        <text class="menu-arrow">›</text>
      </view>

      <view class="menu-item" @click="handleLogout">
        <view class="menu-left">
          <text class="menu-icon">🚪</text>
          <text class="menu-title">退出登录</text>
        </view>
        <text class="menu-arrow">›</text>
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
.account-container {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 60rpx 40rpx;
  display: flex;
  align-items: center;
}

.avatar {
  width: 120rpx;
  height: 120rpx;
  border-radius: 60rpx;
  background-color: rgba(255, 255, 255, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 30rpx;
}

.avatar-text {
  font-size: 48rpx;
  font-weight: bold;
  color: #fff;
}

.user-info {
  flex: 1;
}

.nickname {
  font-size: 36rpx;
  font-weight: bold;
  color: #fff;
  margin-bottom: 10rpx;
}

.username {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.8);
}

.menu-list {
  margin-top: 20rpx;
  background-color: #fff;
}

.menu-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 30rpx 40rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.menu-item:last-child {
  border-bottom: none;
}

.menu-left {
  display: flex;
  align-items: center;
}

.menu-icon {
  font-size: 36rpx;
  margin-right: 20rpx;
}

.menu-title {
  font-size: 30rpx;
  color: #333;
}

.menu-arrow {
  font-size: 48rpx;
  color: #ccc;
}
</style>
