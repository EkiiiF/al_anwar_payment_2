package service

import (
	"context"
)

// ReportService defines business operations for financial reporting.
type ReportService interface {
	GetReport(ctx context.Context, filters map[string]interface{}) (map[string]interface{}, error)
	ExportCSV(ctx context.Context, filters map[string]interface{}) ([]byte, error)
}
