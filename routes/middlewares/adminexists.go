package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/marioarizaj/restaurant-app/server/model"
	"net/http"
)

func AdminExists() gin.HandlerFunc {
	return func(c *gin.Context) {
		err, exists := Model.AdminExists()

		if err != nil && !exists {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
			"error":"Admin already exists",
		})
	}
}

