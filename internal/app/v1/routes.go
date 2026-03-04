package v1

import (
	"retrovisionarios-api/internal/app/v1/events/eventControllers"
	"retrovisionarios-api/internal/app/v1/videos/videoControllers"

	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	EventController *eventControllers.EventController
	VideoController *videoControllers.VideoController
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
		}

		// Agrupamento de Videos
		videos := v1.Group("/videos")
		{
			videos.GET("/", config.VideoController.GetAll)
		}
	}
}
