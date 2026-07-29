package app

import (
	"log"

	"github.com/reinp/event-platform/backend/internal/config"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/server"
)

type App struct {
	cfg *config.Config

	server *server.Server
}

func New() (*App, error) {

	cfg := config.Load()

	db, err := database.Connect(cfg)

	if err != nil {
		return nil, err
	}

	executor := database.NewGormExecutor(db)

	if err := database.Migrate(db); err != nil {
		return nil, err
	}

	deps, err := BuildDependencies(
		cfg,
		executor,
	)

	if err != nil {
		return nil, err
	}

	return &App{
		cfg: cfg,

		server: server.New(
			cfg,
			deps,
		),
	}, nil
}

func (a *App) Start() error {

	log.Println("starting api")

	return a.server.Start(
		":" + a.cfg.Port,
	)
}
