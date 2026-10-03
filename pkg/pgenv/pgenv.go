// Package pgenv は DB_* の環境変数から PostgreSQL の接続 URL を組み立てる
package pgenv

import (
	"fmt"
	"net"
	"net/url"
	"os"
)

// DSN は値をエンコードした接続 URL を返す。パスワードに空白や記号があっても壊れない
func DSN() (string, error) {
	var v [5]string
	for i, k := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_DB"} {
		if v[i] = os.Getenv(k); v[i] == "" {
			return "", fmt.Errorf("required environment variable %s is not set", k)
		}
	}
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(v[2], v[3]),
		Host:     net.JoinHostPort(v[0], v[1]),
		Path:     "/" + v[4],
		RawQuery: "sslmode=disable",
	}
	return u.String(), nil
}
