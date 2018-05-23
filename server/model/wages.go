package Model

import (
	"time"
	"github.com/marioarizaj/restaurant-app/config"
)

type Payments struct {
	Id int `json:"id"`
	EmployeeId int `json:"employeeid"`
	EmployeeName string `json:"employeename"`
	Payed string `json:"payed"`
	Date time.Time `json:"date"`
}

func GetPayments() (error,[]Payments) {
	var payment Payments
	var payments []Payments
	row,err := config.DbCon.Query("SELECT payments.id,employees.id,employees.name,payments.payment,payments.date FROM payments LEFT JOIN employees ON employees.id = payments.employeeid")

	if err != nil {
		return err,nil
	}

	defer row.Close()

	for row.Next() {
		err := row.Scan(&payment.Id,&payment.EmployeeId,&payment.EmployeeName,&payment.Payed,&payment.Date)
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

func GetPaymentsId(id int) (error,[]Payments) {
	var payment Payments
	var payments []Payments
	row,err := config.DbCon.Query("SELECT payments.id,employees.id,employees.name,payments.payment,payments.date FROM payments LEFT JOIN employees ON employees.id = payments.employeeid WHERE payments.employeeid = $1",id)

	if err != nil {
		return err,nil
	}

	defer row.Close()

	for row.Next() {
		err := row.Scan(&payment.Id,&payment.EmployeeId,&payment.EmployeeName,&payment.Payed,&payment.Date)
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
