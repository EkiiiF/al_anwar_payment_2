// Package response — model response untuk pembayaran.
package response

// TransactionResponse — response pembuatan transaksi Midtrans Snap.
type TransactionResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}
