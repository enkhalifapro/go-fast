package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type User struct {
	Id            bson.ObjectId `bson:"_id"`
	UserName      string        `form:"username" json:"username" binding:"required"`
	FirstName     string        `form:"firstName" json:"firstName" binding:"required"`
	LastName      string        `form:"lastName" json:"lastName" binding:"required"`
	Image         string        `form:"image" json:"image"`
	Slug          string
	Password      string        `form:"password" json:"password" binding:"required"`
	Email         string        `form:"email" json:"email" binding:"required"`
	EmailVerified bool
	VerifyToken   string
	RoleId        string `form:"roleId" json:"roleId" binding:"required"`
	Role          *Role `bson:"-"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	UpdaterId     string
}