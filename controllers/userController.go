package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/enkhalifapro/go-fast/models"
	"github.com/enkhalifapro/go-fast/services"
	"github.com/enkhalifapro/go-fast/utilities"
	"github.com/enkhalifapro/go-fast/viewModels"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

type UserController struct {
	userService services.IUserService
	cryptUtil   utilities.ICryptUtil
}

func NewUserController(cryptUtil utilities.ICryptUtil, userService services.IUserService) *UserController {
	controller := UserController{}
	controller.cryptUtil = cryptUtil
	controller.userService = userService
	return &controller
}

func (controller UserController) Root(c *gin.Context) {
	c.JSON(http.StatusOK, true)
}

func (controller UserController) CreateUser(c *gin.Context) {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// e-mail exist
	err, _ = controller.userService.FindOne(&bson.M{"email": user.Email})
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "E-mail address '" + user.Email + "' is already exits"})
		return
	}
	// name exist
	err, _ = controller.userService.FindOne(&bson.M{"username": user.UserName})
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username '" + user.UserName + "' is already exits"})
		return
	}
	passwordLength := len(user.Password)
	if passwordLength < 3 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Passord length should be 3 charachters atleast"})
		return
	}
	err = controller.userService.Insert(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": user.UserName, "email": user.Email})
}

func (controller UserController) ResendVerifyEmail(c *gin.Context) {
	email := c.Param("email")
	err, user := controller.userService.FindOne(&bson.M{"email": email})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	err = controller.userService.ResendVerifyEmail(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"VerifyEmailStatus": "ReSent", "verifyToken": user.VerifyToken})
}

func (controller UserController) Login(c *gin.Context) {
	var loginViewModel viewModels.LoginViewModel
	err := c.Bind(&loginViewModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var loginResult = controller.userService.Login(&loginViewModel)
	if loginResult == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user or password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"sessionToken": loginViewModel.Token, "userId": loginViewModel.UserId, "username": loginViewModel.UserName, "firstName": loginViewModel.FirstName, "lastName": loginViewModel.LastName, "image": loginViewModel.Image})
}

func (controller UserController) VerifyEmail(c *gin.Context) {
	token := c.Param("token")
	err := controller.userService.VerifyEmail(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"emailStatus": "verified"})
}

func (controller UserController) CurrentUser(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	authToken = strings.Replace(authToken, "Bearer ", "", -1)
	fmt.Println(authToken)
	err, user := controller.userService.CurrentUser(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": gin.H{"username": user.UserName, "email": user.Email, "image": user.Image, "firstName": user.FirstName, "lastName": user.LastName}})
}

func (controller UserController) UpdateCurrentUser(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	authToken = strings.Replace(authToken, "Bearer ", "", -1)
	fmt.Println(authToken)
	err, user := controller.userService.CurrentUser(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	var newUser viewModels.UpdateUserViewModel
	err = c.Bind(&newUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	err = controller.userService.UpdateById(user.Id.Hex(), user.Id.Hex(), &newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": gin.H{"username": newUser.UserName, "email": newUser.Email}})
}

func (controller UserController) ChangePassword(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	authToken = strings.Replace(authToken, "Bearer ", "", -1)
	fmt.Println(authToken)
	err, user := controller.userService.CurrentUser(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	var newPassword viewModels.ChangePasswordViewModel
	err = c.Bind(&newPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	hashedPassword := controller.cryptUtil.Bcrypt(newPassword.CurrentPassword)
	if hashedPassword != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid current password"})
		return
	}
	err = controller.userService.ChangePassword(user.Id.Hex(), user.Id.Hex(), newPassword.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "password updated"})
}

func (controller UserController) SendPasswordResetEmail(c *gin.Context) {
	email := c.Param("email")
	fmt.Println(email)
	err, user := controller.userService.FindOne(&bson.M{"email": email})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	err = controller.userService.SendPasswordResetEmail(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Status": "Sent"})
}

func (controller UserController) ResetPassword(c *gin.Context) {
	var resetPassword viewModels.ResetPasswordViewModel
	err := c.Bind(&resetPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	err, userId := controller.userService.ValidatePasswordResetToken(resetPassword.ResetToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	err = controller.userService.ChangePassword(userId, userId, resetPassword.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "password updated"})
}

func (controller UserController) IsAvailableName(c *gin.Context) {
	name := c.Param("name")
	// name exist
	err, _ := controller.userService.FindOne(&bson.M{"username": name})
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"available": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"available": true})
}

func (controller UserController) Logout(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	authToken = strings.Replace(authToken, "Bearer ", "", -1)
	fmt.Println(authToken)
	err := controller.userService.Logout(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"userStatus": "loggedout"})
}
