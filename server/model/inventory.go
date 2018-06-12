package Model

import (
	"github.com/marioarizaj/restaurant-app/config"
	"fmt"
)

type Inventory struct {
	Product
	Quantity int `json:"quantity" binding:"required"`
}

func PurchaseProduct(inventory Inventory) (error,bool) {
	err,exists := checkIfExists(inventory.Product.Id)
	if err != nil {
		return err,false
	}

	if exists {
		_,err := config.DbCon.Query("UPDATE inventory SET quantity = quantity + $1 WHERE productid = $2",inventory.Quantity,inventory.Id)
		if err != nil {
			return err,false
		}
		return nil,true
	}

	_,err = config.DbCon.Query("INSERT INTO inventory VALUES ($1,$2)",inventory.Id,inventory.Quantity)

	if err != nil {
		return err,false
	}

	return nil,true

}

func checkIfExists(productid int) (error,bool) {
	stm, err := config.DbCon.Prepare("SELECT * FROM inventory WHERE productid = $1")

	if err != nil {
		return err,false
	}

	row,err := stm.Exec(productid)

	if err != nil {
		return err,false
	}
	rows,err := row.RowsAffected()
	fmt.Println(rows)

	if err != nil {
		return err,false
	}

	if rows == 0 {
		return nil,false
	}

	return nil,true
}

func GetInventoryAdmin() (error,[]Inventory) {
	var inventory Inventory
	var inventories []Inventory

	rows,err := config.DbCon.Query("SELECT products.id,products.name,supplier.companyname,categories.name,products.buyingprice,products.sellingprice,inventory.quantity " +
																			"FROM inventory" +
																			" LEFT JOIN products ON products.id = inventory.productid" +
																			" LEFT JOIN categories ON categories.id=products.categoryid" +
																			" LEFT JOIN supplier ON supplier.id = products.supplierid")

	if err != nil {
		return err,nil
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&inventory.Id,&inventory.Name,&inventory.SupplierName,&inventory.CategoryName,&inventory.BuyingPrice,&inventory.SellingPrice,&inventory.Quantity)

		if err != nil {
			return err,nil
		}

		inventories = append(inventories,inventory)

		if err = rows.Err(); err != nil {
			return err,nil
		}

	}

	return nil,inventories

}

func GetInventoryServer() (error,[]Inventory) {
	var inventory Inventory
	var inventories []Inventory

	rows,err := config.DbCon.Query("SELECT products.id,products.name,categories.name,sellingprice,inventory.quantity " +
		"FROM products" +
		" LEFT JOIN inventory ON (products.id = inventory.productid) " +
		" LEFT JOIN categories ON (categories.id=products.categoryid)")

	if err != nil {
		return err,nil
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&inventory.Id,&inventory.Name,&inventory.CategoryName,&inventory.SellingPrice,&inventory.Quantity)

		if err != nil {
			return err,nil
		}

		inventories = append(inventories,inventory)

		if err = rows.Err(); err != nil {
			return err,nil
		}

	}

	return nil,inventories

}



