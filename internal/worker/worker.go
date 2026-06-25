package worker

import (
    "context"
    "fmt"
    "os"
    "zheu-system/internal/models"

    "github.com/SherClockHolmes/webpush-go"
    "github.com/joho/godotenv"
    "github.com/redis/go-redis/v9"
    "gorm.io/gorm"
)

func StartNotificationWorker(db *gorm.DB, rdb *redis.Client) {
    fmt.Println("воркер уведомлений запущен и слушает очередь...")
    ctx := context.Background()

    err_env := godotenv.Load()
    if err_env != nil {
        fmt.Println("Предупреждение: Ошибка загрузки .env файла")
    }

    privateKey := os.Getenv("VAPID_PRIVATE_KEY")
    publicKey := os.Getenv("VAPID_PUBLIC_KEY")

    for {
        result, err := rdb.BRPop(ctx, 0, "notifications_queue").Result()
        if err != nil {
            fmt.Printf("Ошибка чтения из Redis: %v\n", err)
            continue
        }

        if len(result) < 2 {
            continue
        }

        requestID := result[1]

        var req models.Request
        if err := db.First(&req, requestID).Error; err != nil {
            fmt.Printf("Заявка ID %s не найдена в базе данных: %v\n", requestID, err)
            continue
        }
        fmt.Printf("Заявка найдена в БД. Автор user_id: %d, Текущий статус: %s\n", req.UserID, req.Status)

        var sub models.Subscription
        if err := db.Where("user_id = ?", req.UserID).First(&sub).Error; err != nil {
            fmt.Printf("Подписка для user_id %d не найдена в БД\n", req.UserID)
            continue
        }
        fmt.Printf("Подписка найдена для user_id: %d\n", req.UserID)

        s := &webpush.Subscription{
            Endpoint: sub.Endpoint,
            Keys: webpush.Keys{P256dh: sub.P256dh, Auth: sub.Auth},
        }

        pushMessage := fmt.Sprintf("Статус вашей заявки №%d изменен на: %s", req.ID, req.Status)

        resp, err := webpush.SendNotification([]byte(pushMessage), s, &webpush.Options{
            Subscriber:      "mailto:ukesaman@gmail.com",
            VAPIDPublicKey:  publicKey,
            VAPIDPrivateKey: privateKey,
            TTL:             30,
        })

        if err != nil {
            fmt.Printf("Не удалось отправить пуш: %v\n", err)
            continue
        }
        
        fmt.Printf("Пуш отправлен! Статус ответа: %d %s\n", resp.StatusCode, resp.Status)
        resp.Body.Close()
    }
}