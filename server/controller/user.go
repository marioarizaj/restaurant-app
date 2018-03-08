package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/marioarizaj/restaurant-app/server/model"
)

func RegisterUser(c *gin.Context) {

	first_name := c.PostForm("first_name")
	last_name := c.PostForm("last_name")
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	password2 := c.PostForm("password2")
	access := 2

	Model.CreateUser(first_name,last_name,email,username,password,password2,access)


}
