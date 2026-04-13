package api

import (
	"net/http"
	"strconv"

	"server/internal/db"

	"github.com/gin-gonic/gin"
)

// RelatedProductResponse 相关产品响应结构
type RelatedProductResponse struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Variety          string  `json:"variety"`
	VarietyCode      *string `json:"variety_code"`
	Location         string  `json:"location"`
	PlantingDate     string  `json:"planting_date"`
	SugarContent     float64 `json:"sugar_content"`
	Weight           float64 `json:"weight"`
	LikeCount        int     `json:"like_count"`
	SimilarityScore  float64 `json:"similarity_score"`
}

// GetRelatedProducts 获取相关产品（同地点）
func GetRelatedProducts(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	// 获取当前产品的地点信息
	var currentLocation string
	err = db.DB.QueryRow(`
		SELECT location FROM planting WHERE id = $1
	`, plantingID).Scan(&currentLocation)
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Planting record not found"})
		return
	}

	if currentLocation == "" {
		c.JSON(http.StatusOK, gin.H{
			"related_products": []gin.H{},
			"count": 0,
			"location": currentLocation,
		})
		return
	}

	// 获取同地点的其他产品
	rows, err := db.DB.Query(`
		SELECT p.id, s.name, s.variety, s.variety_code, p.location, p.planting_date,
		       COALESCE(pq.sugar_content, 0), COALESCE(pq.weight, 0),
		       COALESCE(l.like_count, 0)
		FROM planting p
		JOIN seed_info s ON p.seed_id = s.id
		LEFT JOIN product_quality pq ON p.id = pq.planting_id
		LEFT JOIN (
			SELECT planting_id, COUNT(*) as like_count
			FROM likes
			GROUP BY planting_id
		) l ON p.id = l.planting_id
		WHERE p.location = $1 AND p.id != $2
		ORDER BY p.planting_date DESC
		LIMIT 10
	`, currentLocation, plantingID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get related products"})
		return
	}
	defer rows.Close()

	var relatedProducts []RelatedProductResponse
	for rows.Next() {
		var product RelatedProductResponse
		var likeCount int
		
		if err := rows.Scan(
			&product.ID, &product.Name, &product.Variety, &product.VarietyCode,
			&product.Location, &product.PlantingDate, &product.SugarContent,
			&product.Weight, &likeCount,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan product info"})
			return
		}
		
		product.LikeCount = likeCount
		product.SimilarityScore = 1.0 // 同地点产品相似度设为1.0
		
		relatedProducts = append(relatedProducts, product)
	}

	c.JSON(http.StatusOK, gin.H{
		"related_products": relatedProducts,
		"count": len(relatedProducts),
		"location": currentLocation,
		"current_product_id": plantingID,
	})
}

// GetProductsByLocation 获取指定地点的所有产品
func GetProductsByLocation(c *gin.Context) {
	location := c.Param("location")
	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location is required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// 获取该地点的产品总数
	var totalCount int
	err := db.DB.QueryRow(`
		SELECT COUNT(*) FROM planting WHERE location = $1
	`, location).Scan(&totalCount)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count products"})
		return
	}

	// 获取产品列表
	rows, err := db.DB.Query(`
		SELECT p.id, s.name, s.variety, s.variety_code, p.location, p.planting_date,
		       COALESCE(pq.sugar_content, 0), COALESCE(pq.weight, 0),
		       COALESCE(l.like_count, 0), COALESCE(f.favorite_count, 0)
		FROM planting p
		JOIN seed_info s ON p.seed_id = s.id
		LEFT JOIN product_quality pq ON p.id = pq.planting_id
		LEFT JOIN (
			SELECT planting_id, COUNT(*) as like_count
			FROM likes
			GROUP BY planting_id
		) l ON p.id = l.planting_id
		LEFT JOIN (
			SELECT planting_id, COUNT(*) as favorite_count
			FROM favorites
			GROUP BY planting_id
		) f ON p.id = f.planting_id
		WHERE p.location = $1
		ORDER BY p.planting_date DESC
		LIMIT $2 OFFSET $3
	`, location, limit, offset)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products by location"})
		return
	}
	defer rows.Close()

	var products []gin.H
	for rows.Next() {
		var product struct {
			ID            int     `json:"id"`
			Name          string  `json:"name"`
			Variety       string  `json:"variety"`
			VarietyCode   *string `json:"variety_code"`
			Location      string  `json:"location"`
			PlantingDate  string  `json:"planting_date"`
			SugarContent  float64 `json:"sugar_content"`
			Weight        float64 `json:"weight"`
			LikeCount     int     `json:"like_count"`
			FavoriteCount int     `json:"favorite_count"`
		}
		
		if err := rows.Scan(
			&product.ID, &product.Name, &product.Variety, &product.VarietyCode,
			&product.Location, &product.PlantingDate, &product.SugarContent,
			&product.Weight, &product.LikeCount, &product.FavoriteCount,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan product info"})
			return
		}
		
		products = append(products, gin.H{
			"id":             product.ID,
			"name":           product.Name,
			"variety":        product.Variety,
			"variety_code":   product.VarietyCode,
			"location":       product.Location,
			"planting_date":  product.PlantingDate,
			"sugar_content":  product.SugarContent,
			"weight":         product.Weight,
			"like_count":     product.LikeCount,
			"favorite_count": product.FavoriteCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"products":    products,
		"total":       totalCount,
		"page":        page,
		"limit":       limit,
		"total_pages": (totalCount + limit - 1) / limit,
		"location":    location,
	})
}

// GetSimilarProducts 获取相似产品（基于品种、品质等）
func GetSimilarProducts(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	// 获取当前产品的信息
	var currentProduct struct {
		SeedID      int
		Location    string
		VarietyCode *string
		SugarContent float64
		Weight      float64
	}
	
	err = db.DB.QueryRow(`
		SELECT p.seed_id, p.location, s.variety_code, 
		       COALESCE(pq.sugar_content, 0), COALESCE(pq.weight, 0)
		FROM planting p
		JOIN seed_info s ON p.seed_id = s.id
		LEFT JOIN product_quality pq ON p.id = pq.planting_id
		WHERE p.id = $1
	`, plantingID).Scan(
		&currentProduct.SeedID, &currentProduct.Location, 
		&currentProduct.VarietyCode, &currentProduct.SugarContent,
		&currentProduct.Weight,
	)
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Planting record not found"})
		return
	}

	// 构建相似度查询
	query := `
		SELECT p.id, s.name, s.variety, s.variety_code, p.location, p.planting_date,
		       COALESCE(pq.sugar_content, 0), COALESCE(pq.weight, 0),
		       COALESCE(l.like_count, 0),
		       -- 相似度计算
		       CASE 
		         WHEN p.seed_id = $1 THEN 1.0
		         WHEN s.variety_code = $2 AND $2 IS NOT NULL THEN 0.8
		         WHEN p.location = $3 THEN 0.6
		         ELSE 0.3
		       END as similarity_score
		FROM planting p
		JOIN seed_info s ON p.seed_id = s.id
		LEFT JOIN product_quality pq ON p.id = pq.planting_id
		LEFT JOIN (
			SELECT planting_id, COUNT(*) as like_count
			FROM likes
			GROUP BY planting_id
		) l ON p.id = l.planting_id
		WHERE p.id != $4
		ORDER BY similarity_score DESC, p.planting_date DESC
		LIMIT $5
	`

	rows, err := db.DB.Query(query,
		currentProduct.SeedID, currentProduct.VarietyCode,
		currentProduct.Location, plantingID, limit,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get similar products"})
		return
	}
	defer rows.Close()

	var similarProducts []RelatedProductResponse
	for rows.Next() {
		var product RelatedProductResponse
		
		if err := rows.Scan(
			&product.ID, &product.Name, &product.Variety, &product.VarietyCode,
			&product.Location, &product.PlantingDate, &product.SugarContent,
			&product.Weight, &product.LikeCount, &product.SimilarityScore,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan product info"})
			return
		}
		
		similarProducts = append(similarProducts, product)
	}

	c.JSON(http.StatusOK, gin.H{
		"similar_products": similarProducts,
		"count": len(similarProducts),
		"current_product_id": plantingID,
		"similarity_criteria": gin.H{
			"same_seed": currentProduct.SeedID,
			"same_variety_code": currentProduct.VarietyCode,
			"same_location": currentProduct.Location,
		},
	})
}

// GetLocationStats 获取地点统计信息
func GetLocationStats(c *gin.Context) {
	location := c.Param("location")
	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location is required"})
		return
	}

	var stats struct {
		TotalProducts   int     `json:"total_products"`
		AvgSugarContent float64 `json:"avg_sugar_content"`
		AvgWeight       float64 `json:"avg_weight"`
		TotalLikes      int     `json:"total_likes"`
		TotalFavorites  int     `json:"total_favorites"`
		FirstPlanting   *string `json:"first_planting_date"`
		LastPlanting    *string `json:"last_planting_date"`
	}

	// 获取统计信息
	err := db.DB.QueryRow(`
		SELECT 
			COUNT(DISTINCT p.id) as total_products,
			AVG(pq.sugar_content) as avg_sugar_content,
			AVG(pq.weight) as avg_weight,
			COALESCE(SUM(l.like_count), 0) as total_likes,
			COALESCE(SUM(f.favorite_count), 0) as total_favorites,
			MIN(p.planting_date) as first_planting,
			MAX(p.planting_date) as last_planting
		FROM planting p
		LEFT JOIN product_quality pq ON p.id = pq.planting_id
		LEFT JOIN (
			SELECT planting_id, COUNT(*) as like_count
			FROM likes
			GROUP BY planting_id
		) l ON p.id = l.planting_id
		LEFT JOIN (
			SELECT planting_id, COUNT(*) as favorite_count
			FROM favorites
			GROUP BY planting_id
		) f ON p.id = f.planting_id
		WHERE p.location = $1
		GROUP BY p.location
	`, location).Scan(
		&stats.TotalProducts, &stats.AvgSugarContent, &stats.AvgWeight,
		&stats.TotalLikes, &stats.TotalFavorites,
		&stats.FirstPlanting, &stats.LastPlanting,
	)
	
	if err != nil {
		// 可能没有数据，返回空统计
		stats.TotalProducts = 0
		stats.AvgSugarContent = 0
		stats.AvgWeight = 0
		stats.TotalLikes = 0
		stats.TotalFavorites = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"location_stats": stats,
		"location": location,
	})
}

// SearchLocations 搜索地点
func SearchLocations(c *gin.Context) {
	query := c.DefaultQuery("q", "")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	rows, err := db.DB.Query(`
		SELECT DISTINCT location, COUNT(*) as product_count
		FROM planting 
		WHERE location ILIKE $1
		GROUP BY location
		ORDER BY product_count DESC, location
		LIMIT $2
	`, "%"+query+"%", limit)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search locations"})
		return
	}
	defer rows.Close()

	var locations []gin.H
	for rows.Next() {
		var location string
		var productCount int
		
		if err := rows.Scan(&location, &productCount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan location info"})
			return
		}
		
		locations = append(locations, gin.H{
			"location":      location,
			"product_count": productCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"locations": locations,
		"count": len(locations),
		"search_query": query,
	})
}