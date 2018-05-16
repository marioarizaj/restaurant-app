package Model

import (
	"github.com/marioarizaj/restaurant-app/config"
	"errors"
	"database/sql"
	"fmt"
)

type Restaurant struct {
	Qkr string `json:"qkr" binding:"required"`
	Name string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
	TableNr int `json:"tableNr" binding:"required"`
	Phone string `json:"phonenr" binding:"required"`
}



func CheckIfExists () (error,Restaurant) {
	var qkr string
	var address string
	var phone string
	var tableNr int
	var name string

	err := config.DbCon.QueryRow("SELECT * FROM restorant").Scan(&qkr,&address,&tableNr,&phone,&name)

	if err == sql.ErrNoRows {
		return errors.New("no restaurant"),Restaurant{"","","",0,""}
	}
	if err != nil && err != sql.ErrNoRows {
		return errors.New("internal error"),Restaurant{"","","",0,""}
	}

	return nil , Restaurant{qkr,name,address,tableNr,phone}

}

func CreateRestaurant(res Restaurant) (error,bool) {
	fmt.Println("Entered Function")
	qkr := res.Qkr
	address := res.Address
	tableNr := res.TableNr
	phone := res.Phone


	_,err := config.DbCon.Query("INSERT INTO restorant(qkr,address,tablenr,phone,name) values ($1,$2,$3,$4,$5)",qkr,address,tableNr,phone,res.Phone)

	if err != nil {
		return errors.New("database error"),false
	}

	return nil,true
}
