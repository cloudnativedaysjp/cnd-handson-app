package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// NewHTTPHandler は discovery と JWKS を返す。Istio の RequestAuthentication はこの jwks_uri を参照する
func NewHTTPHandler(jwk map[string]string, issuer string) http.Handler {
	jwks := map[string]any{"keys": []map[string]string{jwk}}
	discovery := map[string]any{
		"issuer":                                issuer,
		"jwks_uri":                              strings.TrimSuffix(issuer, "/") + "/.well-known/jwks.json",
		"id_token_signing_alg_values_supported": []string{"RS256"},
		"response_types_supported":              []string{"token"},
		"subject_types_supported":               []string{"public"},
	}

	mux := http.NewServeMux()
	// パターン（"GET /path"）をそのままスパン名にする
	traced := func(pattern string, h http.HandlerFunc) { mux.Handle(pattern, otelhttp.NewHandler(h, pattern)) }
	traced("GET /.well-known/openid-configuration", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, discovery)
	})
	traced("GET /.well-known/jwks.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=300")
		writeJSON(w, jwks)
	})
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return mux
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
