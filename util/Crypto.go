package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"io"
)

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

// Encrypt encrypts plaintext using a key
func Encrypt(plaintext []byte, key []byte) (ciphertext string, err error) {
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

	return string(gcm.Seal(nonce, nonce, plaintext, nil)[:]), nil
}

// Decrypt decrypts a cipher using a key
func Decrypt(ciphertext []byte, key []byte) (plaintext string, err error) {
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
		nil, )
	return string(arr[:]), err
}
