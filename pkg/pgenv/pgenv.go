// Package pgenv は DB_* の環境変数から PostgreSQL の接続文字列を組み立てる
package pgenv

import (
	"fmt"
	"os"
	"strings"
)

// 値を引用符で囲み、空白や記号があっても壊れないようにする。
// URL 形式にしないのは、Unix ソケットのパスや複数の host / port をそのまま渡せるようにするため
var quote = strings.NewReplacer(`\`, `\\`, `'`, `\'`)

// DSN は key=value 形式の接続文字列を返す
func DSN() (string, error) {
	var v [5]string
	for i, k := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_DB"} {
		if v[i] = os.Getenv(k); v[i] == "" {
			return "", fmt.Errorf("required environment variable %s is not set", k)
		}
		v[i] = "'" + quote.Replace(v[i]) + "'"
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", v[0], v[1], v[2], v[3], v[4]), nil
}
