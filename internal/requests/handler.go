package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"zheu-system/internal/auth"
	"zheu-system/internal/middleware"
	"zheu-system/internal/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func CreateRequestHandler(db *gorm.DB, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
		if !ok {
			http.Error(w, "Не удалось авторизовать пользователя", http.StatusUnauthorized)
			return
		}

		var newRequest models.Request
		err := json.NewDecoder(r.Body).Decode(&newRequest)
		if err != nil {
			http.Error(w, "Не удалось прочитать json", http.StatusBadRequest)
			return
		}
		newRequest.UserID = claims.UserID

		result := db.Create(&newRequest)
		if result.Error != nil {
			http.Error(w, "Ошибка при сохранении в базу", http.StatusInternalServerError)
			return
		}

		err = rdb.LPush(r.Context(), "notifications_queue", newRequest.ID).Err()
		if err != nil {
			http.Error(w, "Ошибка при постановке в очередь", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newRequest)
	}
}


func GetRequestsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(middleware.UserClaimsKey).(*auth.Claims)
		if !ok {
			http.Error(w, "Не удалось получить данные пользователя", http.StatusUnauthorized)
			return
		}
		var userRequests []models.Request
		db.Where("user_id = ?", claims.UserID).Find(&userRequests)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userRequests)
	}
}


type StatusUpdate struct {
	Status string `json:"status"`
}


func UpdateStatusHandler(db *gorm.DB, rdb *redis.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        parts := strings.Split(r.URL.Path, "/")

        if len(parts) < 4 {
            http.Error(w, "ID заявки не указан", http.StatusBadRequest)
            return
        }
        
        id, err := strconv.Atoi(parts[3])
        if err != nil {
            http.Error(w, "Некорректный ID", http.StatusBadRequest)
            return
        }

        var input StatusUpdate
        if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
            http.Error(w, "Неверный json", http.StatusBadRequest)
            return
        }

        result := db.Model(&models.Request{}).Where("id = ?", id).Update("status", input.Status)
        if result.Error != nil {
            http.Error(w, "Ошибка при обновлении в БД", http.StatusInternalServerError)
            return
        }
        
        if result.RowsAffected == 0 {
            http.Error(w, "Заявка не найдена", http.StatusNotFound)
            return
        }

        err = rdb.LPush(r.Context(), "notifications_queue", id).Err()
        if err != nil {
            fmt.Printf("Ошибка постановки в очередь обновлений: %v\n", err)
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"result": "Статус успешно изменен"})
    }
}