package controllers

import (
	"fmt"
	"net/http"
	"retrovisionarios-api/internal/app/v1/events/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventService interface {
	GetAll(year int) ([]models.Event, error)
	Create(event *models.Event) error
}

type EventController struct {
	service EventService
}

func NewEventController(service EventService) *EventController {
	return &EventController{service: service}
}

func (c *EventController) GetAll(ctx *gin.Context) {
	yearStr := ctx.Query("year")
	year := 0

	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}

	events, err := c.service.GetAll(year)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Erro no service Events GetAll",
		})

		fmt.Println(err)

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": events,
	})
}

func (c *EventController) Create(ctx *gin.Context) {
	var event models.Event

	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Payload malformado ou campos inválidos"})
		fmt.Println(err)
		return
	}

	if event.Name == "" || event.Date.IsZero() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Os campos 'date' e 'name' são obrigatórios"})
		return
	}

	if err := c.service.Create(&event); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno ao criar evento"})
		fmt.Println(err)
		return
	}

	ctx.JSON(http.StatusCreated, event)
}
