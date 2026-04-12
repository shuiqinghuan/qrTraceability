package main

import (
    "database/sql"
    "log"
    "os"

    _ "github.com/lib/pq"
    "server/internal/router"
)

func main() {
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

    connStr := "host=" + dbHost + " port=" + dbPort + " user=" + dbUser +
        " password=" + dbPassword + " dbname=" + dbName + " sslmode=disable"

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatal("Failed to ping database:", err)
    }
    log.Println("Database connection established")

    r := router.SetupRouter(db)

    port := os.Getenv("SERVER_PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Server starting on port %s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
