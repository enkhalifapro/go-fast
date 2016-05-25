package main

import (
	"log"

	"github.com/enkhalifapro/go-fast/controllers"
	"github.com/enkhalifapro/go-fast/services"
	"github.com/enkhalifapro/go-fast/utilities"
	"github.com/enkhalifapro/go-fast/utilities/security"
	"github.com/gin-gonic/gin"
)

func GetPort() string {
	confiUtil := utilities.NewConfigUtil()
	port := confiUtil.GetConfig("port")
	if port == "" {
		port = "2020"
		log.Println("[-] No PORT environment variable detected. Setting to ", port)
	}
	return ":" + port
}

func main() {
	// Creates a gin router with default middleware:
	router := gin.Default()

	// utilites helpers
	configUtil := utilities.NewConfigUtil()
	cryptUtil := utilities.NewCryptUtil()
	slugUtil := utilities.NewSlugUtil()

	// users APIs
	userService := services.NewUserService(configUtil, cryptUtil)
	usersController := controllers.NewUserController(cryptUtil, userService)
	router.GET("/", usersController.Root)
	router.POST("/api/v1/signup", usersController.CreateUser)
	router.POST("/api/v1/login", usersController.Login)
	router.POST("/api/v1/verify/email/:email", usersController.ResendVerifyEmail)
	router.GET("/api/v1/verify/email/:token", usersController.VerifyEmail)
	router.GET("/api/v1/users/me", usersController.CurrentUser)
	router.PUT("/api/v1/users/me", usersController.UpdateCurrentUser)
	router.PUT("/api/v1/password/change", usersController.ChangePassword)
	router.POST("/api/v1/password/reset", usersController.ResetPassword)
	router.POST("/api/v1/password/reset/:email", usersController.SendPasswordResetEmail)
	router.GET("/api/v1/author/:name/available", usersController.IsAvailableName)
	router.POST("/api/v1/logout", security.BasicUser, usersController.Logout)

	// roles APIs
	roleService := services.NewRoleService(configUtil, slugUtil)
	roleController := controllers.NewRoleController(userService, roleService)
	router.GET("/api/v1/role", roleController.GetAll)
	router.GET("/api/v1/role/name/:roleName", roleController.GetByRoleName)
	router.POST("/api/v1/role", security.BasicUser, roleController.CreateRole)
	router.PUT("/api/v1/role/name/:roleName", security.BasicUser, roleController.UpdateByName)
	router.DELETE("/api/v1/role/name/:roleName", security.BasicUser, roleController.DeleteByName)

	// Listen and server
	port := GetPort()
	log.Println("Port is " + port)
	router.Run(port)
}
