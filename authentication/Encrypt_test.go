package authentication

import (
	"testing"
)

func TestNewEncryptionKey(t *testing.T) {
	var keys [][]byte
	for i := 0; i < 10000; i++ {
		keys = append(keys, NewEncryptionKey())
	}
	for i := 0; i < 10000; i++ {
		for j, k := range keys {
			if i != j && string(k[:]) == string(keys[i][:]) {
				t.FailNow()
			}
		}
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := NewEncryptionKey()
	test := []byte("abcdefghijklmnopqrstuvwxyz")
	encrypted := Encrypt(test, key)
	decrypted := Decrypt(encrypted, key)
	if string(decrypted[:]) != string(test[:]) {
		t.FailNow()
	}
}
