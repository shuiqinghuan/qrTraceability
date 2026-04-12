package api

import (
	"net/http"
	"strconv"

	"server/internal/db"

	"github.com/gin-gonic/gin"
)

type SeedInfoRequest struct {
	Name        string `json:"name" binding:"required"`
	Variety     string `json:"variety" binding:"required"`
	Description string `json:"description"`
}

func CreateSeedInfo(c *gin.Context) {
	var req SeedInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	var seedID int
	err := db.DB.QueryRow(
		"INSERT INTO seed_info (name, variety, description, created_by) VALUES ($1, $2, $3, $4) RETURNING id",
		req.Name, req.Variety, req.Description, userID,
	).Scan(&seedID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create seed info"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Seed info created successfully",
		"seed": gin.H{
			"id":          seedID,
			"name":        req.Name,
			"variety":     req.Variety,
			"description": req.Description,
			"created_by":  userID,
		},
	})
}

func GetSeedInfo(c *gin.Context) {
	seedID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid seed ID"})
		return
	}

	var seed struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Variety     string `json:"variety"`
		Description string `json:"description"`
		CreatedBy   int    `json:"created_by"`
	}

	err = db.DB.QueryRow(
		"SELECT id, name, variety, description, created_by FROM seed_info WHERE id = $1",
		seedID,
	).Scan(&seed.ID, &seed.Name, &seed.Variety, &seed.Description, &seed.CreatedBy)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seed info not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"seed": seed})
}

func ListSeedInfo(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, name, variety, description, created_by FROM seed_info ORDER BY created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list seed info"})
		return
	}
	defer rows.Close()

	var seeds []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Variety     string `json:"variety"`
		Description string `json:"description"`
		CreatedBy   int    `json:"created_by"`
	}

	for rows.Next() {
		var seed struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Variety     string `json:"variety"`
			Description string `json:"description"`
			CreatedBy   int    `json:"created_by"`
		}
		if err := rows.Scan(&seed.ID, &seed.Name, &seed.Variety, &seed.Description, &seed.CreatedBy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan seed info"})
			return
		}
		seeds = append(seeds, seed)
	}

	c.JSON(http.StatusOK, gin.H{"seeds": seeds})
}

func UpdateSeedInfo(c *gin.Context) {
	seedID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid seed ID"})
		return
	}

	var req SeedInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	// 检查种子信息是否存在且属于当前用户
	var creatorID int
	err = db.DB.QueryRow("SELECT created_by FROM seed_info WHERE id = $1", seedID).Scan(&creatorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seed info not found"})
		return
	}

	if creatorID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	_, err = db.DB.Exec(
		"UPDATE seed_info SET name = $1, variety = $2, description = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4",
		req.Name, req.Variety, req.Description, seedID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update seed info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Seed info updated successfully",
		"seed": gin.H{
			"id":          seedID,
			"name":        req.Name,
			"variety":     req.Variety,
			"description": req.Description,
		},
	})
}

func DeleteSeedInfo(c *gin.Context) {
	seedID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid seed ID"})
		return
	}

	userID, _ := c.Get("userID")

	// 检查种子信息是否存在且属于当前用户
	var creatorID int
	err = db.DB.QueryRow("SELECT created_by FROM seed_info WHERE id = $1", seedID).Scan(&creatorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seed info not found"})
		return
	}

	if creatorID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	_, err = db.DB.Exec("DELETE FROM seed_info WHERE id = $1", seedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete seed info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seed info deleted successfully"})
}