package Model

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"github.com/marioarizaj/restaurant-app/config"
	"errors"
	"net/http"
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

	if len(errorsArray) > 0 {
		fmt.Println(errorsArray)
		return &errorsArray, nil, ""
	}

	err, id := create(fn, ln, un, em, ps, access)

	err, uid := generateSession(id)

	if err != nil {
		return nil, err, ""
	}

	return nil, err, uid

}

func CheckLogin(email , password string) (error,string,int) {
	var id int
	statement,err := config.DbCon.Prepare("SELECT id FROM users WHERE email = $1 AND password = $2")

	if err!=nil {
		return errors.New("Internal Error Occurred"), "",500
	}

		err = statement.QueryRow(email,password).Scan(&id)

	if err!=nil {
		return errors.New("Invalid username or password") , "",406
	}

	err,uid := generateSession(id)
	return nil,uid,200
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

func CheckToken(session string) (error,bool) {
	query , err := config.DbCon.Prepare("SELECT * FROM sessions WHERE uuid = $1")
	if err != nil {
		return err,false
	}

	res,err := query.Exec(session)
	if err != nil {
		return err,false
	}

	rows,err := res.RowsAffected()
	if err != nil {
		return err,false
	}
	if rows != 1 {
		return errors.New("Session expired"),false
	}

	return nil, true

}

func IsAdmin(session string) (error,bool,int) {
	var id int
	var access int
	query , err := config.DbCon.Prepare("SELECT id FROM sessions WHERE uuid = $1")
	if err != nil {
		return err,false,http.StatusInternalServerError
	}

	err = query.QueryRow(session).Scan(&id)

	if err != nil {
		return err,false,http.StatusInternalServerError
	}



	_ = config.DbCon.QueryRow("SELECT type FROM users WHERE id = $1",id).Scan(&access)
	if access != 1 {
		return errors.New("unauthorized"),false,http.StatusUnauthorized
	}


	return nil,true,http.StatusOK

}



func create(fn, ln, un, em, ps string, access int) (error, int) {
	var id int
	err := config.DbCon.QueryRow("INSERT INTO users(first_name,last_name,email,username,password,type) VALUES ($1,$2,$3,$4,$5,$6) RETURNING Id", fn, ln, em, un, ps, access).Scan(&id)
	if err != nil {
		return err, id
	}

	if err != nil {
		return err, id
	}

	return nil, id

}
