package routes

import "github.com/gin-gonic/gin"
import (
	"github.com/marioarizaj/restaurant-app/routes/middlewares"
	"github.com/marioarizaj/restaurant-app/server/controller"
)


func InitializeRoutes(router *gin.Engine){
	user := router.Group("/user")
	{
		user.POST("/register", middlewares.IsGuest(), controller.RegisterUser)
		user.POST("/login", middlewares.IsGuest(),controller.checkLogin)
	}
}