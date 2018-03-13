package authentication

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Username string
	Password []byte
	AuthKeys [][]byte
}

func UserExists(username string, users *mgo.Collection) bool {
	u := FetchUser(username, users)
	return u != nil
}

func CreateUser(username, password string, users *mgo.Collection) (*User, error) {
	if users == nil {
		return nil, errors.New("no collection provided")
	}
	if UserExists(username, users) {
		return nil, errors.New("user already exists")
	}
	u := new(User)
	u.Username = username
	u.Password = Hash([]byte(password))
	u.AuthKeys = make([][]byte, 5)
	users.Insert(u)
	return u, nil
}

func (u *User) Login(password string, users *mgo.Collection) string {
	if Compare(u.Password, []byte(password)) {
		payload := []byte(u.Username)
		encKey := NewEncryptionKey()
		authToken := Encrypt(payload, encKey)
		u.UpdateKeys(encKey, users)
		encoded := string(Encode(authToken)[:])
		return encoded
	}
	return ""
}

func (u *User) Authenticate(encoded string) bool {
	decoded := Decode([]byte(encoded))
	for i := 0; i < 5; i++ {
		if len(u.AuthKeys[i]) > 0 {
			decrypted := Decrypt(decoded, u.AuthKeys[i])
			if string(decrypted[:]) == u.Username {
				return true
			}
		}
	}
	return false
}

func (u *User) UpdateKeys(newKey []byte, users *mgo.Collection) {
	for i := 3; i >= 0; i-- {
		u.AuthKeys[i+1] = u.AuthKeys[i]
	}
	u.AuthKeys[0] = newKey
	qry := bson.M{"username": u.Username}
	update := bson.M{"$set": bson.M{"authkeys": u.AuthKeys}}
	users.Update(qry, update)
}

func FetchUser(username string, users *mgo.Collection) *User {
	qry := bson.M{"username": username}
	var u *User
	users.Find(qry).One(&u)
	return u
}

func FindUserByAuthKey(key string, users *mgo.Collection, matrix Matrix) *User {
	if matrix[key] != nil {
		return matrix[key]
	}
	var usersarr []*User
	users.Find(nil).All(&usersarr)
	for _, u := range usersarr {
		if u.Authenticate(key) {
			matrix[key] = u
			return u
		}
	}
	return nil
}
