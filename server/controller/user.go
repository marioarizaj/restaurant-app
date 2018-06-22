package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/marioarizaj/restaurant-app/server/model"
	"net/http"
	"github.com/gin-gonic/gin/binding"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

func RegisterUser(c *gin.Context) {
	var user Model.User
	if err:=c.MustBindWith(&user,binding.JSON);err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError,gin.H{
			"error":err.Error(),
		})
	}
	err,ex := Model.AdminExists()

	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if !ex {
		err,done := Model.NewCashFlow()
		if err != nil || !done  {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	user.Hiredate = time.Now().Format("2006-01-02")
	err, uuid := Model.CreateUser(user)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK , gin.H{
		"uuid" : uuid,
	})

}

func FirtsAdminExists(c *gin.Context) {
	err,exists := Model.AdminExists()

	if err != nil || !exists {
		c.JSON(http.StatusNoContent,gin.H{
			"message":"Admin does not exists",
		})
		return
	}

	c.JSON(200,gin.H{
		"message":"Admin Exists",
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
		fmt.Println(err)
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	err, uid , user := Model.CheckLogin(login.Username,login.Password)

	if err != nil {
		fmt.Println(err)
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

func IsClockedIn(c *gin.Context) {
	uuid := c.GetHeader("authorization")
	err , user := Model.CurrentUser(uuid)

	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	err,isClock := Model.IsClockedIn(user.Id)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200,isClock)
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
	fmt.Println("here")
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
	fmt.Println("id is :",id)
	intid,err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusNotAcceptable , gin.H{
			"message": err.Error(),
		})
		return
	}

	err,pay := Model.CalculateWage(intid)

	if err != nil {
		c.JSON(http.StatusNotAcceptable , gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"To be paid": pay,
	})

}

func GetPayments(c *gin.Context) {
	err,payments := Model.GetPayments()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{
			"message":"There has been a problem",
		})
	}
	c.JSON(200,payments)
}

func PayEmployee(c *gin.Context){
	type pay struct {
		Id int `json:"id"`
		Bonus float64 `json:"bonus"`
	}
	py := pay{}

	err := c.MustBindWith(&py ,binding.JSON)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message":err.Error(),
		})
		return
	}

	err,payid,payed := Model.PayEmployee(py.Id,py.Bonus)

	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message":err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"payment id": payid,
		"payed": payed,
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
		fmt.Println(err.Error())
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
		c.JSON(http.StatusNoContent,gin.H{
			"message":"No category registered",
		})
		return
	}

	c.JSON(http.StatusOK,categories)
}


func AddProduct(c *gin.Context) {
	var product Model.Product

	err := c.MustBindWith(&product,binding.JSON)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusNotAcceptable,gin.H{
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
		c.JSON(http.StatusNoContent,gin.H{
			"message":"No products registered",
		})
		return
	}

	c.JSON(http.StatusOK,products)
}

func AddInventory(c *gin.Context) {
	var inventory Model.Inventory

	err := c.MustBindWith(&inventory,binding.JSON)

	if err != nil {
		fmt.Println("error",err.Error())
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

	err,add := Model.AddSupplies(inventory.Product.Id,inventory.Quantity)

	if err != nil || !add {
		c.JSON(500,gin.H{
			"message":"There has been a problem with server",
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message":"Inventory successfully added",
	})
}

func RegisterCashFlow(c *gin.Context) {
	err,startdate := Model.GetStartDate()
	if err != nil {
		c.JSON(500,gin.H{
			"message":"Internal server error",
		})
		return
	}
	err,revenue := Model.GetRevenue(startdate)

	err,sup := Model.GetSupplies(startdate)

	err,wag := Model.GetWages(startdate)

	cashflow := Model.CashFlow{
		StartDate: startdate,
		EndDate: time.Now(),
		Revenue: revenue,
		ExpensesSupplies: sup,
		ExpensesWage: wag,
	}

	err,done := Model.RegisterCaashFlow(cashflow)

	if err != nil || !done {
		c.JSON(500,gin.H{
			"message":"Internal server error",
		})
		return
	}

	err,dn := Model.NewCashFlow()

	if err != nil || !dn {
		c.JSON(500,gin.H{
			"message":"Internal server error",
		})
		return
	}


	c.JSON(200,gin.H{
		"message":"Successfully registered message",
	})

}

func GetOrders(c *gin.Context) {
	uuid := c.GetHeader("authorization")
	err , user := Model.CurrentUser(uuid)

	if err != nil {
		c.JSON(500,gin.H{
			"message":err.Error(),
		})
		return
	}

	err,admin,_ := Model.IsAdmin(uuid)

	if admin {
		err,orders := Model.AllOrders()
		if err != nil && err != sql.ErrNoRows{
			fmt.Println(err.Error())
			c.JSON(500,gin.H{
				"message":err.Error(),
			})
			return
		}
		if err == sql.ErrNoRows {
			fmt.Println(err.Error())
			c.JSON(http.StatusNoContent,orders)
		}
		c.JSON(200,orders)
		return
	} else {
		err,orders := Model.ShiftOrders(user.Id)
		if err != nil && err != sql.ErrNoRows{
			fmt.Println(err.Error())
			c.JSON(500,gin.H{
				"message":err.Error(),
			})
			return
		}
		if err == sql.ErrNoRows {
			fmt.Println(err.Error())
			c.JSON(http.StatusNoContent,orders)
		}
		c.JSON(200,orders)
		return
	}

}

func ToBePaidEmployees(c *gin.Context) {
	err,toBePaid := Model.ToBePaid()

	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println(err.Error())
			c.JSON(500,gin.H{
				"message":err.Error(),
			})
			return
		}
	}

	if err == sql.ErrNoRows {
		c.JSON(204,toBePaid)
	}


	c.JSON(200,toBePaid)


}


func AllCashFlows(c *gin.Context) {
	err,allCashFlows := Model.AllCashFlows()

	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println(err.Error())
			c.JSON(500,gin.H{
				"message":err.Error(),
			})
			return
		}
	}


	c.JSON(200,allCashFlows)


}


func GetCurrentCashFlow(c *gin.Context) {
	err,startdate := Model.GetStartDate()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500,gin.H{
			"message":"Internal server error",
		})
		return
	}
	err,revenue := Model.GetRevenue(startdate)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500,gin.H{
			"message":"Internal server error",
		})
		return
	}

	err,sup := Model.GetSupplies(startdate)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500,gin.H{
			"message":"Internal server error",
		})
		return
	}

	err,wag := Model.GetWages(startdate)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500,gin.H{
			"message":"Internal server error",
		})
		return
	}

	cashflow := Model.CashFlow{
		StartDate: startdate,
		Revenue: revenue,
		ExpensesSupplies: sup,
		ExpensesWage: wag,
		Difference:revenue-(sup+wag),
	}

	c.JSON(200,cashflow)

}

func GetCloacks(c *gin.Context) {
	err,clocks := Model.ClockInAndOut()
	if err != nil && err != sql.ErrNoRows{
		fmt.Println(err.Error())
		c.JSON(500,gin.H{
			"message":"Internal server error",
		})
		return
	}

	if err == sql.ErrNoRows {
		c.JSON(200,clocks)
	}

	c.JSON(200,clocks)

}

func DeleteUser(c *gin.Context) {
	type idToDelete struct {
		id int
	}
	var toDelete idToDelete
	err := c.MustBindWith(&toDelete,binding.JSON)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500,gin.H{
			"message":"Internal server error",
		})
	}

	err,done := Model.DeleteUser(toDelete.id)

	if err != nil || !done {
		fmt.Println(err.Error())
		c.JSON(500,gin.H{
			"message":"Internal server error",
		})
	}

	fmt.Println(err.Error())
	c.JSON(http.StatusOK,gin.H{
		"message":"Deleted",
	})

}

func CreateOrder(c *gin.Context) {
	var order Model.Order

	err := c.MustBindWith(&order,binding.JSON)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusNotAcceptable,gin.H{
			"message":err.Error(),
		})
		return
	}

	err,added := Model.CreateOrder(order)

	if err != nil {
		fmt.Println(err.Error())
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


