package model

import "github.com/google/uuid"

// role サービスの roles テーブルと同じ形
type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"type:varchar(255);not null;index"`
	Description string    `gorm:"type:text"`
}
