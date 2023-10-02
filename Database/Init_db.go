package db

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Init_db() (Db, error) {
	database := Db{}
	data, err := Create_DB()
	if err != nil {
		fmt.Println("⚠ ERROR ⚠: could'nt init database ❌ ", err)
		return Db{}, err
	}
	database.Doc = data

	fmt.Println("✅ database has been created successfully")

	database.Tables()
	return database, err
}
