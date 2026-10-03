package db

import (
	"github.com/cloudnativedaysjp/cnd-handson-app/pkg/pgenv"
	"github.com/cloudnativedaysjp/cnd-handson-app/pkg/telemetry"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open() (*gorm.DB, error) {
	dsn, err := pgenv.DSN()
	if err != nil {
		return nil, err
	}
	sqlDB, err := telemetry.OpenPostgres(dsn)
	if err != nil {
		return nil, err
	}
	// TranslateError で一意制約違反を gorm.ErrDuplicatedKey として受け取る
	return gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{TranslateError: true})
}
