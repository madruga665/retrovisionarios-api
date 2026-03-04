package controller

import (
	"context"
	"log/slog"
	"net/http"
	"retrovisionarios-api/internal/app/v1/videos/models"

	"github.com/gin-gonic/gin"
)

type VideoService interface {
	GetAll(ctx context.Context) (map[string][]models.Video, error)
}

type VideoController struct {
	service VideoService
}

func NewVideoController(service VideoService) *VideoController {
	return &VideoController{service: service}
}

// GetAll godoc
// @Summary      Listar videos
// @Description  Retorna uma lista de videos
// @Tags         videos
// @Produce      json
// @Success      200   {object}  map[string][]models.Video
// @Failure      500   {object}  map[string]string
// @Router       /videos [get]
func (c *VideoController) GetAll(ctx *gin.Context) {
	videos, err := c.service.GetAll(ctx.Request.Context())

	if err != nil {
		slog.Error("Erro ao buscar videos", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Erro ao buscar videos",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": videos,
	})
}
