package service

import (
	"time"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/exception"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/request"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/response"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/repository"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthServiceImpl struct {
	DB             *gorm.DB
	Validate       *validator.Validate
	AuthRepository repository.AuthRepository
}

func NewAuthService(db *gorm.DB, validate *validator.Validate, authRepo repository.AuthRepository) AuthService {
	return &AuthServiceImpl{
		DB:             db,
		Validate:       validate,
		AuthRepository: authRepo,
	}
}

func (s *AuthServiceImpl) Login(req request.LoginRequest, ip, userAgent string) (response.LoginResponse, error) {
	if err := s.Validate.Struct(req); err != nil {
		return response.LoginResponse{}, err
	}

	user, err := s.AuthRepository.FindUserByUsername(s.DB, req.Username)
	if err != nil {
		return response.LoginResponse{}, exception.ValidationError{Message: "Username atau password salah"}
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return response.LoginResponse{}, exception.ValidationError{Message: "Username atau password salah"}
	}

	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role.Name)
	if err != nil {
		return response.LoginResponse{}, err
	}

	now := time.Now()
	_ = s.AuthRepository.UpdateLastLogin(s.DB, user.ID, now)

	return response.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		UserID:      user.ID,
		Username:    user.Username,
		Role:        user.Role.Name,
	}, nil
}

func (s *AuthServiceImpl) Logout(token string, userID, ip, userAgent string) error {
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		return err
	}

	return s.AuthRepository.BlacklistToken(s.DB, token, claims.ExpiresAt.Time)
}

func (s *AuthServiceImpl) GetProfile(userID string) (domain.User, error) {
	user, err := s.AuthRepository.FindUserByID(s.DB, userID)
	if err != nil {
		return user, err
	}

	if user.Role.Name == "guardian" {
		var guardian domain.Guardian
		if err := s.DB.Where("user_id = ?", user.ID).First(&guardian).Error; err == nil {
			var address string
			var student domain.Student
			if err := s.DB.Preload("Addresses").Where("id = ?", guardian.StudentID).First(&student).Error; err == nil {
				for _, addr := range student.Addresses {
					if addr.IsPrimary {
						address = addr.AddressLine
						break
					}
				}
			}

			user.FirstName = guardian.Name.FirstName
			user.MiddleName = guardian.Name.MiddleName
			user.LastName = guardian.Name.LastName
			user.Email = guardian.Email
			user.PhoneNumber = guardian.Phone
			user.Address = address
		}
	} else {
		if user.FirstName == "" {
			user.FirstName = user.Username
		}
		if user.Email == "" {
			user.Email = user.Username + "@alanwar.com"
		}
		if user.Gender == "" {
			user.Gender = "L"
		}
	}

	return user, nil
}

func (s *AuthServiceImpl) UpdateProfile(userID string, req request.UpdateProfileRequest) (domain.User, error) {
	if err := s.Validate.Struct(req); err != nil {
		return domain.User{}, err
	}

	user, err := s.AuthRepository.FindUserByID(s.DB, userID)
	if err != nil {
		return domain.User{}, err
	}

	var birthDate *time.Time
	if req.BirthDate != "" {
		t, err := time.Parse("2006-01-02", req.BirthDate)
		if err == nil {
			birthDate = &t
		}
	}

	if user.Role.Name == "guardian" {
		var guardian domain.Guardian
		if err := s.DB.Where("user_id = ?", user.ID).First(&guardian).Error; err != nil {
			return domain.User{}, err
		}

		guardian.Name.FirstName = req.FirstName
		guardian.Name.MiddleName = req.MiddleName
		guardian.Name.LastName = req.LastName
		guardian.Email = req.Email
		guardian.Phone = req.Phone

		if err := s.DB.Save(&guardian).Error; err != nil {
			return domain.User{}, err
		}

		if req.Address != "" {
			var student domain.Student
			if err := s.DB.Preload("Addresses").Where("id = ?", guardian.StudentID).First(&student).Error; err == nil {
				found := false
				for i, addr := range student.Addresses {
					if addr.IsPrimary {
						student.Addresses[i].AddressLine = req.Address
						s.DB.Save(&student.Addresses[i])
						found = true
						break
					}
				}
				if !found && len(student.Addresses) == 0 {
					newAddr := domain.Address{
						ID:          uuid.New().String(),
						StudentID:   student.ID,
						AddressLine: req.Address,
						IsPrimary:   true,
					}
					s.DB.Create(&newAddr)
				}
			}
		}
	} else {
		user.FirstName = req.FirstName
		user.MiddleName = req.MiddleName
		user.LastName = req.LastName
		user.Email = req.Email
		user.PhoneNumber = req.Phone
		if req.Gender != "" {
			user.Gender = req.Gender
		}
		if birthDate != nil {
			user.BirthDate = birthDate
		}
		user.Address = req.Address

		if err := s.DB.Save(&user).Error; err != nil {
			return domain.User{}, err
		}
	}

	return s.GetProfile(userID)
}

func (s *AuthServiceImpl) ChangePassword(userID string, req request.ChangePasswordRequest) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	if req.NewPassword != req.ConfirmPassword {
		return exception.ValidationError{Message: "Konfirmasi password baru tidak cocok"}
	}

	user, err := s.AuthRepository.FindUserByID(s.DB, userID)
	if err != nil {
		return err
	}

	if !utils.CheckPasswordHash(req.CurrentPassword, user.Password) {
		return exception.ValidationError{Message: "Password saat ini salah"}
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return s.DB.Model(&domain.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error
}
