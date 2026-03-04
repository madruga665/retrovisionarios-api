package controller

import (
	"context"
	"log/slog"
	"net/http"
	"retrovisionarios-api/internal/app/v1/events/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventService interface {
	GetAll(ctx context.Context, year int) ([]models.Event, error)
	Create(ctx context.Context, event *models.Event) error
}

type EventController struct {
	service EventService
}

func NewEventController(service EventService) *EventController {
	return &EventController{service: service}
}

// GetAll godoc
// @Summary      Listar eventos
// @Description  Retorna uma lista de eventos, opcionalmente filtrada por ano.
// @Tags         events
// @Produce      json
// @Param        year  query     int  false  "Ano para filtrar eventos"
// @Success      200   {object}  map[string][]models.Event
// @Failure      500   {object}  map[string]string
// @Router       /events [get]
func (c *EventController) GetAll(ctx *gin.Context) {
	yearStr := ctx.Query("year")
	year := 0

	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}

	events, err := c.service.GetAll(ctx.Request.Context(), year)

	if err != nil {
		slog.Error("Erro ao buscar eventos", "error", err, "year", year)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Erro ao buscar eventos",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": events,
	})
}

// Create godoc
// @Summary      Criar evento
// @Description  Cria um novo evento de show.
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        event  body      models.Event  true  "Dados do evento"
// @Success      201    {object}  models.Event
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /events [post]
func (c *EventController) Create(ctx *gin.Context) {
	var event models.Event

	if err := ctx.ShouldBindJSON(&event); err != nil {
		slog.Warn("Falha na validação de payload na criação de evento", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "os campos 'date' e 'name' são obrigatórios e devem estar no formato correto"})
		return
	}

	if err := c.service.Create(ctx.Request.Context(), &event); err != nil {
		slog.Error("Erro interno ao criar evento", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao criar evento"})
		return
	}

	ctx.JSON(http.StatusCreated, event)
}
