package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var key []byte

func InitializeKey() error {
	err := godotenv.Load()
	if err != nil {
		return errors.New("failed to load .env file")
	}

	envKey := os.Getenv("ENCRYPTION_KEY")
	if len(envKey) != 16 {
		return errors.New("encryption key must be 16 bytes long")
	}

	key = []byte(envKey)
	return nil
}
func Encrypt(data string) string {
	block, _ := aes.NewCipher(key)
	nonce := make([]byte, 12)
	aesGCM, _ := cipher.NewGCM(block)
	encrypted := aesGCM.Seal(nil, nonce, []byte(data), nil)
	return base64.StdEncoding.EncodeToString(encrypted)
}

func Decrypt(data string) string {
	block, _ := aes.NewCipher(key)
	nonce := make([]byte, 12)
	aesGCM, _ := cipher.NewGCM(block)
	decoded, _ := base64.StdEncoding.DecodeString(data)
	decrypted, _ := aesGCM.Open(nil, nonce, decoded, nil)
	return string(decrypted)
}
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	return string(hashedPassword), nil
}
