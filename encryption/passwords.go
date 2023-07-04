package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

const keySize = 32 // for AES-256

var additionalBytes = []byte("5BUptaMpkopamfjRe5mbSSNds+U0WbRl")

func encode(b []byte) string { return base64.StdEncoding.EncodeToString(b) }

func buildSecret(secret string) []byte {
	var s = []byte(secret)
	if len(s) < keySize {
		s = append(s, additionalBytes...)
	}
	s = s[:keySize]
	return s
}

// Encrypt method is to encrypt or hide any classified text
func Encrypt(text, secret string) (string, error) {
	block, err := aes.NewCipher(buildSecret(secret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text + encode([]byte(secret)))
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cipher.NewCFBEncrypter(block, iv).XORKeyStream(cipherText[aes.BlockSize:], plainText)
	return encode(cipherText), nil
}

func decode(s string) (data []byte, err error) {
	if data, err = base64.StdEncoding.DecodeString(s); err != nil {
		return
	}
	return
}

// Decrypt method is to extract back the encrypted text
func Decrypt(text, secret string) (string, error) {
	if text == "" {
		return "", nil
	}
	block, err := aes.NewCipher(buildSecret(secret))
	if err != nil {
		return "", err
	}
	var cipherText []byte
	if cipherText, err = decode(text); err != nil {
		return "", err
	}
	if len(cipherText) < aes.BlockSize {
		return "", errors.New("ciphertext block size is too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	plainText := make([]byte, len(cipherText))
	cipher.NewCFBDecrypter(block, iv).XORKeyStream(plainText, cipherText)
	return strings.TrimSuffix(string(plainText), encode([]byte(secret))), nil
}
