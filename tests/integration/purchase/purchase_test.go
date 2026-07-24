package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/reinp/event-platform/backend/internal/appErrors"
	appModels "github.com/reinp/event-platform/backend/internal/models"
	"github.com/reinp/event-platform/backend/internal/models/enum"
	"github.com/reinp/event-platform/backend/internal/requests"

	"github.com/reinp/event-platform/backend/tests/fixtures"
	"github.com/reinp/event-platform/backend/tests/helpers"
)

func setupPurchaseTest(t *testing.T) (
	*appModels.User,
	*appModels.Party,
) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	if err := helpers.CleanDatabase(db); err != nil {
		t.Fatal(err)
	}

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

	return &user, &party
}

func TestUserCanCreatePendingPurchase(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	helpers.CleanDatabase(db)

	user, party := setupPurchaseTest(t)

	service := helpers.NewPurchaseService(db)

	category := appModels.TicketCategory{

		ID: uuid.New(),

		Name: "VIP",

		PartyID: party.ID,

		Price: 5000,
	}

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	purchase, err := service.CreatePurchase(
		context.Background(),
		user.ID,
		party.ID,
		[]requests.PurchaseItemRequest{
			{
				TicketCategoryID: category.ID,
				Quantity:         1,
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if purchase.Status != enum.StatusPending {

		t.Fatalf(
			"expected pending purchase got %s",
			purchase.Status,
		)
	}

	if len(purchase.Items) != 1 {

		t.Fatal(
			"expected one purchase item",
		)
	}
}

func TestPurchaseStoresPriceSnapshot(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	helpers.CleanDatabase(db)

	user, party := setupPurchaseTest(t)

	service := helpers.NewPurchaseService(db)

	category := appModels.TicketCategory{

		ID: uuid.New(),

		Name: "VIP",

		PartyID: party.ID,

		Price: 5000,
	}

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	purchase, err := service.CreatePurchase(
		context.Background(),
		user.ID,
		party.ID,
		[]requests.PurchaseItemRequest{
			{
				TicketCategoryID: category.ID,
				Quantity:         2,
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if purchase.Items[0].UnitPrice != 5000 {

		t.Fatalf(
			"expected snapshot price 5000 got %d",
			purchase.Items[0].UnitPrice,
		)
	}
}

func TestPurchaseRejectsInvalidCategory(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	helpers.CleanDatabase(db)

	user, party := setupPurchaseTest(t)

	service := helpers.NewPurchaseService(db)

	_, err = service.CreatePurchase(
		context.Background(),
		user.ID,
		party.ID,
		[]requests.PurchaseItemRequest{
			{
				TicketCategoryID: uuid.New(),
				Quantity:         1,
			},
		},
	)

	if err == nil {

		t.Fatal(
			"expected invalid category error",
		)
	}
}

func TestPurchaseCalculatesTotal(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	helpers.CleanDatabase(db)

	user, party := setupPurchaseTest(t)

	service := helpers.NewPurchaseService(db)

	category := appModels.TicketCategory{

		ID: uuid.New(),

		Name: "VIP",

		PartyID: party.ID,

		Price: 2500,
	}

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	purchase, err := service.CreatePurchase(
		context.Background(),
		user.ID,
		party.ID,
		[]requests.PurchaseItemRequest{
			{
				TicketCategoryID: category.ID,
				Quantity:         4,
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	expected := int64(10000)

	if purchase.TotalPrice != expected {

		t.Fatalf(
			"expected total %d got %d",
			expected,
			purchase.TotalPrice,
		)
	}
}

func TestCannotPurchaseTicketCategoryFromAnotherParty(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	helpers.CleanDatabase(db)

	user, party := setupPurchaseTest(t)

	otherCategory := fixtures.Category()
	otherCategory.Name = "Concert"

	if err := db.Create(&otherCategory).Error; err != nil {
		t.Fatal(err)
	}

	otherParty := fixtures.PartyWithOrganizer(
		user.ID,
	)

	otherParty.CategoryID = otherCategory.ID

	if err := db.Create(&otherParty).Error; err != nil {
		t.Fatal(err)
	}

	category := appModels.TicketCategory{

		ID: uuid.New(),

		Name: "VIP",

		PartyID: otherParty.ID,

		Price: 5000,
	}

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	service := helpers.NewPurchaseService(db)

	_, err = service.CreatePurchase(
		context.Background(),
		user.ID,
		party.ID,
		[]requests.PurchaseItemRequest{
			{
				TicketCategoryID: category.ID,
				Quantity:         1,
			},
		},
	)

	if err == nil {

		t.Fatal(
			"expected category from another party to fail",
		)
	}

	if err != appErrors.ErrTicketCategoryNotFound {

		t.Fatalf(
			"expected ErrTicketCategoryNotFound got %v",
			err,
		)
	}
}

func TestPurchaseCanBeMarkedPaid(t *testing.T) {

	db, err := helpers.TestDatabase()

	if err != nil {
		t.Fatal(err)
	}

	helpers.CleanDatabase(db)

	user, party := setupPurchaseTest(t)

	service := helpers.NewPurchaseService(db)

	category := appModels.TicketCategory{

		ID: uuid.New(),

		Name: "VIP",

		PartyID: party.ID,

		Price: 5000,
	}

	if err := db.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	purchase, err := service.CreatePurchase(
		context.Background(),
		user.ID,
		party.ID,
		[]requests.PurchaseItemRequest{
			{
				TicketCategoryID: category.ID,
				Quantity:         1,
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	paymentID := "PAYPAL-ORDER-123"

	err = service.AttachPayment(
		context.Background(),
		purchase.ID,
		"paypal",
		paymentID,
	)

	if err != nil {
		t.Fatal(err)
	}

	updatedPurchase, err := service.ConfirmPayment(
		context.Background(),
		paymentID,
	)

	if err != nil {
		t.Fatal(err)
	}

	if updatedPurchase.Status != enum.StatusPaid {

		t.Fatalf(
			"expected purchase status PAID got %s",
			updatedPurchase.Status,
		)
	}

	if updatedPurchase.PaymentProvider != "paypal" {

		t.Fatalf(
			"expected payment provider paypal got %s",
			updatedPurchase.PaymentProvider,
		)
	}

	if updatedPurchase.PaymentID != paymentID {

		t.Fatalf(
			"expected payment id %s got %s",
			paymentID,
			updatedPurchase.PaymentID,
		)
	}
}
