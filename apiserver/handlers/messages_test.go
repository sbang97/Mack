package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/info344-s17/challenges-sbang97/apiserver/models/messages"
	"github.com/info344-s17/challenges-sbang97/apiserver/models/users"
	"github.com/info344-s17/challenges-sbang97/apiserver/sessions"
)

func TestChannelsHandler(t *testing.T) {
	sess, err := mgo.Dial("192.168.99.100:27017")
	defer sess.Close()
	if err != nil {
		t.Fatalf("error dialing Mongo: %v", err)
	}
	store := &messages.MongoStore{
		Session:      sess,
		DatabaseName: "test",
		Channels:     "channels",
		Messages:     "messages",
	}
	sessStore := sessions.NewMemStore(time.Minute)
	userStore := users.NewMemStore()
	ctx := &Context{
		SessionKey:   "secret key",
		SessionStore: sessStore,
		UserStore:    userStore,
		MessageStore: store,
	}
	testChannel := strings.NewReader(`{ 
		"name": "test2",
		"description": "This is a test channel",
		"members": "user1"
	}`)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ctx.ChannelsHandler)
	req, err := http.NewRequest("POST", "/v1/channels", testChannel)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)
	channel, err := ctx.MessageStore.GetChannelByName("test2")
	if nil == channel {
		t.Fatal(err)
	}
	fmt.Printf(channel.Description)
}
