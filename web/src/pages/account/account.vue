<template>
  <view class="game-bg" style="min-height: 100vh;">
    <!-- Header Profile Card -->
    <view class="btn-primary" style="padding: 128rpx 32rpx 96rpx;">
      <view style="max-width: 700rpx; margin: 0 auto; display: flex; align-items: center;">
        <view style="width: 192rpx; height: 192rpx; border-radius: 50%; background: rgba(255, 255, 255, 0.3); backdrop-filter: blur(12px); display: flex; align-items: center; justify-content: center; margin-right: 48rpx; overflow: hidden;">
          <image
            v-if="userInfo?.avatar"
            :src="userInfo.avatar"
            mode="aspectFill"
            style="width: 100%; height: 100%;"
          />
          <text v-else style="font-size: 96rpx; font-weight: bold; color: white;">{{ avatarText }}</text>
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
            :border="false"
            @click="showEditNickname"
            :customStyle="{
              background: 'transparent'
            }"
          >
            <template #icon>
              <text style="font-size: 56rpx; margin-right: 32rpx;">✏️</text>
            </template>
          </u-cell>

          <!-- 只有账号密码登录的用户才显示修改密码选项 -->
          <template v-if="userInfo?.loginType === 'normal'">
            <u-line :color="'rgba(255, 255, 255, 0.1)'" :margin="'16rpx 0'"></u-line>

            <u-cell
              title="修改密码"
              :isLink="true"
              :border="false"
              @click="showChangePassword"
              :customStyle="{
                background: 'transparent'
              }"
            >
              <template #icon>
                <text style="font-size: 56rpx; margin-right: 32rpx;">🔒</text>
              </template>
            </u-cell>

            <u-line :color="'rgba(255, 255, 255, 0.1)'" :margin="'16rpx 0'"></u-line>
          </template>

          <u-cell
            title="退出登录"
            :isLink="true"
            :border="false"
            @click="handleLogout"
            :customStyle="{
              background: 'transparent'
            }"
          >
            <template #icon>
              <text style="font-size: 56rpx; margin-right: 32rpx;">🚪</text>
            </template>
          </u-cell>
        </u-cell-group>
      </view>
    </view>

    <!-- 修改昵称弹窗 -->
    <view v-if="showNicknameModal" class="password-modal-overlay" @click="resetNicknameForm">
      <view class="password-modal-content" @click.stop>
        <!-- 标题 -->
        <view class="modal-title">修改昵称</view>

        <!-- 表单内容 -->
        <view class="modal-form">
          <view class="form-item">
            <view class="form-label">新昵称</view>
            <input
              v-model="nicknameForm.nickname"
              type="text"
              placeholder="请输入新昵称"
              class="form-input"
            />
          </view>

          <!-- 错误提示 -->
          <view v-if="nicknameError" class="error-message">{{ nicknameError }}</view>
        </view>

        <!-- 按钮区域 -->
        <view class="modal-actions">
          <button class="modal-btn cancel-btn" @click.stop="resetNicknameForm">取消</button>
          <button class="modal-btn confirm-btn" @click.stop="handleUpdateNickname">确定</button>
        </view>
      </view>
    </view>

    <!-- 修改密码弹窗 -->
    <view v-if="showPasswordModal" class="password-modal-overlay" @click="resetPasswordForm">
      <view class="password-modal-content" @click.stop>
        <!-- 标题 -->
        <view class="modal-title">修改密码</view>

        <!-- 表单内容 -->
        <view class="modal-form">
          <view class="form-item">
            <view class="form-label">旧密码</view>
            <input
              v-model="passwordForm.oldPassword"
              type="password"
              placeholder="请输入旧密码"
              class="form-input"
            />
          </view>
          <view class="form-item">
            <view class="form-label">新密码</view>
            <input
              v-model="passwordForm.newPassword"
              type="password"
              placeholder="请输入新密码（至少6位）"
              class="form-input"
            />
          </view>
          <view class="form-item">
            <view class="form-label">确认密码</view>
            <input
              v-model="passwordForm.confirmPassword"
              type="password"
              placeholder="请再次输入新密码"
              class="form-input"
            />
          </view>

          <!-- 错误提示 -->
          <view v-if="errorMessage" class="error-message">{{ errorMessage }}</view>
        </view>

        <!-- 按钮区域 -->
        <view class="modal-actions">
          <button class="modal-btn cancel-btn" @click.stop="resetPasswordForm">取消</button>
          <button class="modal-btn confirm-btn" @click.stop="handleChangePassword">确定</button>
        </view>
      </view>
    </view>

    <!-- 退出登录弹窗 -->
    <view v-if="showLogoutModal" class="password-modal-overlay" @click="showLogoutModal = false">
      <view class="password-modal-content" @click.stop>
        <!-- 标题 -->
        <view class="modal-title">退出登录</view>

        <!-- 表单内容 -->
        <view class="modal-form">
          <view style="text-align: center; color: #666; font-size: 28rpx;">
            确定要退出登录吗？
          </view>
        </view>

        <!-- 按钮区域 -->
        <view class="modal-actions">
          <button class="modal-btn cancel-btn" @click.stop="showLogoutModal = false">取消</button>
          <button class="modal-btn confirm-btn" @click.stop="confirmLogout">确定</button>
        </view>
      </view>
    </view>

    <!-- 自定义 Toast -->
    <view v-if="toastVisible" class="custom-toast">
      {{ toastMessage }}
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { getUserInfo as getStorageUserInfo, logout, isLogin } from '@/utils/auth'
import { getUserInfo as fetchUserInfo, updateUserInfo, changePassword } from '@/api/user'
import { setUserInfo } from '@/utils/auth'
import type { User } from '@/types'

const userInfo = ref<User | null>(null)
const showNicknameModal = ref(false)
const showPasswordModal = ref(false)
const showLogoutModal = ref(false)
const nicknameForm = ref({
  nickname: ''
})
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})
const nicknameError = ref('')
const errorMessage = ref('')
const toastVisible = ref(false)
const toastMessage = ref('')

// 自定义 toast 函数
const showToast = (message: string, duration = 2000) => {
  toastMessage.value = message
  toastVisible.value = true
  setTimeout(() => {
    toastVisible.value = false
  }, duration)
}

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
  nicknameForm.value.nickname = userInfo.value?.nickname || ''
  showNicknameModal.value = true
}

const resetNicknameForm = () => {
  nicknameForm.value.nickname = ''
  nicknameError.value = ''
  showNicknameModal.value = false
}

const handleUpdateNickname = async () => {
  const nickname = nicknameForm.value.nickname.trim()

  // 清空之前的错误信息
  nicknameError.value = ''

  // 验证
  if (!nickname) {
    nicknameError.value = '请输入新昵称'
    return
  }

  if (nickname.length > 20) {
    nicknameError.value = '昵称长度不能超过20个字符'
    return
  }

  try {
    uni.showLoading({ title: '修改中...' })
    await updateUserInfo({ nickname })

    // 更新本地信息
    if (userInfo.value) {
      userInfo.value.nickname = nickname
      setUserInfo(userInfo.value)
    }

    uni.hideLoading()
    showToast('修改成功')

    // 重置表单并关闭弹窗
    resetNicknameForm()
  } catch (error) {
    uni.hideLoading()
    nicknameError.value = '修改失败，请稍后重试'
  }
}

const showChangePassword = () => {
  showPasswordModal.value = true
}

const resetPasswordForm = () => {
  passwordForm.value = {
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
  }
  errorMessage.value = ''
  showPasswordModal.value = false
}

const handleChangePassword = async () => {
  const { oldPassword, newPassword, confirmPassword } = passwordForm.value

  // 清空之前的错误信息
  errorMessage.value = ''

  console.log('handleChangePassword called', { oldPassword, newPassword, confirmPassword })

  // 验证
  if (!oldPassword || !newPassword || !confirmPassword) {
    console.log('验证失败：请填写完整信息')
    errorMessage.value = '请填写完整信息'
    return
  }

  if (newPassword.length < 6) {
    console.log('验证失败：新密码长度不能少于6位')
    errorMessage.value = '新密码长度不能少于6位'
    return
  }

  if (newPassword !== confirmPassword) {
    console.log('验证失败：两次输入的新密码不一致')
    errorMessage.value = '两次输入的新密码不一致'
    return
  }

  try {
    console.log('开始调用修改密码 API')
    uni.showLoading({ title: '修改中...' })
    await changePassword({
      old_password: oldPassword,
      new_password: newPassword
    })

    uni.hideLoading()
    showToast('修改成功，请重新登录', 2000)

    // 重置表单并关闭弹窗
    resetPasswordForm()

    // 2秒后退出登录
    setTimeout(() => {
      logout()
    }, 2000)
  } catch (error) {
    console.error('修改密码失败', error)
    uni.hideLoading()
    errorMessage.value = '修改密码失败，请检查旧密码是否正确'
  }
}

const handleLogout = () => {
  showLogoutModal.value = true
}

const confirmLogout = () => {
  showLogoutModal.value = false
  logout()
}
</script>

<style scoped>
/* 使用全局样式 */

/* 覆盖 uView Cell 组件的默认样式 */
:deep(.u-cell__title-text) {
  color: white !important;
}

:deep(.u-cell__value) {
  color: rgba(255, 255, 255, 0.7) !important;
}

:deep(.u-icon) {
  color: white !important;
}

/* 修改密码弹窗样式 */
.password-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.password-modal-content {
  width: 600rpx;
  background: white;
  border-radius: 20rpx;
  overflow: hidden;
}

.modal-title {
  padding: 32rpx;
  text-align: center;
  font-size: 32rpx;
  font-weight: bold;
  border-bottom: 1px solid #f0f0f0;
}

.modal-form {
  padding: 32rpx;
}

.form-item {
  margin-bottom: 32rpx;
}

.form-item:last-child {
  margin-bottom: 0;
}

.form-label {
  color: #666;
  font-size: 28rpx;
  margin-bottom: 16rpx;
}

.form-input {
  width: 100%;
  height: 80rpx;
  padding: 0 24rpx;
  background: #f5f5f5;
  border-radius: 12rpx;
  font-size: 28rpx;
  box-sizing: border-box;
  border: none;
}

.error-message {
  margin-top: 24rpx;
  padding: 20rpx;
  background: #fff3f3;
  border: 1px solid #ffcdd2;
  border-radius: 8rpx;
  color: #d32f2f;
  font-size: 26rpx;
  text-align: center;
  animation: shake 0.3s;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-10rpx); }
  75% { transform: translateX(10rpx); }
}

.modal-actions {
  display: flex;
  border-top: 1px solid #f0f0f0;
}

.modal-btn {
  flex: 1;
  padding: 32rpx;
  text-align: center;
  font-size: 28rpx;
  background: transparent;
  border: none;
  cursor: pointer;
}

.modal-btn::after {
  border: none;
}

.cancel-btn {
  color: #666;
  border-right: 1px solid #f0f0f0;
}

.confirm-btn {
  color: #007aff;
  font-weight: bold;
}

.modal-btn:active {
  background: #f5f5f5;
}

/* 自定义 Toast 样式 */
.custom-toast {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: rgba(0, 0, 0, 0.8);
  color: white;
  padding: 24rpx 48rpx;
  border-radius: 12rpx;
  font-size: 28rpx;
  z-index: 10000;
  animation: fadeInOut 2s ease-in-out;
  max-width: 80%;
  text-align: center;
  word-break: break-word;
}

@keyframes fadeInOut {
  0% { opacity: 0; transform: translate(-50%, -50%) scale(0.9); }
  10% { opacity: 1; transform: translate(-50%, -50%) scale(1); }
  90% { opacity: 1; transform: translate(-50%, -50%) scale(1); }
  100% { opacity: 0; transform: translate(-50%, -50%) scale(0.9); }
}
</style>
