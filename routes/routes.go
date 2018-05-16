package routes

import "github.com/gin-gonic/gin"
import (
	"github.com/marioarizaj/restaurant-app/routes/middlewares"
	"github.com/marioarizaj/restaurant-app/server/controller"
)


func InitializeRoutes(router *gin.Engine){
	router.GET("/restaurant",middlewares.IsGuest(),controller.CheckRestaurant)
	router.POST("/createRestaurant",middlewares.RestaurantExists(),controller.CreateRestaurant)
	router.POST("/firstAdmin",middlewares.IsGuest(),middlewares.AdminExists(),controller.RegisterUser)
	router.POST("/login", middlewares.IsGuest(),controller.CheckLogin)
	user := router.Group("/user",middlewares.IsLoggedIn(),middlewares.SessionValid())
	{
		user.GET("/getCategories",controller.GetCategories)
		user.GET("/currentUser",controller.CurrentUser)
		user.GET("/clockin",middlewares.IsCloackedOut(),controller.ClockIn)
		cloackedIn := router.Group("/clockedin",middlewares.IsCloackedIn())
		{
			cloackedIn.GET("/clockout",controller.ClockOut)
			cloackedIn.GET("/getInventory",controller.GetInventoryServer)
			cloackedIn.POST("/createOrder",controller.CreateOrder)
		}
		adminPrivileges := user.Group("/admin", middlewares.IsAdmin())
		{
			adminPrivileges.POST("/register", controller.RegisterUser)
			adminPrivileges.POST("/editInfo",controller.EditInfo)
			adminPrivileges.POST("/addSupplier",controller.AddSupplier)
			adminPrivileges.GET("/getUsers",controller.GetUsers)
			adminPrivileges.POST("/addCategory",controller.AddCategory)
			adminPrivileges.POST("/addProducts",controller.AddProduct)
			adminPrivileges.GET("/getProducts",controller.GetProducts)
			adminPrivileges.POST("/addInventory",controller.AddInventory)
			adminPrivileges.GET("/getInventory",controller.GetInventoryAdmin)
			adminPrivileges.GET("/getSuppliers",controller.GetSupplier)
			adminPrivileges.GET("/calculateWage",controller.CalculateWage)
		}
	}
}