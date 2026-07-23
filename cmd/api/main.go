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
	executor := database.NewGormExecutor(db)

	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Party{},
		&models.Media{},
		&models.PartyMedia{},
		&models.TicketCategory{},
		&models.Ticket{},
		&models.Purchase{},
		&models.PurchaseItem{},
		&models.PartyMember{},
	); err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_party_ticket_category_name
		ON ticket_categories(name, party_id)
		WHERE deleted_at IS NULL
	`).Error; err != nil {
		log.Fatal(err)
	}

	userRepository := repository.NewUserRepository(db)
	partyRepository := repository.NewPartyRepository(db)
	categoryRepository := repository.NewCategoryRepository(db)
	mediaRepository := repository.NewMediaRepository(db)
	ticketCategoryRepository := repository.NewTicketCategoryRepository(executor)
	ticketRepository := repository.NewTicketRepository(executor)
	partyMemberRepository := repository.NewPartyMemberRepository(db)

	jwt := auth.NewJWT(
		cfg.JWTSecret,
	)

	authService := service.NewAuthService(
		userRepository,
		jwt,
	)

	partyMemberService := service.NewPartyMemberService(
		partyMemberRepository,
	)

	partyService := service.NewPartyService(
		partyRepository,
		partyMemberRepository,
	)

	categoryService := service.NewCategoryService(
		categoryRepository,
	)

	mediaService := service.NewMediaService(
		mediaRepository,
	)

	ticketCategoryService := service.NewTicketCategoryService(
		ticketCategoryRepository,
	)

	ticketService := service.NewTicketService(
		ticketRepository,
		partyRepository,
		ticketCategoryRepository,
		executor,
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
		categoryService,
		mediaService,
		ticketCategoryService,
		ticketService,
		partyMemberService,
	); err != nil {
		log.Fatal(err)
	}
}
