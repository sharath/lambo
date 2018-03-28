package response

import (
	"gopkg.in/mgo.v2"
	"testing"
	"fmt"
)

func TestNewGraph(t *testing.T) {
	s, _ := mgo.Dial("localhost")
	g := NewGraph("Ethereum", s.DB("lambo").C("entries"))
	fmt.Println(g.Data)
}
