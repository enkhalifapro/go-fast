package models

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

type Session struct {
	Id         bson.ObjectId `bson:"_id"`
	Token      string        `form:"userName" json:"userName"`
	ExpiryDate time.Time
	UserId     string
	User       *User
	CreatedAt  bson.MongoTimestamp
}