package util

import "github.com/sharath/lambo/models/intern"
import (
	"gopkg.in/mgo.v2"
)

type Token struct {
	u intern.User
	key string
}

func NewToken(user intern.User, key *mgo.Collection) *Token {
	t := new(Token)
	t.u = user
	type k struct {
		Key string `json:"key" bson:"key"`
	}
	var found k
	key.Find(nil).One(&found)
	t.key = found.Key
	return t
}

func (t *Token) GenToken() string {
	tok := ""
	str := t.u.Username + t.u.Password
	// do some encryption from the t.key collection and return that
	return tok
}