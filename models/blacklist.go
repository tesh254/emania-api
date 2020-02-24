package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Blacklist defines token blacklist structure
type Blacklist struct {
	ID bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Token string `json:"token" bson:"token"`
}

// BlacklistModel defines balclist model
type BlacklistModel struct {}

// FindToken handle getting a single user
func (b *BlacklistModel) FindToken(token string) (doc Blacklist, err error) {
	collection := dbConnect.Use(databaseName, "blacklist")
	err = collection.Find(bson.M{"token": token}).One(&doc)

	return doc, err
}

// Add handle adding token to blacklist
func (b *BlacklistModel) Add(token string) (err error) {
	collection := dbConnect.Use(databaseName, "blacklist")
	err = collection.Insert(bson.M{
		"token": token,
	})

	return err
}