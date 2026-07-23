package enum

type PurchaseStatus string

const (
	StatusPending  PurchaseStatus = "PENDING"
	StatusPaid     PurchaseStatus = "PAID"
	StatusFailed   PurchaseStatus = "FAILED"
	StatusCanceled PurchaseStatus = "CANCELED"
	StatusRefunded PurchaseStatus = "REFUNDED"
)
