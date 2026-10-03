package token_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"testing"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/token"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func encode(t *testing.T, typ string, der []byte) string {
	t.Helper()
	return base64.StdEncoding.EncodeToString(pem.EncodeToMemory(&pem.Block{Type: typ, Bytes: der}))
}

func TestParsePrivateKeyAcceptsPKCS1AndPKCS8(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	pkcs8, err := x509.MarshalPKCS8PrivateKey(key)
	require.NoError(t, err)

	for _, b64 := range []string{
		encode(t, "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(key)),
		encode(t, "PRIVATE KEY", pkcs8),
	} {
		got, err := token.ParsePrivateKey(b64)
		require.NoError(t, err)
		assert.True(t, key.Equal(got))
	}

	_, err = token.ParsePrivateKey("not base64!")
	assert.Error(t, err)
}

func TestIssueSignsRS256WithRegisteredClaims(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	issuer, err := token.NewRS256Issuer(key, "http://idp.test", "handson")
	require.NoError(t, err)
	user := &model.User{ID: uuid.New()}

	signed, exp, err := issuer.Issue(user, []string{"member"})
	require.NoError(t, err)

	claims := &struct {
		jwt.RegisteredClaims
		Roles []string `json:"roles"`
	}{}
	parsed, err := jwt.ParseWithClaims(signed, claims, func(*jwt.Token) (any, error) { return &key.PublicKey, nil },
		jwt.WithValidMethods([]string{"RS256"}), jwt.WithIssuer("http://idp.test"), jwt.WithAudience("handson"), jwt.WithExpirationRequired())
	require.NoError(t, err)
	assert.Equal(t, issuer.JWK()["kid"], parsed.Header["kid"])
	assert.NotEmpty(t, issuer.JWK()["kid"])
	assert.Equal(t, user.ID.String(), claims.Subject)
	assert.Equal(t, exp.Unix(), claims.ExpiresAt.Unix())
	assert.Equal(t, []string{"member"}, claims.Roles)
}

func TestIssueAlwaysIncludesRolesClaim(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	issuer, err := token.NewRS256Issuer(key, "http://idp.test", "handson")
	require.NoError(t, err)

	signed, _, err := issuer.Issue(&model.User{ID: uuid.New()}, nil)
	require.NoError(t, err)
	mc := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(signed, mc, func(*jwt.Token) (any, error) { return &key.PublicKey, nil })
	require.NoError(t, err)
	assert.Equal(t, []any{}, mc["roles"])
}

func TestNewRS256IssuerRequiresIssuerAndAudience(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	_, err = token.NewRS256Issuer(key, "", "handson")
	assert.Error(t, err)
	_, err = token.NewRS256Issuer(key, "http://idp.test", "")
	assert.Error(t, err)
}
