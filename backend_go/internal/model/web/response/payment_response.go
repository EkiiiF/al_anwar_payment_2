package response

type TransactionResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}

type MidtransStatusResponse struct {
	TransactionID     string `json:"transaction_id"`
	OrderID           string `json:"order_id"`
	StatusCode        string `json:"status_code"`
	TransactionStatus string `json:"transaction_status"`
	PaymentType       string `json:"payment_type"`
	GrossAmount       string `json:"gross_amount"`
}
