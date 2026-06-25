package models

import (
	"gorm.io/gorm"
)

type Request struct {
	gorm.Model
	UserID      uint   `json:"user_id"`
	Description string `json:"description"`
	Status      string `json:"status" gorm:"default:open"`
	HouseID     uint   `json:"house_id"`
}