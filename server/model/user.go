package Model

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"github.com/marioarizaj/restaurant-app/config"
)

type User struct {
	FirstName string
	LastName  string
	Email     string
	Username  string
	Password  string
	Access    int
}

//Create User

func CreateUser(fn, ln, em, un, ps, ps2 string, access int,) (*[]string, error, string) {

	validateString(fn)
	validateString(ln)
	validateUsername(un)
	validateEmail(em)
	validatePassword(ps, ps2)

	if len(errors) > 0 {
		fmt.Println(errors)
		return &errors, nil, ""
	}

	err, id := create(fn, ln, un, em, ps, access)

	err, uid := generateSession(id)

	if err != nil {
		return nil, err, ""
	}

	return nil, err, uid

}

func generateSession(id int) (error, string) {
	uuidd, err := uuid.NewV4()

	uniqueValue := uuidd.String()

	fmt.Println(uniqueValue)

	query, err := config.DbCon.Prepare("INSERT INTO sessions(UUID,user_id) VALUES ($1,$2)")

	_,err = query.Exec(uniqueValue, id)

	if err != nil {
		fmt.Println("Error")
		fmt.Println(err)
		return err, ""
	}

	return err, uniqueValue

}



func create(fn, ln, un, em, ps string, access int) (error, int) {
	var id int
	err := config.DbCon.QueryRow("INSERT INTO users(first_name,last_name,email,username,password,age) VALUES ($1,$2,$3,$4,$5,$6) RETURNING Id", fn, ln, em, un, ps, access).Scan(&id)
	if err != nil {
		return err, id
	}

	if err != nil {
		return err, id
	}

	return nil, id

}
