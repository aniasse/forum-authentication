package tools

import (
	"fmt"
	data "forum/Database"
)

func GetName_byID(database data.Db, ID string) (string, string, string, error) {
	//getting the user's name
	condition := fmt.Sprintf("WHERE %s = '%s'", data.Id_user, ID)
	request := fmt.Sprintf("%s, %s, %s", data.Username, data.Name, data.Surname)
	info, errn := database.GetData(request, data.User, condition)
	if errn != nil {
		fmt.Println("⚠ ERROR ⚠ : Couldn't get the username according to the id from database ❌")
		fmt.Printf("⚠ : %v\n", errn)
		return "", "", "", errn
	}

	var username, surname, name string
	for info.Next() {
		err := info.Scan(&username, &name, &surname)
		if err != nil {
			fmt.Printf("⚠ : %v\n", err)
			return "", "", "", err
		}
	}

	return username, name, surname, nil
}

func GetPic_byID(database data.Db, Id string) (string, error) {
	//getting the user's profil image
	condition := fmt.Sprintf("WHERE %s = '%s'", data.Id_user, Id)
	info, errn := database.GetData(data.Pp, data.User, condition)
	if errn != nil {
		fmt.Println("⚠ ERROR ⚠ : Couldn't get the profil photo according to the id from database ❌")
		fmt.Printf("⚠ : %v\n", errn)
		return "", errn
	}

	picture, errpic := data.Getelement(info)
	if errpic != nil {
		fmt.Println("⚠ ERROR ⚠ : Couldn't get the profil image from database ❌")
		fmt.Printf("⚠ : %v\n", errn)
		return "", errpic
	}

	return picture, nil
}

func GetName_bycomment(database data.Db, ID string) (string, error) {
	//getting the user's name
	condition := fmt.Sprintf("WHERE %s = '%s'", data.Id_comment, ID)
	info, errn := database.GetData(data.Username, data.Comment, condition)
	if errn != nil {
		fmt.Println("⚠ ERROR ⚠ : Couldn't get the username according to the commentId from database ❌")
		fmt.Printf("⚠ : %v\n", errn)
		return "", errn
	}

	var username string
	for info.Next() {
		err := info.Scan(&username)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
	}

	return username, nil
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
