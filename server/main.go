package main

import (
	"log"
	"play-tools/api"
	"play-tools/database"
	"play-tools/middleware"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	database.InitDB()

	// 创建Gin引擎
	r := gin.Default()

	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = false // 允许cookie时不能使用AllowAllOrigins
	config.AllowOriginFunc = func(origin string) bool {
		// 允许localhost和127.0.0.1
		if strings.HasPrefix(origin, "http://localhost") ||
				strings.HasPrefix(origin, "http://127.0.0.1") {
				return true
			}
			// 允许局域网IP（192.168.x.x 和 10.x.x.x）
			if strings.HasPrefix(origin, "http://192.168.") ||
				strings.HasPrefix(origin, "http://10.") {
				return true
			}
			return false
	}
	config.AllowCredentials = true // 允许携带cookie
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	r.Use(cors.New(config))

	// API路由组
	apiGroup := r.Group("/api")
	{
		// 用户相关接口（无需认证）
		userGroup := apiGroup.Group("/user")
		{
			userGroup.POST("/register", api.Register)
			userGroup.POST("/login", api.Login)
		}

		// 需要认证的接口
		authGroup := apiGroup.Group("")
		authGroup.Use(middleware.AuthMiddleware())
		{
			// 用户信息
			authGroup.GET("/user/info", api.GetUserInfo)
			authGroup.PUT("/user/update", api.UpdateUserInfo)
			authGroup.PUT("/user/password", api.ChangePassword)

			// 游戏相关接口
			gameGroup := authGroup.Group("/game")
			{
				gameGroup.POST("/room/create", api.CreateRoom)
				gameGroup.POST("/room/join", api.JoinRoom)
				gameGroup.GET("/room/:roomId", api.GetRoomInfo)
				gameGroup.GET("/room/:roomId/settlements", api.GetRoomSettlements)
				gameGroup.POST("/room/:roomId/ready", api.ToggleReady)
				gameGroup.POST("/room/:roomId/start", api.StartGame)
				gameGroup.POST("/transaction/add", api.AddTransaction)
				gameGroup.POST("/room/finish", api.FinishGame)
				gameGroup.GET("/rooms", api.GetMyRooms)
			}
		}
	}

	// 启动服务器
	log.Println("Server is running on :8080")
	log.Println("Access URLs:")
	log.Println("  Local:   http://localhost:8080")
	log.Println("  Network: http://192.168.0.102:8080")
	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
