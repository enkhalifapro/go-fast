package viewModels

import ()

type ChangePasswordViewModel struct {
	CurrentPassword string        `form:"currentPassword" json:"currentPassword" binding:"required"`
	NewPassword     string        `form:"newPassword" json:"newPassword" binding:"required"`
}