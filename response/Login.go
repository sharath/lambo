package response

// Login is the response format for the login endpoint
type Login struct {
	AuthKey string `json:"auth_key"`
}

// NewLogin returns filled response format
func NewLogin(authKey string) *Login {
	l := new(Login)
	l.AuthKey = authKey
	return l
}
