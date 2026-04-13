package api

import (
	"net/http"
	"strconv"

	"server/internal/db"

	"github.com/gin-gonic/gin"
)

type ProductQualityRequest struct {
	PlantingID      int     `json:"planting_id" binding:"required"`
	HarvestStartDate string `json:"harvest_start_date"`
	HarvestEndDate   string `json:"harvest_end_date"`
	SugarContent    float64 `json:"sugar_content"`
	Weight          float64 `json:"weight"`
	TasteDescription string `json:"taste_description"`
	SuitableFor     string  `json:"suitable_for"`
	QualitySummary  string  `json:"quality_summary"`
}

func CreateProductQuality(c *gin.Context) {
	var req ProductQualityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	var qualityID int
	err := db.DB.QueryRow(
		"INSERT INTO product_quality (planting_id, harvest_start_date, harvest_end_date, sugar_content, weight, taste_description, suitable_for, quality_summary, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id",
		req.PlantingID, req.HarvestStartDate, req.HarvestEndDate, req.SugarContent, req.Weight, req.TasteDescription, req.SuitableFor, req.QualitySummary, userID,
	).Scan(&qualityID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product quality record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product quality record created successfully",
		"quality": gin.H{
			"id":               qualityID,
			"planting_id":      req.PlantingID,
			"harvest_start_date": req.HarvestStartDate,
			"harvest_end_date":   req.HarvestEndDate,
			"sugar_content":    req.SugarContent,
			"weight":          req.Weight,
			"taste_description": req.TasteDescription,
			"suitable_for":    req.SuitableFor,
			"quality_summary": req.QualitySummary,
			"created_by":      userID,
		},
	})
}

func GetProductQuality(c *gin.Context) {
	qualityID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quality ID"})
		return
	}

	var quality struct {
		ID               int     `json:"id"`
		PlantingID       int     `json:"planting_id"`
		HarvestStartDate *string `json:"harvest_start_date"`
		HarvestEndDate   *string `json:"harvest_end_date"`
		SugarContent     float64 `json:"sugar_content"`
		Weight           float64 `json:"weight"`
		TasteDescription string  `json:"taste_description"`
		SuitableFor      string  `json:"suitable_for"`
		QualitySummary   string  `json:"quality_summary"`
		CreatedBy        int     `json:"created_by"`
	}

	err = db.DB.QueryRow(
		"SELECT id, planting_id, harvest_start_date, harvest_end_date, sugar_content, weight, taste_description, suitable_for, quality_summary, created_by FROM product_quality WHERE id = $1",
		qualityID,
	).Scan(&quality.ID, &quality.PlantingID, &quality.HarvestStartDate, &quality.HarvestEndDate, &quality.SugarContent, &quality.Weight, &quality.TasteDescription, &quality.SuitableFor, &quality.QualitySummary, &quality.CreatedBy)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product quality record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"quality": quality})
}

func GetProductQualityByPlanting(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("planting_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	var quality struct {
		ID               int     `json:"id"`
		PlantingID       int     `json:"planting_id"`
		HarvestStartDate *string `json:"harvest_start_date"`
		HarvestEndDate   *string `json:"harvest_end_date"`
		SugarContent     float64 `json:"sugar_content"`
		Weight           float64 `json:"weight"`
		TasteDescription string  `json:"taste_description"`
		SuitableFor      string  `json:"suitable_for"`
		QualitySummary   string  `json:"quality_summary"`
		CreatedBy        int     `json:"created_by"`
	}

	err = db.DB.QueryRow(
		"SELECT id, planting_id, harvest_start_date, harvest_end_date, sugar_content, weight, taste_description, suitable_for, quality_summary, created_by FROM product_quality WHERE planting_id = $1",
		plantingID,
	).Scan(&quality.ID, &quality.PlantingID, &quality.HarvestStartDate, &quality.HarvestEndDate, &quality.SugarContent, &quality.Weight, &quality.TasteDescription, &quality.SuitableFor, &quality.QualitySummary, &quality.CreatedBy)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product quality record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"quality": quality})
}

func UpdateProductQuality(c *gin.Context) {
	qualityID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quality ID"})
		return
	}

	var req ProductQualityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	// 检查品质记录是否存在且属于当前用户
	var creatorID int
	err = db.DB.QueryRow("SELECT created_by FROM product_quality WHERE id = $1", qualityID).Scan(&creatorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product quality record not found"})
		return
	}

	if creatorID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	_, err = db.DB.Exec(
		"UPDATE product_quality SET planting_id = $1, harvest_start_date = $2, harvest_end_date = $3, sugar_content = $4, weight = $5, taste_description = $6, suitable_for = $7, quality_summary = $8, updated_at = CURRENT_TIMESTAMP WHERE id = $9",
		req.PlantingID, req.HarvestStartDate, req.HarvestEndDate, req.SugarContent, req.Weight, req.TasteDescription, req.SuitableFor, req.QualitySummary, qualityID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product quality record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product quality record updated successfully",
		"quality": gin.H{
			"id":               qualityID,
			"planting_id":      req.PlantingID,
			"harvest_start_date": req.HarvestStartDate,
			"harvest_end_date":   req.HarvestEndDate,
			"sugar_content":    req.SugarContent,
			"weight":          req.Weight,
			"taste_description": req.TasteDescription,
			"suitable_for":    req.SuitableFor,
			"quality_summary": req.QualitySummary,
		},
	})
}

func DeleteProductQuality(c *gin.Context) {
	qualityID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quality ID"})
		return
	}

	userID, _ := c.Get("userID")

	// 检查品质记录是否存在且属于当前用户
	var creatorID int
	err = db.DB.QueryRow("SELECT created_by FROM product_quality WHERE id = $1", qualityID).Scan(&creatorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product quality record not found"})
		return
	}

	if creatorID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	_, err = db.DB.Exec("DELETE FROM product_quality WHERE id = $1", qualityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product quality record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product quality record deleted successfully"})
}