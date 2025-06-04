package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWTHelper struct {
	secret []byte
}

func NewJWTHelper(secret []byte) *JWTHelper {
	return &JWTHelper{secret: secret}
}

func (h *JWTHelper) ParseToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return h.secret, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}
