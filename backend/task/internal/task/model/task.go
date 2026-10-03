package model

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`                              // UUID型の主キー
	Title       string    `gorm:"type:varchar(255);not null"`                        // VARCHAR型
	Description string    `gorm:"type:text;"`                                        // TEXT型
	Status      string    `gorm:"type:varchar(255);not null"`                        // VARCHAR型
	Start_time  time.Time `gorm:"type:timestamp;not null;default:current_timestamp"` // TIMESTAMP型
	End_time    time.Time `gorm:"type:timestamp;"`                                   // TIMESTAMP型
	Column_id   uuid.UUID `gorm:"type:uuid;not null"`                                // UUID型の外部キー
	Assignee_id uuid.UUID `gorm:"type:uuid;not null"`                                // UUID型の外部キー
	Project_id  uuid.UUID `gorm:"type:uuid;index"`                                   // 既存の行があっても migrate できるよう NOT NULL にしない
}
