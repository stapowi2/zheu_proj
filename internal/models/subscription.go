package models


type Subscription struct {
    ID        uint   `gorm:"primaryKey"`
    UserID    uint   `json:"user_id"`
    Endpoint  string `json:"endpoint"`
    P256dh    string `json:"p256dh"`
    Auth      string `json:"auth"`
}