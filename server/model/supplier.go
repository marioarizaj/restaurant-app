package Model

import (
	"github.com/marioarizaj/restaurant-app/config"
)

type Supplier struct {
	Id int `json:"id"`
	Name string `json:"name" binding:"required"`
	Address string  `json:"address" binding:"required"`
	Phone string `json:"phone" binding:"required"`
}

func CreateSupplier(s Supplier) (error,bool) {
	_,err := config.DbCon.Query("INSERT INTO supplier(companyname,address,phone) values ($1,$2,$3)",s.Name,s.Address,s.Phone)
	if err != nil {
		return err,false
	}
	return nil,true
}

func GetSuppliers() (error,[]Supplier) {
	var supplier Supplier
	var suppliers []Supplier

	row,err := config.DbCon.Query("SELECT * FROM supplier")

	if err != nil {
		return err,nil
	}

	defer row.Close()

	for row.Next() {
		err := row.Scan(&supplier.Id,&supplier.Name,&supplier.Address,&supplier.Phone)

		if err != nil {
			return err,nil
		}

		if err = row.Err();err!=nil {
			return err, nil
		}

		suppliers = append(suppliers,supplier)

	}

	return nil,suppliers
}


