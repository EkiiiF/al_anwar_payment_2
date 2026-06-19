package service

import (
	"context"
)

// ReportService — kontrak business logic laporan keuangan.
type ReportService interface {
	GetReport(ctx context.Context, filters map[string]interface{}) (map[string]interface{}, error)
	ExportCSV(ctx context.Context, filters map[string]interface{}) ([]byte, error)
}
