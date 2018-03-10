package database

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"testing"
)

func TestUserExists(t *testing.T) {
	s, err := mgo.Dial("localhost")
	defer s.Close()
	if err != nil {
		t.FailNow()
	}
	users := s.DB("ot_test").C("users")
	users.DropCollection()
	u1, err := CreateUser("user1", "password", users)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	if !UserExists(u1.Username, users) {
		t.FailNow()
	}
}

func TestUser(t *testing.T) {
	s, err := mgo.Dial("localhost")
	defer s.Close()
	if err != nil {
		t.FailNow()
	}
	users := s.DB("ot_test").C("users")
	users.DropCollection()
	u, _ := CreateUser("ot_user", "password", users)
	token := u.Login("notpassword", users)
	if token != "" {
		fmt.Println("invalid password is working")
		t.FailNow()
	}
	token = u.Login("password", users)
	if token == "" {
		fmt.Println("password isn't working")
		t.FailNow()
	}
	if !u.Authenticate(token) {
		fmt.Println("authentication isn't working")
		t.FailNow()
	}
}

func TestUser_Authenticate(t *testing.T) {
	s, err := mgo.Dial("localhost")
	defer s.Close()
	if err != nil {
		t.FailNow()
	}
	users := s.DB("ot_test").C("users")
	users.DropCollection()
	u, _ := CreateUser("ot_user", "password", users)
	token1 := u.Login("password", users)
	if u.Authenticate("this shouldn't work") {
		fmt.Println("invalid auth key working")
		t.FailNow()
	}
	if !u.Authenticate(token1) {
		fmt.Println("authentication isn't working")
		t.FailNow()
	}
	if len(u.AuthKeys) != 5 {
		fmt.Println("authkey array not working properly")
		t.FailNow()
	}

	ut := FetchUser(u.Username, users)
	if ut == nil {
		fmt.Println("fetch user returning nil")
		t.FailNow()
	}
	if string(ut.AuthKeys[0][:]) == "" {
		fmt.Println("not storing last authkey")
		t.FailNow()
	}
	for i := 0; i < 10; i++ {
		u.Login("password", users)
	}
	ut = FetchUser(u.Username, users)
	if len(ut.AuthKeys) != 5 {
		fmt.Println("authkey array not working properly")
		t.FailNow()
	}
	for i := 0; i < 5; i++ {
		if string(ut.AuthKeys[i][:]) == "" {
			fmt.Println("all authkeys aren't filled")
			t.FailNow()
		}
	}
	if ut.Authenticate(token1) {
		fmt.Println("old token is working")
		t.FailNow()
	}
}
