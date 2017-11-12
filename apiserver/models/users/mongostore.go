package users

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoStore struct {
	Session        *mgo.Session
	DatabaseName   string
	CollectionName string
}

func (ms *MongoStore) GetAll() ([]*User, error) {
	users := []*User{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Find(nil).All(&users)
	if err != nil {
		return nil, err
	}
	return users, err
}

func (ms *MongoStore) GetByID(id UserID) (*User, error) {
	user := &User{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).FindId(id).One(user)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (ms *MongoStore) GetByEmail(email string) (*User, error) {
	user := &User{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Find(bson.M{"email": email}).One(user)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (ms *MongoStore) GetByUserName(name string) (*User, error) {
	user := &User{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Find(bson.M{"username": name}).One(user)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (ms *MongoStore) Insert(newUser *NewUser) (*User, error) {
	user, err := newUser.ToUser()
	if err != nil {
		return nil, err
	}
	if nil == user {
		return nil, fmt.Errorf(".ToUser() returned nil")
	}
	newID := UserID(bson.NewObjectId().Hex())
	user.ID = newID
	err = ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Insert(user)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (ms *MongoStore) Update(update *UserUpdates, currentUser *User) error {
	col := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName)
	currentUser.FirstName = update.FirstName
	currentUser.LastName = update.LastName
	updates := bson.M{"$set": bson.M{"firstname": currentUser.FirstName, "lastname": currentUser.LastName}}
	return col.UpdateId(currentUser.ID, updates)
}
