package refund

import (
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"

	"github.com/reinp/event-platform/backend/internal/handlers"
	"github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/routes"
	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
	testhttp "github.com/reinp/event-platform/backend/tests/helpers/http"
	"github.com/reinp/event-platform/backend/tests/scenarios"
)

func TestOrganizerReceivesHTTP200WhenRefunding(t *testing.T) {

	gin.SetMode(gin.TestMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	helpers.CleanDatabase(db)

	paymentService, fakeGateway :=
		helpers.SetupPaymentService(db)

	handler := handlers.NewPaymentHandler(
		paymentService,
	)

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleOrganizer,
	)

	router := testhttp.RefundRouter(
		handler,
		scenario.Organizer.ID,
	)

	recorder := testhttp.ExecuteRefundRequest(
		router,
		scenario.Purchase.ID,
		"",
	)

	if recorder.Code != http.StatusOK {

		t.Fatalf(
			"expected 200 got %d body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	if !fakeGateway.RefundCalled {

		t.Fatal(
			"expected refund gateway call",
		)
	}
}

func TestAdminReceivesHTTP200WhenRefunding(t *testing.T) {

	gin.SetMode(gin.TestMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	paymentService, fakeGateway :=
		helpers.SetupPaymentService(db)

	handler := handlers.NewPaymentHandler(
		paymentService,
	)

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleAdmin,
	)

	router := testhttp.RefundRouter(
		handler,
		scenario.Actor.ID,
	)

	recorder := testhttp.ExecuteRefundRequest(
		router,
		scenario.Purchase.ID,
		"",
	)

	if recorder.Code != http.StatusOK {

		t.Fatalf(
			"expected 200 got %d body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	if !fakeGateway.RefundCalled {

		t.Fatal(
			"expected refund gateway to be called",
		)
	}

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

	gin.SetMode(gin.TestMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	paymentService, fakeGateway :=
		helpers.SetupPaymentService(db)

	handler := handlers.NewPaymentHandler(
		paymentService,
	)

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleRefunder,
	)

	router := testhttp.RefundRouter(
		handler,
		scenario.Actor.ID,
	)

	recorder := testhttp.ExecuteRefundRequest(
		router,
		scenario.Purchase.ID,
		"",
	)

	if recorder.Code != http.StatusOK {

		t.Fatalf(
			"expected 200 got %d body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	if !fakeGateway.RefundCalled {

		t.Fatal(
			"expected refund gateway to be called",
		)
	}

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

	gin.SetMode(gin.TestMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	paymentService, fakeGateway :=
		helpers.SetupPaymentService(db)

	handler := handlers.NewPaymentHandler(
		paymentService,
	)

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleStaff,
	)

	router := testhttp.RefundRouter(
		handler,
		scenario.Actor.ID,
	)

	recorder := testhttp.ExecuteRefundRequest(
		router,
		scenario.Purchase.ID,
		"",
	)

	if recorder.Code != http.StatusForbidden {

		t.Fatalf(
			"expected 403 got %d body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	if fakeGateway.RefundCalled {

		t.Fatal(
			"staff should not reach refund gateway",
		)
	}

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

	gin.SetMode(gin.TestMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	paymentService, fakeGateway :=
		helpers.SetupPaymentService(db)

	handler := handlers.NewPaymentHandler(
		paymentService,
	)

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleAttendee,
	)

	router := testhttp.RefundRouter(
		handler,
		scenario.Actor.ID,
	)

	recorder := testhttp.ExecuteRefundRequest(
		router,
		scenario.Purchase.ID,
		"",
	)

	if recorder.Code != http.StatusForbidden {

		t.Fatalf(
			"expected 403 got %d body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	if fakeGateway.RefundCalled {

		t.Fatal(
			"attendee should not reach refund gateway",
		)
	}

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

	gin.SetMode(gin.TestMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	paymentService, fakeGateway :=
		helpers.SetupPaymentService(db)

	handler := handlers.NewPaymentHandler(
		paymentService,
	)

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleAttendee,
	)

	outsider := fixtures.User()

	if err := db.Create(
		&outsider,
	).Error; err != nil {

		t.Fatal(err)
	}

	router := testhttp.RefundRouter(
		handler,
		outsider.ID,
	)

	recorder := testhttp.ExecuteRefundRequest(
		router,
		scenario.Purchase.ID,
		"",
	)

	if recorder.Code != http.StatusForbidden {

		t.Fatalf(
			"expected 403 got %d body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	if fakeGateway.RefundCalled {

		t.Fatal(
			"non-member should not reach refund gateway",
		)
	}

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

	gin.SetMode(gin.TestMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

	paymentService, fakeGateway :=
		helpers.SetupPaymentService(db)

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
		handler.Refund,
	)

	recorder := testhttp.ExecuteRefundRequest(
		router,
		scenario.Purchase.ID,
		"",
	)

	if recorder.Code != http.StatusUnauthorized {

		t.Fatalf(
			"expected 401 got %d body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	if fakeGateway.RefundCalled {

		t.Fatal(
			"unauthenticated user should not reach refund gateway",
		)
	}

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
func TestPurchaseNotFoundReturnsHTTP404(t *testing.T) {

	gin.SetMode(gin.TestMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	helpers.CleanDatabase(db)

	authService := helpers.NewAuthService(
		db,
	)

	paymentService, _ :=
		helpers.SetupPaymentService(db)

	router := gin.New()

	// TODO
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

	if err := db.Create(
		&user,
	).Error; err != nil {

		t.Fatal(err)
	}

	token := helpers.CreateAuthToken(
		user.ID,
	)

	missingPurchaseID := uuid.New()

	recorder := testhttp.ExecuteRefundRequest(
		router,
		missingPurchaseID,
		token,
	)

	if recorder.Code != http.StatusNotFound {

		t.Fatalf(
			"expected 404 got %d body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}
}

func TestAlreadyRefundedPurchaseReturnsHTTP400(t *testing.T) {

	gin.SetMode(gin.TestMode)

	db, err := helpers.TestDatabaseSilent()

	if err != nil {
		t.Fatal(err)
	}

	helpers.CleanDatabase(db)

	paymentService, fakeGateway :=
		helpers.SetupPaymentService(db)

	handler := handlers.NewPaymentHandler(
		paymentService,
	)

	scenario := scenarios.CreateRefundScenario(
		t,
		db,
		enum.RoleOrganizer,
	)

	// mark purchase already refunded
	scenario.Purchase.Status =
		enum.PurchaseStatusRefunded

	if err := db.Save(
		scenario.Purchase,
	).Error; err != nil {

		t.Fatal(err)
	}

	router := testhttp.RefundRouter(
		handler,
		scenario.Actor.ID,
	)

	recorder := testhttp.ExecuteRefundRequest(
		router,
		scenario.Purchase.ID,
		"",
	)

	if recorder.Code != http.StatusBadRequest {

		t.Fatalf(
			"expected 400 got %d body=%s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	if !strings.Contains(
		recorder.Body.String(),
		"already refunded",
	) {

		t.Fatalf(
			"expected already refunded error got %s",
			recorder.Body.String(),
		)
	}

	if fakeGateway.RefundCalled {

		t.Fatal(
			"gateway should not be called twice",
		)
	}
}
