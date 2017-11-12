package messages

import (
	"time"

	"github.com/info344-s17/challenges-sbang97/apiserver/models/users"
)

type MessageID string

type Message struct {
	ID        MessageID    `json:"id" bson:"_id"`
	ChannelID ChannelID    `json:"channelID"`
	Body      string       `json:"body"`
	CreatedAt time.Time    `json:"createdAt"`
	CreatorID users.UserID `json:"creatorID"`
	PhotoURL  string       `json:"photoURL"`
	UserName  string       `json:"username"`
	EditedAt  time.Time    `json:"editedAt"`
}

type NewMessage struct {
	ChannelID ChannelID `json:"channelID"`
	Body      string    `json:"body"`
}

type MessageUpdates struct {
	Body string
}

func (newMessage *NewMessage) ToMessage() (*Message, error) {
	msg := &Message{
		Body:      newMessage.Body,
		CreatedAt: time.Now(),
		ChannelID: newMessage.ChannelID,
	}
	return msg, nil
}
