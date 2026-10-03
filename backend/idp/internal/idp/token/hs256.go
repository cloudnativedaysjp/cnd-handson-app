package token

import (
	"errors"
	"time"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/model"
	jwt "github.com/golang-jwt/jwt/v5"
)

const accessTokenTTL = 15 * time.Minute

type HS256Issuer struct {
	key []byte
	now func() time.Time
}

func NewHS256Issuer(key []byte) (*HS256Issuer, error) {
	if len(key) == 0 {
		return nil, errors.New("signing key is empty")
	}
	return &HS256Issuer{key: key, now: time.Now}, nil
}

func (i *HS256Issuer) Issue(user *model.User) (string, time.Time, error) {
	exp := i.now().Add(accessTokenTTL)
	claims := jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		ExpiresAt: jwt.NewNumericDate(exp),
		IssuedAt:  jwt.NewNumericDate(i.now()),
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(i.key)
	return signed, exp, err
}
