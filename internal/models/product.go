package models

import (
	"mime/multipart"
	"time"
)

type UpsertProductRequest struct {
	Name        string                `form:"name" json:"name" binding:"required"`
	Description string                `form:"description" json:"description" binding:"required"`
	Price       float64               `form:"price" json:"price" binding:"required"`
	Category    string                `form:"category" json:"category" binding:"required"`
	Stock       int                   `form:"stock" json:"stock" binding:"required"`
	Filename    string                `json:"filename"`
	Image       *multipart.FileHeader `form:"image" json:"image"`
	ImageFile   *multipart.File       `form:"image_file" json:"image_file"`
}

type Products struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:string;not null;size:500"`
	Description string    `json:"description" gorm:"size:1000"`
	Price       float64   `json:"price" gorm:"default:0.00"`
	Category    string    `json:"category" gorm:"size:500"`
	Stock       int       `json:"stock" gorm:"default:0"`
	ImageName   string    `json:"imageName" gorm:"size:500"`
	ImageURL    string    `json:"imageUrl" gorm:"size:500"`
	CreatedAt   time.Time `json:"createdAt" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"default:CURRENT_TIMESTAMP"`
}
