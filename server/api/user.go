package api

import (
	"fmt"
	"net/http"
	"play-tools/database"
	"play-tools/model"
	"play-tools/utils"

	"github.com/gin-gonic/gin"
)

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// WechatLoginRequest 微信登录请求
type WechatLoginRequest struct {
	Code     string `json:"code" binding:"required"`     // 微信登录 code
	Nickname string `json:"nickname" binding:"required"` // 用户昵称
	Avatar   string `json:"avatar"`                      // 用户头像
}

// Register 用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 检查用户名是否已存在
	var existUser model.User
	if err := database.DB.Where("username = ?", req.Username).First(&existUser).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "用户名已存在",
		})
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "密码加密失败",
		})
		return
	}

	// 创建用户
	user := model.User{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Nickname,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建用户失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
		},
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 查找用户
	var user model.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 生成Token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "生成Token失败",
		})
		return
	}

	// 设置Cookie（7天有效期）
	c.SetCookie(
		"token",                  // cookie名称
		token,                    // cookie值
		7*24*60*60,               // 有效期（秒）
		"/",                      // 路径
		"",                       // 域名（空表示当前域名）
		false,                    // secure（生产环境建议true）
		true,                     // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":        user.ID,
				"username":  user.Username,
				"nickname":  user.Nickname,
				"avatar":    user.Avatar,
				"createdAt": user.CreatedAt,
			},
		},
	})
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	userID, _ := c.Get("userId")

	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
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
			"id":        user.ID,
			"username":  user.Username,
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

	if err := database.DB.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
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

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 查找用户
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "旧密码错误",
		})
		return
	}

	// 密码长度验证
	if len(req.NewPassword) < 6 {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "新密码长度不能少于6位",
		})
		return
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "密码加密失败",
		})
		return
	}

	// 更新密码
	if err := database.DB.Model(&model.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "修改密码失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "修改密码成功",
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
		user = model.User{
			WxOpenID:  session.OpenID,
			WxUnionID: session.UnionID,
			Nickname:  req.Nickname,
			Avatar:    req.Avatar,
			LoginType: "wechat",
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
		if session.UnionID != "" {
			updates["wx_union_id"] = session.UnionID
		}

		database.DB.Model(&user).Updates(updates)
	}

	// 生成 Token
	token, err := utils.GenerateToken(user.ID, user.WxOpenID)
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
				"id":        user.ID,
				"nickname":  user.Nickname,
				"avatar":    user.Avatar,
				"loginType": user.LoginType,
				"createdAt": user.CreatedAt,
			},
		},
	})
}
