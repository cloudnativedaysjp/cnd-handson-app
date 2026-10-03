package pgenv

import (
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDSNKeepsSpecialCharacters(t *testing.T) {
	t.Setenv("DB_HOST", "db")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "app user")
	t.Setenv("DB_PASSWORD", "two words @:/?#%")
	t.Setenv("DB_DB", "my db")

	dsn, err := DSN()
	require.NoError(t, err)
	cfg, err := pgconn.ParseConfig(dsn)
	require.NoError(t, err)
	assert.Equal(t, "db", cfg.Host)
	assert.Equal(t, uint16(5432), cfg.Port)
	assert.Equal(t, "app user", cfg.User)
	assert.Equal(t, "two words @:/?#%", cfg.Password)
	assert.Equal(t, "my db", cfg.Database)
}

func TestDSNRequiresAllVariables(t *testing.T) {
	t.Setenv("DB_HOST", "db")
	t.Setenv("DB_PORT", "")
	_, err := DSN()
	assert.ErrorContains(t, err, "DB_PORT")
}
