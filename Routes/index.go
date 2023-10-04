package Route

import (
	"fmt"
	Err "forum/Authentification"
	db "forum/Database"
	"html/template"
	"net/http"
)

/*
Index parses the the homepage where no interaction is possible
we only display the forum's informations
*/
func Index(w http.ResponseWriter, r *http.Request, database db.Db) {
	//code ajouté
	Err.CheckCookie(w, r, database)
	//fin code
	if r.Method != "GET" {
		fmt.Printf("⚠ ERROR ⚠ : cannot access to that page by get mode must log out to reach it ❌")
		w.WriteHeader(http.StatusMethodNotAllowed)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "405")
		return
	}

	//checking whether the route exists or not
	if r.URL.Path != "/" && r.URL.Path != "/home" && r.URL.Path != "/myprofil" && r.URL.Path != "/filter" {
		fmt.Printf("⚠ ERROR ⚠ parsing --> page not found ❌\n")
		w.WriteHeader(http.StatusNotFound)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "404")
		return
	}

	Display_mngmnt(w, r) // displaying datas
	//--displaying welcoming post when database is empty
	if len(postab) == 0 {
		errwel := postab.Welcome_user(database, "index")
		if errwel != nil {
			fmt.Printf("⚠ INDEX ERRWEL ⚠ :%s ❌", errwel)
			Err.Snippets(w, 500)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		fmt.Println("✔ ✨ welcome post sent ✨")
	}

	//--removing the reactions highlihts
	for i := range commtab {
		commtab[i].SessionReact = ""
	}

	for i := range postab {
		postab[i].SessionReact = ""
	}

	file, errf := template.ParseFiles("templates/index.html")
	if errf != nil {
		//sending metadata about the error to the servor
		fmt.Printf("⚠ ERROR ⚠ parsing --> %v\n", errf)
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}

	//struct to execute
	final := Res{
		Postab: postab,
	}

	//sending data to html
	errexc := file.Execute(w, final)
	if errexc != nil {
		//sending metadata about the error to the servor
		fmt.Printf("⚠ ERROR index ⚠ executing file --> %v\n", errexc)
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}
	fmt.Println("--------------- 🟢🌐 data sent from index -----------------------")

}
