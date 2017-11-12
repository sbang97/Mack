package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/info344-s17/challenges-sbang97/apiserver/models/users"
	"github.com/info344-s17/challenges-sbang97/apiserver/sessions"
)

func TestUsersHandler(t *testing.T) {
	testUser := strings.NewReader(`{ 
		"id": "startingid",
		"email": "test@test.com",
		"password": "password",
		"passwordConf": "password",
		"userName": "tester",
		"firstName": "test",
		"lastName": "test"
	}`)
	sessStore := sessions.NewMemStore(time.Minute)
	userStore := users.NewMemStore()
	ctx := &Context{
		SessionKey:   "secret key",
		SessionStore: sessStore,
		UserStore:    userStore,
	}
	handler := http.HandlerFunc(ctx.UsersHandler)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/v1/users", testUser)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)
	user, err := ctx.UserStore.GetByEmail("test@test.com")
	if err != nil {
		t.Fatal("couldn't find user" + err.Error())
	}
	fmt.Printf("userName:" + user.UserName + "\n")
	req, err = http.NewRequest("GET", "/v1/users", nil)
	handler.ServeHTTP(rr, req)
	if err != nil {
		t.Fatal("failed getting the users:" + err.Error())
	}
}

func TestSessionsHandler(t *testing.T) {
	testCred := strings.NewReader(`{
		"email": "test@test.com",
		"password": "password",
	}`)
	ctx := &Context{
		SessionKey:   "secret key",
		SessionStore: sessions.NewMemStore(time.Minute),
		UserStore:    users.NewMemStore(),
	}
	handler := http.HandlerFunc(ctx.SessionsHandler)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/v1/sessions", testCred)
	if nil != err {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)
}

func SessionsMineHandler(t *testing.T) {
	ctx := &Context{
		SessionKey:   "secret key",
		SessionStore: sessions.NewMemStore(time.Minute),
		UserStore:    users.NewMemStore(),
	}
	handler := http.HandlerFunc(ctx.SessionsHandler)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", "/v1/sessions/mine", nil)
	if nil != err {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)
}
