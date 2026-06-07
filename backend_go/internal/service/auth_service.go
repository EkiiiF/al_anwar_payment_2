package service

import (
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/request"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/response"
)

// AuthService defines business operations for authentication.
type AuthService interface {
	Login(req request.LoginRequest, ip, userAgent string) (response.LoginResponse, error)
	Logout(token string, userID, ip, userAgent string) error
	GetProfile(userID string) (domain.User, error)
	UpdateProfile(userID string, req request.UpdateProfileRequest) (domain.User, error)
	ChangePassword(userID string, req request.ChangePasswordRequest) error
}
