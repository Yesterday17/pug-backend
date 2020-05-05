package auth

import (
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(key *rsa.PrivateKey, uuid string, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": uuid,
		"exp": time.Now().Add(duration).Unix(),
	})

	tokenStr, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func CancelToken(key *rsa.PrivateKey) (string, error) {
	return GenerateToken(key, "", 0)
}
