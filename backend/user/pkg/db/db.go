package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// init 関数でデータベース接続を確立
func init() {
	// 初回接続の試行
	connectDB()
}

// connectDB はデータベースに接続する処理を担当
func connectDB() {
	env, err := requireEnv("DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_DB")
	if err != nil {
		log.Fatalf("Invalid database configuration: %v", err)
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		env["DB_HOST"], env["DB_USER"], env["DB_PASSWORD"], env["DB_DB"], env["DB_PORT"])

	log.Println("Connecting to database...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	DB = db
}

// InitDB は接続が必要な場合に再接続を試みる
func InitDB() (*gorm.DB, error) {
	log.Println("Checking database connection...")
	if DB == nil || !isDBConnected() {
		log.Println("Database connection lost. Retrying...")
		connectDB() // 再接続
	}

	// 接続状態を再確認
	if DB == nil || !isDBConnected() {
		return nil, fmt.Errorf("failed to reconnect to the database")
	}

	return DB, nil
}

// isDBConnected はデータベースが接続されているか確認する関数
func isDBConnected() bool {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Failed to get database connection: %v", err)
		return false
	}
	if err := sqlDB.Ping(); err != nil {
		return false
	}
	return true
}

// requireEnv は未設定の環境変数があればエラーを返す
func requireEnv(keys ...string) (map[string]string, error) {
	vals := make(map[string]string, len(keys))
	for _, k := range keys {
		v := os.Getenv(k)
		if v == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", k)
		}
		vals[k] = v
	}
	return vals, nil
}
