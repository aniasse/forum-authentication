package auth

import (
	"fmt"
	db "forum/Database"
	"net/http"
)

func ComSession_Checker(w http.ResponseWriter, r *http.Request, database db.Db) (bool, *http.Cookie) {
	c, errc := r.Cookie("session_token")
	if errc != nil {
		fmt.Println("pas de cookie session")
		return false, nil

	} else {

		s, err, _ := HelpersBA("sessions",database, "user_id", "WHERE id_session='"+c.Value+"'", "")
		// fmt.Println("here", s, "error", err)
		if err != nil {
			fmt.Println("erreur du serveur", err)
		}
		if s == "" {
			fmt.Println("cookie invalide,affichage de /", s, "verif vide")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return false, nil
		}
	}
	return true, c
}
