package request

type CreateStudentRequest struct {
	NIS             string `json:"nis"              validate:"required,min=3,max=20"`
	FullName        string `json:"full_name"         validate:"required,min=3"`
	Nik             string `json:"nik"`
	BirthDate       string `json:"birth_date"`
	Gender          string `json:"gender"`
	MuhadhorohLevel int    `json:"muhadhoroh_level"`
	CurrentSemester int    `json:"current_semester"`
	StatusTypeID    string `json:"status_type_id"    validate:"required"`

	AddressLine string `json:"address_line"`
	Village     string `json:"village"`
	District    string `json:"district"`
	City        string `json:"city"`
	Province    string `json:"province"`
	Country     string `json:"country"`
	PostalCode  string `json:"postal_code"`

	CreateNewGuardian bool   `json:"create_new_guardian"`
	GuardianName      string `json:"guardian_name"`
	GuardianPhone     string `json:"guardian_phone"`
	GuardianEmail     string `json:"guardian_email"`
	GuardianRelation  string `json:"guardian_relation"`
	GuardianUsername  string `json:"guardian_username"`
	GuardianPassword  string `json:"guardian_password"`

	GuardianID string `json:"guardian_id"`
}

type UpdateStudentRequest struct {
	NIS             string `json:"nis"              validate:"required,min=3,max=20"`
	FullName        string `json:"full_name"         validate:"required,min=3"`
	Nik             string `json:"nik"`
	BirthDate       string `json:"birth_date"`
	Gender          string `json:"gender"`
	MuhadhorohLevel int    `json:"muhadhoroh_level"`
	CurrentSemester int    `json:"current_semester"`
	StatusTypeID    string `json:"status_type_id"    validate:"required"`

	AddressLine string `json:"address_line"`
	Village     string `json:"village"`
	District    string `json:"district"`
	City        string `json:"city"`
	Province    string `json:"province"`
	Country     string `json:"country"`
	PostalCode  string `json:"postal_code"`

	GuardianName     string `json:"guardian_name"`
	GuardianPhone    string `json:"guardian_phone"`
	GuardianEmail    string `json:"guardian_email"`
	GuardianRelation string `json:"guardian_relation"`
	GuardianUsername string `json:"guardian_username"`
	GuardianPassword string `json:"guardian_password"`
}
