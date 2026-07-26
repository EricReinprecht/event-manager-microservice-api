package paypal

type refundRequest struct {
	Amount refundAmount `json:"amount"`
}

type refundAmount struct {
	Value        string `json:"value"`
	CurrencyCode string `json:"currency_code"`
}

type refundResponse struct {
	ID string `json:"id"`
}
