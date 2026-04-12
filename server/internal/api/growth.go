package api

import (
	"net/http"
	"strconv"

	"server/internal/db"

	"github.com/gin-gonic/gin"
)

type GrowthMediaRequest struct {
	PlantingID  int    `json:"planting_id" binding:"required"`
	MediaType   string `json:"media_type" binding:"required,oneof=image video"`
	FilePath    string `json:"file_path" binding:"required"`
	Description string `json:"description"`
}

func CreateGrowthMedia(c *gin.Context) {
	var req GrowthMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	var mediaID int
	err := db.DB.QueryRow(
		"INSERT INTO growth_media (planting_id, media_type, file_path, description, created_by) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		req.PlantingID, req.MediaType, req.FilePath, req.Description, userID,
	).Scan(&mediaID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create growth media"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Growth media created successfully",
		"media": gin.H{
			"id":          mediaID,
			"planting_id": req.PlantingID,
			"media_type":  req.MediaType,
			"file_path":   req.FilePath,
			"description": req.Description,
			"created_by":  userID,
		},
	})
}

func GetGrowthMedia(c *gin.Context) {
	mediaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	var media struct {
		ID          int    `json:"id"`
		PlantingID  int    `json:"planting_id"`
		MediaType   string `json:"media_type"`
		FilePath    string `json:"file_path"`
		Description string `json:"description"`
		CreatedBy   int    `json:"created_by"`
	}

	err = db.DB.QueryRow(
		"SELECT id, planting_id, media_type, file_path, description, created_by FROM growth_media WHERE id = $1",
		mediaID,
	).Scan(&media.ID, &media.PlantingID, &media.MediaType, &media.FilePath, &media.Description, &media.CreatedBy)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Growth media not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"media": media})
}

func ListGrowthMediaByPlanting(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("planting_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	rows, err := db.DB.Query(
		"SELECT id, planting_id, media_type, file_path, description, created_by FROM growth_media WHERE planting_id = $1 ORDER BY created_at DESC",
		plantingID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list growth media"})
		return
	}
	defer rows.Close()

	var mediaList []struct {
		ID          int    `json:"id"`
		PlantingID  int    `json:"planting_id"`
		MediaType   string `json:"media_type"`
		FilePath    string `json:"file_path"`
		Description string `json:"description"`
		CreatedBy   int    `json:"created_by"`
	}

	for rows.Next() {
		var media struct {
			ID          int    `json:"id"`
			PlantingID  int    `json:"planting_id"`
			MediaType   string `json:"media_type"`
			FilePath    string `json:"file_path"`
			Description string `json:"description"`
			CreatedBy   int    `json:"created_by"`
		}
		if err := rows.Scan(&media.ID, &media.PlantingID, &media.MediaType, &media.FilePath, &media.Description, &media.CreatedBy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan growth media"})
			return
		}
		mediaList = append(mediaList, media)
	}

	c.JSON(http.StatusOK, gin.H{"media": mediaList})
}

func DeleteGrowthMedia(c *gin.Context) {
	mediaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	userID, _ := c.Get("userID")

	// 检查媒体是否存在且属于当前用户
	var creatorID int
	err = db.DB.QueryRow("SELECT created_by FROM growth_media WHERE id = $1", mediaID).Scan(&creatorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Growth media not found"})
		return
	}

	if creatorID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	_, err = db.DB.Exec("DELETE FROM growth_media WHERE id = $1", mediaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete growth media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Growth media deleted successfully"})
}