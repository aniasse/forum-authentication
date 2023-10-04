package hdle

import (
	"fmt"
	"log"
	"net/http"

	auth "forum/Authentification"
	db "forum/Database"
	Rt "forum/Routes"
)

func Handlers() {
	// http.Handle("/static/", http.StripPrefix("/static/", fs))
	// http.HandleFunc("/home", Rt.Communication)
	// http.HandleFunc("/myprofil", Rt.Profil)
	// http.HandleFunc("/filter", Rt.Filter)
	// http.HandleFunc("/", Rt.Index)
	// fmt.Println("🌐 server has started at : http://localhost:8080 🟢")
	// http.ListenAndServe(":8080", nil)

	staticHandler := http.FileServer(http.Dir("templates"))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	tab, err := db.Init_db()
	if err != nil {
		log.Fatalln(err)
	}

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			Rt.Index(w, r, tab)
		case "/create":
			Rt.CreateAccountPage(w, r, tab)
		case "/login":
			Rt.LoginPage(w, r, tab)
		case "/logout":
			Rt.LogOutHandler(w, r, tab)
		case "/home":
			Rt.HomeHandler(w, r, tab)
		case "/myprofil":
			Rt.Profil(w, r, tab)
		case "/filter":
			Rt.Filter(w, r, tab)
		case "/index":
			Rt.Indexfilter(w, r, tab)
		default:
			// Rt.Error404Handler(w, r)
			auth.Snippets(w, http.StatusNotFound)
		}
	}))
	// Démarrage du serveur
	fmt.Println("Server start at http://localhost:8080")
	error := http.ListenAndServe(":8080", nil)
	fmt.Println(error)
}
