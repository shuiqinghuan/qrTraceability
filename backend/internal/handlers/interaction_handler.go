package handlers

import (
	"net/http"
	"strconv"

	"qr-traceability/internal/models"
	"qr-traceability/internal/service"

	"github.com/gin-gonic/gin"
)

type InteractionHandler struct {
	service *service.InteractionService
}

func NewInteractionHandler(service *service.InteractionService) *InteractionHandler {
	return &InteractionHandler{service: service}
}

// RecordInteraction 记录用户交互
func (h *InteractionHandler) RecordInteraction(c *gin.Context) {
	var interaction models.UserInteraction
	if err := c.ShouldBindJSON(&interaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取用户IP
	interaction.IP = c.ClientIP()

	if err := h.service.RecordInteraction(&interaction); err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, interaction)
}

// GetInteractionStats 获取交互统计
func (h *InteractionHandler) GetInteractionStats(c *gin.Context) {
	batchID, err := strconv.ParseUint(c.Param("batch_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid batch ID"})
		return
	}

	stats, err := h.service.GetInteractionStats(uint(batchID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
