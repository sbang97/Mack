package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"
	redis "gopkg.in/redis.v5"

	"github.com/info344-s17/challenges-sbang97/apiserver/handlers"
	"github.com/info344-s17/challenges-sbang97/apiserver/middleware"
	"github.com/info344-s17/challenges-sbang97/apiserver/models/messages"
	"github.com/info344-s17/challenges-sbang97/apiserver/models/users"
	"github.com/info344-s17/challenges-sbang97/apiserver/notifier"
	"github.com/info344-s17/challenges-sbang97/apiserver/sessions"
)

const defaultPort = "443"

const (
	apiRoot    = "/v1/"
	apiSummary = apiRoot + "summary"
)

//main is the main entry point for this program
func main() {
	//read and use the following environment variables
	//when initializing and starting your web server
	// PORT - port number to listen on for HTTP requests (if not set, use defaultPort)
	// HOST - host address to respond to (if not set, leave empty, which means any host)
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	if len(port) == 0 {
		port = defaultPort
	}
	tlsCertPath := os.Getenv("TLSCERT")
	tlsKeyPath := os.Getenv("TLSKEY")
	sessKey := os.Getenv("SESSIONKEY")
	dbAddr := os.Getenv("DBADDR")
	msgAddr := os.Getenv("MESSAGEADDR")
	if len(dbAddr) == 0 {
		dbAddr = "192.168.99.100:27017"
	}
	if len(msgAddr) == 0 {
		msgAddr = "192.168.99.100:80"
	}
	sess, err := mgo.Dial(dbAddr)
	if err != nil {
		log.Fatal("error connecting to database", err.Error())
	}
	sess2, err := mgo.Dial(msgAddr)
	if err != nil {
		log.Fatal("error connecting to message database:" + err.Error())
	}
	messageStore := &messages.MongoStore{
		Session:      sess2,
		DatabaseName: "messages",
		Channels:     "channels",
		Messages:     "messages",
	}
	userStore := &users.MongoStore{
		Session:        sess,
		DatabaseName:   "apiserver",
		CollectionName: "users",
	}
	defer sess.Close()

	redisAddr := os.Getenv("REDISADDR")
	if len(redisAddr) == 0 {
		redisAddr = "192.168.99.100:6379"
	}
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	redisStore := sessions.NewRedisStore(client, time.Minute)
	ctx := handlers.Context{
		SessionKey:   sessKey,
		SessionStore: redisStore,
		UserStore:    userStore,
		MessageStore: messageStore,
		Notifier:     notifier.NewNotifier(),
	}
	go ctx.Notifier.Start()
	mux := http.NewServeMux()
	muxWithCors := http.NewServeMux()
	muxWithCors.HandleFunc("/v1/users", ctx.UsersHandler)
	muxWithCors.HandleFunc("/v1/sessions", ctx.SessionsHandler)
	muxWithCors.HandleFunc("/v1/sessions/mine", ctx.SessionsMineHandler)
	muxWithCors.HandleFunc("/v1/users/me", ctx.UsersMeHandler)
	muxWithCors.HandleFunc("/v1/channels", ctx.ChannelsHandler)
	muxWithCors.HandleFunc("/v1/channels/", ctx.SpecificChannelHandler)
	muxWithCors.HandleFunc("/v1/messages", ctx.MessagesHandler)
	muxWithCors.HandleFunc("/v1/messages/", ctx.SpecificMessageHandler)
	muxWithCors.HandleFunc(apiSummary, handlers.SummaryHandler)

	mux.Handle("/", middleware.Adapt(muxWithCors, middleware.CORS("", "", "", "")))
	mux.HandleFunc("/v1/websocket", ctx.WebSocketUpgradeHandler)
	mux.HandleFunc("/v1/chatbot", ctx.ChatBotHandler)

	http.Handle(apiRoot, mux)

	//add your handlers.SummaryHandler function as a handler
	//for the apiSummary route
	//HINT: https://golang.org/pkg/net/http/#HandleFunc

	//start your web server and use log.Fatal() to log
	//any errors that occur if the server can't start
	//HINT: https://golang.org/pkg/net/http/#ListenAndServe
	fmt.Printf("listening at %s...\n", host+":"+port)
	log.Fatal(http.ListenAndServeTLS(host+":"+port, tlsCertPath, tlsKeyPath, nil))
}
