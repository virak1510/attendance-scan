package pkg

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaim struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AuthUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func GenerateToken(data UserClaim, secret string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       data.ID,
		"username": data.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	secretKey := []byte(secret) // Replace with your own secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string, secret string) (*UserClaim, error) {
	claims := &UserClaim{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return claims, nil
}
