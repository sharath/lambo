package authentication

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	test := []byte("abcdefghijklmnopqrstuvwxyz")
	bad := []byte("abcdefghijklmnorstvwxyz")
	if string(Hash(test)[:]) == string(Hash(test)[:]) {
		fmt.Println("same hash")
		t.FailNow()
	}
	if !Compare(Hash(test), test) {
		fmt.Println("correct pass isn't working")
		t.FailNow()
	}
	if Compare(Hash(test), bad) {
		fmt.Println("invalid pass is working")
		t.FailNow()
	}
}
