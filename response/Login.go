package response

type Login struct {
	AuthKey string `json:"auth_key"`
}

func NewLogin(authKey string) *Login {
	l := new(Login)
	l.AuthKey = authKey
	return l
}
