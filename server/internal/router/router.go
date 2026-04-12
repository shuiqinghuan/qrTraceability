package router

import (
    "database/sql"

    "github.com/gin-gonic/gin"
    "server/internal/handlers"
)

func SetupRouter(db *sql.DB) *gin.Engine {
    r := gin.Default()

    // 根路径
    r.GET("/", func(c *gin.Context) {
        c.String(200, "Hello from Go server!")
    })

    r.GET("/health", func(c *gin.Context) {
        c.String(200, "Hello from Go server!")
    })

    // API 路由组
    api := r.Group("/api")
    {
        auth := api.Group("/auth")
        {
            auth.POST("/login", handlers.Login(db))
        }
    }

    return r
}
