package intern

import (
	"errors"
	"github.com/sharath/lambo/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

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
func AuthenticateUser(username string, password string, users *mgo.Collection) (string, string, error) {
	var user User
	users.Find(bson.M{"username": username}).One(&user)
	if util.CompareHash(user.Password, password) {
		authKey, err := user.getAuthKey(users)
		return user.ID, authKey, err
	}
	return "", "", errors.New("invalid login")
}

// VerifyAuthKey returns whether a username authkey pair is valid
func VerifyAuthKey(user string, enc string, users *mgo.Collection) (bool, error) {
	var u User
	var match bool

	err := users.Find(bson.M{"id": user}).One(&u)
	if err != nil {
		return match, errors.New("invalid user")
	}
	for _, key := range u.AuthKeysD {
		k, _ := util.CookieCoding.DecodeString(key)
		decrypt, _ := util.Decrypt(enc, k)
		if decrypt == u.Username {
			match = true
		}
	}
	return match, err
}

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
	u.AuthKeysD[0] = util.CookieCoding.EncodeToString(key)
	users.Update(bson.M{
		"id": u.ID,
	}, u)
	return enc, err
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
