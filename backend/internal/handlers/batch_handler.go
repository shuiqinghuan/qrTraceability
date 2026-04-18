package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"qr-traceability/internal/models"
	"qr-traceability/internal/service"
	"qr-traceability/internal/utils"

	"github.com/gin-gonic/gin"
)

type BatchHandler struct {
	service *service.BatchService
}

func NewBatchHandler(service *service.BatchService) *BatchHandler {
	return &BatchHandler{service: service}
}

// CreateBatch 创建批次
func (h *BatchHandler) CreateBatch(c *gin.Context) {
	var batch models.ProductBatch
	if err := c.ShouldBindJSON(&batch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateBatch(&batch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, batch)
}

// GetBatchByID 根据ID获取批次
func (h *BatchHandler) GetBatchByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid batch ID"})
		return
	}

	batch, err := h.service.GetBatchByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Batch not found"})
		return
	}

	c.JSON(http.StatusOK, batch)
}

// GetBatchByUniqueID 根据唯一ID获取批次
func (h *BatchHandler) GetBatchByUniqueID(c *gin.Context) {
	uniqueID := c.Param("unique_id")
	if uniqueID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unique ID is required"})
		return
	}

	batch, err := h.service.GetBatchByUniqueID(uniqueID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Batch not found"})
		return
	}

	c.JSON(http.StatusOK, batch)
}

// ListBatchesByProductID 根据产品ID列出批次
func (h *BatchHandler) ListBatchesByProductID(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	batches, err := h.service.ListBatchesByProductID(uint(productID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, batches)
}

// UpdateBatch 更新批次
func (h *BatchHandler) UpdateBatch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid batch ID"})
		return
	}

	var batch models.ProductBatch
	if err := c.ShouldBindJSON(&batch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	batch.ID = uint(id)
	if err := h.service.UpdateBatch(&batch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, batch)
}

// DeleteBatch 删除批次
func (h *BatchHandler) DeleteBatch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid batch ID"})
		return
	}

	if err := h.service.DeleteBatch(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Batch deleted successfully"})
}

// AddMedia 添加媒体文件
func (h *BatchHandler) AddMedia(c *gin.Context) {
	batchID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid batch ID"})
		return
	}

	var media models.MediaFile
	if err := c.ShouldBindJSON(&media); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddMedia(uint(batchID), &media); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, media)
}

// UpdateHarvestQuality 更新采收质量
func (h *BatchHandler) UpdateHarvestQuality(c *gin.Context) {
	batchID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid batch ID"})
		return
	}

	var harvest models.HarvestQuality
	if err := c.ShouldBindJSON(&harvest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateHarvestQuality(uint(batchID), &harvest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, harvest)
}

// GenerateQRCode 为批次生成二维码
func (h *BatchHandler) GenerateQRCode(c *gin.Context) {
	batchID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid batch ID"})
		return
	}

	// 获取批次信息
	batch, err := h.service.GetBatchByID(uint(batchID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Batch not found"})
		return
	}

	// 生成产品详情页URL
	// 注意：在实际生产环境中，应该从配置中获取域名
	productURL := fmt.Sprintf("https://example.com/product/%s", batch.UniqueID)

	// 生成二维码
	qrCode, err := utils.GenerateQRCode(productURL, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	// 设置响应头并返回二维码
	c.Header("Content-Type", "image/png")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=qr_%s.png", batch.UniqueID))
	c.Data(http.StatusOK, "image/png", qrCode)
}
