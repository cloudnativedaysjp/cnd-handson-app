package handler_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/handler"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/token"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func get(t *testing.T, h http.Handler, path string, out any) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	if out != nil {
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), out))
	}
	return rec
}

func TestJWKSVerifiesIssuedToken(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	issuer, err := token.NewRS256Issuer(key, "http://idp.test/", "handson")
	require.NoError(t, err)
	h := handler.NewHTTPHandler(issuer.PublicKey(), issuer.KeyID(), "http://idp.test/")

	var disc struct {
		Issuer  string `json:"issuer"`
		JWKSURI string `json:"jwks_uri"`
	}
	assert.Equal(t, http.StatusOK, get(t, h, "/.well-known/openid-configuration", &disc).Code)
	assert.Equal(t, "http://idp.test/", disc.Issuer)
	assert.Equal(t, "http://idp.test/.well-known/jwks.json", disc.JWKSURI)

	var jwks struct {
		Keys []map[string]string `json:"keys"`
	}
	assert.Equal(t, http.StatusOK, get(t, h, "/.well-known/jwks.json", &jwks).Code)
	require.Len(t, jwks.Keys, 1)
	k := jwks.Keys[0]
	assert.Equal(t, map[string]string{"kty": "RSA", "use": "sig", "alg": "RS256", "kid": issuer.KeyID()},
		map[string]string{"kty": k["kty"], "use": k["use"], "alg": k["alg"], "kid": k["kid"]})

	n, err := base64.RawURLEncoding.DecodeString(k["n"])
	require.NoError(t, err)
	e, err := base64.RawURLEncoding.DecodeString(k["e"])
	require.NoError(t, err)
	pub := &rsa.PublicKey{N: new(big.Int).SetBytes(n), E: int(new(big.Int).SetBytes(e).Int64())}

	signed, _, err := issuer.Issue(&model.User{ID: uuid.New()})
	require.NoError(t, err)
	_, err = jwt.Parse(signed, func(*jwt.Token) (any, error) { return pub, nil }, jwt.WithValidMethods([]string{"RS256"}))
	assert.NoError(t, err)
}

func TestHealthz(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	h := handler.NewHTTPHandler(&key.PublicKey, "kid", "http://idp.test")
	assert.Equal(t, http.StatusOK, get(t, h, "/healthz", nil).Code)
	assert.Equal(t, http.StatusNotFound, get(t, h, "/nope", nil).Code)
}
