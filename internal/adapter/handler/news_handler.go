package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shivangsaxena/inshorts-task/internal/core/usecase"
)

type NewsHandler struct {
	useCase *usecase.NewsUseCase
}

func NewNewsHandler(uc *usecase.NewsUseCase) *NewsHandler {
	return &NewsHandler{useCase: uc}
}

func (h *NewsHandler) GetNews(c *gin.Context) {
	query := c.Query("q")
	latStr := c.Query("lat")
	lngStr := c.Query("lng")

	var lat, lng float64
	if latStr != "" && lngStr != "" {
		lat, _ = strconv.ParseFloat(latStr, 64)
		lng, _ = strconv.ParseFloat(lngStr, 64)
	}

	result, err := h.useCase.GetNews(c.Request.Context(), query, lat, lng)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *NewsHandler) GetTrending(c *gin.Context) {
	articles, err := h.useCase.GetTrending(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"articles": articles})
}

func (h *NewsHandler) RecordView(c *gin.Context) {
	id := c.Param("id")
	if err := h.useCase.RecordView(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func RegisterRoutes(r *gin.Engine, h *NewsHandler) {
	api := r.Group("/api/v1")
	{
		api.GET("/news", h.GetNews)
		api.GET("/news/nearby", h.GetNews)
		api.GET("/news/trending", h.GetTrending)
		api.POST("/news/:id/view", h.RecordView)
	}
}
