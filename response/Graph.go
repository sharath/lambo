package response

import (
	"github.com/sharath/lambo/visualization"
	"gopkg.in/mgo.v2"
)

// Graph is a wrapper for the graph url
type GraphResp struct {
	Data string `json:"data" bson:"data"`
}

// NewGraph makes a new graph and returns the url wrapper
func NewGraphResp(name, kind string, entries *mgo.Collection) *GraphResp {
	p := visualization.NewGraph(name, kind, entries)
	return &GraphResp{Data: p}
}
