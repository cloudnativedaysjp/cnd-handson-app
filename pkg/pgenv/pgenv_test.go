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
	t.Setenv("DB_PASSWORD", `two words @:/?#% it's \ ok`)
	t.Setenv("DB_DB", "my db")

	dsn, err := DSN()
	require.NoError(t, err)
	cfg, err := pgconn.ParseConfig(dsn)
	require.NoError(t, err)
	assert.Equal(t, "db", cfg.Host)
	assert.Equal(t, uint16(5432), cfg.Port)
	assert.Equal(t, "app user", cfg.User)
	assert.Equal(t, `two words @:/?#% it's \ ok`, cfg.Password)
	assert.Equal(t, "my db", cfg.Database)
}

func TestDSNKeepsUnixSocketAndMultipleHosts(t *testing.T) {
	t.Setenv("DB_USER", "u")
	t.Setenv("DB_PASSWORD", "p")
	t.Setenv("DB_DB", "d")

	t.Setenv("DB_HOST", "/var/run/postgresql")
	t.Setenv("DB_PORT", "5432")
	dsn, err := DSN()
	require.NoError(t, err)
	cfg, err := pgconn.ParseConfig(dsn)
	require.NoError(t, err)
	assert.Equal(t, "/var/run/postgresql", cfg.Host)

	t.Setenv("DB_HOST", "db1,db2")
	t.Setenv("DB_PORT", "5432,5433")
	dsn, err = DSN()
	require.NoError(t, err)
	cfg, err = pgconn.ParseConfig(dsn)
	require.NoError(t, err)
	assert.Equal(t, "db1", cfg.Host)
	assert.Equal(t, uint16(5432), cfg.Port)
	require.Len(t, cfg.Fallbacks, 1)
	assert.Equal(t, "db2", cfg.Fallbacks[0].Host)
	assert.Equal(t, uint16(5433), cfg.Fallbacks[0].Port)
}

func TestDSNRequiresAllVariables(t *testing.T) {
	t.Setenv("DB_HOST", "db")
	t.Setenv("DB_PORT", "")
	_, err := DSN()
	assert.ErrorContains(t, err, "DB_PORT")
}
