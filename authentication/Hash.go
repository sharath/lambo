package authentication

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash bcrypt hashes a byte array and returns it
func Hash(plaintext []byte) []byte {
	v, err := bcrypt.GenerateFromPassword(plaintext, 10)
	if err != nil {
		return nil
	}
	return v
}

// Compare compares a hash and a password and returns bool
func Compare(hash, unhashed []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, unhashed)
	if err != nil {
		return false
	}
	return true
}
