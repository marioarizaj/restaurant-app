package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/marioarizaj/restaurant-app/routes"
	"github.com/marioarizaj/restaurant-app/config"
	"github.com/marioarizaj/restaurant-app/routes/middlewares"
)

func main() {
	var err error
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())
	config.DbCon, err = config.NewDB("host=localhost port=5432  user=mario password=nil  dbname=restaurant sslmode=disable")

	if err != nil {
		fmt.Println("An error ocurred")
		return
	}

	router.Use(gin.Logger())

	routes.InitializeRoutes(router)

	router.Run(":8080")

}
