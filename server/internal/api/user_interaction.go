package api

import (
	"net/http"
	"strconv"
	"time"

	"server/internal/db"
	"server/internal/utils"

	"github.com/gin-gonic/gin"
)

// 初始化速率限制器，每IP每分钟最多5次点赞请求
var likeLimiter = utils.NewRateLimiter(5, time.Minute)

func AddFavorite(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("planting_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	userID, _ := c.Get("userID")

	// 检查是否已经收藏
	var count int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM favorites WHERE user_id = $1 AND planting_id = $2", userID, plantingID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Already favorited"})
		return
	}

	_, err = db.DB.Exec("INSERT INTO favorites (user_id, planting_id) VALUES ($1, $2)", userID, plantingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Added to favorites"})
}

func RemoveFavorite(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("planting_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	userID, _ := c.Get("userID")

	_, err = db.DB.Exec("DELETE FROM favorites WHERE user_id = $1 AND planting_id = $2", userID, plantingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed from favorites"})
}

func ListFavorites(c *gin.Context) {
	userID, _ := c.Get("userID")

	rows, err := db.DB.Query(`
		SELECT p.id, p.seed_id, p.planting_date, p.transplanting_date, p.location, p.notes 
		FROM planting p
		JOIN favorites f ON p.id = f.planting_id
		WHERE f.user_id = $1
		ORDER BY f.created_at DESC
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list favorites"})
		return
	}
	defer rows.Close()

	var favorites []struct {
		ID                int    `json:"id"`
		SeedID            int    `json:"seed_id"`
		PlantingDate      string `json:"planting_date"`
		TransplantingDate string `json:"transplanting_date"`
		Location          string `json:"location"`
		Notes             string `json:"notes"`
	}

	for rows.Next() {
		var fav struct {
			ID                int    `json:"id"`
			SeedID            int    `json:"seed_id"`
			PlantingDate      string `json:"planting_date"`
			TransplantingDate string `json:"transplanting_date"`
			Location          string `json:"location"`
			Notes             string `json:"notes"`
		}
		if err := rows.Scan(&fav.ID, &fav.SeedID, &fav.PlantingDate, &fav.TransplantingDate, &fav.Location, &fav.Notes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan favorite"})
			return
		}
		favorites = append(favorites, fav)
	}

	c.JSON(http.StatusOK, gin.H{"favorites": favorites})
}

func AddLike(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("planting_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	// 获取用户IP地址
	clientIP := c.ClientIP()

	// 检查速率限制
	if !likeLimiter.Allow(clientIP) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests, please try again later"})
		return
	}

	userID, _ := c.Get("userID")

	// 检查是否已经点赞
	var count int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = $1 AND planting_id = $2", userID, plantingID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Already liked"})
		return
	}

	_, err = db.DB.Exec("INSERT INTO likes (user_id, planting_id) VALUES ($1, $2)", userID, plantingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add like"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Added like"})
}

func RemoveLike(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("planting_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	userID, _ := c.Get("userID")

	_, err = db.DB.Exec("DELETE FROM likes WHERE user_id = $1 AND planting_id = $2", userID, plantingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove like"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed like"})
}

func CheckLikeStatus(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("planting_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	userID, _ := c.Get("userID")

	var count int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = $1 AND planting_id = $2", userID, plantingID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"liked": count > 0})
}