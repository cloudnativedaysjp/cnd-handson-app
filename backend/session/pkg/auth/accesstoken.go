package auth

import (
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// ユーザーIDからアクセストークンを生成する
func GenerateAccessToken(userID uuid.UUID) (token string, exp int64, err error) {
	expiresAt := time.Now().Add(time.Minute * 5).Unix()

	claims := jwt.MapClaims{
		"id":  userID,
		"exp": expiresAt,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := jwtToken.SignedString(secretKey)
	if err != nil {
		return "", int64(0), errors.Wrap(err, "failed to sign JWT")
	}
	return tokenString, expiresAt, nil
}

// アクセストークンを検証する
func ValidateAccessToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate token")
	}
	// トークンが有効か確認
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

// アクセストークンからユーザーIDを取得する
func GetUserIDFromToken(tokenString string) (uuid.UUID, error) {
	token, err := ValidateAccessToken(tokenString)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to validate token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, errors.New("invalid token claims")
	}

	userID, ok := claims["id"].(string)
	if !ok {
		return uuid.Nil, errors.New("user ID not found in token claims")
	}

	return uuid.Parse(userID)
}
