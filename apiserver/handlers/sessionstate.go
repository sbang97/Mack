package handlers

import (
	"time"

	"github.com/info344-s17/challenges-sbang97/apiserver/models/users"
)

type SessionState struct {
	BeganAt    time.Time
	ClientAddr string
	User       *users.User
}
