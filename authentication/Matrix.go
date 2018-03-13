package authentication

type Matrix map[string]*User

func NewAuthenticationMatrix() Matrix {
	m := make(Matrix)
	return m
}
