package utils

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var privateKey = LoadEcdsaPrivateKeyKey()

// GenerateJWT create a new RegisteredClaim and sign with private key
func GenerateJWT(issuer string, userID uint) (string, error) {
	// Create JWT token
	claims := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   strconv.Itoa(int(userID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), // 1 hour
	})

	// Sign with private key
	t, err := claims.SignedString(privateKey)
	return t, err
}
