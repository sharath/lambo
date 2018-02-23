package intern

type User struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"username" bson:"password"`
}