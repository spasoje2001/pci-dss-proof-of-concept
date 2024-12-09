package utils

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Secret key for signing the token (you can store it in environment variables for better security)
var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// Claims - struktura koja predstavlja korisničke podatke u JWT-u
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// GenerateJWT - funkcija za generisanje JWT tokena
func GenerateJWT(userID int, username string, role string) (string, error) {
	// Set the expiration time for the token
	expirationTime := time.Now().Add(24 * time.Hour) // Token će biti validan 24 sata

	// Kreiramo korisničke podatke za JWT (u Claims)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // Postavljamo vreme isteka tokena
			Issuer:    "go-pci-dss",          // Možete koristiti ime svoje aplikacije ili nešto drugo
		},
	}

	// Kreiramo token sa korisničkim podacima i tajnim ključem
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Potpisujemo token sa tajnim ključem
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return signedToken, nil
}
func ValidateJWT(tokenString string) (*Claims, error) {
	// Parsiramo token i proveravamo njegovu validnost
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verifikacija algoritma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Ako je token validan, vraćamo claimove
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token claims")
}

// GetSecretKey - funkcija koja vraća tajni ključ (moguće je koristiti ovo za testiranje)
func GetSecretKey() []byte {
	return secretKey
}
