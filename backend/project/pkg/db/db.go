package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB          *gorm.DB
	initialized = false // データベース初期化状態を追跡
)

// init関数はパッケージのインポート時に一度だけ実行される
func init() {
	// テスト環境の場合は初期化をスキップ
	if isTestEnvironment() {
		log.Println("Test environment detected, skipping PostgreSQL connection")
		return
	}

	// PostgreSQL接続がテスト以外の環境で必要な場合の初期化
	if !initialized {
		initializeDB()
	}
}

// isTestEnvironment はテスト環境かどうかを判定する
func isTestEnvironment() bool {
	// testing.Tが使われているかを確認（テスト実行時の標準的な環境変数）
	if os.Getenv("GOTEST") == "1" {
		return true
	}

	// コマンドライン引数や実行ファイル名からテスト実行を検出
	executable := os.Args[0]
	if strings.HasSuffix(executable, ".test") || strings.Contains(executable, "/_test/") {
		return true
	}

	// テストモードの引数を検出
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-test.") {
			return true
		}
	}

	return false
}

// initializeDB はデータベース接続を初期化する
func initializeDB() {
	// 既に初期化済みならスキップ
	if initialized {
		return
	}

	// .envファイルを読み込む
	loadEnvFile()

	// 初回接続の試行
	connectDB()

	// 初期化状態を更新
	initialized = true
}

// loadEnvFile は.envファイルを読み込む
func loadEnvFile() {
	// カレントディレクトリから親ディレクトリへ順番に.envファイルを検索
	dir, err := os.Getwd()
	if err != nil {
		log.Printf("Failed to get current directory: %v", err)
		return
	}

	// 最大3階層まで遡ってenvファイルを探す
	for i := 0; i < 3; i++ {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			if err := godotenv.Load(envPath); err == nil {
				log.Printf("Loaded .env file from: %s", envPath)
				return
			}
		}
		// 親ディレクトリへ
		dir = filepath.Dir(dir)
	}

	log.Println("Warning: .env file not found")
}

// InitPostgresDB はPostgreSQLデータベースに接続する処理（interface.goから参照されるため公開）
func InitPostgresDB() (*gorm.DB, error) {
	dsn, err := buildDSN()
	if err != nil {
		return nil, err
	}

	log.Println("Connecting to PostgreSQL database...")
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// buildDSN は必須の環境変数から接続文字列を組み立てる。未設定ならエラー
func buildDSN() (string, error) {
	var v [5]string
	for i, k := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_DB"} {
		if v[i] = os.Getenv(k); v[i] == "" {
			return "", fmt.Errorf("required environment variable %s is not set", k)
		}
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", v[0], v[1], v[2], v[3], v[4]), nil
}

// connectDB はデータベースに接続する処理を担当
func connectDB() {
	dsn, err := buildDSN()
	if err != nil {
		log.Printf("Invalid database configuration: %v", err)
		return
	}

	log.Println("Connecting to PostgreSQL database...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// エラーログを出力するだけでパニックは発生させない
		log.Printf("Failed to connect to the database: %v", err)
		return
	}
	log.Println("Successfully connected to database")
	DB = db
}

// InitDB は接続が必要な場合に再接続を試みる
func InitDB() (*gorm.DB, error) {
	// テスト環境ではDB接続を行わない
	if isTestEnvironment() {
		return DB, nil
	}

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

// SetDB は外部からDBインスタンスを設定するためのヘルパー関数
func SetDB(db *gorm.DB) {
	DB = db
	initialized = true
}

// isDBConnected はデータベースが接続されているか確認する関数
func isDBConnected() bool {
	// DBがnilの場合は接続されていない
	if DB == nil {
		return false
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Failed to get database connection: %v", err)
		return false
	}
	if err := sqlDB.Ping(); err != nil {
		log.Printf("Failed to ping database: %v", err)
		return false
	}
	return true
}
