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
