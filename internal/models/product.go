package models

import "time"

type Products struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:string;not null"`
	Description string    `json:"description"`
	Price       float64   `json:"price" gorm:"default:0.00"`
	Category    string    `json:"category"`
	Stock       int       `json:"stock" gorm:"default:0"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"createdAt" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"default:CURRENT_TIMESTAMP"`
}
