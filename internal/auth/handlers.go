package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"zheu-system/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		var user models.User
		if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		if err := user.CheckPassword(req.Password); err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			fmt.Println(err)
			return
		}


		accessToken, refreshToken, err := GenerateTokens(user.ID, user.Username, user.Role)
		if err != nil {
			http.Error(w, "Could not generate tokens", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}


func RefreshHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            RefreshToken string `json:"refresh_token"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "bad request", http.StatusBadRequest)
            return
        }

        claims := &RefreshClaims{} 
        
        token, err := jwt.ParseWithClaims(req.RefreshToken, claims, func(token *jwt.Token) (any, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
            return
        }

        var user models.User
        if err := db.First(&user, claims.UserID).Error; err != nil {
            http.Error(w, "User not found", http.StatusUnauthorized)
            return
        }

        newAccess, newRefresh, err := GenerateTokens(user.ID, user.Username, user.Role)
        if err != nil {
            http.Error(w, "Error generate tokens", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "access_token":  newAccess,
            "refresh_token": newRefresh,
        })
    }
}

func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}


func RegisterHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var u models.User
        json.NewDecoder(r.Body).Decode(&u)
        
        hashedPassword, _ := hashPassword(u.Password)
		u.Password = hashedPassword
        db.Create(&u)
        w.WriteHeader(http.StatusCreated)
    }
}