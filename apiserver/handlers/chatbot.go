package handlers

import (
	"encoding/json"
	"net/http"

	"bytes"

	"github.com/info344-s17/challenges-sbang97/apiserver/sessions"
)

func (ctx *Context) ChatBotHandler(w http.ResponseWriter, r *http.Request) {
	state := &SessionState{}
	_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, state)
	if err != nil {
		http.Error(w, "failed to retrieve state:"+err.Error(), http.StatusBadRequest)
		return
	}
	user, err := json.Marshal(state.User)
	req, err := http.NewRequest("POST", "http://api.sbang9.me/chatbot", bytes.NewBuffer(user))
	req.Header.Set("User", string(user))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	w.Header().Add(headerContentType, contentTypeJSONUTF8)
	encoder := json.NewEncoder(w)
	encoder.Encode(resp)
}
