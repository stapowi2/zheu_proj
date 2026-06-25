package news

import (
	"time"
	"encoding/json"
	"net/http"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"zheu-system/internal/models"
)

func GetNewsHandler(db *gorm.DB, rdb *redis.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cachedNews, err := rdb.Get(r.Context(), "news_cache").Result()
        if err == nil {
            w.Header().Set("Content-Type", "application/json")
            w.Write([]byte(cachedNews))
            return
        }

        var news []models.News
        db.Find(&news)

        data, _ := json.Marshal(news)
        rdb.Set(r.Context(), "news_cache", data, 10*time.Minute)

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(news)
    }
}


func CreateNewsHandler(db *gorm.DB, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
			return
		}

		var newNews models.News
		if err := json.NewDecoder(r.Body).Decode(&newNews); err != nil {
			http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
			return
		}

		if err := db.Create(&newNews).Error; err != nil {
			http.Error(w, "Ошибка при публикации новости", http.StatusInternalServerError)
			return
		}

		rdb.Del(r.Context(), "news_cache")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newNews)
	}
}