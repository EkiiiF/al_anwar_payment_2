package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/repository"
	"gorm.io/gorm"
)

type ReportServiceImpl struct {
	DB         *gorm.DB
	Repository repository.ReportRepository
}

func NewReportService(db *gorm.DB, repo repository.ReportRepository) ReportService {
	return &ReportServiceImpl{DB: db, Repository: repo}
}

func (s *ReportServiceImpl) GetReport(ctx context.Context, filters map[string]interface{}) (map[string]interface{}, error) {
	invoices, err := s.Repository.GetInvoiceReports(s.DB, filters)
	if err != nil {
		return nil, err
	}

	summary, err := s.Repository.GetSummary(s.DB, filters)
	if err != nil {
		return nil, err
	}

	summary["invoices"] = invoices
	return summary, nil
}

func (s *ReportServiceImpl) ExportCSV(ctx context.Context, filters map[string]interface{}) ([]byte, error) {
	invoices, err := s.Repository.GetInvoiceReports(s.DB, filters)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	// Header
	if err := writer.Write([]string{"No. Invoice", "Nama Santri", "NIS", "Bulan", "Tahun", "Nominal", "Status"}); err != nil {
		return nil, fmt.Errorf("gagal menulis header CSV: %w", err)
	}

	for _, inv := range invoices {
		if err := writer.Write([]string{
			inv.InvoiceNumber,
			inv.Student.FullName(),
			inv.Student.StudentNumber,
			fmt.Sprintf("%d", inv.Month),
			fmt.Sprintf("%d", inv.Year),
			fmt.Sprintf("%.2f", inv.AmountDue),
			inv.Status,
		}); err != nil {
			return nil, fmt.Errorf("gagal menulis baris CSV: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("gagal flush CSV: %w", err)
	}
	return buf.Bytes(), nil
}
