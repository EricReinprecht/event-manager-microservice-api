package main

import (
	"log"

	"github.com/reinp/event-platform/backend/internal/auth"
	"github.com/reinp/event-platform/backend/internal/clock"
	"github.com/reinp/event-platform/backend/internal/config"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/payment/paypal"
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

		&models.Ticket{},
		&models.TicketCategory{},
		&models.TicketAccessWindow{},
		&models.TicketScan{},

		&models.Purchase{},
		&models.PurchaseItem{},

		&models.PartyMember{},
		&models.PartyMemberRole{},

		&models.PaymentEvent{},
	); err != nil {
		log.Fatal(err)
	}

	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_party_ticket_category_name
		ON ticket_categories(name, party_id)
		WHERE deleted_at IS NULL
	`).Error; err != nil {
		log.Fatal(err)
	}

	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_ticket_scan_window_active
		ON ticket_scans(ticket_id, ticket_access_window_id)
		WHERE status IN ('PENDING','VERIFIED')
	`).Error; err != nil {
		log.Fatal(err)
	}

	appClock := clock.RealClock{}

	userRepository := repository.NewUserRepository(db)
	partyRepository := repository.NewPartyRepository(executor)
	categoryRepository := repository.NewCategoryRepository(db)
	mediaRepository := repository.NewMediaRepository(db)
	ticketCategoryRepository := repository.NewTicketCategoryRepository(executor)
	ticketRepository := repository.NewTicketRepository(executor)
	purchaseRepository := repository.NewPurchaseRepository(executor)
	partyMemberRepository := repository.NewPartyMemberRepository(executor)
	ticketScanRepository := repository.NewTicketScanRepository(executor)
	ticketAccessWindowRepository := repository.NewTicketAccessWindowRepository(executor)
	paymentEventRepository := repository.NewPaymentEventRepository(executor)
	partyMemberRoleRepository := repository.NewPartyMemberRoleRepository(executor)

	jwt := auth.NewJWT(
		cfg.JWTSecret,
		clock.RealClock{},
	)

	authService := service.NewAuthService(
		userRepository,
		jwt,
	)

	partyMemberService := service.NewPartyMemberService(
		partyMemberRepository,
		partyRepository,
		partyMemberRoleRepository,
	)

	partyService := service.NewPartyService(
		partyRepository,
		partyMemberRepository,
		categoryRepository,
		mediaRepository,
		partyMemberRoleRepository,
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
		partyMemberRepository,
		ticketScanRepository,
		ticketAccessWindowRepository,
		executor,
		appClock,
		cfg.TicketVerificationTTL,
	)

	purchaseService := service.NewPurchaseService(
		purchaseRepository,
		ticketRepository,
	)

	paypalClient := paypal.NewClient(
		cfg.PayPalClientID,
		cfg.PayPalClientSecret,
		cfg.PayPalBaseURL,
		cfg.PayPalReturnURL,
		cfg.PayPalCancelURL,
		cfg.PayPalWebhookID,
	)

	paymentService := service.NewPaymentService(
		purchaseService,
		ticketService,
		partyMemberService,
		paypalClient,
		paymentEventRepository,
		purchaseRepository,
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
		purchaseService,
		paymentService,
		partyMemberService,
	); err != nil {
		log.Fatal(err)
	}
}
