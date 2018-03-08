package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/marioarizaj/restaurant-app/config"
	"github.com/marioarizaj/restaurant-app/server/controller"
	"database/sql"
	"log"
)

var router *gin.Engine
func main() {
	var err error
	router = gin.Default()
	config.DbCon,err = config.NewDB("host=localhost port=5432  user=mario password=nil  dbname=restaurant sslmode=disable")

	if err != nil {
		fmt.Println("An error ocurred")
		return
	}



	router.POST("/register", controller.RegisterUser)

	router.Run(":8080")

	fmt.Println("success")

}

