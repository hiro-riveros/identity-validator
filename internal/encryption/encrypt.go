package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func Encrypt(plainText string, secretKey []byte) (string, error) {
	if len(secretKey) != 16 && len(secretKey) != 24 && len(secretKey) != 32 {
		return "", fmt.Errorf("secret key must be 16, 24, or 32 bytes long")
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.Strict().EncodeToString(cipherText), nil
}

func Decrypt(cypherTextBase64 string, secretKey []byte) (string, error) {
	if len(secretKey) != 16 && len(secretKey) != 24 && len(secretKey) != 32 {
		return "", fmt.Errorf("secret key must be 16, 24, or 32 bytes long")
	}

	cipherText, err := base64.StdEncoding.DecodeString(cypherTextBase64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(cipherText) < nonceSize {
		return "", err
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
