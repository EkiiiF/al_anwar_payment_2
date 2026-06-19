package service

import (
	"context"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/response"
)

// PaymentService — kontrak business logic pembayaran Midtrans.
type PaymentService interface {
	CreateTransaction(ctx context.Context, invoiceIDs []string) (response.TransactionResponse, error)
	HandleNotification(ctx context.Context, notification map[string]interface{}) error
	CheckStatus(ctx context.Context, orderID string) error
}
