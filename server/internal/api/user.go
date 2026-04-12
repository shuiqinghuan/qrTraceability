package api

import (
	"net/http"

	"server/internal/auth"
	"server/internal/db"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
	Role        string `json:"role" binding:"required,oneof=serverseed servergrow servermanager clentcustomer"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查手机号是否已存在
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE phone_number = $1", req.PhoneNumber).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Phone number already registered"})
		return
	}

	// 哈希密码
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}

	// 创建用户
	var userID int
	err = db.DB.QueryRow(
		"INSERT INTO users (phone_number, password_hash, role) VALUES ($1, $2, $3) RETURNING id",
		req.PhoneNumber, hashedPassword, req.Role,
	).Scan(&userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed"})
		return
	}

	// 生成token
	token, err := auth.GenerateToken(userID, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   token,
		"user": gin.H{
			"id":          userID,
			"phone_number": req.PhoneNumber,
			"role":        req.Role,
		},
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找用户
	var userID int
	var hashedPassword, role string
	err := db.DB.QueryRow(
		"SELECT id, password_hash, role FROM users WHERE phone_number = $1",
		req.PhoneNumber,
	).Scan(&userID, &hashedPassword, &role)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid phone number or password"})
		return
	}

	// 验证密码
	if !auth.CheckPasswordHash(req.Password, hashedPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid phone number or password"})
		return
	}

	// 生成token
	token, err := auth.GenerateToken(userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":          userID,
			"phone_number": req.PhoneNumber,
			"role":        role,
		},
	})
}