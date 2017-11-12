package messages

import (
	"fmt"

	"github.com/info344-s17/challenges-sbang97/apiserver/models/users"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoStore struct {
	Session      *mgo.Session
	DatabaseName string
	Channels     string
	Messages     string
}

// GetAllChannels gets all the channels a user has access to or is a member of
func (ms *MongoStore) GetAllChannels(uID users.UserID) ([]*Channel, error) {
	channels := []*Channel{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.Channels).Find(bson.M{"$or": []bson.M{bson.M{"members": uID}, bson.M{"private": false}}}).All(&channels)
	if err != nil {
		return nil, err
	}
	return channels, err
}

//GetRecentMessages returns the last n messages
func (ms *MongoStore) GetMessages(chID ChannelID, n int) ([]*Message, error) {
	messages := []*Message{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.Messages).Find(bson.M{"channelid": chID}).Limit(n).All(&messages)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

//Delete deletes a channel as well as all messages posted to that channel
func (ms *MongoStore) DeleteChannel(chID ChannelID) error {
	err := ms.Session.DB(ms.DatabaseName).C(ms.Channels).RemoveId(chID)
	err2 := ms.Session.DB(ms.DatabaseName).C(ms.Channels).Remove(bson.M{"channelID": chID})
	if err != nil || err2 != nil {
		return err
	}
	return nil
}

//AddUser Adds a user to a channel's Members list
func (ms *MongoStore) AddUser(uID users.UserID, channel *Channel) error {
	col := ms.Session.DB(ms.DatabaseName).C(ms.Channels)
	channel.Members = append(channel.Members, uID)
	update := bson.M{"$set": bson.M{"members": channel.Members}}
	return col.UpdateId(channel.ID, update)
}

//InsertChannel inserts a new channel
func (ms *MongoStore) InsertChannel(newChannel *NewChannel, uID users.UserID) (*Channel, error) {
	channel, err := newChannel.ToChannel()
	if err != nil {
		return nil, fmt.Errorf("ToChannel() returned error:" + err.Error())
	}
	channel.CreatorID = uID
	newID := ChannelID(bson.NewObjectId().Hex())
	channel.ID = newID
	err = ms.Session.DB(ms.DatabaseName).C(ms.Channels).Insert(channel)
	if err != nil {
		return nil, err
	}
	return channel, err
}

//InsertMessage inserts a new message
func (ms *MongoStore) InsertMessage(newMessage *NewMessage, user *users.User) (*Message, error) {
	message, err := newMessage.ToMessage()
	if err != nil {
		return nil, fmt.Errorf("ToMessage() returned error:" + err.Error())
	}
	message.CreatorID = user.ID
	message.UserName = user.UserName
	message.PhotoURL = user.PhotoURL
	newID := MessageID(bson.NewObjectId().Hex())
	message.ID = newID
	err = ms.Session.DB(ms.DatabaseName).C(ms.Messages).Insert(message)
	if err != nil {
		return nil, err
	}
	return message, err
}

//DeleteMessage deletes a message
func (ms *MongoStore) DeleteMessage(mID MessageID) error {
	err := ms.Session.DB(ms.DatabaseName).C(ms.Messages).RemoveId(mID)
	if err != nil {
		return err
	}
	return nil
}

//RemoveUser removes a user from a channel
func (ms *MongoStore) RemoveUser(uID users.UserID, channel *Channel) error {
	col := ms.Session.DB(ms.DatabaseName).C(ms.Channels)
	members := channel.Members
	for i, member := range members {
		if member == uID {
			if i < len(members)-1 {
				members = append(members[:i], members[i+1:]...)
			} else {
				members = members[:i]
			}
		}
	}
	updates := bson.M{"$set": bson.M{"members": members}}
	err := col.UpdateId(channel.ID, updates)
	if err != nil {
		return err
	}
	return nil
}

//Update applies UserUpdates to the currentUser
func (ms *MongoStore) UpdateChannel(update *ChannelUpdates, currentChannel *Channel) (*Channel, error) {
	col := ms.Session.DB(ms.DatabaseName).C(ms.Channels)
	currentChannel.Name = update.Name
	currentChannel.Description = update.Description
	updates := bson.M{"$set": bson.M{"name": currentChannel.Name, "description": currentChannel.Description}}
	err := col.UpdateId(currentChannel.ID, updates)
	if err != nil {
		return nil, err
	}
	return currentChannel, nil
}

//UpdateMessage updates a message's body text
func (ms *MongoStore) UpdateMessage(update *MessageUpdates, currentMessage *Message) (*Message, error) {
	col := ms.Session.DB(ms.DatabaseName).C(ms.Messages)
	currentMessage.Body = update.Body
	updates := bson.M{"$set": bson.M{"body": currentMessage.Body}}
	err := col.UpdateId(currentMessage.ID, updates)
	if err != nil {
		return nil, err
	}
	return currentMessage, nil
}

func (ms *MongoStore) GetChannelByName(name string) (*Channel, error) {
	channel := &Channel{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.Channels).Find(bson.M{"name": name}).One(channel)
	if err != nil {
		return nil, err
	}
	return channel, err
}

func (ms *MongoStore) GetChannelByID(chID ChannelID) (*Channel, error) {
	channel := &Channel{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.Channels).FindId(chID).One(channel)
	if err != nil {
		return nil, err
	}
	return channel, err
}

func (ms *MongoStore) GetMessageByID(msgID MessageID) (*Message, error) {
	msg := &Message{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.Messages).FindId(msgID).One(msg)
	if err != nil {
		return nil, err
	}
	fmt.Printf(string(msg.ID))
	return msg, err
}
