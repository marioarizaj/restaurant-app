package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/marioarizaj/restaurant-app/server/model"
	"net/http"
	"github.com/gin-gonic/gin/binding"
	"database/sql"
	"fmt"
	"strconv"
)

func RegisterUser(c *gin.Context) {
	var user Model.User
	if err:=c.MustBindWith(&user,binding.JSON);err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,gin.H{
			"error":"Internal server error",
		})
	}
	err, uuid := Model.CreateUser(user)
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

func GetSupplier(c *gin.Context) {
	err, suppliers := Model.GetSuppliers()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if len(suppliers) == 0 {
		c.JSON(http.StatusNoContent,gin.H{
			"message":"No suppliers registered",
		})
		return
	}

	c.JSON(http.StatusOK, suppliers)

}

func CurrentUser(c *gin.Context) {
	uuid := c.GetHeader("authorization")

	err,user := Model.CurrentUser(uuid)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotAcceptable,gin.H{
			"message":err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK,user)
}
func AddSupplier(c *gin.Context) {
	var supplier Model.Supplier
	if c.MustBindWith(&supplier,binding.JSON) != nil {
		c.JSON(http.StatusNotAcceptable,
					gin.H{
						"message":"Error ocurred binding",
					})
		return
	}

	err, created := Model.CreateSupplier(supplier)
	if err != nil {
		c.JSON(http.StatusNotAcceptable,
			gin.H{
				"message":err,
			})
		return
	}
	if !created {
		c.JSON(http.StatusNotAcceptable,
			gin.H{
				"message":"Error ocurred",
			})
		return
	}

		c.JSON(http.StatusOK , gin.H {
		"message": "Supplier successfully registered",
	})

}

func GetUsers(c *gin.Context){
	err,users := Model.GetUsers()
	if err != nil {
		c.JSON(http.StatusBadGateway,gin.H{
			"message":"Problem retreiving users",
			"error":err.Error(),
		})
		return
	}
	if len(users) == 0 {
		c.JSON(http.StatusNotFound,gin.H{
			"message":"There are no users",
		})
		return
	}

	c.JSON(http.StatusOK,users)

}

func CheckLogin(c *gin.Context) {
	var login Model.LoginCmd
	err := c.MustBindWith(&login,binding.JSON)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	err, uid , user := Model.CheckLogin(login.Username,login.Password)

	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200 , gin.H {
		"message": "You are successfully logged in " ,
		"uuid": uid,
		"user":user,
	})

}

func EditInfo(c *gin.Context) {
	var userToEdit Model.User

	err := c.MustBindWith(&userToEdit,binding.JSON)

	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	err,edited := Model.EditUser(userToEdit)

	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	if !edited {
		c.JSON(500, gin.H{
			"message": "There has been a problem",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Information successfully edited",
	})
	return

}

func AddCategory(c *gin.Context) {
	var category Model.Category

	err := c.MustBindWith(&category,binding.JSON)

	if err != nil {
		c.JSON(http.StatusBadGateway,gin.H{
			"message":err.Error(),
		})
		return
	}

	err,done := Model.AddCategory(category)

	if err != nil {
		c.JSON(http.StatusBadGateway,gin.H{
			"message":err.Error(),
		})
		return
	}

	if !done {
		c.JSON(http.StatusBadGateway,gin.H{
			"message":"Problem proccessing data",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message":"Category successfully added",
	})
}

func ClockIn(c *gin.Context) {
	uuid := c.GetHeader("authorization")
	err , user := Model.CurrentUser(uuid)

	if err != nil {
		c.JSON(http.StatusNotAcceptable,gin.H{
			"message": err.Error(),
		})
		return
	}

	err,done := Model.ClockIn(user.Id)

	if err != nil {
		c.JSON(http.StatusNotAcceptable,gin.H{
			"message": err.Error(),
		})
		return
	}

	if !done {
		c.JSON(http.StatusNotAcceptable,gin.H{
			"message": "Problem with server",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message": "Successfully logged in",
	})
}

//TODO: PAYMENTS TABLE

func CalculateWage(c *gin.Context) {
	id := c.Query("id")
	intid,err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusNotAcceptable , gin.H{
			"message": "Id should be integer",
		})
		return
	}

	err,pay := Model.CalculateWage(intid)

	if err != nil {
		c.JSON(http.StatusNotAcceptable , gin.H{
			"message": "Id should be integer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"To be paid": pay,
	})

}

func GetInventoryServer(c *gin.Context) {
	err,inventory := Model.GetInventoryServer()
	if err != nil {
		c.JSON(http.StatusNotAcceptable,gin.H{
			"message": err.Error(),
		})
		return
	}

	if len(inventory)==0 || err == sql.ErrNoRows {
		c.JSON(http.StatusBadGateway,gin.H{
			"message":"No inventory registered",
		})
		return
	}

	c.JSON(http.StatusOK,inventory)
	return
}

func GetInventoryAdmin(c *gin.Context) {
	err,inventory := Model.GetInventoryAdmin()
	if err != nil {
		c.JSON(http.StatusNotAcceptable,gin.H{
			"message": err.Error(),
		})
		return
	}

	if len(inventory)==0 || err == sql.ErrNoRows {
		c.JSON(http.StatusBadGateway,gin.H{
			"message":"No inventory registered",
		})
		return
	}

	c.JSON(http.StatusOK,inventory)
	return
}

func ClockOut (c *gin.Context) {
	uuid := c.GetHeader("authorization")
	err , user := Model.CurrentUser(uuid)
	if err != nil {
		c.JSON(http.StatusNotAcceptable,gin.H{
			"message": err.Error(),
		})
		return
	}

	err,done := Model.ClockOut(user.Id)

	if err != nil {
		c.JSON(http.StatusNotAcceptable,gin.H{
			"message": err.Error(),
		})
		return
	}

	if !done {
		c.JSON(http.StatusNotAcceptable,gin.H{
			"message": "Problem with server",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message": "Successfully logged in",
	})

}

func GetCategories(c *gin.Context) {
	err,categories := Model.GetCategories()

	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusBadGateway,gin.H{
			"message":err.Error(),
		})
		return
	}

	if len(categories)==0 || err == sql.ErrNoRows {
		c.JSON(http.StatusBadGateway,gin.H{
			"message":"No category registered",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message":categories,
	})
}


func AddProduct(c *gin.Context) {
	var product Model.Product

	err := c.MustBindWith(&product,binding.JSON)

	if err != nil {
		c.JSON(http.StatusBadGateway,gin.H{
			"message":err.Error(),
		})
		return
	}

	err,done := Model.CreateProduct(product)

	if err != nil {
		c.JSON(http.StatusBadGateway,gin.H{
			"message":err.Error(),
		})
		return
	}

	if !done {
		c.JSON(http.StatusBadGateway,gin.H{
			"message":"Problem proccessing data",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message":"Product successfully added",
	})
}

func GetProducts(c *gin.Context) {
	err,products := Model.GetProducts()

	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusBadGateway,gin.H{
			"message":err.Error(),
		})
		return
	}

	if len(products)==0 || err == sql.ErrNoRows {
		c.JSON(http.StatusBadGateway,gin.H{
			"message":"No category registered",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message":products,
	})
}

func AddInventory(c *gin.Context) {
	var inventory Model.Inventory

	err := c.MustBindWith(&inventory,binding.JSON)

	if err != nil {
		c.JSON(http.StatusNotAcceptable,gin.H{
			"message":err.Error(),
		})
		return
	}

	err,added := Model.PurchaseProduct(inventory)

	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"message":err.Error(),
		})
		return
	}

	if !added {
		c.JSON(http.StatusInternalServerError,gin.H{
			"message":"Can not proccess data",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message":"Inventory successfully added",
	})
}

func CreateOrder(c *gin.Context) {
	var order Model.Order

	err := c.MustBindWith(&order,binding.JSON)

	if err != nil {
		c.JSON(http.StatusNotAcceptable,gin.H{
			"message":err.Error(),
		})
		return
	}

	err,added := Model.CreateOrder(order)

	if err != nil {
		c.JSON(http.StatusNotAcceptable,gin.H{
			"message":err.Error(),
		})
		return
	}

	if !added {
		c.JSON(http.StatusInternalServerError,gin.H{
			"message":"Problem in server",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message":"Order successfully registered",
	})

}


