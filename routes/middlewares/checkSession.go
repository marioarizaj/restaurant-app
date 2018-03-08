package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"errors"
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
