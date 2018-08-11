package response

import (
	"github.com/globalsign/mgo"
	"testing"
)

func TestNewGraph(t *testing.T) {
	s, _ := mgo.Dial("localhost")
	g := NewGraph("Ethereum", s.DB("lambo").C("entries"))
	if g.Data == "" {
		t.FailNow()
	}
}
