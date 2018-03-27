package Model

import (
	"regexp"
	"github.com/marioarizaj/restaurant-app/config"
)

var errorsArray []string

func validateString(str string) {
	r := regexp.MustCompile("^[A-Z]{1}[a-z]{2,20}$")
	valid := r.MatchString(str)
	if !valid {
		errorsArray = append(errorsArray, "String in not valid")
	}
}

func validateUsername(str string) {
	r := regexp.MustCompile("^[a-z0-9]{6,20}$")
	valid := r.MatchString(str)
	emailExists, err := checkIfUnlExists(str)

	if err != nil {
		errorsArray = append(errorsArray, "There has benn an error")
	}

	if emailExists {
		errorsArray = append(errorsArray, "Username is taken")
	}

	if !valid {
		errorsArray = append(errorsArray, "Username in not valid")
	}
}

func validatePassword(pass, pass2 string) {
	if pass != pass2 {
		errorsArray = append(errorsArray, "Passwords do not match")
	}

	r := regexp.MustCompile("^[a-z0-9]{8,20}$")

	valid := r.MatchString(pass)

	if !valid {
		errorsArray = append(errorsArray, "Password in not valid")
	}

}

func validateEmail(email string) {
	emailExists, err := checkIfEmailExists(email)

	if err != nil {
		errorsArray = append(errorsArray, "There has benn an error")
	}

	if emailExists {
		errorsArray = append(errorsArray, "Email is taken")
	}

	r := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	valid := r.MatchString(email)

	if !valid {
		errorsArray = append(errorsArray, "Email is not valid")
	}

}

func checkIfEmailExists(email string) (bool, error) {
	query, err := config.DbCon.Prepare("SELECT * FROM users WHERE email = $1")
	if err != nil {
		return true, err
	}
	res, err := query.Exec(email)
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

func checkIfUnlExists(un string) (bool, error) {
	query, err := config.DbCon.Prepare("SELECT * FROM users WHERE username = $1")
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
