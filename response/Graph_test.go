package response

import (
	"gopkg.in/mgo.v2"
	"testing"
)

func TestNewGraph(t *testing.T) {
	s, _ := mgo.Dial("localhost")
	g := NewGraph("Ethereum", s.DB("lambo").C("entries"))
	if g.Data == "" {
		t.FailNow()
	}
}
