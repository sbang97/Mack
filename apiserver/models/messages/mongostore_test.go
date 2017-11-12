package messages

import (
	"fmt"
	"testing"

	"github.com/info344-s17/challenges-sbang97/apiserver/models/users"

	mgo "gopkg.in/mgo.v2"
)

func TestMongoStore(t *testing.T) {
	sess, err := mgo.Dial("192.168.99.100:27017")
	if err != nil {
		t.Fatalf("error dialing Mongo: %v", err)
	}
	defer sess.Close()
	uID1 := users.UserID("user1")
	uID2 := users.UserID("user2")
	uID3 := users.UserID("user3")
	members := []users.UserID{uID1, uID2, uID3}
	members2 := []users.UserID{uID2}
	store := &MongoStore{
		Session:      sess,
		DatabaseName: "test",
		Channels:     "channels",
		Messages:     "messages",
	}

	nc := &NewChannel{
		Name:        "TEST Channel",
		Description: "this channel is a channel used for testing purposes",
		Members:     members,
	}

	nc2 := &NewChannel{
		Name:        "TEST Channel2",
		Description: "this channel is a channel used for testing purposes2",
		Members:     members2,
	}

	channel, err := store.InsertChannel(nc, uID1)
	if err != nil {
		t.Errorf("error inserting channel: %v\n", err)
	}
	channel2, err := store.InsertChannel(nc2, uID2)
	if err != nil {
		t.Errorf("error inserting channel: %v\n", err)
	}
	fmt.Printf(channel.Description + "\n")
	fmt.Printf("userIDs=%v, %v, %v\n", channel.Members[0], channel.Members[1], channel.Members[2])

	fmt.Printf(channel2.Name + "\n")
	fmt.Printf("userIDs=%v\n", channel2.Members[0])
	channels, err := store.GetAllChannels(uID1)
	if err != nil {
		fmt.Printf("error " + err.Error())
	}
	fmt.Printf("There should be 2 channels, there are: %v channels \n", len(channels))
	nm := &NewMessage{
		ChannelID: channel.ID,
		Body:      "This is a test message, please disregard",
	}
	nm2 := &NewMessage{
		ChannelID: channel.ID,
		Body:      "This is another test message, please disregard",
	}
	msg1, _ := store.InsertMessage(nm, uID1)
	msg2, _ := store.InsertMessage(nm2, uID2)
	err = store.DeleteChannel(channel.ID)
	if err != nil {
		fmt.Printf("error " + err.Error())
	}
	channels, err = store.GetAllChannels(uID1)
	if err != nil {
		fmt.Printf("error " + err.Error())
	}
	fmt.Printf("There should be 1 channel, there is: %v channel \n", len(channels))
	err = store.DeleteMessage(msg1.ID)
	if err != nil {
		fmt.Printf("error " + err.Error())
	}
	msgs, err := store.GetMessages(channel.ID, 10)
	if err != nil {
		fmt.Printf("error " + err.Error())
	}
	fmt.Printf("DELETE: There should now be 1 msg, there is currnetly %v msg\n", len(msgs))
	err = store.DeleteMessage(msg2.ID)
	if err != nil {
		fmt.Printf("error " + err.Error())
	}
	msgs, err = store.GetMessages(channel.ID, 10)
	if err != nil {
		fmt.Printf("error " + err.Error())
	}
	for _, msg := range msgs {
		fmt.Printf(msg.Body + "\n")
	}

	fmt.Printf("DELETE: There should now be 0 msgs, there are currnetly %v msgs\n", len(msgs))
	sess.DB(store.DatabaseName).C(store.Channels).RemoveAll(nil)
	sess.DB(store.DatabaseName).C(store.Messages).RemoveAll(nil)
}
