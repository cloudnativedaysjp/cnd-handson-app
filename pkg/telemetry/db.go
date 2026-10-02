package telemetry

import (
	"database/sql"

	"github.com/XSAM/otelsql"
	_ "github.com/jackc/pgx/v5/stdlib" // "pgx" ドライバ
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// OpenPostgres は otelsql で計装した *sql.DB を返す。クエリの引数はスパンに入らない
func OpenPostgres(dsn string) (*sql.DB, error) {
	attrs := otelsql.WithAttributes(semconv.DBSystemNamePostgreSQL)
	db, err := otelsql.Open("pgx", dsn, attrs)
	if err != nil {
		return nil, err
	}
	if _, err := otelsql.RegisterDBStatsMetrics(db, attrs); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
