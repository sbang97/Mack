package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/info344-s17/challenges-sbang97/apiserver/sessions"
)

func (ctx *Context) WebSocketUpgradeHandler(w http.ResponseWriter, r *http.Request) {
	state := &SessionState{}
	_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, state)
	if err != nil {
		http.Error(w, "failed to retrieve state", http.StatusBadRequest)
	}
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	ctx.Notifier.AddClient(conn)
}
