package main

import (
	"log"

	"github.com/reinp/event-platform/backend/internal/auth"
	"github.com/reinp/event-platform/backend/internal/config"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/repository"
	"github.com/reinp/event-platform/backend/internal/server"
	"github.com/reinp/event-platform/backend/internal/service"
)

func main() {

	cfg := config.Load()

	db, err := database.Connect(cfg)

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Party{},
		&models.Media{},
		&models.PartyMedia{},
	)

	if err != nil {
		log.Fatal(err)
	}

	userRepository := repository.NewUserRepository(db)
	partyRepository := repository.NewPartyRepository(db)

	jwt := auth.NewJWT(
		cfg.JWTSecret,
	)

	authService := service.NewAuthService(
		userRepository,
		jwt,
	)

	partyService := service.NewPartyService(
		partyRepository,
	)

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatal(err)
	}

	defer sqlDB.Close()

	log.Println("Database connected")

	if err := server.Start(
		":"+cfg.Port,
		authService,
		partyService,
	); err != nil {
		log.Fatal(err)
	}
}
