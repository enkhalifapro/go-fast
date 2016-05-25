package security

import (
	"net/http"
	"strings"

	"github.com/enkhalifapro/go-fast/services"
	"github.com/enkhalifapro/go-fast/utilities"
	"github.com/gin-gonic/gin"
)

func validUser(sessionToken string) bool {
	configUtil := utilities.NewConfigUtil()
	cryptUtil := utilities.NewCryptUtil()
	sessionService := services.NewSessionService(configUtil, cryptUtil)
	isValid := sessionService.Valid(sessionToken)
	return isValid
}

func BasicUser(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	authToken = strings.Replace(authToken, "Bearer ", "", -1)
	if validUser(authToken) == false {
		//c.AbortWithError(http.StatusInternalServerError, errors.New("Invalid user session"))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthoirzed user"})
		c.Abort()
		return
	}
	c.Next()
}
