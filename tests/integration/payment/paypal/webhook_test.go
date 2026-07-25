package paypal

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/reinp/event-platform/backend/internal/database"
	"github.com/reinp/event-platform/backend/internal/handlers"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/payment/paypal"
	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

func TestPayPalWebhookCheckoutOrderApprovedTriggersCapture(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
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

	eventID := "WH-CHECKOUT-APPROVED-" + uuid.New().String()

	orderID := "ORDER-TEST-" + uuid.New().String()

	body := fmt.Appendf(
		nil,
		`
	{
		"id": "%s",

		"event_type": "CHECKOUT.ORDER.APPROVED",

		"resource": {
			"id": "%s"
		}
	}
	`,
		eventID,
		orderID,
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/paypal/webhook",
		bytes.NewReader(body),
	)

	req.Header.Set("PAYPAL-TRANSMISSION-ID", "test-id")
	req.Header.Set("PAYPAL-TRANSMISSION-TIME", "2026-01-01T00:00:00Z")
	req.Header.Set("PAYPAL-TRANSMISSION-SIG", "test-signature")
	req.Header.Set("PAYPAL-CERT-URL", "https://example.com/cert")
	req.Header.Set("PAYPAL-AUTH-ALGO", "SHA256withRSA")

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)

	c.Request = req

	handler.Webhook(c)

	if w.Code != http.StatusOK {

		t.Fatalf(
			"expected status 200, got %d body: %s",
			w.Code,
			w.Body.String(),
		)
	}

	if !fakeGateway.CaptureOrderCalled {

		t.Fatal(
			"expected capture order to be called",
		)
	}

	if fakeGateway.CapturedOrderID != orderID {

		t.Fatalf(
			"expected %s, got %s",
			orderID,
			fakeGateway.CapturedOrderID,
		)
	}
}

func TestPayPalWebhookCheckoutOrderApprovedInvalidOrderID(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	failingGateway := &helpers.FailingPaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		failingGateway,
	)

	handler := handlers.NewPaymentHandler(
		paymentService,
	)

	body := []byte(`
	{
		"id": "WH-INVALID-ORDER-TEST-001",

		"event_type": "CHECKOUT.ORDER.APPROVED",

		"resource": {
			"id": "INVALID_ORDER_ID"
		}
	}
	`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/paypal/webhook",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-ID",
		"test-id",
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-TIME",
		"2026-01-01T00:00:00Z",
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-SIG",
		"test-signature",
	)

	req.Header.Set(
		"PAYPAL-CERT-URL",
		"https://example.com/cert",
	)

	req.Header.Set(
		"PAYPAL-AUTH-ALGO",
		"SHA256withRSA",
	)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)

	c.Request = req

	handler.Webhook(c)

	if w.Code != http.StatusInternalServerError {

		t.Fatalf(
			"expected status 500, got %d body: %s",
			w.Code,
			w.Body.String(),
		)
	}
}

func TestPayPalWebhookCaptureCompletedConfirmsPayment(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
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

	user := fixtures.User()

	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	party := fixtures.Party()

	party.CategoryID = category.ID
	party.OrganizerID = user.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.StatusPending,
	)

	paymentID := "PAYPAL-ORDER-WEBHOOK-123"

	purchase.PaymentProvider = "paypal"

	purchase.PaymentID = paymentID

	if err := db.Save(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	ticketCategory := fixtures.TicketCategory(
		party.ID,
	)

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	item := models.PurchaseItem{

		ID: uuid.New(),

		PurchaseID: purchase.ID,

		TicketCategoryID: ticketCategory.ID,

		Quantity: 1,

		UnitPrice: ticketCategory.Price,
	}

	if err := db.Create(&item).Error; err != nil {
		t.Fatal(err)
	}

	eventID := "WH-CAPTURE-COMPLETED-" + uuid.New().String()

	body := []byte(`
	{
		"id": "` + eventID + `",

		"event_type": "PAYMENT.CAPTURE.COMPLETED",

		"resource": {

			"id": "CAPTURE-123",

			"supplementary_data": {

				"related_ids": {

					"order_id": "` + paymentID + `"

				}
			}
		}
	}
	`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/paypal/webhook",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-ID",
		"test-id",
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-TIME",
		"2026-01-01T00:00:00Z",
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-SIG",
		"test-signature",
	)

	req.Header.Set(
		"PAYPAL-CERT-URL",
		"https://example.com/cert",
	)

	req.Header.Set(
		"PAYPAL-AUTH-ALGO",
		"SHA256withRSA",
	)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)

	c.Request = req

	handler.Webhook(c)

	if w.Code != http.StatusOK {

		t.Fatalf(
			"expected status 200, got %d body: %s",
			w.Code,
			w.Body.String(),
		)
	}

	if fakeGateway.CaptureOrderCalled {

		t.Fatal(
			"capture should not be called for completed payment",
		)
	}

	var updatedPurchase models.Purchase

	err = db.
		First(
			&updatedPurchase,
			"id = ?",
			purchase.ID,
		).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if updatedPurchase.Status != enum.StatusPaid {

		t.Fatalf(
			"expected purchase status PAID, got %s",
			updatedPurchase.Status,
		)
	}

	var tickets []models.Ticket

	err = db.
		Where(
			"user_id = ?",
			user.ID,
		).
		Find(&tickets).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if len(tickets) != 1 {

		t.Fatalf(
			"expected 1 generated ticket, got %d",
			len(tickets),
		)
	}
}

func TestPayPalWebhookCaptureCompletedGeneratesTickets(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
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

	user := fixtures.User()

	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	party := fixtures.Party()

	party.CategoryID = category.ID
	party.OrganizerID = user.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.StatusPending,
	)

	paymentID := "PAYPAL-ORDER-TICKET-123"

	purchase.PaymentProvider = "paypal"
	purchase.PaymentID = paymentID

	if err := db.Save(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	ticketCategory := fixtures.TicketCategory(
		party.ID,
	)

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	item := models.PurchaseItem{

		ID: uuid.New(),

		PurchaseID: purchase.ID,

		TicketCategoryID: ticketCategory.ID,

		Quantity: 1,

		UnitPrice: ticketCategory.Price,
	}

	if err := db.Create(&item).Error; err != nil {
		t.Fatal(err)
	}

	eventID := "WH-TICKET-GENERATION-123"

	body := fmt.Appendf(
		nil,
		`
	{
		"id": "%s",

		"event_type": "PAYMENT.CAPTURE.COMPLETED",

		"resource": {

			"id": "CAPTURE-123",

			"supplementary_data": {

				"related_ids": {

					"order_id": "%s"

				}
			}
		}
	}
	`,
		eventID,
		paymentID,
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/paypal/webhook",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-ID",
		"test-id",
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-TIME",
		"2026-01-01T00:00:00Z",
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-SIG",
		"test-signature",
	)

	req.Header.Set(
		"PAYPAL-CERT-URL",
		"https://example.com/cert",
	)

	req.Header.Set(
		"PAYPAL-AUTH-ALGO",
		"SHA256withRSA",
	)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)

	c.Request = req

	handler.Webhook(c)

	if w.Code != http.StatusOK {

		t.Fatalf(
			"expected status 200 got %d body: %s",
			w.Code,
			w.Body.String(),
		)
	}

	var tickets []models.Ticket

	err = db.
		Where(
			"user_id = ?",
			user.ID,
		).
		Find(&tickets).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if len(tickets) != 1 {

		t.Fatalf(
			"expected 1 generated ticket, got %d",
			len(tickets),
		)
	}

	if fakeGateway.CaptureOrderCalled {

		t.Fatal(
			"capture should not be called for completed payment",
		)
	}
}

func TestPayPalWebhookCaptureCompletedMarksEventProcessed(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
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

	user := fixtures.User()

	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	party := fixtures.Party()

	party.CategoryID = category.ID
	party.OrganizerID = user.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.StatusPending,
	)

	paymentID := "PAYPAL-ORDER-EVENT-PROCESSED-123"

	purchase.PaymentProvider = "paypal"
	purchase.PaymentID = paymentID

	if err := db.Save(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	ticketCategory := fixtures.TicketCategory(
		party.ID,
	)

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	item := models.PurchaseItem{

		ID: uuid.New(),

		PurchaseID: purchase.ID,

		TicketCategoryID: ticketCategory.ID,

		Quantity: 1,

		UnitPrice: ticketCategory.Price,
	}

	if err := db.Create(&item).Error; err != nil {
		t.Fatal(err)
	}

	eventID := "WH-MARK-PROCESSED-123"

	body := fmt.Appendf(
		nil,
		`
	{
		"id": "%s",

		"event_type": "PAYMENT.CAPTURE.COMPLETED",

		"resource": {

			"id": "CAPTURE-123",

			"supplementary_data": {

				"related_ids": {

					"order_id": "%s"

				}
			}
		}
	}
	`,
		eventID,
		paymentID,
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/paypal/webhook",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-ID",
		"test-id",
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-TIME",
		"2026-01-01T00:00:00Z",
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-SIG",
		"test-signature",
	)

	req.Header.Set(
		"PAYPAL-CERT-URL",
		"https://example.com/cert",
	)

	req.Header.Set(
		"PAYPAL-AUTH-ALGO",
		"SHA256withRSA",
	)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)

	c.Request = req

	handler.Webhook(c)

	if w.Code != http.StatusOK {

		t.Fatalf(
			"expected status 200 got %d body: %s",
			w.Code,
			w.Body.String(),
		)
	}

	var event models.PaymentEvent

	err = db.
		Where(
			"event_id = ?",
			eventID,
		).
		First(&event).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if !event.Processed {

		t.Fatal(
			"expected payment event to be marked processed",
		)
	}

	if event.Type != "PAYMENT.CAPTURE.COMPLETED" {

		t.Fatalf(
			"expected event type PAYMENT.CAPTURE.COMPLETED got %s",
			event.Type,
		)
	}
}

func TestPayPalWebhookDuplicateEventDoesNotDuplicatePayment(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
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

	user := fixtures.User()

	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	party := fixtures.Party()

	party.CategoryID = category.ID
	party.OrganizerID = user.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.StatusPending,
	)

	paymentID := "PAYPAL-DUPLICATE-ORDER-123"

	purchase.PaymentProvider = "paypal"
	purchase.PaymentID = paymentID

	if err := db.Save(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	ticketCategory := fixtures.TicketCategory(
		party.ID,
	)

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	item := models.PurchaseItem{

		ID: uuid.New(),

		PurchaseID: purchase.ID,

		TicketCategoryID: ticketCategory.ID,

		Quantity: 1,

		UnitPrice: ticketCategory.Price,
	}

	if err := db.Create(&item).Error; err != nil {
		t.Fatal(err)
	}

	eventID := "WH-DUPLICATE-EVENT-123"

	body := fmt.Appendf(
		nil,
		`
	{
		"id": "%s",

		"event_type": "PAYMENT.CAPTURE.COMPLETED",

		"resource": {

			"id": "CAPTURE-123",

			"supplementary_data": {

				"related_ids": {

					"order_id": "%s"

				}
			}
		}
	}
	`,
		eventID,
		paymentID,
	)

	sendWebhook := func() *httptest.ResponseRecorder {

		req := httptest.NewRequest(
			http.MethodPost,
			"/paypal/webhook",
			bytes.NewReader(body),
		)

		req.Header.Set(
			"PAYPAL-TRANSMISSION-ID",
			"test-id",
		)

		req.Header.Set(
			"PAYPAL-TRANSMISSION-TIME",
			"2026-01-01T00:00:00Z",
		)

		req.Header.Set(
			"PAYPAL-TRANSMISSION-SIG",
			"test-signature",
		)

		req.Header.Set(
			"PAYPAL-CERT-URL",
			"https://example.com/cert",
		)

		req.Header.Set(
			"PAYPAL-AUTH-ALGO",
			"SHA256withRSA",
		)

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = req

		handler.Webhook(c)

		return w
	}

	first := sendWebhook()

	if first.Code != http.StatusOK {

		t.Fatalf(
			"first webhook expected 200 got %d body %s",
			first.Code,
			first.Body.String(),
		)
	}

	second := sendWebhook()

	if second.Code != http.StatusOK {

		t.Fatalf(
			"duplicate webhook expected 200 got %d body %s",
			second.Code,
			second.Body.String(),
		)
	}

	var events []models.PaymentEvent

	err = db.
		Where(
			"event_id = ?",
			eventID,
		).
		Find(&events).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if len(events) != 1 {

		t.Fatalf(
			"expected 1 payment event, got %d",
			len(events),
		)
	}

	var tickets []models.Ticket

	err = db.
		Where(
			"user_id = ?",
			user.ID,
		).
		Find(&tickets).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if len(tickets) != 1 {

		t.Fatalf(
			"expected 1 ticket, got %d",
			len(tickets),
		)
	}
}

func TestPayPalWebhookUnprocessedEventRetrySucceeds(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	user := fixtures.User()

	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	party := fixtures.Party()

	party.CategoryID = category.ID
	party.OrganizerID = user.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	paymentID := "PAYPAL-RETRY-ORDER-" + uuid.New().String()

	eventID := "WH-RETRY-EVENT-" + uuid.New().String()

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.StatusPending,
	)

	purchase.PaymentProvider = "paypal"
	purchase.PaymentID = paymentID

	if err := db.Save(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	ticketCategory := fixtures.TicketCategory(
		party.ID,
	)

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	item := models.PurchaseItem{

		ID: uuid.New(),

		PurchaseID: purchase.ID,

		TicketCategoryID: ticketCategory.ID,

		Quantity: 1,

		UnitPrice: ticketCategory.Price,
	}

	if err := db.Create(&item).Error; err != nil {
		t.Fatal(err)
	}

	createRequest := func(body []byte) *http.Request {

		req := httptest.NewRequest(
			http.MethodPost,
			"/paypal/webhook",
			bytes.NewReader(body),
		)

		req.Header.Set(
			"PAYPAL-TRANSMISSION-ID",
			"test-id",
		)

		req.Header.Set(
			"PAYPAL-TRANSMISSION-TIME",
			"2026-01-01T00:00:00Z",
		)

		req.Header.Set(
			"PAYPAL-TRANSMISSION-SIG",
			"test-signature",
		)

		req.Header.Set(
			"PAYPAL-CERT-URL",
			"https://example.com/cert",
		)

		req.Header.Set(
			"PAYPAL-AUTH-ALGO",
			"SHA256withRSA",
		)

		return req
	}

	// ==================================================
	// FIRST ATTEMPT
	// invalid order id -> ConfirmPayment fails
	// ==================================================

	firstBody := fmt.Appendf(
		nil,
		`
{
	"id": "%s",

	"event_type": "PAYMENT.CAPTURE.COMPLETED",

	"resource": {

		"id": "CAPTURE-123",

		"supplementary_data": {

			"related_ids": {

				"order_id": "INVALID-ORDER-ID"

			}
		}
	}
}
`,
		eventID,
	)

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		&helpers.FakePaymentGateway{},
	)

	handler := handlers.NewPaymentHandler(
		paymentService,
	)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)

	c.Request = createRequest(firstBody)

	handler.Webhook(c)

	if w.Code != http.StatusOK {

		t.Fatalf(
			"expected 200 because unknown order is ignored, got %d body %s",
			w.Code,
			w.Body.String(),
		)
	}

	var event models.PaymentEvent

	err = db.
		Where(
			"event_id = ?",
			eventID,
		).
		First(&event).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if event.Processed {

		t.Fatal(
			"expected event to be unprocessed after failure",
		)
	}

	// ==================================================
	// SECOND ATTEMPT
	// same event id, correct order id -> succeeds
	// ==================================================

	secondBody := fmt.Appendf(
		nil,
		`
{
	"id": "%s",

	"event_type": "PAYMENT.CAPTURE.COMPLETED",

	"resource": {

		"id": "CAPTURE-123",

		"supplementary_data": {

			"related_ids": {

				"order_id": "%s"

			}
		}
	}
}
`,
		eventID,
		paymentID,
	)

	w = httptest.NewRecorder()

	c, _ = gin.CreateTestContext(w)

	c.Request = createRequest(secondBody)

	handler.Webhook(c)

	if w.Code != http.StatusOK {

		t.Fatalf(
			"expected retry 200, got %d body %s",
			w.Code,
			w.Body.String(),
		)
	}

	err = db.
		Where(
			"event_id = ?",
			eventID,
		).
		First(&event).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if !event.Processed {

		t.Fatal(
			"expected event processed after retry",
		)
	}

	var updatedPurchase models.Purchase

	err = db.
		First(
			&updatedPurchase,
			"id = ?",
			purchase.ID,
		).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if updatedPurchase.Status != enum.StatusPaid {

		t.Fatalf(
			"expected PAID after retry, got %s",
			updatedPurchase.Status,
		)
	}
}

func TestPaymentEvent_CannotDuplicateWebhook(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
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

	// -----------------------------
	// Create purchase
	// -----------------------------

	user := fixtures.User()

	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	party := fixtures.Party()

	party.CategoryID = category.ID
	party.OrganizerID = user.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.StatusPending,
	)

	paymentID := "PAYPAL-DUPLICATE-" + uuid.New().String()

	purchase.PaymentProvider = "paypal"
	purchase.PaymentID = paymentID

	if err := db.Save(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	ticketCategory := fixtures.TicketCategory(
		party.ID,
	)

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	item := models.PurchaseItem{

		ID: uuid.New(),

		PurchaseID: purchase.ID,

		TicketCategoryID: ticketCategory.ID,

		Quantity: 1,

		UnitPrice: ticketCategory.Price,
	}

	if err := db.Create(&item).Error; err != nil {
		t.Fatal(err)
	}

	// -----------------------------
	// Webhook payload
	// -----------------------------

	eventID := "WH-DUPLICATE-" + uuid.New().String()

	body := fmt.Appendf(
		nil,
		`
{
	"id": "%s",

	"event_type": "PAYMENT.CAPTURE.COMPLETED",

	"resource": {

		"id": "CAPTURE-123",

		"supplementary_data": {

			"related_ids": {

				"order_id": "%s"

			}
		}
	}
}
`,
		eventID,
		paymentID,
	)

	createRequest := func() *http.Request {

		req := httptest.NewRequest(
			http.MethodPost,
			"/paypal/webhook",
			bytes.NewReader(body),
		)

		req.Header.Set(
			"PAYPAL-TRANSMISSION-ID",
			"test-id",
		)

		req.Header.Set(
			"PAYPAL-TRANSMISSION-TIME",
			"2026-01-01T00:00:00Z",
		)

		req.Header.Set(
			"PAYPAL-TRANSMISSION-SIG",
			"test-signature",
		)

		req.Header.Set(
			"PAYPAL-CERT-URL",
			"https://example.com/cert",
		)

		req.Header.Set(
			"PAYPAL-AUTH-ALGO",
			"SHA256withRSA",
		)

		return req
	}

	// -----------------------------
	// First webhook
	// -----------------------------

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)

	c.Request = createRequest()

	handler.Webhook(c)

	if w.Code != http.StatusOK {

		t.Fatalf(
			"first webhook failed: %d body %s",
			w.Code,
			w.Body.String(),
		)
	}

	// verify event stored once

	var count int64

	err = db.
		Model(&models.PaymentEvent{}).
		Where(
			"event_id = ?",
			eventID,
		).
		Count(&count).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {

		t.Fatalf(
			"expected 1 payment event, got %d",
			count,
		)
	}

	// -----------------------------
	// Second identical webhook
	// -----------------------------

	w = httptest.NewRecorder()

	c, _ = gin.CreateTestContext(w)

	c.Request = createRequest()

	handler.Webhook(c)

	if w.Code != http.StatusOK {

		t.Fatalf(
			"duplicate webhook failed: %d body %s",
			w.Code,
			w.Body.String(),
		)
	}

	// still only one event

	count = 0

	err = db.
		Model(&models.PaymentEvent{}).
		Where(
			"event_id = ?",
			eventID,
		).
		Count(&count).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {

		t.Fatalf(
			"expected no duplicate payment events, got %d",
			count,
		)
	}

	// verify purchase was paid only once

	var updatedPurchase models.Purchase

	err = db.
		First(
			&updatedPurchase,
			"id = ?",
			purchase.ID,
		).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if updatedPurchase.Status != enum.StatusPaid {

		t.Fatalf(
			"expected PAID, got %s",
			updatedPurchase.Status,
		)
	}
}

func TestWebhook_RejectsInvalidSignature(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	invalidSignatureGateway := &helpers.InvalidSignaturePaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		invalidSignatureGateway,
	)

	handler := handlers.NewPaymentHandler(
		paymentService,
	)

	body := []byte(`
{
	"id": "WH-INVALID-SIGNATURE-001",

	"event_type": "PAYMENT.CAPTURE.COMPLETED",

	"resource": {

		"id": "CAPTURE-123",

		"supplementary_data": {

			"related_ids": {

				"order_id": "ORDER123"

			}
		}
	}
}
`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/paypal/webhook",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-ID",
		"invalid-id",
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-TIME",
		"2026-01-01T00:00:00Z",
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-SIG",
		"invalid-signature",
	)

	req.Header.Set(
		"PAYPAL-CERT-URL",
		"https://example.com/cert",
	)

	req.Header.Set(
		"PAYPAL-AUTH-ALGO",
		"SHA256withRSA",
	)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)

	c.Request = req

	handler.Webhook(c)

	if w.Code != http.StatusUnauthorized {

		t.Fatalf(
			"expected 401 invalid signature, got %d body: %s",
			w.Code,
			w.Body.String(),
		)
	}

	if !strings.Contains(
		w.Body.String(),
		"invalid paypal signature",
	) {

		t.Fatalf(
			"unexpected response: %s",
			w.Body.String(),
		)
	}
}

func TestWebhook_RejectsMissingTransmissionHeaders(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	invalidSignatureGateway := &helpers.InvalidSignaturePaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		invalidSignatureGateway,
	)

	handler := handlers.NewPaymentHandler(
		paymentService,
	)

	body := []byte(`
{
	"id": "WH-MISSING-HEADERS-001",

	"event_type": "PAYMENT.CAPTURE.COMPLETED",

	"resource": {

		"id": "CAPTURE-123",

		"supplementary_data": {

			"related_ids": {

				"order_id": "ORDER123"

			}
		}
	}
}
`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/paypal/webhook",
		bytes.NewReader(body),
	)

	// Intentionally do NOT set:
	//
	// PAYPAL-TRANSMISSION-ID
	// PAYPAL-TRANSMISSION-TIME
	// PAYPAL-TRANSMISSION-SIG
	// PAYPAL-CERT-URL
	// PAYPAL-AUTH-ALGO

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)

	c.Request = req

	handler.Webhook(c)

	if w.Code != http.StatusUnauthorized {

		t.Fatalf(
			"expected 401 for missing headers, got %d body: %s",
			w.Code,
			w.Body.String(),
		)
	}

	if !strings.Contains(
		w.Body.String(),
		"invalid paypal signature",
	) {

		t.Fatalf(
			"expected invalid signature response, got %s",
			w.Body.String(),
		)
	}
}

func TestWebhook_RejectsModifiedPayload(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	invalidSignatureGateway := &helpers.InvalidSignaturePaymentGateway{}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		invalidSignatureGateway,
	)

	handler := handlers.NewPaymentHandler(
		paymentService,
	)

	// Payload was modified after signing

	body := []byte(`
{
	"id": "WH-MODIFIED-PAYLOAD-001",

	"event_type": "PAYMENT.CAPTURE.COMPLETED",

	"resource": {

		"id": "CAPTURE-123",

		"supplementary_data": {

			"related_ids": {

				"order_id": "ORDER-MODIFIED"

			}
		}
	}
}
`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/paypal/webhook",
		bytes.NewReader(body),
	)

	// These headers simulate PayPal headers,
	// but signature verification will fail because
	// the payload was modified.

	req.Header.Set(
		"PAYPAL-TRANSMISSION-ID",
		"transmission-123",
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-TIME",
		"2026-01-01T00:00:00Z",
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-SIG",
		"old-valid-signature",
	)

	req.Header.Set(
		"PAYPAL-CERT-URL",
		"https://example.com/cert",
	)

	req.Header.Set(
		"PAYPAL-AUTH-ALGO",
		"SHA256withRSA",
	)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)

	c.Request = req

	handler.Webhook(c)

	if w.Code != http.StatusUnauthorized {

		t.Fatalf(
			"expected 401 for modified payload, got %d body: %s",
			w.Code,
			w.Body.String(),
		)
	}

	if !strings.Contains(
		w.Body.String(),
		"invalid paypal signature",
	) {

		t.Fatalf(
			"expected invalid signature error, got %s",
			w.Body.String(),
		)
	}
}

func TestPayPalWebhookReplayProtection(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	err = helpers.CleanDatabase(db)

	if err != nil {
		t.Fatal(err)
	}

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	gateway := &helpers.FakePaymentGateway{
		Order: &paypal.Order{
			ID: "ORDER-REPLAY-001",
		},
	}

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		gateway,
	)

	handler := handlers.NewPaymentHandler(
		paymentService,
	)

	// ----------------------------
	// Create purchase setup
	// ----------------------------

	user := fixtures.User()

	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	category := fixtures.Category()

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	party := fixtures.PartyWithOrganizer(
		user.ID,
	)

	party.CategoryID = category.ID

	if err := db.Create(&party).Error; err != nil {
		t.Fatal(err)
	}

	purchase := helpers.CreatePurchase(
		t,
		db,
		&user,
		&party,
		enum.StatusPending,
	)

	ticketCategory := fixtures.TicketCategory(
		party.ID,
	)

	if err := db.Create(&ticketCategory).Error; err != nil {
		t.Fatal(err)
	}

	item := models.PurchaseItem{

		ID: uuid.New(),

		PurchaseID: purchase.ID,

		TicketCategoryID: ticketCategory.ID,

		Quantity: 1,

		UnitPrice: ticketCategory.Price,
	}

	if err := db.Create(&item).Error; err != nil {
		t.Fatal(err)
	}

	purchase.PaymentID = "ORDER-REPLAY-001"

	purchase.PaymentProvider = "paypal"

	if err := db.Save(&purchase).Error; err != nil {
		t.Fatal(err)
	}

	// ----------------------------
	// Webhook payload
	// ----------------------------

	body := []byte(`
{
	"id": "WH-REPLAY-001",

	"event_type": "PAYMENT.CAPTURE.COMPLETED",

	"resource": {

		"id": "CAPTURE-REPLAY-001",

		"supplementary_data": {

			"related_ids": {

				"order_id": "ORDER-REPLAY-001"

			}
		}
	}
}
`)

	sendWebhook := func() int {

		req := httptest.NewRequest(
			http.MethodPost,
			"/paypal/webhook",
			bytes.NewReader(body),
		)

		req.Header.Set(
			"PAYPAL-TRANSMISSION-ID",
			"replay-transmission-001",
		)

		req.Header.Set(
			"PAYPAL-TRANSMISSION-TIME",
			"2026-01-01T00:00:00Z",
		)

		req.Header.Set(
			"PAYPAL-TRANSMISSION-SIG",
			"valid-signature",
		)

		req.Header.Set(
			"PAYPAL-CERT-URL",
			"https://example.com/cert",
		)

		req.Header.Set(
			"PAYPAL-AUTH-ALGO",
			"SHA256withRSA",
		)

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Request = req

		handler.Webhook(c)

		return w.Code
	}

	// ----------------------------
	// First delivery
	// ----------------------------

	code := sendWebhook()

	if code != http.StatusOK {

		t.Fatalf(
			"first webhook failed with status %d",
			code,
		)
	}

	// ----------------------------
	// Replay same webhook
	// ----------------------------

	code = sendWebhook()

	if code != http.StatusOK {

		t.Fatalf(
			"replayed webhook failed with status %d",
			code,
		)
	}

	// ----------------------------
	// Verify only one event stored
	// ----------------------------

	var events []models.PaymentEvent

	err = db.
		Where(
			"event_id = ?",
			"WH-REPLAY-001",
		).
		Find(&events).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if len(events) != 1 {

		t.Fatalf(
			"expected exactly one payment event, got %d",
			len(events),
		)
	}

	// ----------------------------
	// Verify tickets not duplicated
	// ----------------------------

	var tickets []models.Ticket

	err = db.
		Where(
			"user_id = ?",
			user.ID,
		).
		Find(&tickets).
		Error

	if err != nil {
		t.Fatal(err)
	}

	if len(tickets) != 1 {

		t.Fatalf(
			"expected exactly one ticket, got %d",
			len(tickets),
		)
	}
}

func TestPayPalWebhookUnknownOrderDoesNotCreatePayment(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	// --------------------------------
	// Setup payment service
	// --------------------------------

	executor := database.NewGormExecutor(db)

	purchaseService := helpers.NewPurchaseService(db)

	ticketService := helpers.NewTicketService(db)

	paymentService := helpers.NewPaymentService(
		executor,
		purchaseService,
		ticketService,
		&helpers.FakePaymentGateway{
			Order: &paypal.Order{
				ID: "ORDER-DOES-NOT-EXIST",
			},
		},
	)

	handler := handlers.NewPaymentHandler(
		paymentService,
	)

	// --------------------------------
	// Unknown PayPal order payload
	// --------------------------------

	body := []byte(`
{
	"id": "WH-UNKNOWN-ORDER-001",

	"event_type": "PAYMENT.CAPTURE.COMPLETED",

	"resource": {

		"id": "CAPTURE-UNKNOWN",

		"supplementary_data": {

			"related_ids": {

				"order_id": "ORDER-DOES-NOT-EXIST"

			}
		}
	}
}
`)

	req := httptest.NewRequest(
		http.MethodPost,
		"/paypal/webhook",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-ID",
		"transmission-unknown-order",
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-TIME",
		"2026-01-01T00:00:00Z",
	)

	req.Header.Set(
		"PAYPAL-TRANSMISSION-SIG",
		"valid-signature",
	)

	req.Header.Set(
		"PAYPAL-CERT-URL",
		"https://example.com/cert",
	)

	req.Header.Set(
		"PAYPAL-AUTH-ALGO",
		"SHA256withRSA",
	)

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)

	c.Request = req

	// --------------------------------
	// Execute webhook
	// --------------------------------

	handler.Webhook(c)

	// --------------------------------
	// Assertions
	// --------------------------------

	if w.Code != http.StatusOK {

		t.Fatalf(
			"expected 200 for unknown order webhook, got %d body: %s",
			w.Code,
			w.Body.String(),
		)
	}

	// --------------------------------
	// No purchase should exist
	// --------------------------------

	var purchases []models.Purchase

	if err := db.Find(
		&purchases,
	).Error; err != nil {

		t.Fatal(err)
	}

	if len(purchases) != 0 {

		t.Fatalf(
			"expected no purchases, got %d",
			len(purchases),
		)
	}

	// --------------------------------
	// No tickets generated
	// --------------------------------

	var tickets []models.Ticket

	if err := db.Find(
		&tickets,
	).Error; err != nil {

		t.Fatal(err)
	}

	if len(tickets) != 0 {

		t.Fatalf(
			"expected no tickets, got %d",
			len(tickets),
		)
	}

	// --------------------------------
	// Event should still be stored
	// --------------------------------
	// PayPal replay protection requires unknown events
	// to be recorded. We only ignore the business action.

	var events []models.PaymentEvent

	if err := db.
		Where(
			"event_id = ?",
			"WH-UNKNOWN-ORDER-001",
		).
		Find(&events).
		Error; err != nil {

		t.Fatal(err)
	}

	if len(events) != 1 {

		t.Fatalf(
			"expected webhook event to be stored once, got %d",
			len(events),
		)
	}
}
