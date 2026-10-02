package db

import "gorm.io/gorm"

func Migrate(db *gorm.DB, models ...any) error {
	return db.AutoMigrate(models...)
}
