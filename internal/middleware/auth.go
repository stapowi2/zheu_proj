package middleware

import (
	"net/http"
	"strings"
	"github.com/golang-jwt/jwt/v5"
	"zheu-system/internal/auth"
	"context"
)

type ContextKey string
const UserClaimsKey ContextKey = "userClaims"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid auth format", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(parts[1], &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret_key"), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserClaimsKey, token.Claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}


func RoleMiddleware(requiredRole string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(UserClaimsKey).(*auth.Claims)
			
			if claims.Role != requiredRole && claims.Role != "admin" {
				http.Error(w, "У вас недостаточно прав", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}
	}
}