package intern

type User struct {
	ID       string   `json:"id" bson:"id"`
	Username string   `json:"username" bson:"username"`
	Password string   `json:"username" bson:"password"`
	AuthKey  [5]string `json:"auth_key" json:"auth_key"`
}
