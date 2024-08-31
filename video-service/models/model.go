package models

import "time"

type Video struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Size      int64     `gorm:"column:size" json:"size"`
	Type      string    `gorm:"column:type" json:"type"`
	S3Url     string    `gorm:"column:s3_url" json:"s3_url"`
}
