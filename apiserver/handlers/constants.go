package handlers

import (
	"github.com/info344-s17/challenges-sbang97/apiserver/models/messages"
	"github.com/info344-s17/challenges-sbang97/apiserver/models/users"
)

const (
	charsetUTF8         = "charset=utf-8"
	contentTypeJSON     = "application/json"
	contentTypeJSONUTF8 = contentTypeJSON + "; " + charsetUTF8
	headerContentType   = "Content-Type"
	newUserCreated      = "new user"
	newChannelCreated   = "new channel"
	channelUpdated      = "channel updated"
	channelDeleted      = "channel deleted"
	userJoinedChannel   = "joined channel"
	userLeftChannel     = "left channel"
	newMessage          = "new message"
	messageUpdated      = "message updated"
	messageDeleted      = "message deleted"
)

type UserEvent struct {
	Type string      `json:"type"`
	User *users.User `json:"user"`
}

type MessageEvent struct {
	Type    string            `json:"type"`
	Message *messages.Message `json:"message"`
}

type ChannelEvent struct {
	Type    string            `json:"type"`
	Channel *messages.Channel `json:"channel"`
}
