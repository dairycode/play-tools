package main

import (
	"log"
	"play-tools/api"
	"play-tools/config"
	"play-tools/database"
	"play-tools/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置文件
	if err := config.Init(); err != nil {
		panic(err)
	}

	// 初始化数据库
	if err := database.Init(); err != nil {
		panic(err)
	}

	// 创建 Gin 引擎
	r := gin.Default()

	// 配置 CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = false // 允许 cookie 时不能使用 AllowAllOrigins
	corsConfig.AllowOriginFunc = func(origin string) bool { return true }
	corsConfig.AllowCredentials = true // 允许携带 cookie
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	r.Use(cors.New(corsConfig))

	// API 路由组
	apiGroup := r.Group("/api")
	{
		// 用户相关接口（无需认证）
		userGroup := apiGroup.Group("/user")
		{
			userGroup.POST("/register", api.Register)
			userGroup.POST("/login", api.Login)
			userGroup.POST("/wechat-login", api.WechatLogin) // 微信小程序登录
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
				gameGroup.POST("/room/:roomId/leave", api.LeaveRoom)
				gameGroup.GET("/room/current", api.GetCurrentRoom)
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
	serverAddr := config.AppConfig.Server.Host + ":" + config.AppConfig.Server.Port
	log.Printf("Server is running on %s\n", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		panic(err)
	}
}
