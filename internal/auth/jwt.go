package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("secret_key")

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	UserID   uint   `json:"user_id"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}


func GenerateTokens(userID uint, username, role string) (string, string, error) {
	accessClaims := &Claims{
		Username: username,
		Role:     role,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(150 * time.Minute)),
		},
	}
	accessToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(jwtKey)

	refreshClaims := &RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	refreshToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(jwtKey)

	return accessToken, refreshToken, nil
}