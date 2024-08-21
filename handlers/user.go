package handlers

import (
	"loan-fund/models"
	"time"

	"github.com/golang-jwt/jwt"
)

var JWTKey = []byte("secret_key")

func GenerateJWT(customer models.User) (string, error) {
	claims := &jwt.MapClaims{
		"customer-id": customer.ID,
		"username":    customer.Username,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}