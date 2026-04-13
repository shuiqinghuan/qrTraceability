package main

import (
    "log"
    "os"

    "server/internal/auth"
    "server/internal/db"
    "server/internal/router"
)

func main() {
    // 从环境变量获取数据库配置
    dbHost := os.Getenv("DB_HOST")
    if dbHost == "" {
        dbHost = "127.0.0.1"
    }
    dbPort := os.Getenv("DB_PORT")
    if dbPort == "" {
        dbPort = "5432"
    }
    dbUser := os.Getenv("DB_USER")
    if dbUser == "" {
        dbUser = "qruser"
    }
    dbPassword := os.Getenv("DB_PASSWORD")
    if dbPassword == "" {
        dbPassword = "6180680"
    }
    dbName := os.Getenv("DB_NAME")
    if dbName == "" {
        dbName = "qrtraceability"
    }

    // 初始化数据库连接
    err := db.InitDB(dbHost, dbPort, dbUser, dbPassword, dbName)
    if err != nil {
        log.Fatal("Failed to initialize database:", err)
    }
    defer db.CloseDB()

    // 初始化JWT密钥
    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        jwtSecret = "your_jwt_secret_key_for_production_change_this"
        log.Println("Warning: Using default JWT secret. For production, set JWT_SECRET environment variable.")
    }
    auth.InitJWTSecret(jwtSecret)

    // 设置路由
    r := router.SetupRouter(db.DB)

    // 启动服务器
    port := os.Getenv("SERVER_PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Server starting on port %s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
