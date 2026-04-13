package router

import (
    "database/sql"

    "github.com/gin-gonic/gin"
    "server/internal/api"
    "server/internal/handlers"
    "server/internal/middleware"
)

func SetupRouter(db *sql.DB) *gin.Engine {
    r := gin.Default()

    // 根路径
    r.GET("/", func(c *gin.Context) {
        c.String(200, "农产品二维码追溯系统 API")
    })

    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "healthy",
            "service": "qr-traceability",
        })
    })

    // API 路由组
    apiGroup := r.Group("/api")
    {
        // 认证路由
        auth := apiGroup.Group("/auth")
        {
            auth.POST("/login", handlers.Login(db))
        }

        // 公共产品API（无需认证）
        products := apiGroup.Group("/products")
        {
            products.GET("/", api.ListProducts)
            products.GET("/:id", api.GetProductDetail)
            products.GET("/:id/media", api.GetProductMedia)
            products.GET("/:id/quality", api.GetProductQuality)
            products.GET("/:id/stats", api.GetProductStats)
            products.GET("/:id/related", api.GetRelatedProducts)
            products.GET("/:id/similar", api.GetSimilarProducts)
            products.GET("/:id/likes", api.GetLikeCount)
            products.GET("/:id/like-status", api.GetLikeStatus)
            products.POST("/:id/like", api.LikeProduct)
        }

        // 地点相关API
        locations := apiGroup.Group("/locations")
        {
            locations.GET("/:location/products", api.GetProductsByLocation)
            locations.GET("/:location/stats", api.GetLocationStats)
            locations.GET("/search", api.SearchLocations)
        }

        // 标签API
        tags := apiGroup.Group("/tags")
        {
            tags.GET("/search", api.SearchTags)
        }

        // 需要用户认证的API
        authenticated := apiGroup.Group("")
        authenticated.Use(middleware.AuthMiddleware())
        {
            // 用户相关API
            user := authenticated.Group("/user")
            {
                user.POST("/likes/:id/remove", api.RemoveLike)
            }

            // 种植标签管理（需要权限）
            plantingTags := authenticated.Group("/planting/:id/tags")
            {
                plantingTags.GET("/", api.GetPlantingTags)
                plantingTags.POST("/", api.AddPlantingTag)
                plantingTags.POST("/batch", api.BatchAddPlantingTags)
                plantingTags.PUT("/:tagId", api.UpdatePlantingTag)
                plantingTags.DELETE("/:tagId", api.DeletePlantingTag)
            }
        }

        // 管理员API
        admin := apiGroup.Group("/admin")
        {
            admin.POST("/login", api.AdminLogin)
            
            // 需要管理员认证的路由
            adminAuth := admin.Group("")
            adminAuth.Use(middleware.AdminMiddleware())
            {
                // 产品管理
                adminAuth.GET("/products", api.AdminListProducts)
                adminAuth.POST("/products", api.AdminCreateProduct)
                adminAuth.GET("/products/:id", api.AdminGetProductDetail)
                adminAuth.PUT("/products/:id", api.AdminUpdateProduct)
                adminAuth.DELETE("/products/:id", api.AdminDeleteProduct)
                
                // 媒体管理
                adminAuth.POST("/media/upload", api.AdminUploadMedia)
                
                // 品质管理
                adminAuth.POST("/quality", api.AdminUpdateQuality)
                
                // 系统管理
                adminAuth.POST("/likes/:id/reset-ip", api.ResetIPLikeRestriction)
            }
        }

        // 现有API路由（保持向后兼容）
        legacy := apiGroup.Group("")
        legacy.Use(middleware.AuthMiddleware())
        {
            // 种子信息管理
            legacy.POST("/seed", api.CreateSeedInfo)
            legacy.GET("/seed", api.ListSeedInfo)
            legacy.GET("/seed/:id", api.GetSeedInfo)
            legacy.PUT("/seed/:id", api.UpdateSeedInfo)
            legacy.DELETE("/seed/:id", api.DeleteSeedInfo)
            
            // 种植管理
            // 这里可以添加现有的种植管理API
            
            // 品质管理（更新后的）
            legacy.POST("/quality", api.CreateProductQuality)
            legacy.GET("/quality/:id", api.GetProductQuality)
            legacy.GET("/quality/planting/:planting_id", api.GetProductQualityByPlanting)
            legacy.PUT("/quality/:id", api.UpdateProductQuality)
            legacy.DELETE("/quality/:id", api.DeleteProductQuality)
        }
    }

    return r
}
