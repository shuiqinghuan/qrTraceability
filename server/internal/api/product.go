package api

import (
	"net/http"
	"strconv"
	"time"

	"server/internal/db"

	"github.com/gin-gonic/gin"
)

// ProductDetailResponse 产品详情响应结构
type ProductDetailResponse struct {
	ID          int                    `json:"id"`
	SeedInfo    SeedInfoDetail         `json:"seed_info"`
	Planting    PlantingDetail         `json:"planting"`
	Media       []GrowthMediaDetail    `json:"media"`
	Quality     ProductQualityDetail   `json:"quality"`
	Tags        []PlantingTagDetail    `json:"tags"`
	LikeCount   int                    `json:"like_count"`
	FavoriteCount int                  `json:"favorite_count"`
}

// SeedInfoDetail 种子信息详情
type SeedInfoDetail struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Variety     string `json:"variety"`
	VarietyCode string `json:"variety_code"`
	Description string `json:"description"`
}

// PlantingDetail 种植信息详情
type PlantingDetail struct {
	ID               int    `json:"id"`
	PlantingDate     string `json:"planting_date"`
	TransplantingDate string `json:"transplanting_date"`
	Location         string `json:"location"`
	Notes            string `json:"notes"`
}

// GrowthMediaDetail 生长媒体详情
type GrowthMediaDetail struct {
	ID          int    `json:"id"`
	MediaType   string `json:"media_type"`
	FilePath    string `json:"file_path"`
	Description string `json:"description"`
}

// ProductQualityDetail 产品品质详情
type ProductQualityDetail struct {
	ID               int     `json:"id"`
	HarvestStartDate *string `json:"harvest_start_date"`
	HarvestEndDate   *string `json:"harvest_end_date"`
	SugarContent     float64 `json:"sugar_content"`
	Weight           float64 `json:"weight"`
	TasteDescription string  `json:"taste_description"`
	SuitableFor      string  `json:"suitable_for"`
	QualitySummary   string  `json:"quality_summary"`
}

// PlantingTagDetail 种植标签详情
type PlantingTagDetail struct {
	ID      int    `json:"id"`
	TagName string `json:"tag_name"`
}

// GetProductDetail 获取产品详情
func GetProductDetail(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	// 获取种植信息
	var planting PlantingDetail
	err = db.DB.QueryRow(`
		SELECT id, planting_date, transplanting_date, location, notes 
		FROM planting WHERE id = $1
	`, plantingID).Scan(
		&planting.ID, &planting.PlantingDate, &planting.TransplantingDate, 
		&planting.Location, &planting.Notes,
	)
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Planting record not found"})
		return
	}

	// 获取种子信息
	var seedInfo SeedInfoDetail
	err = db.DB.QueryRow(`
		SELECT s.id, s.name, s.variety, s.variety_code, s.description
		FROM seed_info s
		JOIN planting p ON p.seed_id = s.id
		WHERE p.id = $1
	`, plantingID).Scan(
		&seedInfo.ID, &seedInfo.Name, &seedInfo.Variety, 
		&seedInfo.VarietyCode, &seedInfo.Description,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get seed info"})
		return
	}

	// 获取媒体信息
	var media []GrowthMediaDetail
	rows, err := db.DB.Query(`
		SELECT id, media_type, file_path, description
		FROM growth_media 
		WHERE planting_id = $1
		ORDER BY created_at
	`, plantingID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get media info"})
		return
	}
	defer rows.Close()
	
	for rows.Next() {
		var m GrowthMediaDetail
		if err := rows.Scan(&m.ID, &m.MediaType, &m.FilePath, &m.Description); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan media info"})
			return
		}
		media = append(media, m)
	}

	// 获取品质信息
	var quality ProductQualityDetail
	err = db.DB.QueryRow(`
		SELECT id, harvest_start_date, harvest_end_date, sugar_content, 
		       weight, taste_description, suitable_for, quality_summary
		FROM product_quality 
		WHERE planting_id = $1
	`, plantingID).Scan(
		&quality.ID, &quality.HarvestStartDate, &quality.HarvestEndDate,
		&quality.SugarContent, &quality.Weight, &quality.TasteDescription,
		&quality.SuitableFor, &quality.QualitySummary,
	)
	
	// 如果找不到品质信息，创建空结构
	if err != nil {
		quality = ProductQualityDetail{}
	}

	// 获取标签信息
	var tags []PlantingTagDetail
	tagRows, err := db.DB.Query(`
		SELECT id, tag_name
		FROM planting_tags 
		WHERE planting_id = $1
		ORDER BY tag_name
	`, plantingID)
	
	if err == nil {
		defer tagRows.Close()
		for tagRows.Next() {
			var tag PlantingTagDetail
			if err := tagRows.Scan(&tag.ID, &tag.TagName); err != nil {
				break
			}
			tags = append(tags, tag)
		}
	}

	// 获取点赞和收藏统计
	var likeCount, favoriteCount int
	db.DB.QueryRow("SELECT COUNT(*) FROM likes WHERE planting_id = $1", plantingID).Scan(&likeCount)
	db.DB.QueryRow("SELECT COUNT(*) FROM favorites WHERE planting_id = $1", plantingID).Scan(&favoriteCount)

	response := ProductDetailResponse{
		ID:          plantingID,
		SeedInfo:    seedInfo,
		Planting:    planting,
		Media:       media,
		Quality:     quality,
		Tags:        tags,
		LikeCount:   likeCount,
		FavoriteCount: favoriteCount,
	}

	c.JSON(http.StatusOK, gin.H{
		"product": response,
	})
}

// GetProductMedia 获取产品媒体信息
func GetProductMedia(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	var media []GrowthMediaDetail
	rows, err := db.DB.Query(`
		SELECT id, media_type, file_path, description
		FROM growth_media 
		WHERE planting_id = $1
		ORDER BY created_at
	`, plantingID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get media info"})
		return
	}
	defer rows.Close()
	
	for rows.Next() {
		var m GrowthMediaDetail
		if err := rows.Scan(&m.ID, &m.MediaType, &m.FilePath, &m.Description); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan media info"})
			return
		}
		media = append(media, m)
	}

	c.JSON(http.StatusOK, gin.H{
		"media": media,
	})
}

// GetProductQuality 获取产品品质信息
func GetProductQuality(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	var quality ProductQualityDetail
	err = db.DB.QueryRow(`
		SELECT id, harvest_start_date, harvest_end_date, sugar_content, 
		       weight, taste_description, suitable_for, quality_summary
		FROM product_quality 
		WHERE planting_id = $1
	`, plantingID).Scan(
		&quality.ID, &quality.HarvestStartDate, &quality.HarvestEndDate,
		&quality.SugarContent, &quality.Weight, &quality.TasteDescription,
		&quality.SuitableFor, &quality.QualitySummary,
	)
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product quality record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"quality": quality,
	})
}

// ListProducts 获取产品列表
func ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// 获取产品总数
	var totalCount int
	err := db.DB.QueryRow(`
		SELECT COUNT(*) FROM planting
	`).Scan(&totalCount)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count products"})
		return
	}

	// 获取产品列表
	rows, err := db.DB.Query(`
		SELECT p.id, p.planting_date, p.transplanting_date, p.location,
		       s.name, s.variety, s.variety_code,
		       COALESCE(pq.sugar_content, 0), COALESCE(pq.weight, 0)
		FROM planting p
		JOIN seed_info s ON p.seed_id = s.id
		LEFT JOIN product_quality pq ON p.id = pq.planting_id
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list products"})
		return
	}
	defer rows.Close()

	var products []gin.H
	for rows.Next() {
		var product struct {
			ID               int     `json:"id"`
			PlantingDate     string  `json:"planting_date"`
			TransplantingDate *string `json:"transplanting_date"`
			Location         string  `json:"location"`
			Name             string  `json:"name"`
			Variety          string  `json:"variety"`
			VarietyCode      *string `json:"variety_code"`
			SugarContent     float64 `json:"sugar_content"`
			Weight           float64 `json:"weight"`
		}
		
		if err := rows.Scan(
			&product.ID, &product.PlantingDate, &product.TransplantingDate,
			&product.Location, &product.Name, &product.Variety, 
			&product.VarietyCode, &product.SugarContent, &product.Weight,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan product info"})
			return
		}
		
		products = append(products, gin.H{
			"id":                product.ID,
			"planting_date":     product.PlantingDate,
			"transplanting_date": product.TransplantingDate,
			"location":          product.Location,
			"name":              product.Name,
			"variety":           product.Variety,
			"variety_code":      product.VarietyCode,
			"sugar_content":     product.SugarContent,
			"weight":            product.Weight,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"products":   products,
		"total":      totalCount,
		"page":       page,
		"limit":      limit,
		"total_pages": (totalCount + limit - 1) / limit,
	})
}

// GetProductStats 获取产品统计信息
func GetProductStats(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	var stats struct {
		LikeCount     int `json:"like_count"`
		FavoriteCount int `json:"favorite_count"`
		MediaCount    int `json:"media_count"`
	}

	db.DB.QueryRow("SELECT COUNT(*) FROM likes WHERE planting_id = $1", plantingID).Scan(&stats.LikeCount)
	db.DB.QueryRow("SELECT COUNT(*) FROM favorites WHERE planting_id = $1", plantingID).Scan(&stats.FavoriteCount)
	db.DB.QueryRow("SELECT COUNT(*) FROM growth_media WHERE planting_id = $1", plantingID).Scan(&stats.MediaCount)

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}