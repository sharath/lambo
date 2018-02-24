package intern

import (
	"gopkg.in/mgo.v2"
	"strconv"
	"gopkg.in/mgo.v2/bson"
	"errors"
	"github.com/sharath/lambo/util"
)

type User struct {
	ID       string    `json:"id" bson:"id"`
	Username string    `json:"username" bson:"username"`
	Password string    `json:"password" bson:"password"`
	AuthKey  [5]string `json:"auth_key" json:"auth_key"`
}

func GenerateUserID(users *mgo.Collection) string {
	count, _ := users.Count()
	return strconv.Itoa(count + 1)
}

func ValidNewUsername(users *mgo.Collection, username string) bool {
	count, _ := users.Find(bson.M{"username": username}).Count()
	if count != 0 {
		return false
	}
	return true
}

func ValidPassword(password string) bool {
	return !(len(password) < 7)
}

func CreateUser(username string, password string, users *mgo.Collection) (*User, error) {
	u := new(User)
	if !ValidNewUsername(users, username) {
		return u, errors.New("invalid username")
	}
	if !ValidPassword(password) {
		return u, errors.New("invalid password")
	}

	u.ID = GenerateUserID(users)
	u.Username = username
	u.Password = password
	users.Insert(u)
	return u, nil
}

func AuthenticateUser(username string, password string, users *mgo.Collection) (string, error) {
	var user User
	users.Find(bson.M{"username": username}).One(&user)
	if user.Password == password {
		var err error
		payload := []byte(username + password)
		key, err := util.NewEncryptionKey()
		enc, err := util.Encrypt(payload, key)
		return string(enc), err
	} else {
		return "", errors.New("invalid login")
	}
}
