package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"

	"strings"

	"github.com/enkhalifapro/go-fast/models"
	"github.com/enkhalifapro/go-fast/services"
)

type RoleController struct {
	userService services.IUserService
	roleService services.IRoleService
}

func NewRoleController(userService services.IUserService, roleService services.IRoleService) *RoleController {
	controller := RoleController{}
	controller.userService = userService
	controller.roleService = roleService
	return &controller
}

func (controller RoleController) GetAll(c *gin.Context) {
	err, roles := controller.roleService.Find(&bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

func (controller RoleController) GetByRoleName(c *gin.Context) {
	roleName := c.Param("roleName")
	err, role := controller.roleService.FindOne(&bson.M{"slug": roleName})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

func (r RoleController) CreateRole(c *gin.Context) {
	var role models.Role
	err := c.Bind(&role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// name exist
	err, _ = r.roleService.FindOne(&bson.M{"name": role.Name})
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role name '" + role.Name + "' is already exits"})
		return
	}
	authToken := c.Request.Header.Get("Authorization")
	authToken = strings.Replace(authToken, "Bearer ", "", -1)
	err, user := r.userService.CurrentUser(authToken)
	err = r.roleService.Insert(user.Id.Hex(), &role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"createdRole": role})
}

func (controller RoleController) UpdateByName(c *gin.Context) {
	roleName := c.Param("roleName")
	var role models.Role
	err := c.Bind(&role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authToken := c.Request.Header.Get("Authorization")
	authToken = strings.Replace(authToken, "Bearer ", "", -1)
	err, user := controller.userService.CurrentUser(authToken)
	err = controller.roleService.UpdateByName(user.Id.Hex(), roleName, &role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"updatedRole": role})
}

func (r RoleController) DeleteByName(c *gin.Context) {
	roleName := c.Param("roleName")
	err := r.roleService.DeleteByName(roleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deletedRole": roleName})
}
