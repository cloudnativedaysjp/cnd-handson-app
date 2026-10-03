package token

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/model"
	jwt "github.com/golang-jwt/jwt/v5"
)

const accessTokenTTL = 15 * time.Minute

type claims struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

type RS256Issuer struct {
	key      *rsa.PrivateKey
	kid      string
	issuer   string
	audience string
	now      func() time.Time
}

func NewRS256Issuer(key *rsa.PrivateKey, issuer, audience string) (*RS256Issuer, error) {
	if issuer == "" || audience == "" {
		return nil, errors.New("issuer and audience are required")
	}
	kid, err := thumbprint(&key.PublicKey)
	if err != nil {
		return nil, err
	}
	return &RS256Issuer{key: key, kid: kid, issuer: issuer, audience: audience, now: time.Now}, nil
}

func (i *RS256Issuer) Issue(user *model.User, roles []string) (string, time.Time, error) {
	now := i.now()
	exp := now.Add(accessTokenTTL)
	if roles == nil {
		roles = []string{}
	}
	c := claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    i.issuer,
			Audience:  jwt.ClaimStrings{i.audience},
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Roles: roles,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	t.Header["kid"] = i.kid
	signed, err := t.SignedString(i.key)
	return signed, exp, err
}

func (i *RS256Issuer) PublicKey() *rsa.PublicKey { return &i.key.PublicKey }

func (i *RS256Issuer) KeyID() string { return i.kid }

// ParsePrivateKey は base64 でエンコードした PEM（PKCS#1 / PKCS#8）を読む
func ParsePrivateKey(b64 string) (*rsa.PrivateKey, error) {
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("signing key is not base64: %w", err)
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, errors.New("signing key is not PEM")
	}
	if k, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return k, nil
	}
	k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse signing key: %w", err)
	}
	rsaKey, ok := k.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("signing key is not RSA")
	}
	return rsaKey, nil
}

// kid は RFC 7638 の JWK thumbprint。鍵を差し替えると kid も変わる
func thumbprint(pub *rsa.PublicKey) (string, error) {
	b, err := json.Marshal(struct {
		E   string `json:"e"`
		Kty string `json:"kty"`
		N   string `json:"n"`
	}{
		E:   b64url(big.NewInt(int64(pub.E)).Bytes()),
		Kty: "RSA",
		N:   b64url(pub.N.Bytes()),
	})
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return b64url(sum[:]), nil
}

func b64url(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }
