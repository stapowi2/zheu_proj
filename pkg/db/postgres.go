package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
	"fmt"
)

func InitDB(dsn string) *gorm.DB {
	var db *gorm.DB
    var err error

	for i := 0; i < 10; i++ {
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err == nil {
            return db
        }
        fmt.Println("База еще не готова, жду 2 секунды...")
        time.Sleep(2 * time.Second)
    }
    
    log.Fatal("Не удалось подключиться к базе данных после 10 попыток")
    return nil
	
}