package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/saulfrancisco-ruizacevedo/event-weaver-svc/internal/infrastructure/di"
)

const BasePath string = "event-weaver-svc"

func RegisterRoutes(app *di.App, router *gin.Engine) {
	apiGroup := router.Group(fmt.Sprintf("/%s", BasePath))
	{
		app.SpecController.RegisterRoutes(apiGroup)
	}
}
