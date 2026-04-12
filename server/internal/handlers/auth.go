package handlers

import (
    "database/sql"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
    PhoneNumber string `json:"phone_number"`
    Password    string `json:"password"`
}

func Login(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req LoginRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
            return
        }

        log.Printf("Login attempt for phone: %s", req.PhoneNumber)

        var user struct {
            ID       int
            Phone    string
            Password string
            Role     string
        }

        query := `SELECT id, phone_number, password, role FROM users WHERE phone_number = $1`
        err := db.QueryRow(query, req.PhoneNumber).Scan(&user.ID, &user.Phone, &user.Password, &user.Role)
        if err != nil {
            log.Printf("User not found: %v", err)
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid phone number or password"})
            return
        }

        log.Printf("User found: ID=%d, Role=%s", user.ID, user.Role)
        log.Printf("Stored hash: %s", user.Password)
        log.Printf("Input password: %s", req.Password)

        err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
        if err != nil {
            log.Printf("Password mismatch: %v", err)
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid phone number or password"})
            return
        }

        log.Printf("Login successful for user: %s", user.Phone)

        c.JSON(http.StatusOK, gin.H{
            "message": "Login successful",
            "user": gin.H{
                "id":           user.ID,
                "phone_number": user.Phone,
                "role":         user.Role,
            },
        })
    }
}
