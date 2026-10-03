package db

import (
	"fmt"
	"os"

	"github.com/cloudnativedaysjp/cnd-handson-app/pkg/telemetry"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open() (*gorm.DB, error) {
	var v [5]string
	for i, k := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_DB"} {
		if v[i] = os.Getenv(k); v[i] == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", k)
		}
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", v[0], v[1], v[2], v[3], v[4])
	sqlDB, err := telemetry.OpenPostgres(dsn)
	if err != nil {
		return nil, err
	}
	return gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
}
