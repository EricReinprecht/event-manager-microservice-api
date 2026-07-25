package enum

type PurchaseStatus string

const (
	PruchaseStatusPending PurchaseStatus = "PENDING"

	PurchaseStatusPaid PurchaseStatus = "PAID"

	PurchaseStatusFailed PurchaseStatus = "FAILED"

	PurchaseStatusCanceled PurchaseStatus = "CANCELED"

	PurchaseStatusRefunded PurchaseStatus = "REFUNDED"
)
