package response

import "github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"

type GuardianDashboardStats struct {
	Student        domain.Student   `json:"student"`
	TotalUnpaid    float64          `json:"total_unpaid"`
	UnpaidCount    int              `json:"unpaid_count"`
	TotalInvoices  int              `json:"total_invoices"`
	LastPayment    *domain.Payment  `json:"last_payment"`
	RecentPayments []domain.Payment `json:"recent_payments"`
}
