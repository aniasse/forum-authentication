package auth

import (
	"database/sql"
	"fmt"
	db "forum/Database"
)

type User struct {
	Id       string
	Username string
	Name     string
	Email    string
	Password string
	Pp       string
	Pc       string
}

func GetDatafromBA(tab *sql.DB, data, attribute, table string) bool {
	var response bool
	selectSQL := "SELECT " + attribute + " FROM " + table + "  ;"

	rows, err := tab.Query(selectSQL)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&attribute)
		if err != nil {
			fmt.Println(err)
			return false
		}
		if attribute == data {
			response = true
		}
	}
	return response
}

func GetElementOfOneUser(db *sql.DB, username string) (user User, response bool) {
	rows, err := db.Query("SELECT id_user,name,email,pp,pc FROM users WHERE username='" + username + "';")

	var id_user, name, email, pp, pc string
	if err != nil {
		fmt.Println(err, "1")
		return user, false
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id_user, &name, &email, &pp, &pc)
		if err != nil {
			fmt.Println(err, "2")
			return user, false
		}
	}
	user = User{Id: id_user, Username: username, Name: name, Email: email, Pp: pp, Pc: pc}
	return user, true
}

func HelpersBA(from string, tab db.Db, attribute, condition, compare string) (string, error, bool) {
	result := ""
	response := false
	rows, errorrows := tab.GetData(attribute, from, condition)
	if errorrows != nil {
		// _, _, confirmemail := auth.HelpersBA(tab, "username", "", username)
		return result, errorrows, response
	}

	for rows.Next() {

		// var password string
		err := rows.Scan(&attribute)
		if err != nil {
			// fmt.Println(err)
			return result, err, response
		}
		if attribute == compare {
			response = true
		}
		result = attribute
	}
	return result, nil, response
}
