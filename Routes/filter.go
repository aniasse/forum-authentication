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
	//-------------- retrieving datas ---------------//
	//--1
	
	errGetPost := postab.GetPost_data(database)
	if errGetPost != nil {
		auth.Snippets(w, 500)
		return
	}
	//--2
	errGetComm := commtab.GetComment_data(database)
	if errGetComm != nil {
		auth.Snippets(w, 500)
		return
	}
	//--3
	errectabcomm := reactab_com.GetReact_comdata(database)
	if errectabcomm != nil {
		auth.Snippets(w, 500)
		return
	}
	//--4
	categos, err := Com.GetPost_categories(database)
	if err != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't get categories data from database\n")
		auth.Snippets(w, 500)
		return
	}
	//--5
	errectab := reactab.Get_reacPosts_data(database)
	if errectab != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't get reaction for display a from database\n")
		auth.Snippets(w, 500)
		return
	}

	//storing user's name and profil image in structures
	for i := range postab {
		username, name, surname, errGN := tools.GetName_byID(database, postab[i].UserId)
		Profil, errprof := tools.GetPic_byID(database, postab[i].UserId)

		if errprof != nil || errGN != nil {
			//sending metadata about the error to the servor
			auth.Snippets(w, 500)
			return
		}
		postab[i].Profil = Profil
		postab[i].Username = username
		postab[i].Name = name
		postab[i].Surname = surname
	}

	for i := range commtab {
		username, name, surname, errGN := tools.GetName_byID(database, commtab[i].UserId)
		Profil, errprof := tools.GetPic_byID(database, commtab[i].UserId)

		if errprof != nil || errGN != nil {
			//sending metadata about the error to the servor
			auth.Snippets(w, 500)
			return
		}
		commtab[i].Profil = Profil
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
