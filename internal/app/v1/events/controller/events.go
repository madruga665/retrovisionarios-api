package controller

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"retrovisionarios-api/internal/app/v1/events/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type EventService interface {
	GetAll(ctx context.Context, year int, showDeleted bool) ([]models.Event, error)
	GetByID(ctx context.Context, id int) (*models.Event, error)
	Create(ctx context.Context, event *models.Event) error
	Update(ctx context.Context, event *models.UpdateEventRequest) error
	Delete(ctx context.Context, id int) error
}

type EventController struct {
	service EventService
}

func NewEventController(service EventService) *EventController {
	return &EventController{service: service}
}

// Update godoc
// @Summary      Atualizar evento
// @Description  Atualiza as propriedades de um evento existente.
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        id     path      int           true  "ID do evento"
// @Param        event  body      models.Event  true  "Dados do evento para atualizar"
// @Success      200    {object}  models.Event
// @Failure      400    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /events/{id} [patch]
func (c *EventController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var event models.UpdateEventRequest
	if err := ctx.ShouldBindJSON(&event); err != nil {
		slog.Warn("Falha na validação de payload na atualização de evento", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "dados do evento inválidos"})
		return
	}

	event.ID = id

	if err := c.service.Update(ctx.Request.Context(), &event); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Evento não encontrado."})
			return
		}
		slog.Error("Erro interno ao atualizar evento", "error", err, "id", id)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao atualizar evento"})
		return
	}

	updatedEvent, err := c.service.GetByID(ctx, event.ID)

	ctx.JSON(http.StatusOK, updatedEvent)
}

// GetAll godoc
// @Summary      Listar eventos
// @Description  Retorna uma lista de eventos, opcionalmente filtrada por ano e status de deleção.
// @Tags         events
// @Produce      json
// @Param        year     query     int   false  "Ano para filtrar eventos"
// @Param        deleted  query     bool  false  "Se true, lista também eventos deletados"
// @Success      200      {object}  map[string][]models.Event
// @Failure      500      {object}  map[string]string
// @Router       /events [get]
func (c *EventController) GetAll(ctx *gin.Context) {
	yearStr := ctx.Query("year")
	year := 0

	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}

	showDeleted, _ := strconv.ParseBool(ctx.DefaultQuery("deleted", "false"))

	events, err := c.service.GetAll(ctx.Request.Context(), year, showDeleted)

	if err != nil {
		slog.Error("Erro ao buscar eventos", "error", err, "year", year, "showDeleted", showDeleted)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Erro ao buscar eventos",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": events,
	})
}

// Delete godoc
// @Summary      Deletar evento
// @Description  Realiza um soft delete de um evento.
// @Tags         events
// @Param        id   path      int  true  "ID do evento"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /events/{id} [delete]
func (c *EventController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		slog.Error("Erro ao deletar evento", "error", err, "id", id)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar evento"})
		return
	}

	ctx.Status(http.StatusNoContent)
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
