package response

import (
	"testing"
	"gopkg.in/mgo.v2"
	"fmt"
)

func TestNewGraph(t *testing.T) {
	s, _ := mgo.Dial("localhost")
	g := NewGraph("Ethereum",s.DB("lambo").C("entries"))
	fmt.Println(g)
}