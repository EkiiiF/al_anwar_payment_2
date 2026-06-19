package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"time"

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

	// Write metadata rows
	exporter, _ := filters["exporter"].(string)
	if exporter == "" {
		exporter = "Sistem"
	}
	yearVal, _ := filters["year"].(string)
	monthVal, _ := filters["month"].(string)

	writer.Write([]string{"LAPORAN KEUANGAN PONDOK PESANTREN AL-ANWAR"})
	writer.Write([]string{"Nama Laporan", "Laporan Pemasukan Santri (Lunas)"})
	writer.Write([]string{"Dibuat Oleh", exporter})
	
	// Convert server time to local or standard display format
	importTime := time.Now().Format("02 January 2006 15:04:05 MST")
	writer.Write([]string{"Waktu Ekspor", importTime})

	if yearVal != "" {
		if monthVal != "" {
			writer.Write([]string{"Periode", fmt.Sprintf("Bulan %s Tahun %s H", monthVal, yearVal)})
		} else {
			writer.Write([]string{"Periode", fmt.Sprintf("Tahun %s H", yearVal)})
		}
	}
	writer.Write([]string{""}) // Empty row as separator

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
