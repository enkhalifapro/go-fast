package viewModels

import ()

type ResetPasswordViewModel struct {
	ResetToken  string        `form:"resetToken" json:"resetToken" binding:"required"`
	NewPassword string        `form:"newPassword" json:"newPassword" binding:"required"`
}