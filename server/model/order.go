package Model

import (
	"github.com/marioarizaj/restaurant-app/config"
	"errors"
	"time"
)

type OrderDetail struct {
	ProductId int `json:"productid"`
	Quantity int `json:"quantity"`
}

type Order struct {
	Id int `json:"id"`
	EmployeeId int `json:"employeeid"`
	OrderDetails []OrderDetail `json:"orderdetails"`
	TableId int `json:"table"`
	Date time.Time `json:"date"`
	TotalPrice float32 `json:"totalprice"`
}

func CreateOrder (order Order) (error,bool) {
	err,priceCheck := checkPrice(order.OrderDetails)
	if err != nil {
		return err,false
	}

	if priceCheck != order.TotalPrice {
		return errors.New("price mismatch"),false
	}

	date := time.Now()

	stm := config.DbCon.QueryRow("INSERT INTO orders(employeeid,time,tableid,totalprice) VALUES ($1,$2,$3,$4) RETURNING Id",order.EmployeeId,date,order.TableId,order.TotalPrice)

	err = stm.Scan(&order.Id)

	if err != nil {
		return err,false
	}

	err,created := enterDetails(order.Id,order.OrderDetails)

	if err != nil {
		return err,false
	}

	if !created {
		return errors.New("problem entering data in database"),false
	}

	err,substracted := decreaseInventory(order.OrderDetails)

	if err != nil {
		return nil,false
	}

	if !substracted {
		return errors.New("server error"),false
	}

	return nil,true
}

func decreaseInventory(details []OrderDetail) (error,bool) {
	stm , err := config.DbCon.Prepare("UPDATE inventory SET quantity = quantity - $1 WHERE productid = $2")
	if err != nil {
		return err,false
	}

	for i:=0;i<len(details);i++ {
		_,err := stm.Exec(details[i].Quantity,details[i].ProductId)
		if err != nil {
			return err,false
		}
	}
	return nil,true
}

func enterDetails(id int, detail []OrderDetail) (error,bool) {
	stm , err := config.DbCon.Prepare("INSERT INTO orderdetails(orderid,productid,quantity) VALUES ($1,$2,$3)")
	if err != nil {
		return err,false
	}
	for i:=0;i<len(detail);i++ {
		_,err := stm.Exec(id,detail[i].ProductId,detail[i].Quantity)
		if err != nil {
			return err,false
		}
	}

	return nil,true
}

func checkPrice(orderdt []OrderDetail) (error,float32) {
	var totprice float32
	var price float32
	var tempPrice float32

	for i:=0;i<len(orderdt);i++ {
		stmt := config.DbCon.QueryRow("SELECT sellingprice FROM products WHERE id = $1",orderdt[i].ProductId)
		err := stmt.Scan(&price)
		if err != nil {
			return err,0
		}
		tempPrice = float32(orderdt[i].Quantity) * price
		totprice = totprice + tempPrice
	}

	return nil,totprice
}

func GetRevenue(date time.Time) (error,float64) {
	var tempPrice float64
	var totprice float64
	res,err := config.DbCon.Query("SELECT totalprice FROM orders WHERE time > $1",date)
	if err != nil {
		return err,0
	}
	defer res.Close()
	for res.Next() {
		err = res.Scan(&tempPrice)
		if err != nil {
			return err,0
		}
		totprice += tempPrice
		if res.Err() != nil {
			return res.Err(),0
		}
	}
	return nil,totprice
}

func ShiftOrders(id int) (error,[]Order) {
	var order Order
	var orders []Order

	row,err := config.DbCon.Query("SELECT orders.time, orders.tableid,employees.id,order.id, orders.totalprice " +
																			"	FROM orders " +
																			"LEFT JOIN employees ON employee.id = orders.employeeid " +
																			"LEFT JOIN hoursworked ON employee.id = hoursworked.userid WHERE id = $1 AND hoursworked.cloak_out IS NULL ",id)

	if err != nil {
		return err,nil
	}

	defer row.Close()

	for row.Next() {
		err := row.Scan(&order.Date, &order.TableId,&order.EmployeeId,&order.Id, &order.TotalPrice)

		if err != nil {
			return err,nil
		}

		orders = append(orders,order)
	}

	return nil,orders
}

func AllOrders() (error,[]Order) {
	var order Order
	var orders []Order

	row,err := config.DbCon.Query("SELECT orders.time, orders.tableid, employees.id,orders.id, orders.totalprice FROM orders LEFT JOIN employees ON employees.id = orders.employeeid")

	if err != nil {
		return err,nil
	}

	defer row.Close()

	for row.Next() {
		err := row.Scan(&order.Date, &order.TableId,&order.EmployeeId,&order.Id, &order.TotalPrice)

		if err != nil {
			return err,nil
		}

		orders = append(orders,order)
	}

	return nil,orders
}