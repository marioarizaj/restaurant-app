package Model

import (
	"time"
	"github.com/marioarizaj/restaurant-app/config"
)

type Supplies struct {
	Id int `json:"suppliesId"`
	ProductId int `json:"productsId"`
	Quantity int `json:"quantity"`
	TotPrice float64 `json:"totPrice"`
	DatePurchased time.Time `json:"datePurchased"`
}

func AddSupplies(productid int, quantity int)(error,bool) {
	date := time.Now()
	err,price := getProductBuyingPrice(productid)
	if err != nil {
		return err,false
	}
	totPrice := price * float64(quantity)
	_,err = config.DbCon.Query("INSERT INTO supplies(productid,quantity,date,totprice) VALUES ($1,$2,$3,$4)",productid,quantity,date,totPrice)
	if err != nil {
		return err,false
	}
	return nil,true
}

func getProductBuyingPrice(id int) (error,float64) {
	var buyingPrice float64
	query := config.DbCon.QueryRow("SELECT buyingprice FROM products WHERE id = $1",id)
	err := query.Scan(&buyingPrice)
	if err != nil {
		return err,0
	}
	return nil,buyingPrice

}


func GetSupplies(date time.Time) (error,float64) {
	var tempPrice float64
	var totprice float64
	res,err := config.DbCon.Query("SELECT totprice FROM supplies WHERE date > $1",date)
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
