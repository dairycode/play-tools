<template>
  <view class="game-bg">
    <view class="status-bar" :style="{ height: statusBarHeight + 'px' }"></view>

    <view style="padding: 32rpx; padding-top: 160rpx; position: relative; z-index: 10;">
      <!-- Title -->
      <view style="text-align: center; margin-bottom: 64rpx;">
        <view class="text-glow" style="font-size: 56rpx; font-weight: bold; color: white; margin-bottom: 16rpx;">
          完善个人信息
        </view>
        <view class="text-white-70" style="font-size: 28rpx;">设置您的头像和昵称</view>
      </view>

      <!-- Form Card -->
      <view class="glass-card animate-fade-in" style="padding: 48rpx;">
        <!-- Avatar -->
        <view style="display: flex; flex-direction: column; align-items: center; margin-bottom: 48rpx;">
          <view style="margin-bottom: 16rpx; color: rgba(255, 255, 255, 0.9); font-size: 28rpx;">
            点击上传头像
          </view>
          <!-- #ifdef MP-WEIXIN -->
          <button
            class="avatar-button"
            open-type="chooseAvatar"
            @chooseavatar="onChooseAvatar"
          >
            <image
              v-if="avatarUrl"
              :src="avatarUrl"
              mode="aspectFill"
              style="width: 160rpx; height: 160rpx; border-radius: 80rpx;"
            />
            <view
              v-else
              class="glass-card"
              style="width: 160rpx; height: 160rpx; border-radius: 80rpx; display: flex; align-items: center; justify-content: center;"
            >
              <text style="font-size: 80rpx;">👤</text>
            </view>
          </button>
          <!-- #endif -->
          <!-- #ifndef MP-WEIXIN -->
          <view
            class="glass-card"
            style="width: 160rpx; height: 160rpx; border-radius: 80rpx; display: flex; align-items: center; justify-content: center;"
            @click="chooseImage"
          >
            <image
              v-if="avatarUrl"
              :src="avatarUrl"
              mode="aspectFill"
              style="width: 160rpx; height: 160rpx; border-radius: 80rpx;"
            />
            <text v-else style="font-size: 80rpx;">👤</text>
          </view>
          <!-- #endif -->
        </view>

        <!-- Nickname -->
        <view style="margin-bottom: 32rpx;">
          <view style="margin-bottom: 16rpx; color: rgba(255, 255, 255, 0.9); font-size: 28rpx;">
            昵称
          </view>
          <!-- #ifdef MP-WEIXIN -->
          <input
            type="nickname"
            v-model="nickname"
            placeholder="请输入昵称"
            maxlength="20"
            @blur="onNicknameInput"
            :style="{
              background: 'rgba(255, 255, 255, 0.1)',
              border: '2rpx solid rgba(255, 255, 255, 0.2)',
              borderRadius: '12rpx',
              padding: '24rpx 32rpx',
              fontSize: '32rpx',
              color: '#ffffff'
            }"
            placeholder-style="color: rgba(255, 255, 255, 0.5)"
          />
          <!-- #endif -->
          <!-- #ifndef MP-WEIXIN -->
          <u-input
            v-model="nickname"
            placeholder="请输入昵称"
            :maxlength="20"
            color="#ffffff"
            :customStyle="{
              background: 'rgba(255, 255, 255, 0.1)',
              border: '2rpx solid rgba(255, 255, 255, 0.2)',
              borderRadius: '12rpx',
              padding: '24rpx 32rpx',
              fontSize: '32rpx'
            }"
            placeholderStyle="color: rgba(255, 255, 255, 0.5)"
          />
          <!-- #endif -->
        </view>

        <!-- Submit Button -->
        <u-button
          type="primary"
          text="完成"
          @click="handleSubmit"
          :customStyle="{
            width: '100%',
            marginBottom: '24rpx',
            fontSize: '32rpx',
            fontWeight: 'bold'
          }"
          :custom-class="'btn-primary'"
        ></u-button>

        <!-- Skip Button -->
        <view
          class="text-white-70"
          style="text-align: center; font-size: 28rpx; cursor: pointer;"
          @click="handleSkip"
        >
          暂时跳过
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { updateUserInfo, uploadAvatar as uploadAvatarApi } from '@/api/user'
import { getUserInfo, setUserInfo } from '@/utils/auth'

const statusBarHeight = ref(0)
const nickname = ref('')
const avatarUrl = ref('')
const redirectUrl = ref('/pages/index/index')

// 获取状态栏高度
uni.getSystemInfo({
  success: (res) => {
    statusBarHeight.value = res.statusBarHeight || 0
  }
})

// 获取当前用户信息和重定向路径
onLoad((options: any) => {
  if (options.redirect) {
    redirectUrl.value = decodeURIComponent(options.redirect)
  }

  // 获取当前登录用户信息作为默认值
  const userInfo = getUserInfo()
  if (userInfo) {
    nickname.value = userInfo.nickname || ''
    avatarUrl.value = userInfo.avatar || ''
  }
})

// 选择头像（微信小程序）
const onChooseAvatar = (e: any) => {
  const { avatarUrl: tempPath } = e.detail
  console.log('选择头像，临时路径:', tempPath)
  avatarUrl.value = tempPath  // 先显示临时图片

  // 立即上传到服务器获取永久URL
  uploadAvatar(tempPath)
}

// 选择头像(H5/其他平台)
const chooseImage = () => {
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      const tempPath = res.tempFilePaths[0]
      if (tempPath) {
        avatarUrl.value = tempPath
        uploadAvatar(tempPath)
      }
    }
  })
}

// 上传头像到服务器
const uploadAvatar = async (filePath: string) => {
  console.log('开始上传头像，文件路径:', filePath)
  
  // 验证文件路径
  if (!filePath) {
    console.error('文件路径为空')
    uni.showToast({
      title: '文件路径无效',
      icon: 'none'
    })
    return
  }

  try {
    uni.showLoading({
      title: '上传中...'
    })

    console.log('调用上传API...')
    const res = await uploadAvatarApi(filePath)
    console.log('上传API返回:', res)
    
    // 将临时路径替换为服务器返回的公开URL
    if (res.data && res.data.url) {
      avatarUrl.value = res.data.url
      console.log('头像上传成功，公开URL:', res.data.url)
      
      uni.hideLoading()
      uni.showToast({
        title: '上传成功',
        icon: 'success',
        duration: 1500
      })
    } else {
      console.error('服务器未返回URL:', res)
      uni.hideLoading()
      uni.showToast({
        title: '上传失败：未获取到URL',
        icon: 'none'
      })
    }
  } catch (error) {
    uni.hideLoading()
    console.error('上传头像异常:', error)
    // 打印完整的错误对象
    console.error('错误详情:', JSON.stringify(error))
    uni.showToast({
      title: '上传失败，请重试',
      icon: 'none'
    })
  }
}


// 昵称输入
const onNicknameInput = (e: any) => {
  // 微信小程序 type="nickname" 会自动处理昵称
  if (e.detail && e.detail.value) {
    nickname.value = e.detail.value
  }
}

// 提交
const handleSubmit = async () => {
  if (!nickname.value.trim()) {
    uni.showToast({
      title: '请输入昵称',
      icon: 'none'
    })
    return
  }

  try {
    uni.showLoading({
      title: '保存中...'
    })

    // 更新用户信息（只保存昵称，头像已在上传时保存）
    await updateUserInfo({
      nickname: nickname.value.trim()
    })

    // 更新本地存储的用户信息
    const userInfo = getUserInfo()
    if (userInfo) {
      userInfo.nickname = nickname.value.trim()
      // avatarUrl 已在上传时更新，同步到本地存储
      if (avatarUrl.value) {
        userInfo.avatar = avatarUrl.value
      }
      setUserInfo(userInfo)
    }

    uni.hideLoading()
    uni.showToast({
      title: '保存成功',
      icon: 'success'
    })

    // 跳转
    setTimeout(() => {
      navigateToHome()
    }, 1000)
  } catch (error) {
    uni.hideLoading()
  }
}

// 跳过
const handleSkip = () => {
  uni.showModal({
    title: '提示',
    content: '确定暂时跳过吗？您可以稍后在个人中心修改',
    success: (res) => {
      if (res.confirm) {
        navigateToHome()
      }
    }
  })
}

// 跳转到首页
const navigateToHome = () => {
  uni.switchTab({
    url: redirectUrl.value,
    fail: () => {
      uni.redirectTo({
        url: redirectUrl.value,
        fail: () => {
          uni.reLaunch({
            url: redirectUrl.value
          })
        }
      })
    }
  })
}
</script>

<style scoped>
.status-bar {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  z-index: 999;
}

.avatar-button {
  background: none;
  border: none;
  padding: 0;
  margin: 0;
  line-height: 1;
}

.avatar-button::after {
  border: none;
}
</style>
