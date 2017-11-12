package messages

import "github.com/info344-s17/challenges-sbang97/apiserver/models/users"

//Store represents an abstract store for model.User objects.
//This interface is used by the HTTP handlers to insert new users,
//get users, and update users. This interface can be implemented
//for any persistent database you want (e.g., MongoDB, PostgreSQL, etc.)
type Store interface {
	//GetAllChannel returns the channels that the user is allowed to see
	GetAllChannels(uID users.UserID) ([]*Channel, error)

	//GetRecentMessages returns the last n messages
	GetMessages(chID ChannelID, n int) ([]*Message, error)

	//Delete deletes a channel as well as all messages posted to that channel
	DeleteChannel(chID ChannelID) error

	//AddUser Adds a user to a channel's Members list
	AddUser(uID users.UserID, channel *Channel) error

	//InsertChannel inserts a new channel
	InsertChannel(newChannel *NewChannel, uID users.UserID) (*Channel, error)

	//InsertMessage inserts a new message
	InsertMessage(newMessage *NewMessage, user *users.User) (*Message, error)

	//DeleteMessage removes a message
	DeleteMessage(mID MessageID) error

	//RemoveUser removes a user from a channel
	RemoveUser(uID users.UserID, channel *Channel) error

	//Update applies UserUpdates to the currentUser
	UpdateChannel(updates *ChannelUpdates, currentChannel *Channel) (*Channel, error)

	//UpdateMessage updates a message's body text
	UpdateMessage(updates *MessageUpdates, message *Message) (*Message, error)

	GetChannelByID(chID ChannelID) (*Channel, error)

	GetChannelByName(name string) (*Channel, error)

	GetMessageByID(msgID MessageID) (*Message, error)
}
