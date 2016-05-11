package viewModels

import ()

type UpdateUserViewModel struct {
	UserName  string        `form:"username" json:"username" binding:"required"`
	FirstName string        `form:"firstName" json:"firstName" binding:"required"`
	LastName  string        `form:"lastName" json:"lastName" binding:"required"`
	Email     string        `form:"email" json:"email" binding:"required"`
	Image     string        `form:"image" json:"image"`
	RoleId    string        `form:"roleId" json:"roleId" binding:"required"`
}