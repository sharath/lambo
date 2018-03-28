package response

import (
	"gopkg.in/mgo.v2"
	"github.com/sharath/lambo/database"
)

type Graph struct {
	Data string `json:"data" bson:"data"`
}

func NewGraph(name string, entries *mgo.Collection) *Graph {
	var all []*database.MongoEntry
	g := new(Graph)
	entries.Find(nil).All(&all)
	// todo
	return g
}
