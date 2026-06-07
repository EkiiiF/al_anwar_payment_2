package request

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateProfileRequest struct {
	FirstName   string `json:"first_name" validate:"required"`
	MiddleName  string `json:"middle_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Gender      string `json:"gender"`
	BirthDate   string `json:"birth_date"`
	Address     string `json:"address"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}
