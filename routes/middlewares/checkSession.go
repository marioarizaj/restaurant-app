package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"errors"
	"github.com/marioarizaj/restaurant-app/server/model"
)

func IsGuest() gin.HandlerFunc {

	return func(c *gin.Context) {
		fmt.Println(c.GetHeader("authorization"),"here")
		if c.GetHeader("authorization") == "" {
			c.Next()
			return
		}
		c.AbortWithError(http.StatusUnauthorized, errors.New("not authorized"))
		c.JSON(http.StatusUnauthorized,"There has been an error")

	}
}

func IsLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId := c.GetHeader("authorization")

	if sessionId == "" {
		c.AbortWithError(http.StatusUnauthorized, errors.New("You are not logged in"))
		return
	}

	err , isLogged := Model.CheckToken(sessionId)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized , err)
	}
	if !isLogged {
		c.AbortWithError(http.StatusUnauthorized , errors.New("You are not logged in"))
	}
	c.Next()
	return

	}

}

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context){
		sessionId := c.GetHeader("authorization")
		err, isAdmin,status := Model.IsAdmin(sessionId)
		if err != nil {
			c.AbortWithError(status, err)
			return
		}
		if !isAdmin {
			c.AbortWithError(status , errors.New("You are not logged in"))
			return
		}
		c.Next()
		return

	}
}
