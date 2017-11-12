package handlers

import (
	"github.com/info344-s17/challenges-sbang97/apiserver/models/messages"
	"github.com/info344-s17/challenges-sbang97/apiserver/models/users"
	"github.com/info344-s17/challenges-sbang97/apiserver/notifier"
	"github.com/info344-s17/challenges-sbang97/apiserver/sessions"
)

//Context holds all the shared values that
//multiple HTTP Handlers will need
type Context struct {
	SessionKey   string
	SessionStore sessions.Store
	UserStore    users.Store
	MessageStore messages.Store
	Notifier     *notifier.Notifier
}
