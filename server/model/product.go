package Model

import "github.com/marioarizaj/restaurant-app/config"

type Product struct {
	Id int `json:"id"`
	Name string `json:"productname"`
	SupplierId int `json:"supplierid"`
	SupplierName string `json:"suppliername"`
	CategoryId int `json:"categoryid"`
	CategoryName string `json:"categoryname"`
	BuyingPrice float32 `json:"buyingprice"`
	SellingPrice float32 `json:"sellingprice"`
}

func CreateProduct (product Product) (error,bool) {
	_,err := config.DbCon.Query("INSERT into products(name,supplierid,categoryid,sellingprice,buyingprice) VALUES ($1,$2,$3,$4,$5)",product.Name,product.SupplierId,product.CategoryId,product.SellingPrice,product.BuyingPrice)

	if err != nil {
		return err,false
	}

	return nil,true

}

func GetProducts () (error,[]Product) {
	var product Product
	var products []Product

	rows,err := config.DbCon.Query("SELECT products.id,products.name,supplier.companyname,categories.name,buyingprice,sellingprice " +
																				"FROM products" +
																				" LEFT JOIN categories ON (categories.id=products.categoryid)" +
																				" LEFT JOIN supplier ON (supplier.id = products.supplierid)")

	if err != nil {
		return err,nil
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&product.Id,&product.Name,&product.SupplierName,&product.CategoryName,&product.BuyingPrice,&product.SellingPrice)

		if err != nil {
			return err,nil
		}

		products = append(products, product)

		if err = rows.Err(); err != nil {
			return err,nil
		}

	}

	return nil,products

}
