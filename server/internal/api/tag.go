package api

import (
	"net/http"
	"strconv"

	"server/internal/db"

	"github.com/gin-gonic/gin"
)

// TagRequest 标签请求结构
type TagRequest struct {
	TagName string `json:"tag_name" binding:"required"`
}

// AddPlantingTag 添加种植标签
func AddPlantingTag(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	var req TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查种植记录是否存在
	var plantingExists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM planting WHERE id = $1)", plantingID).Scan(&plantingExists)
	if err != nil || !plantingExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Planting record not found"})
		return
	}

	// 插入标签
	var tagID int
	err = db.DB.QueryRow(
		"INSERT INTO planting_tags (planting_id, tag_name) VALUES ($1, $2) RETURNING id",
		plantingID, req.TagName,
	).Scan(&tagID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add tag"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tag added successfully",
		"tag": gin.H{
			"id":        tagID,
			"planting_id": plantingID,
			"tag_name":  req.TagName,
		},
	})
}

// GetPlantingTags 获取种植标签列表
func GetPlantingTags(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	rows, err := db.DB.Query(`
		SELECT id, tag_name
		FROM planting_tags 
		WHERE planting_id = $1
		ORDER BY tag_name
	`, plantingID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tags"})
		return
	}
	defer rows.Close()

	var tags []gin.H
	for rows.Next() {
		var tag struct {
			ID      int    `json:"id"`
			TagName string `json:"tag_name"`
		}
		
		if err := rows.Scan(&tag.ID, &tag.TagName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan tag info"})
			return
		}
		
		tags = append(tags, gin.H{
			"id":       tag.ID,
			"tag_name": tag.TagName,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"tags": tags,
		"count": len(tags),
	})
}

// DeletePlantingTag 删除种植标签
func DeletePlantingTag(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	tagID, err := strconv.Atoi(c.Param("tagId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	// 检查标签是否存在且属于该种植记录
	var exists bool
	err = db.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM planting_tags WHERE id = $1 AND planting_id = $2)",
		tagID, plantingID,
	).Scan(&exists)
	
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	// 删除标签
	_, err = db.DB.Exec("DELETE FROM planting_tags WHERE id = $1", tagID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tag deleted successfully",
	})
}

// UpdatePlantingTag 更新种植标签
func UpdatePlantingTag(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	tagID, err := strconv.Atoi(c.Param("tagId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	var req TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查标签是否存在且属于该种植记录
	var exists bool
	err = db.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM planting_tags WHERE id = $1 AND planting_id = $2)",
		tagID, plantingID,
	).Scan(&exists)
	
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	// 检查新标签名是否已存在（排除当前标签）
	var duplicateExists bool
	err = db.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM planting_tags WHERE planting_id = $1 AND tag_name = $2 AND id != $3)",
		plantingID, req.TagName, tagID,
	).Scan(&duplicateExists)
	
	if err == nil && duplicateExists {
		c.JSON(http.StatusConflict, gin.H{"error": "Tag name already exists"})
		return
	}

	// 更新标签
	_, err = db.DB.Exec(
		"UPDATE planting_tags SET tag_name = $1 WHERE id = $2",
		req.TagName, tagID,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tag updated successfully",
		"tag": gin.H{
			"id":        tagID,
			"planting_id": plantingID,
			"tag_name":  req.TagName,
		},
	})
}

// BatchAddPlantingTags 批量添加种植标签
func BatchAddPlantingTags(c *gin.Context) {
	plantingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid planting ID"})
		return
	}

	var req struct {
		Tags []string `json:"tags" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查种植记录是否存在
	var plantingExists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM planting WHERE id = $1)", plantingID).Scan(&plantingExists)
	if err != nil || !plantingExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Planting record not found"})
		return
	}

	// 批量插入标签
	var addedTags []gin.H
	for _, tagName := range req.Tags {
		if tagName == "" {
			continue
		}
		
		var tagID int
		err := db.DB.QueryRow(
			"INSERT INTO planting_tags (planting_id, tag_name) VALUES ($1, $2) ON CONFLICT (planting_id, tag_name) DO NOTHING RETURNING id",
			plantingID, tagName,
		).Scan(&tagID)
		
		if err == nil {
			addedTags = append(addedTags, gin.H{
				"id":        tagID,
				"planting_id": plantingID,
				"tag_name":  tagName,
			})
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tags added successfully",
		"added_count": len(addedTags),
		"tags": addedTags,
	})
}

// SearchTags 搜索标签
func SearchTags(c *gin.Context) {
	query := c.DefaultQuery("q", "")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	rows, err := db.DB.Query(`
		SELECT DISTINCT tag_name, COUNT(*) as usage_count
		FROM planting_tags 
		WHERE tag_name ILIKE $1
		GROUP BY tag_name
		ORDER BY usage_count DESC, tag_name
		LIMIT $2
	`, "%"+query+"%", limit)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search tags"})
		return
	}
	defer rows.Close()

	var tags []gin.H
	for rows.Next() {
		var tagName string
		var usageCount int
		
		if err := rows.Scan(&tagName, &usageCount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan tag info"})
			return
		}
		
		tags = append(tags, gin.H{
			"tag_name":    tagName,
			"usage_count": usageCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"tags": tags,
		"count": len(tags),
	})
}