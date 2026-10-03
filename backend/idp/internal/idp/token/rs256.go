package token

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
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
	jwk      map[string]string
	issuer   string
	audience string
}

func NewRS256Issuer(key *rsa.PrivateKey, issuer, audience string) (*RS256Issuer, error) {
	if issuer == "" || audience == "" {
		return nil, errors.New("issuer and audience are required")
	}
	b64 := base64.RawURLEncoding.EncodeToString
	n := b64(key.N.Bytes())
	e := b64(big.NewInt(int64(key.E)).Bytes())
	// kid は RFC 7638 の JWK thumbprint（必須メンバーを辞書順に並べた JSON の SHA-256）
	sum := sha256.Sum256(fmt.Appendf(nil, `{"e":%q,"kty":"RSA","n":%q}`, e, n))
	jwk := map[string]string{"kty": "RSA", "use": "sig", "alg": "RS256", "kid": b64(sum[:]), "n": n, "e": e}
	return &RS256Issuer{key: key, jwk: jwk, issuer: issuer, audience: audience}, nil
}

func (i *RS256Issuer) Issue(user *model.User, roles []string) (string, time.Time, error) {
	now := time.Now()
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
	t.Header["kid"] = i.jwk["kid"]
	signed, err := t.SignedString(i.key)
	return signed, exp, err
}

// JWK は JWKS に載せる公開鍵
func (i *RS256Issuer) JWK() map[string]string { return i.jwk }

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
