package controllers

import (
	"fmt"
	"net/http"
	"retrovisionarios-api/internal/app/v1/events/models"

	"github.com/gin-gonic/gin"
)

type EventService interface {
	GetAll() ([]models.Event, error)
}

type EventController struct {
	service EventService
}

func NewEventController(service EventService) *EventController {
	return &EventController{service: service}
}

func (c *EventController) GetAll(ctx *gin.Context) {
	events, err := c.service.GetAll()

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
