package Model

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"github.com/marioarizaj/restaurant-app/config"
	"errors"
	"net/http"
	"database/sql"
	"time"
)

type User struct {
	Id int `json:"id"`
	FirstName string `json:"name" binding:"required"`
	LastName  string  `json:"surname" binding:"required"`
	Username  string `json:"username"`
	Password  string `json:"password" binding:"required"`
	Type    int `json:"type"`
	Birthday string `json:"birthdate" binding:"required"`
	Hiredate string `json:"hiredate" binding:"required"`
	Address string `json:"address" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	Notes string `json:"notes"`
	HourlyWage float32 `json:"wage"`
}

type LoginCmd struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

//Create User

func AdminExists() (error,bool) {
	err := config.DbCon.QueryRow("SELECT * FROM employees").Scan()
	if err != nil {
		if err == sql.ErrNoRows {
			return err,false
		}
	}
	return nil,true

}

func CalculateWage (id int) (error,float64) {
	var wage float64
	stm := config.DbCon.QueryRow("SELECT hourlywage FROM employees WHERE id = $1",id)

	err := stm.Scan(&wage)

	if err != nil {
		return err,0
	}

	err,hours := calculateHours(id)

	if err != nil {
		return err,0
	}

	pay := hours * wage

	return nil,pay
}


func PayEmployee(id int) (error,int,float64) {
	err,wage := CalculateWage(id)
	var pid int
	if err != nil {
		return err,0,0
	}

	stm := config.DbCon.QueryRow("INSERT INTO payments(employeeid,payment,date) VALUES ($1,$2,$3)",id,wage,time.Now())

	err = stm.Scan(&pid)

	if err != nil {
		return err,0,0
	}

	return nil,pid,wage
}

func calculateHours(id int) (error,float64) {
	var clockins time.Time
	var clockouts time.Time
	var duration float64

	var lastPay time.Time

	stm := config.DbCon.QueryRow("SELECT MAX(date) FROM payments WHERE employeeid = $1 ",id)

	err := stm.Scan(&lastPay)

	if err != nil {
		return err,0
	}


	row,err := config.DbCon.Query("SELECT cloak_in,cloak_out FROM hoursworked WHERE userid = $1 AND clock_in > $2 AND clock_out > $2 ",id,lastPay)

	if err != nil {
		return err,0
	}

	defer row.Close()

	for row.Next() {
		err := row.Scan(&clockins,&clockouts)
		if err != nil {
			return err,0
		}
		tempdur := clockouts.Sub(clockins)
		duration += tempdur.Hours()

		if err = row.Err(); err!= nil {
			return err,0
		}
	}

	return nil,duration
}


func CreateUser(user User) (error,string) {
	/*validateString(user.FirstName)
	validateString(user.LastName)
	validateUsername(user.Username)
	validatePassword(user.Password,user.Password)*/

	err, id := create(user.FirstName, user.LastName, user.FirstName+"."+user.LastName, user.Password, user.Birthday, user.Hiredate, user.Address, user.Phone, user.Notes, user.Type,user.HourlyWage)

	if err != nil {
		return err,""
	}

	err, uid := generateSession(id)

	if err != nil {
		return err, ""
	}

	return err, uid
}
func CheckLogin(username , password string) (error,string,User) {
	var user User
	statement := config.DbCon.QueryRow("SELECT id,firstname,lastname,typeid,birthdate,hiredate,address,phone,notes,username,hourlywage FROM employees WHERE username = $1 AND password = $2",username,password)

	err := statement.Scan(&user.Id,&user.FirstName,&user.LastName,&user.Type,&user.Birthday,&user.Hiredate,&user.Address,&user.Phone,&user.Notes,&user.Username,&user.HourlyWage)

	if err == sql.ErrNoRows {
		return errors.New("invalid username or password"),"",User{}
	}

	if err != nil {
		return err,"",User{}
	}

	err,uid := generateSession(user.Id)
	return nil,uid,user
}

func CheckSession(uuid string) (error,bool){
	var dateexp string
	err := config.DbCon.QueryRow("SELECT dateexpired FROM session WHERE uuid=$1",uuid).Scan(&dateexp)
	if err == sql.ErrNoRows {
		return err,false
	}
	date,err := time.Parse(time.RFC3339,dateexp)

	if err!=nil {
		return err,false
	}

	if err != nil {
		return err,false
	}

	dtnow := time.Now().UTC().Format(time.RFC3339)
	dt,err := time.Parse(time.RFC3339,dtnow)
	if err != nil {
		return err,false
	}
	valid := dateComparision(date.UTC(),dt.Add(2*time.Hour))

	if !valid {
		return nil,false
	}
  err,updated :=	UpdateSession(uuid)

  if err != nil {
  	fmt.Println(err)
  	return err,false
	}

  if !updated {
  	return err,false
	}
	return nil , true

}

func UpdateSession(uuid string) (error,bool) {
	dateUpdated := time.Now()
	dateExpired := time.Now().AddDate(0,0,1)
	stm,err := config.DbCon.Prepare("UPDATE session SET dateupdated=$1 , dateexpired=$2 WHERE uuid=$3")

	if err != nil {
		return err,false
	}
	_,err = stm.Exec(dateUpdated,dateExpired,uuid)

	if err != nil {
		return err,false
	}
	return nil,true
}

func dateComparision(date1,date2 time.Time) bool{
	if date1.After(date2) {
		return true
	}
	fmt.Println(date1," is not before  ",date2)
	return false
}

func generateSession(id int) (error, string) {
	uuidd, err := uuid.NewV4()

	uniqueValue := uuidd.String()

	dateUpdated := time.Now()

	dateExpired := time.Now().AddDate(0,0,1)

	fmt.Println(uniqueValue)

	err , deleted := clearSessions(id)

	if err != nil {
		return err,""
	}

	if !deleted {
		return errors.New("an error ocurred"),""
	}

	query, err := config.DbCon.Prepare("INSERT INTO session(UUID,userid,dateupdated,dateexpired) VALUES ($1,$2,$3,$4)")

	_,err = query.Exec(uniqueValue, id, dateUpdated, dateExpired)

	if err != nil {
		fmt.Println("Error")
		fmt.Println(err)
		return err, ""
	}

	return err, uniqueValue

}

func clearSessions (id int) (error,bool) {
	_,err := config.DbCon.Query("DELETE FROM session WHERE userid=$1",id)
	if err != nil {
		return err,false
	}

	return nil,true
}

func CheckToken(session string) (error,bool) {
	query , err := config.DbCon.Prepare("SELECT * FROM session WHERE uuid = $1")
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
		return errors.New("session expired"),false
	}

	return nil, true

}

func IsAdmin(session string) (error,bool,int) {
	var id int
	var access int
	query , err := config.DbCon.Prepare("SELECT userid FROM session WHERE uuid = $1")
	if err != nil {
		return err,false,http.StatusInternalServerError
	}

	err = query.QueryRow(session).Scan(&id)

	if err != nil {
		return err,false,http.StatusInternalServerError
	}

	err = config.DbCon.QueryRow("SELECT typeid FROM employees WHERE id = $1",id).Scan(&access)
	fmt.Println("Access",access)
	if err!= nil {
		return err,false,http.StatusUnauthorized
	}
	if access != 1 {
		return errors.New("unauthorized"),false,http.StatusUnauthorized
	}


	return nil,true,http.StatusOK

}



func create(fn, ln, un, ps, bday, hday,adr,phone,notes string, Type int, wage float32) (error, int) {
	fmt.Printf("%t %t %t %t %t %t %t %t %t %t", fn, ln, un, ps, bday, hday, adr,phone, notes, Type)
	var id int
	err := config.DbCon.QueryRow("INSERT INTO employees(firstname,lastname,username,password,typeid,birthdate,hiredate,address,phone,notes,hourlywage) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING Id", fn, ln, un, ps,Type,bday,hday,adr,phone,notes,wage).Scan(&id)
	if err != nil {
		return err, id
	}

	if err != nil {
		return err, id
	}

	return nil, id

}

func CurrentUser(uuid string) (error,User){
	var user User
	row := config.DbCon.QueryRow("SELECT employees.id,employees.firstname,employees.lastname,employees.typeid,employees.birthdate,employees.hiredate,employees.address,employees.phone,employees.notes,employees.username,employees.hourlywage FROM employees,session WHERE session.uuid = $1 AND employees.id = session.userid",uuid)
	err := row.Scan(&user.Id,&user.FirstName,&user.LastName,&user.Type,&user.Birthday,&user.Hiredate,&user.Address,&user.Phone,&user.Notes,&user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("no row affected"),user
		}
		return err,user
	}

	return nil,user
}

func GetUsers() (error,[]User) {
	rows,err := config.DbCon.Query("SELECT * FROM employees")
	var users []User
	var user User
	if err != nil {
		fmt.Println(err)
		return err,nil
	}

	var notes sql.NullString

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&user.Id,&user.FirstName,&user.LastName,&user.Type,&user.Birthday,&user.Hiredate,&user.Address,&user.Phone,&notes,&user.Username,&user.Password,&user.HourlyWage)
		if err != nil {
			return err,nil
		}
		if notes.Valid {
			user.Notes = notes.String
		}
		user.Birthday = user.Birthday[0:10]
		user.Hiredate = user.Hiredate[0:10]
		users = append(users, user)
		if err = rows.Err();err != nil {
			return err,nil
		}
	}

	return nil,users
}

func EditUser(user User) (error,bool) {
	_,err := config.DbCon.Query(	"UPDATE employees SET " +
																			"firstname=$1,lastname=$2,username=$3,password=$4," +
																			"typeid=$5,birthdate=$6,hiredate=$7,address=$8,phone=$9,notes=$10,hourlywage=$11" +
																			"WHERE id=$12",user.FirstName,user.LastName,user.Username,user.Password,user.Type,user.Birthday,user.Hiredate,user.Address,user.Phone,user.Notes,user.HourlyWage,user.Id )

	if err != nil {
		return err,false
	}

	return nil,true

}

func ClockIn(id int) (error,bool) {
	date := time.Now().UTC()
	_,err := config.DbCon.Query("INSERT INTO hoursworked(userid,cloak_in) VALUES ($1,$2)",id,date)
	if err != nil {
		return err,false
	}
	return nil,true

}

func IsClockedIn(id int) (error,bool) {
	stm,err := config.DbCon.Prepare("SELECT * FROM hoursworked WHERE userid = $1 AND cloak_out IS NULL")
	if err != nil {
		return err,false
	}

	row,err := stm.Exec(id)

	if err != nil {
		return err,false
	}
	nr,err := row.RowsAffected()
	if err != nil {
		return err,false
	}

	if nr == 0 || nr > 1 {
		return errors.New("not clocked in"),false
	}

	return nil,true

}

func ClockOut (id int)(error,bool) {
	date := time.Now().UTC()
	_,err := config.DbCon.Query("UPDATE hoursworked SET cloak_out = $1 WHERE userid = $2 AND cloak_out IS NULL ",date,id)
	if err != nil {
		return err,false
	}
	return nil,true
}
