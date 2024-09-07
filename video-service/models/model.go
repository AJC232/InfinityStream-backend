package models

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	ID            string         `gorm:"column:id;primaryKey" json:"id"`
	Title         string         `gorm:"column:title,size:255" json:"title"`
	Description   string         `gorm:"column:description,type:text" json:"description"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at,index" json:"deleted_at"`
	Type          string         `gorm:"column:type" json:"type"`
	VideoUrl      string         `gorm:"column:video_url" json:"video_url"`
	CoverPhotoURL string         `gorm:"column:cover_photo_url" json:"cover_photo_url"`
	Category      string         `gorm:"column:category" json:"category"`
	IsPremium     bool           `gorm:"column:is_premium" json:"is_premium"`
}
