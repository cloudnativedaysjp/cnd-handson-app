package model

import "github.com/google/uuid"

// session サービスの refresh_tokens テーブルと同じ形。Token は bcrypt のハッシュ
type RefreshToken struct {
	UserID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Token  string    `gorm:"not null"`
	Exp    int64     `gorm:"not null"`
}
