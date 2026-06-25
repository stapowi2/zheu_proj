package notifications

import (
	"encoding/json"
	"net/http"
	"zheu-system/internal/auth"
	"zheu-system/internal/middleware"
	"zheu-system/internal/models"

	"gorm.io/gorm"
)

func RegisterNotificationHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
		var sub models.Subscription
		if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		sub.UserID = claims.UserID
		db.Create(&sub)
		w.WriteHeader(http.StatusCreated)
	}
}


func UnsubscribeHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        claims := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
        
        result := db.Where("user_id = ?", claims.UserID).Delete(&models.Subscription{})
        
        if result.Error != nil {
            http.Error(w, "Ошибка при удалении", http.StatusInternalServerError)
            return
        }
        
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Отписка успешна"))
    }
}