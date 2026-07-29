package server

import (
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/dependencies"
	"github.com/reinp/event-platform/backend/internal/routes"
)

func Register(
	router *gin.Engine,
	deps *dependencies.Container,
) {

	routes.Register(
		router,
		deps,
	)
}
