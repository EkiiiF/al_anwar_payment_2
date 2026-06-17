package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/config"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/response"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/repository"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/utils"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

type PaymentServiceImpl struct {
	DB           *gorm.DB
	PaymentRepo  repository.PaymentRepository
	GuardianRepo repository.GuardianRepository
	SnapClient   snap.Client
	CoreClient   coreapi.Client
}

func NewPaymentService(db *gorm.DB, paymentRepo repository.PaymentRepository, guardianRepo repository.GuardianRepository, cfg *config.Config) PaymentService {
	svc := &PaymentServiceImpl{
		DB:           db,
		PaymentRepo:  paymentRepo,
		GuardianRepo: guardianRepo,
	}

	var env midtrans.EnvironmentType
	if cfg.MidtransIsProd {
		env = midtrans.Production
	} else {
		env = midtrans.Sandbox
	}
	svc.SnapClient.New(cfg.MidtransServerKey, env)
	svc.CoreClient.New(cfg.MidtransServerKey, env)

	return svc
}

func (s *PaymentServiceImpl) CheckStatus(ctx context.Context, orderID string) error {
	resp, err := s.CoreClient.CheckTransaction(orderID)
	if err != nil {
		log.Printf("[Midtrans] CheckTransaction error for OrderID %s: %v", orderID, err)
		
		// Cek jika pembayaran ada di database dan berstatus pending
		var payments []domain.Payment
		if findErr := s.DB.Where("external_id = ?", orderID).Find(&payments).Error; findErr == nil && len(payments) > 0 {
			status := ""
			if payments[0].TransactionStatus != nil {
				status = *payments[0].TransactionStatus
			}
			if status == "pending" || status == "pending_payment" {
				timeDiff := time.Since(payments[0].CreatedAt)
				if timeDiff > 10*time.Minute {
					log.Printf("[Midtrans] Transaksi %s berumur %v dan tidak ditemukan di Midtrans. Mengubah status ke expire.", orderID, timeDiff)
					notification := map[string]interface{}{
						"order_id":           orderID,
						"transaction_status": "expire",
						"payment_type":       "-",
						"fraud_status":       "accept",
						"gross_amount":       "0",
					}
					return s.HandleNotification(ctx, notification)
				}
			}
		}
		return fmt.Errorf("gagal mengecek status ke Midtrans: %w", err)
	}

	notification := map[string]interface{}{
		"order_id":           resp.OrderID,
		"transaction_status": resp.TransactionStatus,
		"payment_type":       resp.PaymentType,
		"fraud_status":       resp.FraudStatus,
		"gross_amount":       resp.GrossAmount,
	}

	return s.HandleNotification(ctx, notification)
}

func (s *PaymentServiceImpl) CreateTransaction(ctx context.Context, invoiceIDs []string) (response.TransactionResponse, error) {
	if len(invoiceIDs) == 0 {
		return response.TransactionResponse{}, fmt.Errorf("pilih minimal satu tagihan")
	}

	var totalAmount float64
	var itemDetails []midtrans.ItemDetails
	var invoices []domain.Invoice

	for _, id := range invoiceIDs {
		invoice, err := s.PaymentRepo.FindInvoiceByID(s.DB, id)
		if err != nil {
			return response.TransactionResponse{}, fmt.Errorf("tagihan ID %s tidak ditemukan", id)
		}

		if invoice.Status == "paid" {
			return response.TransactionResponse{}, fmt.Errorf("tagihan %s sudah lunas", invoice.InvoiceNumber)
		}
		if invoice.Status == "pending" {
			return response.TransactionResponse{}, fmt.Errorf("tagihan %s sedang dalam proses pembayaran (pending)", invoice.InvoiceNumber)
		}

		totalAmount += invoice.AmountDue
		invoices = append(invoices, invoice)
		itemDetails = append(itemDetails, midtrans.ItemDetails{
			ID:    invoice.ID,
			Price: int64(invoice.AmountDue),
			Qty:   1,
			Name:  fmt.Sprintf("Tagihan %s %d/%d", invoice.Student.FullName(), invoice.Month, invoice.Year),
		})
	}

	// Ambil info guardian untuk customer details
	var guardianEmail, guardianName, guardianPhone string
	if len(invoices) > 0 && len(invoices[0].Student.Guardians) > 0 {
		g := invoices[0].Student.Guardians[0]
		guardianEmail = g.Email
		guardianName = g.FullName()
		guardianPhone = g.Phone
	}

	orderID := fmt.Sprintf("PAY-%d-%s", time.Now().Unix(), utils.GenerateID()[:8])

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(totalAmount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: guardianName,
			Email: guardianEmail,
			Phone: guardianPhone,
		},
		Items: &itemDetails,
	}

	snapResp, snapErr := s.SnapClient.CreateTransaction(snapReq)
	if snapErr != nil {
		return response.TransactionResponse{}, fmt.Errorf("gagal membuat transaksi ke Midtrans: %w", snapErr)
	}

	// Buat Payment record per invoice
	statusInitial := "pending_payment"
	for _, inv := range invoices {
		payment := domain.Payment{
			ID:                utils.GenerateID(),
			InvoiceID:         inv.ID,
			ExternalID:        &orderID,
			SnapToken:         &snapResp.Token,
			AmountPaid:        inv.AmountDue,
			TransactionStatus: &statusInitial,
			RawResponse:       "{}",
			CreatedAt:         time.Now(),
		}
		if err := s.PaymentRepo.CreatePayment(s.DB, &payment); err != nil {
			return response.TransactionResponse{}, fmt.Errorf("gagal menyimpan payment: %w", err)
		}
		if err := s.PaymentRepo.UpdateInvoiceStatus(s.DB, inv.ID, "pending"); err != nil {
			return response.TransactionResponse{}, fmt.Errorf("gagal memperbarui status invoice ke pending: %w", err)
		}
	}

	return response.TransactionResponse{
		Token:       snapResp.Token,
		RedirectURL: snapResp.RedirectURL,
	}, nil
}

func (s *PaymentServiceImpl) HandleNotification(ctx context.Context, notification map[string]interface{}) error {
	orderID, ok := notification["order_id"].(string)
	if !ok {
		return fmt.Errorf("invalid notification: missing order_id")
	}

	return s.DB.Transaction(func(tx *gorm.DB) error {
		payments, err := s.PaymentRepo.FindPaymentByExternalID(tx, orderID)
		if err != nil || len(payments) == 0 {
			return fmt.Errorf("transaksi tidak ditemukan: %s", orderID)
		}

		transactionStatus, _ := notification["transaction_status"].(string)
		fraudStatus, _ := notification["fraud_status"].(string)
		paymentType, _ := notification["payment_type"].(string)

		log.Printf("[Midtrans] Webhook received - OrderID: %s, Status: %s", orderID, transactionStatus)

		raw, _ := json.Marshal(notification)
		rawStr := string(raw)

		for i := range payments {
			p := &payments[i]
			p.RawResponse = rawStr
			p.PaymentMethod = &paymentType

			switch transactionStatus {
			case "capture":
				if fraudStatus == "accept" {
					s.markPaymentSuccessWithTx(tx, p)
				}
			case "settlement":
				s.markPaymentSuccessWithTx(tx, p)
			case "deny", "expire", "cancel":
				s.markPaymentFailureWithTx(tx, p, transactionStatus)
			case "pending":
				status := "pending"
				p.TransactionStatus = &status
				if err := s.PaymentRepo.UpdateInvoiceStatus(tx, p.InvoiceID, "pending"); err != nil {
					log.Printf("[WARN] Gagal update invoice status: %v", err)
				}
			}

			if err := s.PaymentRepo.UpdatePayment(tx, p); err != nil {
				return fmt.Errorf("gagal update payment: %w", err)
			}
		}

		// Kirim notifikasi ke guardian
		s.sendPaymentNotificationWithTx(ctx, tx, payments, transactionStatus)

		return nil
	})
}

func (s *PaymentServiceImpl) markPaymentSuccessWithTx(tx *gorm.DB, payment *domain.Payment) {
	status := "settlement"
	payment.TransactionStatus = &status
	now := time.Now()
	payment.PaymentDate = &now
	if err := s.PaymentRepo.UpdateInvoiceStatus(tx, payment.InvoiceID, "paid"); err != nil {
		log.Printf("[WARN] Gagal update invoice ke paid: %v", err)
	}
}

func (s *PaymentServiceImpl) markPaymentFailureWithTx(tx *gorm.DB, payment *domain.Payment, status string) {
	payment.TransactionStatus = &status
	if err := s.PaymentRepo.UpdateInvoiceStatus(tx, payment.InvoiceID, "unpaid"); err != nil {
		log.Printf("[WARN] Gagal reset invoice ke unpaid: %v", err)
	}
}

func (s *PaymentServiceImpl) sendPaymentNotificationWithTx(ctx context.Context, tx *gorm.DB, payments []domain.Payment, transactionStatus string) {
	if len(payments) == 0 {
		return
	}

	payment := payments[0]

	studentName := payment.Invoice.Student.FullName()

	// Cari UserID dari guardian santri ini
	var userID string
	if len(payment.Invoice.Student.Guardians) > 0 {
		userID = payment.Invoice.Student.Guardians[0].UserID
	}

	if userID == "" {
		log.Printf("[WARN] Tidak dapat menemukan UserID guardian untuk student %s", payment.Invoice.StudentID)
		return
	}

	var totalAmountPaid float64
	for _, p := range payments {
		totalAmountPaid += p.AmountPaid
	}

	var notifyTitle, notifyMsg string
	var shouldNotify bool

	switch transactionStatus {
	case "capture", "settlement":
		notifyTitle = "Pembayaran Berhasil"
		notifyMsg = fmt.Sprintf("Alhamdulillah! Pembayaran untuk %s sebesar %s telah berhasil. Terima kasih.",
			studentName, utils.FormatRupiah(totalAmountPaid))
		shouldNotify = true
	case "deny", "expire", "cancel":
		notifyTitle = "Pembayaran Gagal/Dibatalkan"
		notifyMsg = fmt.Sprintf("Mohon maaf, pembayaran untuk %s sebesar %s gagal atau kedaluwarsa.",
			studentName, utils.FormatRupiah(totalAmountPaid))
		shouldNotify = true
	}

	if shouldNotify {
		// 1. Kirim Notifikasi
		notifGuardian := &domain.Notification{
			ID:        utils.GenerateID(),
			UserID:    userID,
			Title:     notifyTitle,
			Message:   notifyMsg,
			IsRead:    false,
			CreatedAt: time.Now(),
		}
		if err := s.GuardianRepo.CreateNotification(ctx, tx, notifGuardian); err != nil {
			log.Printf("[WARN] Gagal simpan notifikasi pembayaran untuk guardian: %v", err)
		}

		// 2. Kirim Notifikasi ke SEMUA Super User (Bendahara/Admin)
		var superUsers []domain.User
		// Cari user yang memiliki role super_user
		tx.Joins("JOIN roles ON roles.id = users.role_id").
			Where("roles.name = ?", "super_user").
			Find(&superUsers)

		for _, admin := range superUsers {
			// Pesan khusus untuk admin
			var adminMsg string
			if transactionStatus == "capture" || transactionStatus == "settlement" {
				adminMsg = fmt.Sprintf("Pembayaran Masuk: %s telah membayar tagihan sebesar %s.",
					studentName, utils.FormatRupiah(totalAmountPaid))
			} else {
				adminMsg = fmt.Sprintf("Pembayaran Gagal: Upaya pembayaran untuk %s sebesar %s gagal.",
					studentName, utils.FormatRupiah(totalAmountPaid))
			}

			notifAdmin := &domain.Notification{
				ID:        utils.GenerateID(),
				UserID:    admin.ID,
				Title:     notifyTitle,
				Message:   adminMsg,
				IsRead:    false,
				CreatedAt: time.Now(),
			}
			if err := s.GuardianRepo.CreateNotification(ctx, tx, notifAdmin); err != nil {
				log.Printf("[WARN] Gagal simpan notifikasi pembayaran untuk admin %s: %v", admin.ID, err)
			}
		}
	}
}
