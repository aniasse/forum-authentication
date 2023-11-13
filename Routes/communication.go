package Route

import (
	"fmt"
	"html/template"
	"net/http"

	auth "forum/Authentification"
	Com "forum/Communication"
	db "forum/Database"
	tools "forum/tools"
)

type Res struct {
	CurrentN     string
	CurrentSN    string
	CurrentUN    string
	CurrentPP    string
	CurrentCover string
	Postab       Com.Posts
	Empty        bool
}

var (
	postab      Com.Posts    // posts local storage
	commtab     Com.Comments // comments local storage
	reactab     Com.Reacts   //posts reactions local storage
	reactab_com Com.ReactC   // comments reactions local storage
	database    db.Db        //database local storage
	errd        error        // manage errors
)

var Id_user string

/*
Communications handles user's posts, comments and reactions
it can only be reached by using method POST or GET
*/
func Communication(w http.ResponseWriter, r *http.Request, Id string, redirect string) {
	Id_user = Id
	//!--checking the http request
	if r.Method != "POST" && r.Method != "GET" {
		fmt.Printf("⚠ ERROR ⚠ : cannot access to that page by with mode other than GET & POST must log out to reach it ❌")
		auth.Snippets(w, 405)
		return
	}

	GetAll_fromDB(w) //display all values in the forum database
	fmt.Println("postab size ->> ", len(postab))
	StatusCode := ProcessData(w, r, redirect) //Process datas received fromn client request
	if StatusCode != 200 {
		auth.Snippets(w, StatusCode)
		return
	}

	file, errf := template.ParseFiles("templates/home.html", "templates/head.html", "templates/navbar.html", "templates/main.html", "templates/footer.html")
	if errf != nil {
		//sending metadata about the error to the servor
		fmt.Printf("⚠ ERROR ⚠ parsing --> %v\n", errf)
		auth.Snippets(w, 500)
		return
	}
	// user's name
	current_username, current_surname, current_name, errGN := tools.GetName_byID(database, Id_user)
	if errGN != nil {
		//sending metadata about the error to the servor
		auth.Snippets(w, 500)
		return
	}

	// code
	current_pp, _, errpp := auth.HelpersBA("users", database, "pp", " WHERE id_user='"+Id_user+"'", "")
	current_cover, _, errcover := auth.HelpersBA("users", database, "pc", " WHERE id_user='"+Id_user+"'", "")
	// handle error
	if errpp || errcover {
		fmt.Println("error pp,", errpp, " error cover", errcover)
		auth.Snippets(w, http.StatusInternalServerError)
	}
	//end
	//returning "empty" signal to show postab is empty
	//(there 's no result after filter)

	//struct to execute
	final := Res{
		CurrentUN:    current_username,
		CurrentSN:    current_surname,
		CurrentN:     current_name,
		CurrentPP:    current_pp,
		CurrentCover: current_cover,
		Postab:       postab,
	}

	//sending data to html
	errexc := file.Execute(w, final)
	if errexc != nil {
		//sending metadata about the error to the servor
		fmt.Printf("⚠ ERROR ⚠ executing file --> %v\n", errexc)
		auth.Snippets(w, 500)
		return
	}
	fmt.Println("--------------- 🟢🌐 home data sent -----------------------") //debug

}
