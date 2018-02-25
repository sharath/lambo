package intern

import (
	"errors"
	"github.com/sharath/lambo/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"fmt"
	"encoding/base64"
)

// User represents the MongoDB model for login/authentication
type User struct {
	ID        string    `json:"id" bson:"id"`
	Username  string    `json:"username" bson:"username"`
	Password  string    `json:"password" bson:"password"`
	AuthKeysD [5]string `json:"auth_key" json:"auth_key"`
}

func (u *User) getAuthKey(users *mgo.Collection) (string, error) {
	var err error
	payload := u.Username
	key, err := util.NewEncryptionKey()
	if err != nil {
		return "", err
	}
	enc, err := util.Encrypt(payload, key)
	if err != nil {
		return "", err
	}
	for i := len(u.AuthKeysD) - 1; i > 0; i-- {
		u.AuthKeysD[i] = u.AuthKeysD[i-1]
	}
	u.AuthKeysD[0] = base64.StdEncoding.EncodeToString(key)
	users.Update(bson.M{
		"id": u.ID,
	}, u)
	fmt.Println(u.AuthKeysD)
	return enc, err
}

func VerifyAuthKey(id string, enc string, users *mgo.Collection) (bool, error) {
	var user User
	var match bool

	err := users.Find(bson.M{"id": id}).One(&user)
	if err != nil {
		return match, errors.New("invalid id")
	}
	for _, key := range user.AuthKeysD {
		k, _ := base64.StdEncoding.DecodeString(key)
		decrypt, _ := util.Decrypt(enc, k)
		if decrypt == user.Password {
			match = true
		}
	}
	return match, err
}

func generateUserID(users *mgo.Collection) string {
	count, _ := users.Count()
	return strconv.Itoa(count + 1)
}

func validNewUsername(users *mgo.Collection, username string) bool {
	count, _ := users.Find(bson.M{"username": username}).Count()
	if count != 0 {
		return false
	}
	return true
}

func validPassword(password string) bool {
	return !(len(password) < 7)
}

// CreateUser makes a new user from a username and password and adds it to MongoDB
func CreateUser(username string, password string, users *mgo.Collection) (*User, error) {
	u := new(User)
	if !validNewUsername(users, username) {
		return u, errors.New("invalid username")
	}
	if !validPassword(password) {
		return u, errors.New("invalid password")
	}
	u.ID = generateUserID(users)
	u.Username = username
	u.Password = util.Hash(password)
	if password == "" {
		return u, errors.New("invalid password")
	}
	users.Insert(u)
	return u, nil
}

// AuthenticateUser checks a username/password to see if it's valid
func AuthenticateUser(username string, password string, users *mgo.Collection) (string, error) {
	var user User
	users.Find(bson.M{"username": username}).One(&user)
	if util.CompareHash(user.Password, password) {
		return user.getAuthKey(users)
	}
	return "", errors.New("invalid login")
}
