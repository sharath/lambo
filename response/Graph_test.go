package response

import (
	"encoding/base64"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"testing"
)

func TestNewGraph(t *testing.T) {
	s, _ := mgo.Dial("localhost")
	g := NewGraph("Ethereum", s.DB("lambo").C("entries"))
	image, _ := base64.StdEncoding.DecodeString(g.Data)
	ioutil.WriteFile("test.png", image, 0777)
}
