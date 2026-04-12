package main

import (
	"fmt"
	"log"
	"os"

	"server/internal/api"
	"server/internal/auth"
	"server/internal/db"
	"server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "plantation"
	}

	if err := db.InitDB(host, port, user, password, dbname); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	// 初始化JWT
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key"
	}
	auth.InitJWTSecret(jwtSecret)

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由
	r := gin.Default()

	// 配置CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// 认证路由
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", api.Register)
		authGroup.POST("/login", api.Login)
	}

	// 需要认证的路由
	protectedGroup := r.Group("/api")
	protectedGroup.Use(middleware.AuthMiddleware())
	{
		// 种子信息管理 - 仅serverseed角色
		seedGroup := protectedGroup.Group("/seed")
		seedGroup.Use(middleware.RoleMiddleware("serverseed"))
		{
			seedGroup.POST("", api.CreateSeedInfo)
			seedGroup.GET("", api.ListSeedInfo)
			seedGroup.GET("/:id", api.GetSeedInfo)
			seedGroup.PUT("/:id", api.UpdateSeedInfo)
			seedGroup.DELETE("/:id", api.DeleteSeedInfo)
		}

		// 种植管理 - 仅serverseed角色
	plantingGroup := protectedGroup.Group("/planting")
	plantingGroup.Use(middleware.RoleMiddleware("serverseed"))
	{
		plantingGroup.POST("", api.CreatePlanting)
		plantingGroup.GET("", api.ListPlanting)
		plantingGroup.GET("/:id", api.GetPlanting)
		plantingGroup.PUT("/:id", api.UpdatePlanting)
		plantingGroup.DELETE("/:id", api.DeletePlanting)
		plantingGroup.GET("/:id/qrcode", api.GenerateQRCode)
	}

		// 生长媒体管理 - 仅servergrow角色
		growthGroup := protectedGroup.Group("/growth")
		growthGroup.Use(middleware.RoleMiddleware("servergrow"))
		{
			growthGroup.POST("", api.CreateGrowthMedia)
			growthGroup.GET("/planting/:planting_id", api.ListGrowthMediaByPlanting)
			growthGroup.GET("/:id", api.GetGrowthMedia)
			growthGroup.DELETE("/:id", api.DeleteGrowthMedia)
		}

		// 产品品质管理 - 仅servermanager角色
		qualityGroup := protectedGroup.Group("/quality")
		qualityGroup.Use(middleware.RoleMiddleware("servermanager"))
		{
			qualityGroup.POST("", api.CreateProductQuality)
			qualityGroup.GET("/planting/:planting_id", api.GetProductQualityByPlanting)
			qualityGroup.GET("/:id", api.GetProductQuality)
			qualityGroup.PUT("/:id", api.UpdateProductQuality)
			qualityGroup.DELETE("/:id", api.DeleteProductQuality)
		}

		// 用户交互 - 仅clentcustomer角色
		userGroup := protectedGroup.Group("/user")
		userGroup.Use(middleware.RoleMiddleware("clentcustomer"))
		{
			// 收藏功能
			userGroup.POST("/favorite/:planting_id", api.AddFavorite)
			userGroup.DELETE("/favorite/:planting_id", api.RemoveFavorite)
			userGroup.GET("/favorites", api.ListFavorites)

			// 点赞功能
			userGroup.POST("/like/:planting_id", api.AddLike)
			userGroup.DELETE("/like/:planting_id", api.RemoveLike)
			userGroup.GET("/like/:planting_id", api.CheckLikeStatus)
		}

		// 产品查询 - 所有角色都可以访问
		protectedGroup.GET("/products", api.ListPlanting)
		protectedGroup.GET("/products/:id", api.GetPlanting)
	}

	// 启动服务器
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	log.Printf("Server starting on port %s", serverPort)
	if err := r.Run(fmt.Sprintf(":%s", serverPort)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}