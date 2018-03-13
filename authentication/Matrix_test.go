package authentication

import (
	"testing"
)

func TestNewAuthenticationMatrix(t *testing.T) {
	u := new(User)
	u.Username = "test"
	m := NewAuthenticationMatrix()
	m["test"] = u
	if m["test"].Username != u.Username {
		t.FailNow()
	}
}
