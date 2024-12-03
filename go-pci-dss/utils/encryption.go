package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

var key = []byte("examplekey123456") // 16 bytes key

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
