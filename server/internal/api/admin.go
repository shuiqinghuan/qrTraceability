package api

import (
	"net/http"
	"strconv"
	"time"

	"server/internal/db"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AdminLoginRequest 管理员登录请求
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AdminLoginResponse 管理员登录响应
type AdminLoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	Admin   gin.H  `json:"admin,omitempty"`
}

// AdminProductRequest 管理员产品请求
type AdminProductRequest struct {
	SeedID           int     `json:"seed_id" binding:"required"`
	PlantingDate     string  `json:"planting_date" binding:"required"`
	TransplantingDate *string `json:"transplanting_date"`
	Location         string  `json:"location" binding:"required"`
	Notes            string  `json:"notes"`
	Tags             []string `json:"tags"`
}

// AdminMediaUploadRequest 媒体上传请求
type AdminMediaUploadRequest struct {
	PlantingID  int    `json:"planting_id" binding:"required"`
	MediaType   string `json:"media_type" binding:"required"`
	Description string `json:"description"`
}

// AdminQualityRequest 管理员品质请求
type AdminQualityRequest struct {
	PlantingID      int     `json:"planting_id" binding:"required"`
	HarvestStartDate string `json:"harvest_start_date"`
	HarvestEndDate   string `json:"harvest_end_date"`
	SugarContent    float64 `json:"sugar_content"`
	Weight          float64 `json:"weight"`
	TasteDescription string `json:"taste_description"`
	SuitableFor     string  `json:"suitable_for"`
	QualitySummary  string  `json:"quality_summary"`
}

// AdminLogin 管理员登录
func AdminLogin(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询管理员用户
	var admin struct {
		ID           int    `json:"id"`
		Username     string `json:"username"`
		PasswordHash string `json:"-"`
		Role         string `json:"role"`
	}
	
	err := db.DB.QueryRow(`
		SELECT id, username, password_hash, role
		FROM admin_users 
		WHERE username = $1
	`, req.Username).Scan(&admin.ID, &admin.Username, &admin.PasswordHash, &admin.Role)
	
	if err != nil {
		response := AdminLoginResponse{
			Success: false,
			Message: "用户名或密码错误",
		}
		c.JSON(http.StatusUnauthorized, gin.H{"response": response})
		return
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password))
	if err != nil {
		response := AdminLoginResponse{
			Success: false,
			Message: "用户名或密码错误",
		}
		c.JSON(http.StatusUnauthorized, gin.H{"response": response})
		return
	}

	// 生成JWT令牌（这里简化处理，实际应该使用JWT）
	token := generateAdminToken(admin.ID, admin.Username, admin.Role)

	response := AdminLoginResponse{
		Success: true,
		Message: "登录成功",
		Token:   token,
		Admin: gin.H{
			"id":       admin.ID,
			"username": admin.Username,
			"role":     admin.Role,
		},
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}

// generateAdminToken 生成管理员令牌（简化版）
func generateAdminToken(adminID int, username, role string) string {
	// 这里应该使用JWT生成安全的令牌
	// 简化处理：返回一个组合字符串
	timestamp := time.Now().Unix()
	return "admin_" + username + "_" + strconv.FormatInt(timestamp, 10) + "_" + role
}

// AdminCreateProduct 管理员创建产品
func AdminCreateProduct(c *gin.Context) {
	var req AdminProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证种子信息是否存在
	var seedExists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM seed_info WHERE id = $1)", req.SeedID).Scan(&seedExists)
	if err != nil || !seedExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Seed info not found"})
		return
	}

	// 获取管理员ID
	adminID, exists := c.Get("adminID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin authentication required"})
		return
	}

	// 创建种植记录
	var plantingID int
	err = db.DB.QueryRow(`
		INSERT INTO planting (seed_id, planting_date, transplanting_date, location, notes, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, req.SeedID, req.PlantingDate, req.TransplantingDate, req.Location, req.Notes, adminID).Scan(&plantingID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	// 添加标签
	if len(req.Tags) > 0 {
		for _, tag := range req.Tags {
			if tag != "" {
				db.DB.Exec(`
					INSERT INTO planting_tags (planting_id, tag_name)
					VALUES ($1, $2)
					ON CONFLICT (planting_id, tag_name) DO NOTHING
				`, plantingID, tag)
			}
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"product": gin.H{
			"id":                plantingID,
			"seed_id":           req.SeedID,
			"planting_date":     req.PlantingDate,
			"transplanting_date": req.TransplantingDate,
			"location":          req.Location,
			"notes":             req.Notes,
			"tags":              req.Tags,
			"created_by":        adminID,
		},
	})
}

// AdminUpdateProduct 管理员更新产品
func AdminUpdateProduct(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	var req AdminProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查产品是否存在
	var exists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM planting WHERE id = $1)", plantingID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// 验证种子信息是否存在
	var seedExists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM seed_info WHERE id = $1)", req.SeedID).Scan(&seedExists)
	if err != nil || !seedExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Seed info not found"})
		return
	}

	// 更新种植记录
	_, err = db.DB.Exec(`
		UPDATE planting 
		SET seed_id = $1, planting_date = $2, transplanting_date = $3, 
		    location = $4, notes = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
	`, req.SeedID, req.PlantingDate, req.TransplantingDate, req.Location, req.Notes, plantingID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	// 更新标签（先删除所有标签，再添加新标签）
	_, err = db.DB.Exec("DELETE FROM planting_tags WHERE planting_id = $1", plantingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tags"})
		return
	}

	// 添加新标签
	if len(req.Tags) > 0 {
		for _, tag := range req.Tags {
			if tag != "" {
				db.DB.Exec(`
					INSERT INTO planting_tags (planting_id, tag_name)
					VALUES ($1, $2)
				`, plantingID, tag)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"product": gin.H{
			"id":                plantingID,
			"seed_id":           req.SeedID,
			"planting_date":     req.PlantingDate,
			"transplanting_date": req.TransplantingDate,
			"location":          req.Location,
			"notes":             req.Notes,
			"tags":              req.Tags,
		},
	})
}

// AdminDeleteProduct 管理员删除产品
func AdminDeleteProduct(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	// 检查产品是否存在
	var exists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM planting WHERE id = $1)", plantingID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// 开始事务
	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// 删除相关数据
	_, err = tx.Exec("DELETE FROM planting_tags WHERE planting_id = $1", plantingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tags"})
		return
	}

	_, err = tx.Exec("DELETE FROM growth_media WHERE planting_id = $1", plantingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete media"})
		return
	}

	_, err = tx.Exec("DELETE FROM product_quality WHERE planting_id = $1", plantingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete quality records"})
		return
	}

	_, err = tx.Exec("DELETE FROM likes WHERE planting_id = $1", plantingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete likes"})
		return
	}

	_, err = tx.Exec("DELETE FROM favorites WHERE planting_id = $1", plantingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete favorites"})
		return
	}

	_, err = tx.Exec("DELETE FROM ip_like_restrictions WHERE planting_id = $1", plantingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete IP restrictions"})
		return
	}

	// 删除种植记录
	_, err = tx.Exec("DELETE FROM planting WHERE id = $1", plantingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
		"product_id": plantingID,
	})
}

// AdminListProducts 管理员获取产品列表
func AdminListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	search := c.DefaultQuery("search", "")
	location := c.DefaultQuery("location", "")
	seedID, _ := strconv.Atoi(c.Query("seed_id"))

	// 构建查询条件
	whereClause := "WHERE 1=1"
	params := []interface{}{}
	paramCount := 0

	if search != "" {
		paramCount++
		whereClause += " AND (s.name ILIKE $" + strconv.Itoa(paramCount) + 
			" OR s.variety ILIKE $" + strconv.Itoa(paramCount) + 
			" OR p.location ILIKE $" + strconv.Itoa(paramCount) + ")"
		params = append(params, "%"+search+"%")
	}

	if location != "" {
		paramCount++
		whereClause += " AND p.location = $" + strconv.Itoa(paramCount)
		params = append(params, location)
	}

	if seedID > 0 {
		paramCount++
		whereClause += " AND p.seed_id = $" + strconv.Itoa(paramCount)
		params = append(params, seedID)
	}

	// 获取产品总数
	countQuery := `
		SELECT COUNT(*)
		FROM planting p
		JOIN seed_info s ON p.seed_id = s.id
	` + whereClause

	var totalCount int
	err := db.DB.QueryRow(countQuery, params...).Scan(&totalCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count products"})
		return
	}

	// 获取产品列表
	paramCount = 0
	paramsWithLimit := make([]interface{}, len(params))
	copy(paramsWithLimit, params)
	
	paramCount = len(params)
	paramsWithLimit = append(paramsWithLimit, limit, offset)

	query := `
		SELECT p.id, p.planting_date, p.transplanting_date, p.location, p.notes,
		       s.id as seed_id, s.name as seed_name, s.variety, s.variety_code,
		       u.phone_number as creator,
		       COALESCE(pq.sugar_content, 0), COALESCE(pq.weight, 0),
		       COALESCE(l.like_count, 0), COALESCE(f.favorite_count, 0),
		       COALESCE(m.media_count, 0)
		FROM planting p
		JOIN seed_info s ON p.seed_id = s.id
		JOIN users u ON p.created_by = u.id
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
		LEFT JOIN (
			SELECT planting_id, COUNT(*) as media_count
			FROM growth_media
			GROUP BY planting_id
		) m ON p.id = m.planting_id
	` + whereClause + `
		ORDER BY p.created_at DESC
		LIMIT $` + strconv.Itoa(paramCount+1) + ` OFFSET $` + strconv.Itoa(paramCount+2)

	rows, err := db.DB.Query(query, paramsWithLimit...)
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
			Notes            string  `json:"notes"`
			SeedID           int     `json:"seed_id"`
			SeedName         string  `json:"seed_name"`
			Variety          string  `json:"variety"`
			VarietyCode      *string `json:"variety_code"`
			Creator          string  `json:"creator"`
			SugarContent     float64 `json:"sugar_content"`
			Weight           float64 `json:"weight"`
			LikeCount        int     `json:"like_count"`
			FavoriteCount    int     `json:"favorite_count"`
			MediaCount       int     `json:"media_count"`
		}
		
		if err := rows.Scan(
			&product.ID, &product.PlantingDate, &product.TransplantingDate,
			&product.Location, &product.Notes, &product.SeedID, &product.SeedName,
			&product.Variety, &product.VarietyCode, &product.Creator,
			&product.SugarContent, &product.Weight, &product.LikeCount,
			&product.FavoriteCount, &product.MediaCount,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan product info"})
			return
		}
		
		products = append(products, gin.H{
			"id":                product.ID,
			"planting_date":     product.PlantingDate,
			"transplanting_date": product.TransplantingDate,
			"location":          product.Location,
			"notes":             product.Notes,
			"seed_id":           product.SeedID,
			"seed_name":         product.SeedName,
			"variety":           product.Variety,
			"variety_code":      product.VarietyCode,
			"creator":           product.Creator,
			"sugar_content":     product.SugarContent,
			"weight":            product.Weight,
			"like_count":        product.LikeCount,
			"favorite_count":    product.FavoriteCount,
			"media_count":       product.MediaCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"products":    products,
		"total":       totalCount,
		"page":        page,
		"limit":       limit,
		"total_pages": (totalCount + limit - 1) / limit,
		"filters": gin.H{
			"search":   search,
			"location": location,
			"seed_id":  seedID,
		},
	})
}

// AdminGetProductDetail 管理员获取产品详情
func AdminGetProductDetail(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	// 使用产品详情API获取基本信息
	// 这里可以调用GetProductDetail函数或直接查询
	var productDetail ProductDetailResponse
	// 简化处理，实际应该查询数据库
	
	c.JSON(http.StatusOK, gin.H{
		"product": productDetail,
	})
}

// AdminUploadMedia 管理员上传媒体文件
func AdminUploadMedia(c *gin.Context) {
	var req AdminMediaUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查种植记录是否存在
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM planting WHERE id = $1)", req.PlantingID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Planting record not found"})
		return
	}

	// 验证媒体类型
	if req.MediaType != "image" && req.MediaType != "video" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media type. Must be 'image' or 'video'"})
		return
	}

	// 获取管理员ID
	adminID, exists := c.Get("adminID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin authentication required"})
		return
	}

	// 这里应该处理文件上传，简化处理
	// 实际应该保存文件到服务器并返回文件路径
	filePath := "/uploads/" + req.MediaType + "/" + strconv.Itoa(req.PlantingID) + "_" + strconv.FormatInt(time.Now().Unix(), 10)

	// 插入媒体记录
	var mediaID int
	err = db.DB.QueryRow(`
		INSERT INTO growth_media (planting_id, media_type, file_path, description, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, req.PlantingID, req.MediaType, filePath, req.Description, adminID).Scan(&mediaID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save media record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Media uploaded successfully",
		"media": gin.H{
			"id":          mediaID,
			"planting_id": req.PlantingID,
			"media_type":  req.MediaType,
			"file_path":   filePath,
			"description": req.Description,
			"created_by":  adminID,
		},
	})
}

// AdminUpdateQuality 管理员更新品质信息
func AdminUpdateQuality(c *gin.Context) {
	var req AdminQualityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查种植记录是否存在
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM planting WHERE id = $1)", req.PlantingID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Planting record not found"})
		return
	}

	// 获取管理员ID
	adminID, exists := c.Get("adminID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin authentication required"})
		return
	}

	// 检查是否已存在品质记录
	var qualityID int
	err = db.DB.QueryRow(`
		SELECT id FROM product_quality WHERE planting_id = $1
	`, req.PlantingID).Scan(&qualityID)
	
	if err != nil {
		// 不存在，创建新记录
		err = db.DB.QueryRow(`
			INSERT INTO product_quality (planting_id, harvest_start_date, harvest_end_date, 
				sugar_content, weight, taste_description, suitable_for, quality_summary, created_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id
		`, req.PlantingID, req.HarvestStartDate, req.HarvestEndDate,
			req.SugarContent, req.Weight, req.TasteDescription,
			req.SuitableFor, req.QualitySummary, adminID).Scan(&qualityID)
	} else {
		// 已存在，更新记录
		_, err = db.DB.Exec(`
			UPDATE product_quality 
			SET harvest_start_date = $1, harvest_end_date = $2,
				sugar_content = $3, weight = $4, taste_description = $5,
				suitable_for = $6, quality_summary = $7, updated_at = CURRENT_TIMESTAMP
			WHERE id = $8
		`, req.HarvestStartDate, req.HarvestEndDate,
			req.SugarContent, req.Weight, req.TasteDescription,
			req.SuitableFor, req.QualitySummary, qualityID)
	}
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save quality record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quality information saved successfully",
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