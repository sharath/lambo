package intern

import (
	"gopkg.in/mgo.v2"
	"strconv"
	"gopkg.in/mgo.v2/bson"
	"errors"
	"github.com/sharath/lambo/util"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string    `json:"id" bson:"id"`
	Username string    `json:"username" bson:"username"`
	Password string    `json:"password" bson:"password"`
	AuthKeys [5]string `json:"auth_key" json:"auth_key"`
}

func hash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hash)
}

func (u *User) getAuthKey(users *mgo.Collection) (string, error) {
	var err error
	payload := []byte(u.Username + u.Password)
	key, err := util.NewEncryptionKey()
	if err != nil {
		return "", err
	}
	enc, err := util.Encrypt(payload, key)
	if err != nil {
		return "", err
	}
	for i := range u.AuthKeys {
		if i != 0 {
			u.AuthKeys[i] = u.AuthKeys[i-1]
		}
	}
	u.AuthKeys[0] = string(key)
	users.Update(bson.M{
		"id": u.ID,
	}, bson.M{
		"id":        u.ID,
		"auth_keys": u.AuthKeys,
	})
	return string(enc), err
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
	u.Password = hash(password)
	if password == "" {
		return u, errors.New("invalid password")
	}
	users.Insert(u)
	return u, nil
}

func AuthenticateUser(username string, password string, users *mgo.Collection) (string, error) {
	var user User
	users.Find(bson.M{"username": username}).One(&user)
	if user.Password == hash(password) {
		return user.getAuthKey(users)
	} else {
		return "", errors.New("invalid login")
	}
}
