package response

type MonthlyPaymentStat struct {
	Month int     `json:"month"`
	Year  int     `json:"year"`
	Total float64 `json:"total"`
}

type SemesterPaymentStat struct {
	SemesterName      string  `json:"semester_name"`
	AcademicYearLabel string  `json:"academic_year_label"`
	Total             float64 `json:"total"`
	InvoiceCount      int     `json:"invoice_count"`
	PaidCount         int     `json:"paid_count"`
	UnpaidCount       int     `json:"unpaid_count"`
}

type HijriMonthInfo struct {
	HijriMonth        int    `json:"hijri_month"`
	HijriMonthName    string `json:"hijri_month_name"`
	HijriYear         int    `json:"hijri_year"`
	Semester          int    `json:"semester"`
	SemesterName      string `json:"semester_name"`
	AcademicYearLabel string `json:"academic_year_label"`
	IsExamMonth       bool   `json:"is_exam_month"`
	IsRegistration    bool   `json:"is_registration"`
}

type SuperUserDashboardStats struct {
	TotalStudents   int64                 `json:"total_students"`
	UnpaidInvoices  int64                 `json:"unpaid_invoices"`
	PaidInvoices    int64                 `json:"paid_invoices"`
	TotalIncomeMo   float64               `json:"total_income_mo"`
	MonthlyPayments []MonthlyPaymentStat  `json:"monthly_payments"`
	SemesterStats   []SemesterPaymentStat `json:"semester_stats"`
	CurrentHijri    HijriMonthInfo        `json:"current_hijri"`
	AvailableYears  []int                 `json:"available_years"`
}
