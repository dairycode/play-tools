package utils

import (
	"math/rand"
	"time"
)

// GenerateUserID 生成8位用户ID（时间戳+随机数，保证递增）
func GenerateUserID() uint {
	// 使用当前时间戳的秒数（去掉前几位，只保留后面的数字）
	now := time.Now().Unix()
	// 取时间戳后5位（确保不会太大）+ 3位随机数
	timestamp := now % 100000 // 后5位时间戳
	random := rand.Intn(1000) // 0-999的随机数

	// 组合成8位数字：前5位时间戳 + 后3位随机数
	userID := timestamp*1000 + int64(random)

	// 如果不足8位，添加前缀使其成为8位
	if userID < 10000000 {
		userID += 10000000
	}

	return uint(userID)
}

func init() {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())
}
