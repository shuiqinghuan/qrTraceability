package api

import (
	"fmt"
	"net/http"
	"strconv"

	"server/internal/db"
	"server/internal/utils"

	"github.com/gin-gonic/gin"
)

type PlantingRequest struct {
	SeedID          int    `json:"seed_id" binding:"required"`
	PlantingDate    string `json:"planting_date" binding:"required"`
	TransplantingDate string `json:"transplanting_date"`
	Location        string `json:"location"`
	Notes           string `json:"notes"`
}

func CreatePlanting(c *gin.Context) {
	var req PlantingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	var plantingID int
	err := db.DB.QueryRow(
		"INSERT INTO planting (seed_id, planting_date, transplanting_date, location, notes, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		req.SeedID, req.PlantingDate, req.TransplantingDate, req.Location, req.Notes, userID,
	).Scan(&plantingID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create planting record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Planting record created successfully",
		"planting": gin.H{
			"id":                plantingID,
			"seed_id":          req.SeedID,
			"planting_date":    req.PlantingDate,
			"transplanting_date": req.TransplantingDate,
			"location":        req.Location,
			"notes":           req.Notes,
			"created_by":      userID,
		},
	})
}

func GetPlanting(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	var planting struct {
		ID                int    `json:"id"`
		SeedID            int    `json:"seed_id"`
		PlantingDate      string `json:"planting_date"`
		TransplantingDate string `json:"transplanting_date"`
		Location          string `json:"location"`
		Notes             string `json:"notes"`
		CreatedBy         int    `json:"created_by"`
	}

	err = db.DB.QueryRow(
		"SELECT id, seed_id, planting_date, transplanting_date, location, notes, created_by FROM planting WHERE id = $1",
		plantingID,
	).Scan(&planting.ID, &planting.SeedID, &planting.PlantingDate, &planting.TransplantingDate, &planting.Location, &planting.Notes, &planting.CreatedBy)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Planting record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"planting": planting})
}

func ListPlanting(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, seed_id, planting_date, transplanting_date, location, notes, created_by FROM planting ORDER BY created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list planting records"})
		return
	}
	defer rows.Close()

	var plantings []struct {
		ID                int    `json:"id"`
		SeedID            int    `json:"seed_id"`
		PlantingDate      string `json:"planting_date"`
		TransplantingDate string `json:"transplanting_date"`
		Location          string `json:"location"`
		Notes             string `json:"notes"`
		CreatedBy         int    `json:"created_by"`
	}

	for rows.Next() {
		var planting struct {
			ID                int    `json:"id"`
			SeedID            int    `json:"seed_id"`
			PlantingDate      string `json:"planting_date"`
			TransplantingDate string `json:"transplanting_date"`
			Location          string `json:"location"`
			Notes             string `json:"notes"`
			CreatedBy         int    `json:"created_by"`
		}
		if err := rows.Scan(&planting.ID, &planting.SeedID, &planting.PlantingDate, &planting.TransplantingDate, &planting.Location, &planting.Notes, &planting.CreatedBy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan planting record"})
			return
		}
		plantings = append(plantings, planting)
	}

	c.JSON(http.StatusOK, gin.H{"plantings": plantings})
}

func UpdatePlanting(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	var req PlantingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	// 检查种植记录是否存在且属于当前用户
	var creatorID int
	err = db.DB.QueryRow("SELECT created_by FROM planting WHERE id = $1", plantingID).Scan(&creatorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Planting record not found"})
		return
	}

	if creatorID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	_, err = db.DB.Exec(
		"UPDATE planting SET seed_id = $1, planting_date = $2, transplanting_date = $3, location = $4, notes = $5, updated_at = CURRENT_TIMESTAMP WHERE id = $6",
		req.SeedID, req.PlantingDate, req.TransplantingDate, req.Location, req.Notes, plantingID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update planting record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Planting record updated successfully",
		"planting": gin.H{
			"id":                plantingID,
			"seed_id":          req.SeedID,
			"planting_date":    req.PlantingDate,
			"transplanting_date": req.TransplantingDate,
			"location":        req.Location,
			"notes":           req.Notes,
		},
	})
}

func DeletePlanting(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	userID, _ := c.Get("userID")

	// 检查种植记录是否存在且属于当前用户
	var creatorID int
	err = db.DB.QueryRow("SELECT created_by FROM planting WHERE id = $1", plantingID).Scan(&creatorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Planting record not found"})
		return
	}

	if creatorID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	_, err = db.DB.Exec("DELETE FROM planting WHERE id = $1", plantingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete planting record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Planting record deleted successfully"})
}

// GenerateQRCode 为产品生成二维码
func GenerateQRCode(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	// 检查种植记录是否存在
	var count int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM planting WHERE id = $1", plantingID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Planting record not found"})
		return
	}

	// 生成二维码
	qrcodePath, err := utils.GenerateQRCode(plantingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "QR code generated successfully",
		"qrcode_path": qrcodePath,
		"product_url": fmt.Sprintf("http://localhost:3000/product/%d", plantingID),
	})
}