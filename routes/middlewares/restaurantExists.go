package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/marioarizaj/restaurant-app/server/model"
	"fmt"
)

func RestaurantExists() gin.HandlerFunc {
	return func(c *gin.Context) {
		err,_ := Model.CheckIfExists()

		if err != nil {
			if err.Error() == "no restaurant" {
				fmt.Println("Restaurant does not exists")
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError , gin.H{
				"error":"Internal server error",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest , gin.H{
			"error":"Restaurant already exists",
		})
		return
	}
}
