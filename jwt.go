package tools

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	Data interface{}
	jwt.StandardClaims
}

func GenerateToken(key string, data interface{}) (string, error) {
	expire := time.Now().Add(time.Hour * 24)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Data: data,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
		},
	})
	return token.SignedString([]byte(key))
}

func ValidToken(tokenStr string, key string) (interface{}, error) {
	var claim Claims
	token, err := jwt.ParseWithClaims(tokenStr, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return nil, errors.New("auth failure")
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Data, nil
	}
	return nil, errors.New("auth failure")
}
