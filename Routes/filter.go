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

func Profil(w http.ResponseWriter, r *http.Request, database db.Db) {
	//code ajouté
	c, errc := r.Cookie("session_token")
	if errc != nil {
		fmt.Println("pas de cookie session")
		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else {

		s, err, _ := auth.HelpersBA(database, "username", "WHERE usersession='"+c.Value+"'", "")
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

	Id_user, _, _ := auth.HelpersBA(database, "id_user", "WHERE usersession='"+c.Value+"'", "")

	//fin code
	//checking the http request
	if r.Method != "GET" {
		fmt.Printf("⚠ ERROR ⚠ : cannot access to that page by with mode other than GET must log out to reach it ❌")
		w.WriteHeader(http.StatusMethodNotAllowed)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "405")
		return
	}
	fmt.Println("in profil")
	//-------------- retrieving datas ---------------//
	//--1
	errGetPost := postab.GetPost_data(database)
	if errGetPost != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}
	//--2
	errGetComm := commtab.GetComment_data(database)
	if errGetComm != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}
	//--3
	errectabcomm := reactab_com.GetReact_comdata(database)
	if errectabcomm != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't get comments reaction for display from database\n")
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}
	//--4
	categos, err := Com.GetPost_categories(database)
	if err != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't get categories data from database\n")
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}
	//--5
	errectab := reactab.Get_reacPosts_data(database)
	if errectab != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't get reaction for display a from database\n")
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}
	//--------------------------------------------------------------------//
	// storing the session's id
	for i := range postab {
		postab[i].SessionId = Id_user
	}
	for i := range commtab {
		commtab[i].SessionId = Id_user
	}

	//storing user's name in structures
	for i := range postab {
		username, name, surname := tools.GetName_byID(database, postab[i].UserId)
		postab[i].Username = username
		postab[i].Name = name
		postab[i].Surname = surname
	}

	for i := range commtab {
		username, name, surname := tools.GetName_byID(database, commtab[i].UserId)
		commtab[i].Username = username
		commtab[i].Name = name
		commtab[i].Surname = surname
	}

	//storing the reactions in corresponding comments
	for i := range commtab {
		for j := range reactab_com {
			if commtab[i].CommentId == reactab_com[j].CommentId {
				switch reactab_com[j].Reaction {
				case true:
					commtab[i].Likecomm = append(commtab[i].Likecomm, "true")
					if reactab_com[j].UserId == Id_user {
						commtab[i].SessionReact = "true"
					}

				case false:
					commtab[i].Dislikecomm = append(commtab[i].Dislikecomm, "false")
					if reactab_com[j].UserId == Id_user {
						commtab[i].SessionReact = "false"
					}
				}
			}
		}
	}

	//storing the comments in corresponding posts
	for i := range postab {
		for j := range commtab {
			if postab[i].PostId == commtab[j].PostId {
				postab[i].Comment_tab = append(postab[i].Comment_tab, commtab[j])
			}
		}
	}

	//storing the categories in corresponding posts
	for i := range postab {
		for j := range categos {
			if postab[i].PostId == categos[j].PostId {
				postab[i].Categorie = append(postab[i].Categorie, categos[j].Category)
			}
		}
	}

	//storing the reactions in corresponding posts
	for i := range postab {
		for j := range reactab {
			if postab[i].PostId == reactab[j].PostId {
				switch reactab[j].Reaction {
				case true:
					postab[i].Like = append(postab[i].Like, "true")
					if reactab[j].UserId == Id_user {
						postab[i].SessionReact = "true"
					}

				case false:
					postab[i].Dislike = append(postab[i].Dislike, "false")
					if reactab[j].UserId == Id_user {
						postab[i].SessionReact = "false"
					}
				}
			}
		}
	}

	//--------retrieving form values ----------
	fmt.Println("--------------------------------------------")
	fmt.Println("             Profil form values             ")
	fmt.Println("--------------------------------------------")

	choice := r.URL.Query().Get("filter")
	fmt.Println("[INFO] choice: ", choice) //debug

	if choice != "posts" && choice != "favorites" && choice != "comments" {
		fmt.Printf("⚠ ERROR ⚠ parsing --> bad request ❌\n")
		w.WriteHeader(http.StatusNotFound)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "400")
		return
	}

	var newtab Com.Posts

	switch {
	case choice == "favorites":
		for _, v := range postab {
			for _, j := range reactab {
				if (v.PostId == j.PostId) && (j.UserId == Id_user) && j.Reaction {
					newtab = append(newtab, v)
				}
			}
		}

	case choice == "posts":
		for _, v := range postab {
			if v.UserId == Id_user {
				newtab = append(newtab, v)
			}
		}

	case choice == "comments":
		for _, v := range postab {
			for _, j := range v.Comment_tab {
				if j.UserId == Id_user {
					newtab = append(newtab, v)
					break
				}
			}
		}
	}

	username, name, surname := tools.GetName_byID(database, Id_user)
	file, errf := template.ParseFiles("templates/profil.html", "templates/head.html", "templates/navbar.html", "templates/main.html", "templates/footer.html")
	if errf != nil {
		//sending metadata about the error to the servor
		fmt.Printf("⚠ ERROR ⚠ parsing profil.html--> %v\n", errf)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		w.WriteHeader(http.StatusInternalServerError)
		error_file.Execute(w, "500")
		return
	}
	//returning "empty" signal to show postab is empty
	//(there 's no result after filter)
	var empty bool
	if len(newtab) == 0 {
		empty = true
	}
	//users name and surname
	//struct to execute
	finalex := struct {
		CurrentN  string
		CurrentSN string
		CurrentUN string
		Postab    Com.Posts
		Empty     bool
	}{
		CurrentN:  name,
		CurrentSN: surname,
		CurrentUN: username,
		Postab:    newtab,
		Empty:     empty,
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

func Filter(w http.ResponseWriter, r *http.Request, database db.Db) {
	//code ajouté
	c, errc := r.Cookie("session_token")
	if errc != nil {
		fmt.Println("pas de cookie session")
		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else {

		s, err, _ := auth.HelpersBA(database, "username", "WHERE usersession='"+c.Value+"'", "")
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

	//-------------- retrieving datas ---------------//
	//--1
	errGetPost := postab.GetPost_data(database)
	if errGetPost != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}

	//--2
	errGetComm := commtab.GetComment_data(database)
	if errGetComm != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}

	//--3
	errectabcomm := reactab_com.GetReact_comdata(database)
	if errectabcomm != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't get comments reaction for display from database\n")
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}

	//--4
	categos, err := Com.GetPost_categories(database)
	if err != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't get categories data from database\n")
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}
	//--5
	errectab := reactab.Get_reacPosts_data(database)
	if errectab != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't get reaction for display a from database\n")
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}
	//--------------------------------------------------------------------//

	// storing the session's id
	for i := range postab {
		postab[i].SessionId = Id_user
	}
	for i := range commtab {
		commtab[i].SessionId = Id_user
	}

	//storing user's name in structures
	for i := range postab {
		username, name, surname := tools.GetName_byID(database, postab[i].UserId)
		postab[i].Username = username
		postab[i].Name = name
		postab[i].Surname = surname
	}

	for i := range commtab {
		username, name, surname := tools.GetName_byID(database, commtab[i].UserId)
		commtab[i].Username = username
		commtab[i].Name = name
		commtab[i].Surname = surname
	}

	//storing the reactions in corresponding comments
	for i := range commtab {
		for j := range reactab_com {
			if commtab[i].CommentId == reactab_com[j].CommentId {
				switch reactab_com[j].Reaction {
				case true:
					commtab[i].Likecomm = append(commtab[i].Likecomm, "true")
					commtab[i].SessionReact = "true"
				case false:
					commtab[i].Dislikecomm = append(commtab[i].Dislikecomm, "false")
					commtab[i].SessionReact = "false"
				}
			}
		}
	}

	//storing the comments in corresponding posts
	for i := range postab {
		for j := range commtab {
			if postab[i].PostId == commtab[j].PostId {
				postab[i].Comment_tab = append(postab[i].Comment_tab, commtab[j])
			}
		}
	}

	//storing the categories in corresponding posts
	for i := range postab {
		for j := range categos {
			if postab[i].PostId == categos[j].PostId {
				postab[i].Categorie = append(postab[i].Categorie, categos[j].Category)
			}
		}
	}

	//storing the reactions in corresponding posts
	for i := range postab {
		for j := range reactab {
			if postab[i].PostId == reactab[j].PostId {
				switch reactab[j].Reaction {
				case true:
					postab[i].Like = append(postab[i].Like, "true")
					postab[i].SessionReact = "true"
				case false:
					postab[i].Dislike = append(postab[i].Dislike, "false")
					postab[i].SessionReact = "false"
				}
			}
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
	username, name, surname := tools.GetName_byID(database, Id_user)
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
		CurrentN  string
		CurrentSN string
		CurrentUN string
		Postab    Com.Posts
		Empty     bool
	}{
		CurrentN:  name,
		CurrentSN: surname,
		CurrentUN: username,
		Postab:    newtab,
		Empty:     empty,
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
	//-------------- retrieving datas ---------------//
	//--1
	errGetPost := postab.GetPost_data(database)
	if errGetPost != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}
	//--2
	errGetComm := commtab.GetComment_data(database)
	if errGetComm != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}
	//--3
	errectabcomm := reactab_com.GetReact_comdata(database)
	if errectabcomm != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't get comments reaction for display from database\n")
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}
	//--4
	categos, err := Com.GetPost_categories(database)
	if err != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't get categories data from database\n")
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}
	//--5
	errectab := reactab.Get_reacPosts_data(database)
	if errectab != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't get reaction for display a from database\n")
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}
	//--------------------------------------------------------------------//
	// storing the session's id
	for i := range postab {
		postab[i].SessionId = Id_user
	}
	for i := range commtab {
		commtab[i].SessionId = Id_user
	}

	//storing user's name in structures
	for i := range postab {
		username, name, surname := tools.GetName_byID(database, postab[i].UserId)
		postab[i].Username = username
		postab[i].Name = name
		postab[i].Surname = surname
	}

	for i := range commtab {
		username, name, surname := tools.GetName_byID(database, commtab[i].UserId)
		commtab[i].Username = username
		commtab[i].Name = name
		commtab[i].Surname = surname
	}

	//storing the reactions in corresponding comments
	for i := range commtab {
		for j := range reactab_com {
			if commtab[i].CommentId == reactab_com[j].CommentId {
				switch reactab_com[j].Reaction {
				case true:
					commtab[i].Likecomm = append(commtab[i].Likecomm, "true")
				case false:
					commtab[i].Dislikecomm = append(commtab[i].Dislikecomm, "false")
				}
			}
		}
	}

	//storing the comments in corresponding posts
	for i := range postab {
		for j := range commtab {
			if postab[i].PostId == commtab[j].PostId {
				postab[i].Comment_tab = append(postab[i].Comment_tab, commtab[j])
			}
		}
	}

	//storing the categories in corresponding posts
	for i := range postab {
		for j := range categos {
			if postab[i].PostId == categos[j].PostId {
				postab[i].Categorie = append(postab[i].Categorie, categos[j].Category)
			}
		}
	}

	//storing the reactions in corresponding posts
	for i := range postab {
		for j := range reactab {
			if postab[i].PostId == reactab[j].PostId {
				switch reactab[j].Reaction {
				case true:
					postab[i].Like = append(postab[i].Like, "true")
				case false:
					postab[i].Dislike = append(postab[i].Dislike, "false")
				}
			}
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
	username, name, surname := tools.GetName_byID(database, Id_user)
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
		CurrentN  string
		CurrentSN string
		CurrentUN string
		Postab    Com.Posts
		Empty     bool
	}{
		CurrentN:  name,
		CurrentSN: surname,
		CurrentUN: username,
		Postab:    newtab,
		Empty:     empty,
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
