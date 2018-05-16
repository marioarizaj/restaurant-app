package Model

import "github.com/marioarizaj/restaurant-app/config"

type Category struct {
	Id int `json:"id"`
	Name string `json:"name" binding:"required"`
}

func AddCategory(cat Category) (error,bool) {
	_,err := config.DbCon.Query("INSERT INTO categories(name) VALUES ($1)",cat.Name)
	if err != nil {
		return err,false
	}
	return nil,true
}

func GetCategories() (error,[]Category) {
	rows,err := config.DbCon.Query("SELECT * FROM categories")
	var categories []Category
	var category Category
	if err != nil {
		return err,nil
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(&category.Id,&category.Name)
		categories = append(categories, category)
		if err = rows.Err(); err != nil {
			return err,nil
		}
	}

	return nil,categories
}