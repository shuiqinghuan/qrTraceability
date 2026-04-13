package api

import (
	"net/http"
	"strconv"
	"time"

	"server/internal/db"

	"github.com/gin-gonic/gin"
)

// LikeRequest 点赞请求结构
type LikeRequest struct {
	PlantingID int `json:"planting_id" binding:"required"`
}

// LikeResponse 点赞响应结构
type LikeResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	LikeCount  int    `json:"like_count"`
	CanLike    bool   `json:"can_like"`
	Cooldown   int    `json:"cooldown_seconds,omitempty"`
}

// getClientIP 获取客户端IP地址
func getClientIP(c *gin.Context) string {
	// 尝试从X-Forwarded-For获取（如果有代理）
	forwarded := c.GetHeader("X-Forwarded-For")
	if forwarded != "" {
		// 取第一个IP（客户端真实IP）
		ips := splitIPs(forwarded)
		if len(ips) > 0 {
			return ips[0]
		}
	}
	
	// 尝试从X-Real-IP获取
	realIP := c.GetHeader("X-Real-IP")
	if realIP != "" {
		return realIP
	}
	
	// 使用远程地址
	return c.ClientIP()
}

// splitIPs 分割IP地址字符串
func splitIPs(ipString string) []string {
	var ips []string
	start := 0
	for i, ch := range ipString {
		if ch == ',' {
			ip := ipString[start:i]
			ips = append(ips, ip)
			start = i + 1
		}
	}
	if start < len(ipString) {
		ips = append(ips, ipString[start:])
	}
	return ips
}

// LikeProduct 点赞产品（带IP限制）
func LikeProduct(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	// 获取客户端IP
	clientIP := getClientIP(c)
	if clientIP == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to determine client IP"})
		return
	}

	// 检查种植记录是否存在
	var plantingExists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM planting WHERE id = $1)", plantingID).Scan(&plantingExists)
	if err != nil || !plantingExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Planting record not found"})
		return
	}

	// 检查IP限制
	canLike, cooldownSeconds, err := checkIPLikeRestriction(clientIP, plantingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check IP restriction"})
		return
	}

	if !canLike {
		response := LikeResponse{
			Success:   false,
			Message:   "点赞过于频繁，请稍后再试",
			LikeCount: 0,
			CanLike:   false,
			Cooldown:  cooldownSeconds,
		}
		c.JSON(http.StatusTooManyRequests, gin.H{"response": response})
		return
	}

	// 获取用户ID（如果已登录）
	var userID interface{}
	userID, exists := c.Get("userID")
	if !exists {
		userID = nil
	}

	// 记录点赞
	var likeID int
	if userID != nil {
		// 已登录用户：使用用户ID记录点赞
		err = db.DB.QueryRow(`
			INSERT INTO likes (user_id, planting_id) 
			VALUES ($1, $2) 
			ON CONFLICT (user_id, planting_id) DO NOTHING 
			RETURNING id
		`, userID, plantingID).Scan(&likeID)
	} else {
		// 未登录用户：只记录IP限制
		likeID = 0
	}

	// 更新IP限制记录
	err = updateIPLikeRestriction(clientIP, plantingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update IP restriction"})
		return
	}

	// 获取当前点赞总数
	var likeCount int
	db.DB.QueryRow("SELECT COUNT(*) FROM likes WHERE planting_id = $1", plantingID).Scan(&likeCount)

	response := LikeResponse{
		Success:   true,
		Message:   "点赞成功",
		LikeCount: likeCount,
		CanLike:   false, // 点赞后暂时不能再次点赞
		Cooldown:  300,   // 5分钟冷却时间
	}

	c.JSON(http.StatusOK, gin.H{
		"response": response,
	})
}

// checkIPLikeRestriction 检查IP点赞限制
func checkIPLikeRestriction(ipAddress string, plantingID int) (bool, int, error) {
	var likeCount int
	var lastLikeTime *time.Time
	
	err := db.DB.QueryRow(`
		SELECT like_count, last_like_time 
		FROM ip_like_restrictions 
		WHERE ip_address = $1 AND planting_id = $2
	`, ipAddress, plantingID).Scan(&likeCount, &lastLikeTime)
	
	if err != nil {
		// 没有记录，可以点赞
		return true, 0, nil
	}

	// 检查冷却时间（5分钟）
	if lastLikeTime != nil {
		cooldownDuration := 5 * time.Minute
		nextAllowedTime := lastLikeTime.Add(cooldownDuration)
		
		if time.Now().Before(nextAllowedTime) {
			cooldownSeconds := int(nextAllowedTime.Sub(time.Now()).Seconds())
			return false, cooldownSeconds, nil
		}
	}

	// 检查每日点赞限制（10次）
	if likeCount >= 10 {
		// 检查是否是新的一天
		if lastLikeTime != nil {
			now := time.Now()
			lastDay := lastLikeTime.Truncate(24 * time.Hour)
			currentDay := now.Truncate(24 * time.Hour)
			
			if lastDay.Equal(currentDay) {
				// 同一天，已达到限制
				// 计算到第二天0点的秒数
				nextDay := currentDay.Add(24 * time.Hour)
				cooldownSeconds := int(nextDay.Sub(now).Seconds())
				return false, cooldownSeconds, nil
			}
		}
	}

	return true, 0, nil
}

// updateIPLikeRestriction 更新IP点赞限制记录
func updateIPLikeRestriction(ipAddress string, plantingID int) error {
	now := time.Now()
	
	// 检查是否是新的一天
	var lastLikeTime *time.Time
	var likeCount int
	
	err := db.DB.QueryRow(`
		SELECT last_like_time, like_count 
		FROM ip_like_restrictions 
		WHERE ip_address = $1 AND planting_id = $2
	`, ipAddress, plantingID).Scan(&lastLikeTime, &likeCount)
	
	if err != nil {
		// 没有记录，创建新记录
		_, err = db.DB.Exec(`
			INSERT INTO ip_like_restrictions (ip_address, planting_id, like_count, last_like_time)
			VALUES ($1, $2, 1, $3)
		`, ipAddress, plantingID, now)
		return err
	}

	// 有记录，检查是否是新的一天
	if lastLikeTime != nil {
		lastDay := lastLikeTime.Truncate(24 * time.Hour)
		currentDay := now.Truncate(24 * time.Hour)
		
		if !lastDay.Equal(currentDay) {
			// 新的一天，重置计数
			likeCount = 0
		}
	}

	// 更新记录
	_, err = db.DB.Exec(`
		UPDATE ip_like_restrictions 
		SET like_count = $1, last_like_time = $2
		WHERE ip_address = $3 AND planting_id = $4
	`, likeCount+1, now, ipAddress, plantingID)
	
	return err
}

// GetLikeStatus 获取点赞状态
func GetLikeStatus(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	// 获取客户端IP
	clientIP := getClientIP(c)
	if clientIP == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to determine client IP"})
		return
	}

	// 检查种植记录是否存在
	var plantingExists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM planting WHERE id = $1)", plantingID).Scan(&plantingExists)
	if err != nil || !plantingExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Planting record not found"})
		return
	}

	// 检查IP限制
	canLike, cooldownSeconds, err := checkIPLikeRestriction(clientIP, plantingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check IP restriction"})
		return
	}

	// 获取当前点赞总数
	var likeCount int
	db.DB.QueryRow("SELECT COUNT(*) FROM likes WHERE planting_id = $1", plantingID).Scan(&likeCount)

	// 检查用户是否已点赞（如果已登录）
	var userLiked bool
	userID, exists := c.Get("userID")
	if exists && userID != nil {
		err = db.DB.QueryRow(`
			SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND planting_id = $2)
		`, userID, plantingID).Scan(&userLiked)
		if err != nil {
			userLiked = false
		}
	}

	response := LikeResponse{
		Success:   true,
		Message:   "获取点赞状态成功",
		LikeCount: likeCount,
		CanLike:   canLike,
		Cooldown:  cooldownSeconds,
	}

	c.JSON(http.StatusOK, gin.H{
		"response":   response,
		"user_liked": userLiked,
		"ip_address": clientIP,
	})
}

// GetLikeCount 获取点赞统计
func GetLikeCount(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	// 检查种植记录是否存在
	var plantingExists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM planting WHERE id = $1)", plantingID).Scan(&plantingExists)
	if err != nil || !plantingExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Planting record not found"})
		return
	}

	// 获取点赞总数
	var likeCount int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM likes WHERE planting_id = $1", plantingID).Scan(&likeCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get like count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"planting_id": plantingID,
		"like_count": likeCount,
	})
}

// RemoveLike 取消点赞
func RemoveLike(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	// 获取用户ID（必须已登录）
	userID, exists := c.Get("userID")
	if !exists || userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User must be logged in to remove likes"})
		return
	}

	// 删除点赞记录
	_, err = db.DB.Exec(`
		DELETE FROM likes 
		WHERE user_id = $1 AND planting_id = $2
	`, userID, plantingID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove like"})
		return
	}

	// 获取更新后的点赞总数
	var likeCount int
	db.DB.QueryRow("SELECT COUNT(*) FROM likes WHERE planting_id = $1", plantingID).Scan(&likeCount)

	response := LikeResponse{
		Success:   true,
		Message:   "取消点赞成功",
		LikeCount: likeCount,
		CanLike:   true,
	}

	c.JSON(http.StatusOK, gin.H{
		"response": response,
	})
}

// ResetIPLikeRestriction 重置IP点赞限制（管理员功能）
func ResetIPLikeRestriction(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	ipAddress := c.Query("ip")
	if ipAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IP address is required"})
		return
	}

	// 删除IP限制记录
	_, err = db.DB.Exec(`
		DELETE FROM ip_like_restrictions 
		WHERE ip_address = $1 AND planting_id = $2
	`, ipAddress, plantingID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset IP restriction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "IP点赞限制已重置",
		"ip_address": ipAddress,
		"planting_id": plantingID,
	})
}