package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	UserID   string `json:"sub"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

func InitJWTSecret(secret string) {
	if secret == "" {
		panic("JWT_SECRET tidak boleh kosong — set di .env")
	}
	jwtSecret = []byte(secret)
}

func GenerateJWT(userID, username, role string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "al-anwar-payment",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("gagal membuat token: %w", err)
	}
	return tokenString, nil
}

func ValidateJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("signing method tidak valid: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("token tidak valid: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("klaim token tidak valid")
	}
	return claims, nil
}
