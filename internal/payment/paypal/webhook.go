package paypal

type WebhookHeaders struct {
	TransmissionID string

	TransmissionTime string

	TransmissionSig string

	CertURL string

	AuthAlgo string
}
