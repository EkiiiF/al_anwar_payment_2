package request

type GenerateInvoiceRequest struct {
	Month      int `json:"month"       validate:"required,min=1,max=12"`
	Year       int `json:"year"        validate:"required,min=2024"`
	HijriMonth int `json:"hijri_month"`
	HijriYear  int `json:"hijri_year"`
}

type CreateCategoryRequest struct {
	Name        string  `json:"name"        validate:"required,min=2"`
	BaseAmount  float64 `json:"base_amount" validate:"required,min=0"`
	IsFixed     bool    `json:"is_fixed"`
	IsActive    bool    `json:"is_active"`
	Description string  `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        string  `json:"name"        validate:"required,min=2"`
	BaseAmount  float64 `json:"base_amount" validate:"required,min=0"`
	IsFixed     bool    `json:"is_fixed"`
	IsActive    bool    `json:"is_active"`
	Description string  `json:"description"`
}

type CreateStatusTypeRequest struct {
	Name               string  `json:"name"                validate:"required,min=2"`
	DiscountPercentage float64 `json:"discount_percentage" validate:"min=0,max=100"`
	IsActiveBilling    *bool   `json:"is_active_billing"   validate:"required"`
	Description        string  `json:"description"`
}

type UpdateStatusTypeRequest struct {
	Name               string  `json:"name"                validate:"required,min=2"`
	DiscountPercentage float64 `json:"discount_percentage" validate:"min=0,max=100"`
	IsActiveBilling    *bool   `json:"is_active_billing"   validate:"required"`
	Description        string  `json:"description"`
}
