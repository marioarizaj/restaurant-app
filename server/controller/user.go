package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/marioarizaj/restaurant-app/server/model"
	"net/http"
)

func RegisterUser(c *gin.Context) {
	first_name := c.PostForm("first_name")
	last_name := c.PostForm("last_name")
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	password2 := c.PostForm("password2")
	access := 2

	errors, err, uuid := Model.CreateUser(first_name,last_name,email,username,password,password2,access)
	if errors != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": errors,
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": errors,
		})
		return
	}

	c.JSON(http.StatusOK , gin.H{
		"message" : uuid,
	})


}

func CheckLogin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")



}
