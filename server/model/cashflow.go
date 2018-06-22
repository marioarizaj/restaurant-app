package Model

import (
	"time"
	"github.com/marioarizaj/restaurant-app/config"
)

type CashFlow struct {
	Id int `json:"id"`
	StartDate time.Time `json:"startDate"`
	EndDate time.Time `json:"endDate"`
	Revenue float64 `json:"revenue"`
	ExpensesWage float64 `json:"expensesWage"`
	ExpensesSupplies float64 `json:"expensesSupplies"`
	Difference float64 `json:"difference"`
}

func RegisterCaashFlow(cashFlow CashFlow) (error,bool) {
	_,err := config.DbCon.Query("UPDATE cashflows SET end_date=$1,revenue=$2,expenses=$3 WHERE end_date IS NULL AND revenue IS NULL AND expenses IS NULL ",cashFlow.EndDate,cashFlow.Revenue,cashFlow.ExpensesWage+cashFlow.ExpensesSupplies)

	if err != nil {
		return err,false
	}
	return nil,true
}

func NewCashFlow() (error,bool){
	date := time.Now()
	_,err := config.DbCon.Query("INSERT INTO cashflows(start_date) VALUES ($1)",date)

	if err != nil {
		return err,false
	}

	return nil,true
}

func GetStartDate() (error,time.Time) {
	var tm time.Time

	stm := config.DbCon.QueryRow("SELECT start_date FROM cashflows WHERE end_date IS NULL")

	err := stm.Scan(&tm)

	if err != nil {
		return err,time.Time{}
	}
	return nil,tm
}

func AllCashFlows() (error, []CashFlow) {
	var cashflow CashFlow
	var cashflows []CashFlow
	row,err := config.DbCon.Query("SELECT id, start_date,end_date,revenue,expenses FROM cashflows WHERE end_date IS NOT NULL")

	if err != nil {
		return err,nil
	}
	defer row.Close()

	for row.Next() {
		err := row.Scan(&cashflow.Id,&cashflow.StartDate,&cashflow.EndDate,&cashflow.Revenue,&cashflow.ExpensesWage)

		if err != nil {
			return err,nil
		}

		cashflow.Difference = cashflow.Revenue - cashflow.ExpensesWage

		cashflows = append(cashflows,cashflow)
	}

	return nil,cashflows
}
