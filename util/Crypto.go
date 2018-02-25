package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"io"
)

// TOOD: fix UTF16(?) encoding for the strings

// Hash returns a hash from a string
func Hash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hash)
}

// CompareHash checks a hash and a string to see if they're the same
func CompareHash(hash, check string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(check))
	return err == nil
}

// NewEncryptionKey Generates a random encryption a key
func NewEncryptionKey() (string, error) {
	key := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key[:])
	return string(key[:]), err
}

// Encrypt encrypts plaintext using a key and returns base64 version
func Encrypt(input string, k string) (string, error) {
	plaintext := []byte(input)
	key := []byte(k)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}
	ciphertext := base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, plaintext, nil))
	return ciphertext, nil
}

// Decrypt decrypts a cipher using a key
func Decrypt(input string, k string) (plaintext string, err error) {
	ciphertext, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}
	key := []byte(k)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("malformed ciphertext")
	}

	arr, err := gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil)
	return string(arr[:]), err
}
