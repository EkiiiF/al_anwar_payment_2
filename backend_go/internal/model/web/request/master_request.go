// Package request — model request untuk master data dan invoice.
package request

// GenerateInvoiceRequest — request generate invoice bulanan.
type GenerateInvoiceRequest struct {
	Month      int `json:"month"       validate:"required,min=1,max=12"`
	Year       int `json:"year"        validate:"required,min=2024"`
	HijriMonth int `json:"hijri_month"`
	HijriYear  int `json:"hijri_year"`
}

// CreateCategoryRequest — request buat/update kategori tagihan.
type CreateCategoryRequest struct {
	Name        string   `json:"name"        validate:"required,min=2"`
	BaseAmount  float64  `json:"base_amount" validate:"required,min=0"`
	IsFixed     bool     `json:"is_fixed"`
	IsActive    bool     `json:"is_active"`
	IsSemester  bool     `json:"is_semester"`
	Description string   `json:"description"`
	ApplyToAll  bool     `json:"apply_to_all"`
	StudentIDs  []string `json:"student_ids"`
}

// CreateStatusTypeRequest — request buat/update tipe status tagihan.
type CreateStatusTypeRequest struct {
	Name               string  `json:"name"                validate:"required,min=2"`
	DiscountPercentage float64 `json:"discount_percentage" validate:"min=0,max=100"`
	IsActiveBilling    *bool   `json:"is_active_billing"   validate:"required"`
	Description        string  `json:"description"`
}

// GenerateSemesterInvoiceRequest — request generate invoice semester.
type GenerateSemesterInvoiceRequest struct {
	Semester  int `json:"semester"   validate:"required,min=1,max=2"`
	HijriYear int `json:"hijri_year" validate:"required,min=1"`
}

// UpdateSettingRequest — request update konfigurasi sistem.
type UpdateSettingRequest struct {
	Key   string `json:"key"   validate:"required"`
	Value string `json:"value" validate:"required"`
}

// UpdateUserRoleRequest — request ubah role user.
type UpdateUserRoleRequest struct {
	RoleID string `json:"role_id" validate:"required"`
}
