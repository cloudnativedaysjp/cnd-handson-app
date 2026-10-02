package model

import (
	"time"

	"github.com/google/uuid"
)

// user サービスの users テーブルと同じ形にして、データをそのまま引き継ぐ
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name         string    `gorm:"type:varchar(255)"`
	Email        string    `gorm:"type:varchar(255);unique;not null"`
	PasswordHash string    `gorm:"type:text;not null"`
	RoleID       uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt    time.Time `gorm:"type:timestamp;not null;default:current_timestamp"`
	UpdatedAt    time.Time `gorm:"type:timestamp;not null"`
}
