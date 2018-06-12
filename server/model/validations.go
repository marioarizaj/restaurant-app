package Model

import (
	"regexp"
	"github.com/marioarizaj/restaurant-app/config"
	"errors"
)


func validateString(str string) (error) {
	r := regexp.MustCompile("^[A-Z]{1}[a-z]{2,20}$")
	valid := r.MatchString(str)
	if !valid {
		return errors.New("name or surname is not valid")
	}
	return nil
}

func validateUsername(str string) (error){
	r := regexp.MustCompile("^[a-z0-9]{6,20}$")
	valid := r.MatchString(str)
	unExists, err := checkIfUnlExists(str)

	if err != nil {
		return errors.New("there has benn an error")
	}

	if unExists {
		return errors.New("username is taken")
	}

	if !valid {
		errors.New("username in not valid")
	}
	return nil
}

func validatePassword(pass string) (error) {
	r := regexp.MustCompile("^[a-z0-9]{8,20}$")

	valid := r.MatchString(pass)

	if !valid {
		return errors.New("password in not valid")
	}
	return nil
}


func checkIfUnlExists(un string) (bool, error) {
	query, err := config.DbCon.Prepare("SELECT * FROM employees WHERE username = $1")
	if err != nil {
		return true, err
	}
	res, err := query.Exec(un)
	if err != nil {
		return true, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return true, err
	}

	if rows > 0 {
		return true, nil
	}

	return false, nil

}
