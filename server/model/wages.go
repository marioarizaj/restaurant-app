package Model

import (
	"time"
	"github.com/marioarizaj/restaurant-app/config"
)

type Payments struct {
	Id int `json:"id"`
	EmployeeId int `json:"employeeid"`
	EmployeeName string `json:"employeename"`
	EmployeeSurname string `json:"surname"`
	HoursWorked float64 `json:"hrworked"`
	HourlyWage float64 `json:"hourlywage"`
	Payed float64 `json:"payed"`
	Bonus float64 `json:"bonus"`
	Date time.Time `json:"date"`
}

func GetPayments() (error,[]Payments) {
	var payment Payments
	var payments []Payments
	row,err := config.DbCon.Query("SELECT payments.id,employees.id,employees.firstname,employees.lastname,employees.hourlywage,payments.payment+payments.bonus,payments.date FROM payments LEFT JOIN employees ON employees.id = payments.employeeid")

	if err != nil {
		return err,nil
	}

	defer row.Close()

	for row.Next() {
		err := row.Scan(&payment.Id,&payment.EmployeeId,&payment.EmployeeName,&payment.EmployeeSurname,&payment.HourlyWage,&payment.Payed,&payment.Date)
		payment.HoursWorked = (payment.Payed - payment.Bonus)/payment.HourlyWage
		if err != nil {
			return err,nil
		}

		if row.Err() != nil {
			return row.Err(),nil
		}
		payments = append(payments,payment)
	}

	return nil,payments
}

func GetWages(date time.Time) (error,float64) {
	var tempPrice float64
	var totprice float64
	res,err := config.DbCon.Query("SELECT payment+bonus FROM payments WHERE date > $1",date)
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