package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type PasswordResetToken struct {
	Id          bson.ObjectId `bson:"_id"`
	UserId      string
	Token       string
	IsValid     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UpdaterId   string
}