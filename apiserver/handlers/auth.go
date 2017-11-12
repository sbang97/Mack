package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/info344-s17/challenges-sbang97/apiserver/models/users"
	"github.com/info344-s17/challenges-sbang97/apiserver/sessions"
)

func (ctx *Context) UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		newUser := &users.NewUser{}
		if r.Body == nil {
			http.Error(w, "Post must contain a user and cannot be empty", http.StatusBadRequest)
		}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(newUser); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if err := newUser.Validate(); err != nil {
			http.Error(w, "error validating user:"+err.Error(), http.StatusBadRequest)
			return
		}
		u, err := ctx.UserStore.GetByEmail(newUser.Email)
		if u != nil {
			http.Error(w, "user already exists", http.StatusBadRequest)
			return
		}
		u, err = ctx.UserStore.GetByUserName(newUser.UserName)
		if u != nil {
			http.Error(w, "user already exists", http.StatusBadRequest)
			return
		}
		user, err := ctx.UserStore.Insert(newUser)
		if err != nil {
			http.Error(w, "error inserting user:"+err.Error(), http.StatusInternalServerError)
			return
		}
		state := SessionState{
			BeganAt:    time.Now(),
			ClientAddr: r.RemoteAddr,
			User:       user,
		}
		_, err = sessions.BeginSession(ctx.SessionKey, ctx.SessionStore, state, w)
		if err != nil {
			http.Error(w, "error beginning session: "+err.Error(), http.StatusInternalServerError)
			return
		}
		Encode(w, contentTypeJSONUTF8, headerContentType, user)
	case "GET":
		users, err := ctx.UserStore.GetAll()
		if err != nil {
			http.Error(w, "error getting users", http.StatusInternalServerError)
			return
		}
		w.Header().Add(headerContentType, contentTypeJSONUTF8)
		encoder := json.NewEncoder(w)
		encoder.Encode(users)
	}
}

func (ctx *Context) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		newSess := &users.Credentials{}
		if err := decoder.Decode(newSess); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		user, err := ctx.UserStore.GetByEmail(newSess.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err = user.Authenticate(newSess.Password); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		state := SessionState{
			BeganAt:    time.Now(),
			ClientAddr: r.RemoteAddr,
			User:       user,
		}
		_, err = sessions.BeginSession(ctx.SessionKey, ctx.SessionStore, state, w)
		if err != nil {
			http.Error(w, "error beginning session", http.StatusInternalServerError)
		}

		Encode(w, contentTypeJSONUTF8, headerContentType, user)
	default:
		http.Error(w, "must be a POST method", http.StatusBadRequest)
		return
	}
}

func (ctx *Context) SessionsMineHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		sessions.EndSession(r, ctx.SessionKey, ctx.SessionStore)
	default:
		http.Error(w, "must be a DELETE method", http.StatusBadRequest)
		return
	}
}

func (ctx *Context) UsersMeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		state := &SessionState{}
		_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, state)
		if err != nil {
			http.Error(w, "failed to retrieve state", http.StatusBadRequest)
		} else {
			Encode(w, contentTypeJSONUTF8, headerContentType, state.User)
		}
	case "PATCH":
		decoder := json.NewDecoder(r.Body)
		update := &users.UserUpdates{}
		if err := decoder.Decode(update); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		state := &SessionState{}
		sID, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, state)
		if err != nil {
			http.Error(w, "error updating the user", http.StatusInternalServerError)
		}
		if err = ctx.UserStore.Update(update, state.User); err != nil {
			http.Error(w, "error updating the user", http.StatusInternalServerError)
		}
		err = ctx.SessionStore.Save(sID, state)
		if err != nil {
			http.Error(w, "error updating the user", http.StatusInternalServerError)
		}

	}
}

func Encode(w http.ResponseWriter, ctype string, hctype string, user *users.User) {
	w.Header().Add(headerContentType, contentTypeJSONUTF8)
	encoder := json.NewEncoder(w)
	encoder.Encode(user)
}
