package main

import (
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"

	"zheu-system/internal/auth"
	"zheu-system/internal/middleware"
	"zheu-system/internal/models"
	"zheu-system/internal/news"
	"zheu-system/internal/requests"
	"zheu-system/internal/worker"
	"zheu-system/pkg/db"
    "zheu-system/internal/notifications"
)

func main() {
    dsn := "host=db user=user password=password dbname=zheu port=5432 sslmode=disable"
    database := db.InitDB(dsn)

    err := database.AutoMigrate(&models.User{}, &models.Request{}, &models.News{}, &models.Subscription{})
    if err != nil {
        fmt.Println("Ошибка миграции бд: ", err)
    }

    

    rdb := redis.NewClient(&redis.Options{
        Addr: "redis:6379",
    })

    go worker.StartNotificationWorker(database, rdb)

    http.Handle("/", http.FileServer(http.Dir("./static")))
    http.HandleFunc("/auth/login", auth.LoginHandler(database))
    http.HandleFunc("/auth/register", auth.RegisterHandler(database))
    http.HandleFunc("/auth/refresh", auth.RefreshHandler(database))

    http.HandleFunc("/news", news.GetNewsHandler(database, rdb))
    http.HandleFunc("/news/create", middleware.AuthMiddleware(
        middleware.RoleMiddleware("admin")(news.CreateNewsHandler(database, rdb)),
    ))

    http.HandleFunc("/requests", middleware.AuthMiddleware(requests.CreateRequestHandler(database, rdb)))
    http.HandleFunc("/requests/list", middleware.AuthMiddleware(requests.GetRequestsHandler(database)))
    http.HandleFunc("/requests/status/", middleware.AuthMiddleware(
        middleware.RoleMiddleware("employee")(requests.UpdateStatusHandler(database, rdb)),
    ))


    http.HandleFunc("/notifications/register", middleware.AuthMiddleware(
        notifications.RegisterNotificationHandler(database),
    ))
    http.HandleFunc("/notifications/unsubscribe", middleware.AuthMiddleware(
        notifications.UnsubscribeHandler(database),
    ))


    http.ListenAndServe(":8080", nil)
}