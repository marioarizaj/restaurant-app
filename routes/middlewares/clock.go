package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/marioarizaj/restaurant-app/server/model"
	"fmt"
)

func IsCloackedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId := c.GetHeader("authorization")
		err,user := Model.CurrentUser(sessionId)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized , gin.H{
				"error":err.Error(),
			})
			return
		}
		err , isValid := Model.IsClockedIn(user.Id)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized , gin.H{
				"error":err.Error(),
			})
			return
		}
		if !isValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized , gin.H{
				"error":"Not clocked in",
			})
		}
		c.Next()
		return

	}

}

func IsCloackedOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId := c.GetHeader("authorization")
		err,user := Model.CurrentUser(sessionId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized , gin.H{
				"error":err.Error(),
			})
			return
		}
		err , isValid := Model.IsClockedIn(user.Id)
		if err != nil && err.Error() != "not clocked in" {
			c.AbortWithStatusJSON(http.StatusUnauthorized , gin.H{
				"error":err.Error(),
			})
			return
		}
		if isValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized , gin.H{
				"error":"Not clocked out",
			})
		}
		c.Next()
		return

	}

}
