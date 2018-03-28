package authentication

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// NewEncryptionKey generates a random encryption key
func NewEncryptionKey() []byte {
	var key = make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		panic(err)
	}
	return key
}

// Encrypt encrypts an input using the key and returns the encrypted result
func Encrypt(plaintext []byte, key []byte) (ciphertext []byte) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil
	}

	return gcm.Seal(nonce, nonce, plaintext, nil)
}

// Decrypt decrypts an input using the key and returns the plaintext result
func Decrypt(ciphertext []byte, key []byte) (plaintext []byte) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil
	}

	p, err := gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
	if err != nil {
		return nil
	}
	return p
}
