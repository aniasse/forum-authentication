package Route

import (
	"fmt"
	auth "forum/Authentication"
	Com "forum/Communication"
	db "forum/Database"
	tools "forum/tools"
	"html/template"
	"net/http"
	"path"
)

func Profil(w http.ResponseWriter, r *http.Request, database db.Db) {
	//checking session
	HaveSession, session := auth.ComSession_Checker(w, r, database)
	if !HaveSession {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	Id_user, _, _ := auth.HelpersBA("sessions", database, "user_id", "WHERE id_session='"+session.Value+"'", "")

	//checking the http request
	if r.Method != "GET" && r.Method != "POST" {
		fmt.Printf("⚠ ERROR ⚠ : cannot access to that page by with mode other than GET must log out to reach it ❌")
		auth.Snippets(w, 405)
		return
	}
	fmt.Println("in profil")

	//--------retrieving form values ----------
	fmt.Println("--------------------------------------------")
	fmt.Println("             Profil form value             ")
	fmt.Println("--------------------------------------------")
	choice := path.Base(r.URL.Path)

	fmt.Println("[INFO] profil filter choice: ", choice) //debug
	if choice != "posts" {
		fmt.Printf("⚠ ERROR ⚠ parsing --> bad request ❌\n")
		auth.Snippets(w, 400)
		return
	}

	GetAll_fromDB(w)
	UploadImageUser(w, r, Id_user)
	StatusCode := ProcessData(w, r, "/myprofil/"+choice)
	if StatusCode != 200 {
		auth.Snippets(w, StatusCode)
		return
	}
	var newtab Com.Posts
	for _, v := range postab {
		if v.UserId == Id_user {
			newtab = append(newtab, v)
		}
	}

	username, name, surname, errGN := tools.GetName_byID(database, Id_user)
	if errGN != nil {
		//sending metadata about the error to the servor
		auth.Snippets(w, 500)
		return
	}

	//code
	current_pp, _, errpp := auth.HelpersBA("users",database, "pp", " WHERE id_user='"+Id_user+"'", "")
	current_cover, _, errcover := auth.HelpersBA("users",database, "pc", " WHERE id_user='"+Id_user+"'", "")
	//handle error
	if errpp || errcover {
		fmt.Println("error pp,", errpp, " error cover", errcover)
		auth.Snippets(w, http.StatusInternalServerError)
	}
	//end
	file, errf := template.ParseFiles("templates/profil.html", "templates/head.html", "templates/navbar.html", "templates/main.html", "templates/footer.html")
	if errf != nil {
		//sending metadata about the error to the servor
		auth.Snippets(w, 500)
		return
	}

	//returning "empty" signal to show postab is empty (there 's no result after filter)
	var empty bool
	if len(newtab) == 0 {
		empty = true
	}
	//users name and surname
	//struct to execute
	finalex := Res{
		CurrentN:     name,
		CurrentSN:    surname,
		CurrentUN:    username,
		CurrentPP:    current_pp,
		CurrentCover: current_cover,
		Postab:       newtab,
		Empty:        empty,
	}

	//sending data to html
	errexc := file.Execute(w, finalex)
	if errexc != nil {
		//sending metadata about the error to the servor
		fmt.Printf("⚠ ERROR ⚠ executing in profil --> %v\n", errexc)
		http.Error(w, "⚠ INTERNAL SERVER ERROR ⚠", http.StatusInternalServerError)
		return
	}
	fmt.Println("--------------- 🟢🌐 profil data sent -----------------------") //debug
}
