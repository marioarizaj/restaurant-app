package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/marioarizaj/restaurant-app/server/model"
	"net/http"
	"strconv"
)

func RegisterUser(c *gin.Context) {
	first_name := c.PostForm("first_name")
	last_name := c.PostForm("last_name")
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	password2 := c.PostForm("password2")
	userType := c.PostForm("type")
	tp , err := strconv.Atoi(userType)
	errors, err, uuid := Model.CreateUser(first_name,last_name,email,username,password,password2,tp)
	if errors != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": errors,
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
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

	err , uid, status := Model.CheckLogin(email,password)

	if err != nil {
		c.JSON(status, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(status , gin.H {
		"message": "You are successfully logged in " ,
		"uuid": uid,
	})


}
