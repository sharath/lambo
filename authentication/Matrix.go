package authentication

// Matrix stores frequently used authkeys
type Matrix map[string]*User

// NewAuthenticationMatrix returns a new matrix
func NewAuthenticationMatrix() Matrix {
	m := make(Matrix)
	return m
}
