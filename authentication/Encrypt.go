package authentication

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func NewEncryptionKey() []byte {
	var key = make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		panic(err)
	}
	return key
}

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
