package service

import (
	"context"
	"fmt"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/response"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/repository"
	"gorm.io/gorm"
)

type GuardianServiceImpl struct {
	DB           *gorm.DB
	GuardianRepo repository.GuardianRepository
}

func NewGuardianService(db *gorm.DB, guardianRepo repository.GuardianRepository) GuardianService {
	return &GuardianServiceImpl{DB: db, GuardianRepo: guardianRepo}
}

func (s *GuardianServiceImpl) findStudentByUserID(ctx context.Context, userID string) (domain.Student, error) {
	guardian, err := s.GuardianRepo.FindGuardianByUserID(ctx, s.DB, userID)
	if err != nil {
		return domain.Student{}, fmt.Errorf("data wali tidak ditemukan: %w", err)
	}
	return guardian.Student, nil
}

func (s *GuardianServiceImpl) GetDashboardStats(ctx context.Context, userID string) (response.GuardianDashboardStats, error) {
	student, err := s.findStudentByUserID(ctx, userID)
	if err != nil {
		return response.GuardianDashboardStats{}, err
	}

	invoices, err := s.GuardianRepo.FindInvoicesByStudentID(ctx, s.DB, student.ID)
	if err != nil {
		return response.GuardianDashboardStats{}, fmt.Errorf("gagal mengambil tagihan: %w", err)
	}

	if student.MuhadhorohLevel == 0 {
		var filtered []domain.Invoice
		for _, inv := range invoices {
			if inv.Category.Name != "Syahriyyah Muhadhoroh" {
				filtered = append(filtered, inv)
			}
		}
		invoices = filtered
	}

	payments, err := s.GuardianRepo.FindPaymentsByStudentID(ctx, s.DB, student.ID)
	if err != nil {
		return response.GuardianDashboardStats{}, fmt.Errorf("gagal mengambil pembayaran: %w", err)
	}

	var totalUnpaid float64
	var unpaidCount int
	for _, inv := range invoices {
		if inv.Status == "unpaid" {
			totalUnpaid += inv.AmountDue
			unpaidCount++
		}
	}

	var lastPayment *domain.Payment
	if len(payments) > 0 {
		lastPayment = &payments[0]
	}

	return response.GuardianDashboardStats{
		Student:        student,
		TotalUnpaid:    totalUnpaid,
		UnpaidCount:    unpaidCount,
		TotalInvoices:  len(invoices),
		LastPayment:    lastPayment,
		RecentPayments: payments,
	}, nil
}

func (s *GuardianServiceImpl) GetInvoices(ctx context.Context, userID string) ([]domain.Invoice, error) {
	student, err := s.findStudentByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	invoices, err := s.GuardianRepo.FindInvoicesByStudentID(ctx, s.DB, student.ID)
	if err != nil {
		return nil, err
	}

	if student.MuhadhorohLevel == 0 {
		var filtered []domain.Invoice
		for _, inv := range invoices {
			if inv.Category.Name != "Syahriyyah Muhadhoroh" {
				filtered = append(filtered, inv)
			}
		}
		return filtered, nil
	}

	return invoices, nil
}

func (s *GuardianServiceImpl) GetPayments(ctx context.Context, userID string) ([]domain.Payment, error) {
	student, err := s.findStudentByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.GuardianRepo.FindPaymentsByStudentID(ctx, s.DB, student.ID)
}

func (s *GuardianServiceImpl) GetNotifications(ctx context.Context, userID string) ([]domain.Notification, error) {
	return s.GuardianRepo.FindNotificationsByUserID(ctx, s.DB, userID)
}

func (s *GuardianServiceImpl) MarkNotificationRead(ctx context.Context, notificationID string, userID string) error {
	return s.GuardianRepo.MarkNotificationAsRead(ctx, s.DB, notificationID, userID)
}
