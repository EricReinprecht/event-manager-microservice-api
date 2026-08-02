package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/config"
	"github.com/reinp/event-platform/backend/internal/dependencies"
	"github.com/reinp/event-platform/backend/internal/i18n"
	"github.com/reinp/event-platform/backend/internal/middleware"
)

type Server struct {
	router *gin.Engine
}

func New(
	cfg *config.Config,
	deps *dependencies.Container,
) (*Server, error) {

	router := gin.Default()
	router.Static("/uploads", "./uploads")

	translationRegistry, err := i18n.NewRegistry(
		"en",
		"de",
	)

	if err != nil {
		return nil, err
	}

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
		middleware.Language(
			translationRegistry,
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
	}, nil
}

func (s *Server) Start(
	address string,
) error {

	return s.router.Run(address)
}
