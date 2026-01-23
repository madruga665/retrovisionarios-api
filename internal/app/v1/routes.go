package v1

import (
	"retrovisionarios-api/internal/app/v1/events/controllers"

	"github.com/gin-gonic/gin"
)

func EventRoutes(r *gin.Engine, c *controllers.EventController) {
	v1 := r.Group("/v1")
	{
		v1.GET("/events", c.GetAll)
	}
}
