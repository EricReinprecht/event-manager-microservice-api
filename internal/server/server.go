package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/config"
	"github.com/reinp/event-platform/backend/internal/dependencies"
	"github.com/reinp/event-platform/backend/internal/middleware"
)

type Server struct {
	router *gin.Engine
}

func New(
	cfg *config.Config,
	deps *dependencies.Container,
) *Server {

	router := gin.Default()

	globalLimiter := middleware.NewIPRateLimiter(
		20,
		50,
	)

	router.Use(
		middleware.SecurityHeaders(),
	)

	router.Use(
		cors.New(
			middleware.CORS(
				cfg.CORSAllowedOrigins,
			),
		),
	)

	router.Use(
		globalLimiter.Middleware(),
	)

	Register(
		router,
		deps,
	)

	return &Server{
		router: router,
	}
}

func (s *Server) Start(
	address string,
) error {

	return s.router.Run(address)
}
