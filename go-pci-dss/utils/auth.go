package utils

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pquerna/otp/totp"
)

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func GenerateJWT(userID int, username string, role string) (string, error) {
	var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "go-pci-dss",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Println("ovo je secret key ", secretKey)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return signedToken, nil
}
func ValidateJWT(tokenString string) (*Claims, error) {
	var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token claims")
}

// Funkcija za generisanje TOTP sekreta
func GenerateTOTPSecret() (string, error) {
	// Generišemo sekert za korisnika
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "go-pci-dss",
		AccountName: "user@example.com",
	})
	if err != nil {
		return "", err
	}
	return secret.Secret(), nil
}

// Funkcija za verifikaciju TOTP koda
func ValidateTOTP(secret string, otp string) (bool, error) {
	valid := totp.Validate(otp, secret)
	return valid, nil
}
