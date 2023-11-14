package Route

import (
	"fmt"
	"html/template"
	"net/http"

	auth "forum/Authentication"
	Com "forum/Communication"
	db "forum/Database"
	tools "forum/tools"
)

func Filter(w http.ResponseWriter, r *http.Request, database db.Db) {
	//code ajouté
	c, errc := r.Cookie("session_token")
	if errc != nil {
		fmt.Println("pas de cookie session")
		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else {
		s, err, _ := auth.HelpersBA("sessions", database, "user_id", "WHERE id_session='"+c.Value+"'", "")
		// fmt.Println("here", s, "error", err)
		if err != nil {
			fmt.Println("erreur du serveur", err)
		}
		if s == "" {
			fmt.Println("cookie invalide,affichage de /", s, "verif vide")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
	//fin code

	//-------------- retrieving datas ---------------/
	GetAll_fromDB(w)

	//--------retrieving form values ----------
	fmt.Println("--------------------------------------------")
	fmt.Println("             Filter form values             ")
	fmt.Println("--------------------------------------------")

	categorie := r.URL.Query().Get("filter")
	if categorie == "art" {
		categorie = "art & culture"
	}
	if categorie != "art & culture" && categorie != "education" && categorie != "sport" && categorie != "cinema" && categorie != "health" && categorie != "others" {
		fmt.Printf("⚠ ERROR ⚠ filtering --> bad request ❌\n")
		w.WriteHeader(http.StatusNotFound)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "400")
		return
	}
	fmt.Println("[INFO] categorie choice: ", categorie) //debug
	var newtab Com.Posts

	for _, v := range postab {
		for _, j := range v.Categorie {
			if j == categorie {
				newtab = append(newtab, v)
				break
			}
		}
	}

	file, errf := template.ParseFiles("templates/home.html", "templates/head.html", "templates/navbar.html", "templates/main.html", "templates/footer.html")
	if errf != nil {
		//sending metadata about the error to the servor
		fmt.Printf("⚠ ERROR ⚠ parsing home.html--> %v\n", errf)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		w.WriteHeader(http.StatusInternalServerError)
		error_file.Execute(w, "500")
		return
	}
	username, name, surname, errGN := tools.GetName_byID(database, Id_user)
	if errGN != nil {
		//sending metadata about the error to the servor
		auth.Snippets(w, 500)
		return
	}
	current_pp, _, errpp := auth.HelpersBA("users",database, "pp", " WHERE id_user='"+Id_user+"'", "")
	current_cover, _, errcover := auth.HelpersBA("users",database, "pc", " WHERE id_user='"+Id_user+"'", "")
	//handle error
	if errpp || errcover {
		fmt.Println("error pp,", errpp, " error cover", errcover)
		auth.Snippets(w, http.StatusInternalServerError)
	}
	fmt.Println("creds ", username, name, surname)

	//returning "empty" signal to show postab is empty
	//(there 's no result after filter)
	var empty bool
	if len(newtab) == 0 {
		empty = true
	}
	fmt.Println("empty bool filter -> ", empty)

	//users name and surname
	//struct to execute
	final := struct {
		CurrentN     string
		CurrentSN    string
		CurrentUN    string
		CurrentPP    string
		CurrentCover string
		Postab       Com.Posts
		Empty        bool
	}{
		CurrentN:     name,
		CurrentSN:    surname,
		CurrentUN:    username,
		CurrentPP:    current_pp,
		CurrentCover: current_cover,
		Postab:       newtab,
		Empty:        empty,
	}

	//sending data to html
	errexc := file.Execute(w, final)
	if errexc != nil {
		//sending metadata about the error to the servor
		fmt.Printf("⚠ ERROR ⚠ executing in home --> %v\n", errexc)
		http.Error(w, "⚠ INTERNAL SERVER ERROR ⚠", http.StatusInternalServerError)
		return
	}
	fmt.Println("--------------- 🟢🌐 filter data sent -----------------------") //debug

}

func Indexfilter(w http.ResponseWriter, r *http.Request, database db.Db) {
	auth.CheckCookie(w, r, database)

	GetAll_fromDB(w) // displaying datas
	//--removing the reactions highlihts

	for i := range postab {
		postab[i].SessionReact = ""
	}
	for i := range postab {
		for v := range postab[i].Comment_tab {
			postab[i].Comment_tab[v].SessionReact = ""
		}
	}

	//--------retrieving form values ----------
	fmt.Println("--------------------------------------------")
	fmt.Println("             Filter form values             ")
	fmt.Println("--------------------------------------------")

	categorie := r.URL.Query().Get("filter")
	if categorie == "art" {
		categorie = "art & culture"
	}
	fmt.Println("[INFO] categorie choice: ", categorie) //debug
	if categorie != "art & culture" && categorie != "education" && categorie != "sport" && categorie != "cinema" && categorie != "health" && categorie != "others" {
		fmt.Printf("⚠ ERROR ⚠ parsing --> bad request ❌\n")
		w.WriteHeader(http.StatusNotFound)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "400")
		return
	}
	var newtab Com.Posts

	for _, v := range postab {
		for _, j := range v.Categorie {
			if j == categorie {
				newtab = append(newtab, v)
				break
			}
		}
	}

	file, errf := template.ParseFiles("templates/index.html", "templates/footer.html", "templates/navbar.html", "templates/head.html")
	if errf != nil {
		//sending metadata about the error to the servor
		fmt.Printf("⚠ ERROR ⚠ parsing home.html--> %v\n", errf)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		w.WriteHeader(http.StatusInternalServerError)
		error_file.Execute(w, "500")
		return
	}
	username, name, surname, errGN := tools.GetName_byID(database, Id_user)
	if errGN != nil {
		//sending metadata about the error to the servor
		auth.Snippets(w, 500)
		return
	}
	//code
	current_pp, _, errpp := auth.HelpersBA("users", database, "pp", " WHERE id_user='"+Id_user+"'", "")
	current_cover, _, errcover := auth.HelpersBA("users", database, "pc", " WHERE id_user='"+Id_user+"'", "")
	//handle error
	if errpp || errcover {
		fmt.Println("error pp,", errpp, " error cover", errcover)
		auth.Snippets(w, http.StatusInternalServerError)
	}
	//end
	fmt.Println("creds ", username, name, surname)
	//returning "empty" signal to show postab is empty
	//(there 's no result after filter)
	var empty bool
	if len(newtab) == 0 {
		empty = true
	}
	//users name and surname
	//struct to execute
	final := struct {
		CurrentN     string
		CurrentSN    string
		CurrentUN    string
		CurrentPP    string
		CurrentCover string
		Postab       Com.Posts
		Empty        bool
	}{
		CurrentN:     name,
		CurrentSN:    surname,
		CurrentUN:    username,
		CurrentPP:    current_pp,
		CurrentCover: current_cover,
		Postab:       newtab,
		Empty:        empty,
	}

	//sending data to html
	errexc := file.Execute(w, final)
	if errexc != nil {
		//sending metadata about the error to the servor
		fmt.Printf("⚠ ERROR ⚠ executing in home --> %v\n", errexc)
		http.Error(w, "⚠ INTERNAL SERVER ERROR ⚠", http.StatusInternalServerError)
		return
	}
	fmt.Println("--------------- 🟢🌐 filter data sent -----------------------") //debug

}
