package messages

import (
	"time"

	"github.com/info344-s17/challenges-sbang97/apiserver/models/users"
)

type ChannelID string

type Channel struct {
	ID          ChannelID      `json:"id" bson:"_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"createdAt"`
	CreatorID   users.UserID   `json:"userID"`
	Members     []users.UserID `json:"members"`
	Private     bool           `json:"private"`
}

type NewChannel struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Members     []users.UserID `json:"members"`
	Private     bool           `json:"private"`
}

type ChannelUpdates struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (newCh *NewChannel) ToChannel() (*Channel, error) {
	channel := &Channel{
		Name:        newCh.Name,
		Description: newCh.Description,
		CreatedAt:   time.Now(),
		Members:     newCh.Members,
		Private:     newCh.Private,
	}
	return channel, nil
}
