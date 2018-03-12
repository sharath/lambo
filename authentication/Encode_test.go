package authentication

import (
	"testing"
)

func TestEncode(t *testing.T) {
	test1 := []byte("abcdefghijklmnopqrstuvwxyz")
	correct := []byte("YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXo*")
	result := Encode(test1)
	if string(result[:]) != string(correct[:]) {
		t.FailNow()
	}
}

func TestDecode(t *testing.T) {
	test1 := []byte("YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXo*")
	correct := []byte("abcdefghijklmnopqrstuvwxyz")
	result := Decode(test1)
	if string(result[:]) != string(correct[:]) {
		t.FailNow()
	}
}
