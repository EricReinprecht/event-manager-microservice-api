package refund

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/handlers"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/routes"
	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func TestOrganizerReceivesHTTP200WhenRefunding(t *testing.T) {

	gin.SetMode(
		gin.TestMode,
	)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

	handler := handlers.NewPaymentHandler(
		paymentService,
	)

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleOrganizer,
	)

	router := gin.New()

	router.POST(
		"/api/purchases/:id/refund",
		func(c *gin.Context) {

			// temporary auth simulation
			c.Set(
				"userID",
				scenario.Organizer.ID,
			)

			handler.Refund(c)
		},
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/purchases/"+scenario.Purchase.ID.String()+"/refund",
		bytes.NewBuffer(nil),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		req,
	)

	if recorder.Code != http.StatusOK {

		t.Fatalf(
			"expected status 200 got %d body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	if !fakeGateway.RefundCalled {

		t.Fatal(
			"expected refund gateway to be called",
		)
	}
}

func TestAdminReceivesHTTP200WhenRefunding(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

	paymentHandler := handlers.NewPaymentHandler(
		paymentService,
	)

	// -------------------------
	// Scenario
	// -------------------------

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleAdmin,
	)

	router := gin.New()

	router.POST(
		"/api/purchases/:id/refund",
		func(c *gin.Context) {

			c.Set(
				"userID",
				scenario.Actor.ID,
			)

			paymentHandler.Refund(c)
		},
	)

	// -------------------------
	// Execute HTTP request
	// -------------------------

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/purchases/"+scenario.Purchase.ID.String()+"/refund",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		req,
	)

	// -------------------------
	// Verify response
	// -------------------------

	if recorder.Code != http.StatusOK {

		t.Fatalf(
			"expected HTTP 200 got %d body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	// -------------------------
	// Verify gateway
	// -------------------------

	if !fakeGateway.RefundCalled {

		t.Fatal(
			"expected refund gateway to be called",
		)
	}

	// -------------------------
	// Verify purchase
	// -------------------------

	var updated models.Purchase

	if err := db.First(
		&updated,
		"id = ?",
		scenario.Purchase.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	if updated.Status != enum.PurchaseStatusRefunded {

		t.Fatalf(
			"expected REFUNDED got %s",
			updated.Status,
		)
	}
}

func TestRefunderReceivesHTTP200WhenRefunding(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

	paymentHandler := handlers.NewPaymentHandler(
		paymentService,
	)

	// -------------------------
	// Scenario
	// -------------------------

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleRefunder,
	)

	// -------------------------
	// Router
	// -------------------------

	router := gin.New()

	router.POST(
		"/api/purchases/:id/refund",
		func(c *gin.Context) {

			c.Set(
				"userID",
				scenario.Actor.ID,
			)

			paymentHandler.Refund(c)
		},
	)

	// -------------------------
	// Execute request
	// -------------------------

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/purchases/"+scenario.Purchase.ID.String()+"/refund",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		req,
	)

	// -------------------------
	// Verify HTTP
	// -------------------------

	if recorder.Code != http.StatusOK {

		t.Fatalf(
			"expected HTTP 200 got %d body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	// -------------------------
	// Verify gateway
	// -------------------------

	if !fakeGateway.RefundCalled {

		t.Fatal(
			"expected refund gateway to be called",
		)
	}

	// -------------------------
	// Verify purchase
	// -------------------------

	var updated models.Purchase

	if err := db.First(
		&updated,
		"id = ?",
		scenario.Purchase.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	if updated.Status != enum.PurchaseStatusRefunded {

		t.Fatalf(
			"expected REFUNDED got %s",
			updated.Status,
		)
	}
}

func TestStaffReceivesHTTP403WhenRefunding(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

	paymentHandler := handlers.NewPaymentHandler(
		paymentService,
	)

	// -------------------------
	// Scenario
	// -------------------------

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleStaff,
	)

	// -------------------------
	// Router
	// -------------------------

	router := gin.New()

	router.POST(
		"/api/purchases/:id/refund",
		func(c *gin.Context) {

			c.Set(
				"userID",
				scenario.Actor.ID,
			)

			paymentHandler.Refund(c)
		},
	)

	// -------------------------
	// Execute request
	// -------------------------

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/purchases/"+scenario.Purchase.ID.String()+"/refund",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		req,
	)

	// -------------------------
	// Verify HTTP 403
	// -------------------------

	if recorder.Code != http.StatusForbidden {

		t.Fatalf(
			"expected HTTP 403 got %d body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	// -------------------------
	// Verify gateway untouched
	// -------------------------

	if fakeGateway.RefundCalled {

		t.Fatal(
			"staff should not reach refund gateway",
		)
	}

	// -------------------------
	// Verify purchase unchanged
	// -------------------------

	var updated models.Purchase

	if err := db.First(
		&updated,
		"id = ?",
		scenario.Purchase.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	if updated.Status != enum.PurchaseStatusPaid {

		t.Fatalf(
			"expected PAID got %s",
			updated.Status,
		)
	}
}

func TestAttendeeReceivesHTTP403WhenRefunding(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

	paymentHandler := handlers.NewPaymentHandler(
		paymentService,
	)

	// -------------------------
	// Scenario
	// -------------------------

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleAttendee,
	)

	// -------------------------
	// Router
	// -------------------------

	router := gin.New()

	router.POST(
		"/api/purchases/:id/refund",
		func(c *gin.Context) {

			c.Set(
				"userID",
				scenario.Actor.ID,
			)

			paymentHandler.Refund(c)
		},
	)

	// -------------------------
	// Execute request
	// -------------------------

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/purchases/"+scenario.Purchase.ID.String()+"/refund",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		req,
	)

	// -------------------------
	// Verify HTTP 403
	// -------------------------

	if recorder.Code != http.StatusForbidden {

		t.Fatalf(
			"expected HTTP 403 got %d body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	// -------------------------
	// Verify gateway untouched
	// -------------------------

	if fakeGateway.RefundCalled {

		t.Fatal(
			"attendee should not reach refund gateway",
		)
	}

	// -------------------------
	// Verify purchase unchanged
	// -------------------------

	var updated models.Purchase

	if err := db.First(
		&updated,
		"id = ?",
		scenario.Purchase.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	if updated.Status != enum.PurchaseStatusPaid {

		t.Fatalf(
			"expected PAID got %s",
			updated.Status,
		)
	}
}

func TestNonMemberReceivesHTTP403WhenRefunding(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

	paymentHandler := handlers.NewPaymentHandler(
		paymentService,
	)

	// -------------------------
	// Create valid refund target
	// -------------------------

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleAttendee,
	)

	// -------------------------
	// Create outsider user
	// -------------------------

	outsider := fixtures.User()

	if err := db.Create(
		&outsider,
	).Error; err != nil {

		t.Fatal(err)
	}

	// -------------------------
	// Router
	// -------------------------

	router := gin.New()

	router.POST(
		"/api/purchases/:id/refund",
		func(c *gin.Context) {

			c.Set(
				"userID",
				outsider.ID,
			)

			paymentHandler.Refund(c)
		},
	)

	// -------------------------
	// Execute request
	// -------------------------

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/purchases/"+scenario.Purchase.ID.String()+"/refund",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		req,
	)

	// -------------------------
	// Verify HTTP 403
	// -------------------------

	if recorder.Code != http.StatusForbidden {

		t.Fatalf(
			"expected HTTP 403 got %d body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	// -------------------------
	// Verify gateway untouched
	// -------------------------

	if fakeGateway.RefundCalled {

		t.Fatal(
			"non-member should not reach refund gateway",
		)
	}

	// -------------------------
	// Verify purchase unchanged
	// -------------------------

	var updated models.Purchase

	if err := db.First(
		&updated,
		"id = ?",
		scenario.Purchase.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	if updated.Status != enum.PurchaseStatusPaid {

		t.Fatalf(
			"expected PAID got %s",
			updated.Status,
		)
	}
}

func TestUnauthenticatedUserReceivesHTTP401WhenRefunding(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	fakeGateway := &helpers.FakePaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		fakeGateway,
	)

	paymentHandler := handlers.NewPaymentHandler(
		paymentService,
	)

	// -------------------------
	// Scenario
	// -------------------------

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleOrganizer,
	)

	// -------------------------
	// Router
	// -------------------------

	router := gin.New()

	router.POST(
		"/api/purchases/:id/refund",
		paymentHandler.Refund,
	)

	// -------------------------
	// Execute request
	// -------------------------

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/purchases/"+scenario.Purchase.ID.String()+"/refund",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		req,
	)

	// -------------------------
	// Verify HTTP 401
	// -------------------------

	if recorder.Code != http.StatusUnauthorized {

		t.Fatalf(
			"expected HTTP 401 got %d body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	// -------------------------
	// Verify gateway untouched
	// -------------------------

	if fakeGateway.RefundCalled {

		t.Fatal(
			"unauthenticated user should not reach refund gateway",
		)
	}

	// -------------------------
	// Verify purchase unchanged
	// -------------------------

	var updated models.Purchase

	if err := db.First(
		&updated,
		"id = ?",
		scenario.Purchase.ID,
	).Error; err != nil {

		t.Fatal(err)
	}

	if updated.Status != enum.PurchaseStatusPaid {

		t.Fatalf(
			"expected PAID got %s",
			updated.Status,
		)
	}
}

func TestInvalidPurchaseIDReturnsHTTP400(t *testing.T) {

	router := gin.New()

	paymentHandler := handlers.NewPaymentHandler(
		nil,
	)

	router.POST(
		"/api/purchases/:id/refund",
		paymentHandler.Refund,
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/purchases/not-a-valid-uuid/refund",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		req,
	)

	if recorder.Code != http.StatusBadRequest {

		t.Fatalf(
			"expected 400 got %d body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}
}

func TestPurchaseNotFoundReturnsHTTP404(t *testing.T) {

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	authService := helpers.NewAuthService(
		db,
	)

	paymentService := helpers.NewPaymentService(
		database.NewGormExecutor(db),
		helpers.NewPurchaseService(db),
		helpers.NewTicketService(db),
		&helpers.FakePaymentGateway{},
	)

	router := gin.New()

	routes.Register(
		router,
		authService,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		paymentService,
		nil,
	)

	user := fixtures.User()

	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	token := helpers.CreateAuthToken(
		user.ID,
	)

	if err != nil {
		t.Fatal(err)
	}

	missingPurchaseID := uuid.New()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/purchases/"+missingPurchaseID.String()+"/refund",
		nil,
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(
		recorder,
		req,
	)

	if recorder.Code != http.StatusNotFound {

		t.Fatalf(
			"expected 404 got %d body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}
}
