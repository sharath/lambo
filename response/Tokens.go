package response

import (
	"github.com/sharath/lambo/database"
	"gopkg.in/mgo.v2"
)

type Tokens []string

func NewTokens(entries *mgo.Collection) Tokens {
	var all []*database.MongoEntry
	entries.Find(nil).All(&all)
	var t Tokens
	if len(all) < 1 {
		return t
	}
	for _, i := range all[0].Tokens {
		t = append(t, i.Name)
	}
	return t
}
