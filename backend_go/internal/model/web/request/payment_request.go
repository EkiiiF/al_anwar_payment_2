// Package request — model request untuk transaksi pembayaran.
package request

// CreateTransactionRequest — request pembuatan transaksi Midtrans.
type CreateTransactionRequest struct {
	InvoiceIDs []string `json:"invoice_ids" validate:"required,min=1"`
}
