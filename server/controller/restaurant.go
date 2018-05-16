package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/marioarizaj/restaurant-app/server/model"
	"net/http"
	"github.com/gin-gonic/gin/binding"
	"fmt"
)

func CheckRestaurant(c *gin.Context) {
	err,restaurant := Model.CheckIfExists()
	if err != nil {
		if err.Error() == "no restaurant" {
			c.JSON(http.StatusNoContent , gin.H{
				"restaurant": restaurant,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"restaurant": restaurant,
		})
		return
	}

	c.JSON(http.StatusOK , gin.H {
		"restaurant": restaurant,
	})
}

func CreateRestaurant(c *gin.Context){
	fmt.Println("Entered function")
	var restaurant Model.Restaurant
	if err := c.MustBindWith(&restaurant,binding.JSON); err!=nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusNonAuthoritativeInfo)
		return
	}
	err, created := Model.CreateRestaurant(restaurant)
	if err != nil {
		c.AbortWithStatus(500)
		return
	} else if !created {
		c.AbortWithStatus(500)
		return
	}
	c.JSON(201,restaurant)

}