package api

import (
	"fmt"
	"net/http"
	"path/filepath"
	"play-tools/config"
	"play-tools/database"
	"play-tools/model"
	"play-tools/utils"

	"github.com/gin-gonic/gin"
)

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// WechatLoginRequest 微信登录请求
type WechatLoginRequest struct {
	Code     string `json:"code" binding:"required"`     // 微信登录 code
	Nickname string `json:"nickname" binding:"required"` // 用户昵称
	Avatar   string `json:"avatar"`                      // 用户头像
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	userID, _ := c.Get("userId")

	var user model.User
	if err := database.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"userId":    user.UserID,
			"nickname":  user.Nickname,
			"avatar":    user.Avatar,
			"createdAt": user.CreatedAt,
		},
	})
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if err := database.DB.Model(&model.User{}).Where("user_id = ?", userID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// WechatLogin 微信小程序登录
func WechatLogin(c *gin.Context) {
	var req WechatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 调用微信接口获取 openid
	session, err := utils.GetWechatOpenID(req.Code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  fmt.Sprintf("微信登录失败: %v", err),
		})
		return
	}

	// 查找是否已存在该 openid 的用户
	var user model.User
	result := database.DB.Where("wx_open_id = ?", session.OpenID).First(&user)

	if result.Error != nil {
		// 用户不存在，创建新用户
		// 生成8位用户ID
		userID := utils.GenerateUserID()

		// 检查生成的 UserID 是否已存在，如果存在则重新生成
		for {
			var existUser model.User
			if err := database.DB.Where("user_id = ?", userID).First(&existUser).Error; err != nil {
				// 不存在，可以使用
				break
			}
			// 存在，重新生成
			userID = utils.GenerateUserID()
		}

		user = model.User{
			UserID:   userID,
			WxOpenID: session.OpenID,
			Nickname: req.Nickname,
			Avatar:   req.Avatar,
		}

		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "创建用户失败",
			})
			return
		}
	} else {
		// 用户已存在，更新昵称和头像
		updates := map[string]interface{}{
			"nickname": req.Nickname,
		}
		if req.Avatar != "" {
			updates["avatar"] = req.Avatar
		}

		database.DB.Model(&user).Updates(updates)
	}

	// 生成 Token，使用 UserID
	token, err := utils.GenerateToken(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "生成Token失败",
		})
		return
	}

	// 设置 Cookie（7天有效期）
	c.SetCookie(
		"token",
		token,
		7*24*60*60,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"userId":    user.UserID,
				"nickname":  user.Nickname,
				"avatar":    user.Avatar,
				"createdAt": user.CreatedAt,
			},
		},
	})
}

// UploadAvatar 上传头像
func UploadAvatar(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未登录",
		})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "未找到上传的文件",
		})
		return
	}

	// 验证文件
	maxSize := config.AppConfig.Upload.MaxSize
	if err := utils.ValidateImageFile(file, maxSize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	// 使用 userId + 文件扩展名作为文件名
	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%d%s", userID.(uint), ext)

	// 构建保存路径
	avatarPath := config.AppConfig.Upload.AvatarPath
	destPath := fmt.Sprintf("%s/%s", avatarPath, newFileName)

	// 保存文件
	if err := utils.SaveUploadedFile(file, destPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  fmt.Sprintf("保存文件失败: %v", err),
		})
		return
	}

	// 构建公开访问的URL
	avatarURL := fmt.Sprintf("/uploads/avatars/%s", newFileName)

	// 更新用户头像信息到数据库（可选）
	database.DB.Model(&model.User{}).Where("user_id = ?", userID).Update("avatar", avatarURL)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "上传成功",
		"data": gin.H{
			"url": avatarURL,
		},
	})
}
