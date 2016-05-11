package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Role struct {
	Id        bson.ObjectId `bson:"_id"`
	Name      string        `bson:"name" form:"name" json:"name"`
	Slug      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UpdaterId   string
}