package enum

type RefundRequestStatus string

const (
	RefundRequestStatusPending RefundRequestStatus = "PENDING"

	RefundRequestStatusApproved RefundRequestStatus = "APPROVED"

	RefundRequestStatusRejected RefundRequestStatus = "REJECTED"

	RefundRequestStatusCompleted RefundRequestStatus = "COMPLETED"
)
