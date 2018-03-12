package authentication

import "github.com/sharath/lambo/database"

type Matrix map[string]*database.User

func NewAuthenticationMatrix() Matrix {
	m := make(Matrix)
	return m
}
