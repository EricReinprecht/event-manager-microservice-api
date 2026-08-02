package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/dependencies"
	"github.com/reinp/event-platform/backend/internal/handlers"
	"github.com/reinp/event-platform/backend/internal/routes/constants"
)

func registerMediaRoutes(
	router *gin.RouterGroup,
	deps *dependencies.Container,
) {

	handler := handlers.NewMediaHandler(
		deps.MediaService,
	)

	router.POST(
		constants.MediaUpload,
		handler.Upload,
	)
}
