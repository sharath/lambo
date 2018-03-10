package authentication

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(plaintext []byte) []byte {
	v, err := bcrypt.GenerateFromPassword(plaintext, 14)
	if err != nil {
		return nil
	}
	return v
}

func Compare(hash, unhashed []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, unhashed)
	if err != nil {
		return false
	}
	return true
}
