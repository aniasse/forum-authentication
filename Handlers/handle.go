package hdle

import (
	"fmt"
	"net/http"

	auth "forum/Authentication"
	db "forum/Database"
	Rt "forum/Routes"
)

func Handlers() {
	staticHandler := http.FileServer(http.Dir("templates"))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	tab, err := db.Init_db()
	if err != nil {
		fmt.Println(err)
		return
	}

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/": //default page
			Rt.Index(w, r, tab)

		case "/create": //create account page
			Rt.CreateAccountPage(w, r, tab)

		case "/auth/google/login": // googleAuth login page
			Rt.HandleGoogleLogin(w, r, tab)

		case "/auth/google/callback": //googleAuth response url
			Rt.HandleCallback(w, r, tab)

		case "/auth/github/login": // githubAuth login page
			Rt.HandleGitHubLogin(w, r, tab)

		case "/auth/github/callback": //githubAuth response url
			Rt.HandleGitHubCallback(w, r, tab)

		case "/login": //login page
			Rt.LoginPage(w, r, tab)

		case "/logout": //logout page
			Rt.LogOutHandler(w, r, tab)

		case "/home": //home page
			Rt.HomeHandler(w, r, tab)

		case "/myprofil/posts": //filtered created post page
			Rt.Profil(w, r, tab)

		case "/myprofil/favorites": //filtered liked post page
			Rt.Profil_fav(w, r, tab)

		case "/myprofil/comments": //filtered commented post page
			Rt.Profil_comment(w, r, tab)

		case "/filter": //filtered post by categorie page for registered
			Rt.Filter(w, r, tab)

		case "/index": //filtered post by categorie page for non-registered
			Rt.Indexfilter(w, r, tab)

		default: // page does not exist
			auth.Snippets(w, http.StatusNotFound)
		}
	}))

	// Launchinh server
	fmt.Println("📡----------------------------------------------------📡")
	fmt.Println("|                                                    |")
	fmt.Println("| 🌐 Server has started at \033[32mhttp://localhost:8080\033[0m 🟢  |")
	fmt.Println("|                                                    |")
	fmt.Println("📡----------------------------------------------------📡")
	errr := http.ListenAndServe(":8080", nil)
	if errr != nil {
		fmt.Printf("Erreur de serveur HTTP : %s\n", errr)
	}
}
