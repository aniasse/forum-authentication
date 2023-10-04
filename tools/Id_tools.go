package tools

import (
	"fmt"
	data "forum/Database"
	"log"
)

func GetName_byID(database data.Db, ID string) (string, string, string) {
	//getting the user's name
	condition := fmt.Sprintf("WHERE %s = '%s'", data.Id_user, ID)
	request := fmt.Sprintf("%s, %s, %s", data.Username, data.Name, data.Surname)
	info, errn := database.GetData(request, data.User, condition)
	if errn != nil {
		fmt.Println("⚠ ERROR ⚠ : Couldn't get the username according to the id from database ❌")
		log.Fatalf("⚠ : %v\n", errn)
	}

	var username, surname, name string
	for info.Next() {
		err := info.Scan(&username, &name, &surname)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return username, name, surname
}

func GetName_bycomment(database data.Db, ID string) string {
	//getting the user's name
	condition := fmt.Sprintf("WHERE %s = '%s'", data.Id_comment, ID)
	info, errn := database.GetData(data.Username, data.Comment, condition)
	if errn != nil {
		fmt.Println("⚠ ERROR ⚠ : Couldn't get the username according to the commentId from database ❌")
		log.Fatalf("⚠ : %v\n", errn)
	}

	var username string
	for info.Next() {
		err := info.Scan(&username)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return username
}

func IsnotExist_user(id string, database data.Db) bool {
	Condition := fmt.Sprintf("WHERE %s = '%s'", data.Id_user, id)
	got, _ := database.GetData(data.Email, data.User, Condition)
	stored, _ := data.Getelement(got)
	fmt.Println("stored user in database: ", stored)
	if stored == "" {
		fmt.Printf("✖ Id n°= %s doesn't exist in Database\n", id)
		return true
	}
	return false
}

func IsnotExist_Post(id string, database data.Db) bool {
	Condition := fmt.Sprintf("WHERE %s = '%s'", data.Id_post, id)
	got, _ := database.GetData(data.Description, data.Post, Condition)
	stored, _ := data.Getelement(got)
	fmt.Println("stored post content in database: ", stored)
	if stored == "" {
		fmt.Printf("✖ Post n°= %s doesn't exist in Database\n", id)
		return true
	}
	return false
}

func IsnotExist_Comment(id string, database data.Db) bool {
	Condition := fmt.Sprintf("WHERE %s = '%s'", data.Id_comment, id)
	gotcomm, _ := database.GetData(data.Content, data.Comment, Condition)
	stored, _ := data.Getelement(gotcomm)
	fmt.Println("stored comm to reply content  in data base: ", stored)
	if stored == "" {
		fmt.Printf("✖ comment n°= %s doesn't exist in Database\n", id)
		return true
	}
	return false
}
