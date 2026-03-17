package v1

import (
	controller1 "retrovisionarios-api/internal/app/v1/events/controller"
	"retrovisionarios-api/internal/app/v1/videos/controller"

	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	EventController *controller1.EventController
	VideoController *controller.VideoController
}

// SetupRoutes centraliza o agrupamento v1
func SetupRoutes(r *gin.Engine, config RouterConfig) {
	v1 := r.Group("/v1")
	{
		// Agrupamento de Events
		events := v1.Group("/events")
		{
			events.GET("/", config.EventController.GetAll)
			events.POST("/", config.EventController.Create)
			events.DELETE("/:id", config.EventController.Delete)
		}

		// Agrupamento de Videos
		videos := v1.Group("/videos")
		{
			videos.GET("/", config.VideoController.GetAll)
		}
	}
}
