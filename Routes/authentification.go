package Route

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"

	auth "forum/Authentification"
	db "forum/Database"
)

// the credentials structure stores the data of the logged in user
type Credentials struct {
	id       string
	Name     string
	Surname  string
	Username string
	Email    string
	Password string
}
type Create struct {
	Surname    string
	Name       string
	Username   string
	Email      string
	Password   string
	Confirmpwd string
}
type Register struct {
	Username string
	Password string
	Message  string
}
type Message struct {
	Errormessage string
	CreateForm   Create
}

// CreateAccountPage manages the account creation and once successful redirects
//
//	the user to their home page, otherwise it displays an error page
func CreateAccountPage(w http.ResponseWriter, r *http.Request, tab db.Db) {
	messageE := ""
	// check if the user is not already logged in to be able to access this page
	auth.CheckCookie(w, r, tab)
	// method verification
	if r.Method == "GET" {
		fmt.Println("here")
		auth.DisplayFile(w, "templates/createacount.html")
		return
	} else if r.Method == "POST" {
		// retrieving query data
		surname := strings.TrimSpace(r.FormValue("surname"))
		name := strings.TrimSpace(r.FormValue("name"))
		username := strings.TrimSpace(r.FormValue("username"))
		email := strings.TrimSpace(r.FormValue("email"))
		password := strings.TrimSpace(r.FormValue("password"))
		confirmpwd := strings.TrimSpace(r.FormValue("confirmpwd"))

		formcreate := Create{Surname: surname, Name: name, Username: username, Email: email, Password: password, Confirmpwd: confirmpwd}

		// fmt.Println(name, username, email, password)
		//check that the fields are not empty
		if auth.FieldsLimited(name, 2, 15) && auth.FieldsLimited(surname, 2, 15) && auth.FieldsLimited(username, 2, 15) && auth.FieldsLimited(email, 10, 133) && auth.FieldsLimited(password, 8, 15) && auth.FieldsLimited(confirmpwd, 8, 15) {
			// if name != "" && username != "" && email != "" && password != "" && surname != "" && confirmpwd != "" {
			// fmt.Println(name, username, password, email)
			if auth.NotAllow(name) || auth.NotAllow(username) || auth.NotAllow(surname) || auth.NotAllow(email) || auth.NotAllow(password) {
				fmt.Println("esq amna apostorof")
				messageE = "Character \"'\" not allowed"
				message := Message{Errormessage: messageE, CreateForm: formcreate}
				auth.DisplayFilewithexecute(w, "templates/createacount.html", message, http.StatusBadRequest)
				return
			}
			creds := &Credentials{}
			//check that the email and username have not already been used
			validemail, right := auth.ValidMailAddress(email)
			if !right {
				fmt.Println("mauvais format d'email", validemail)
				messageE = "Bad email format"
				// auth.Snippets(w, http.StatusBadRequest)
				message := Message{Errormessage: messageE, CreateForm: formcreate}
				auth.DisplayFilewithexecute(w, "templates/createacount.html", message, http.StatusBadRequest)
				return
			}

			email = validemail
			_, _, confirmemail := auth.HelpersBA("users", tab, "email", " WHERE email='"+email+"'", email)
			_, _, confirmusername := auth.HelpersBA("users", tab, "username", " WHERE username='"+username+"'", username)

			if confirmemail || confirmusername {
				messageE = "email/username already used"
				message := Message{Errormessage: messageE, CreateForm: formcreate}
				auth.DisplayFilewithexecute(w, "templates/createacount.html", message, http.StatusBadRequest)
				fmt.Println("⚠ ERROR ⚠:❌  email ", email, "ou username existant", username)
				return
			}

			if password != confirmpwd {
				fmt.Println("mots de passe ne match pas")
				messageE = "Incorrect password confirmation"
				message := Message{Errormessage: messageE, CreateForm: formcreate}
				auth.DisplayFilewithexecute(w, "templates/createacount.html", message, http.StatusBadRequest)
				return
			}

			// password hash
			hashpassword, errorhash := auth.HashPassword(password)
			if errorhash != nil {
				fmt.Println("error hash")
				auth.Snippets(w, http.StatusInternalServerError)
				return
			}

			// store current user information
			newid, err := uuid.NewV4()
			if err != nil {
				fmt.Println("erreur avec le uuid niveau create account")
				auth.Snippets(w, http.StatusInternalServerError)
				return
			}
			creds = &Credentials{Name: name, Username: username, Email: email, Password: hashpassword, id: newid.String(), Surname: surname}
			//save user in database
			// fmt.Println("creds", creds)
			values := "('" + creds.id + "','" + creds.Email + "','" + creds.Name + "','" + creds.Username + "','" + creds.Surname + "','" + creds.Password + "','../static/front-tools/images/profil.jpeg','../static/front-tools/images/mur.png')"
			attributes := "(id_user,email,name,username,surname, password,pp,pc)"
			error := tab.INSERT(db.User, attributes, values)
			if error != nil {
				fmt.Println("something wrong")
				fmt.Println("error", error)
				auth.Snippets(w, http.StatusInternalServerError)
				return

			}
			valuesession := "('" + creds.id + "')"
			attributessession := "(user_id)"
			errorsession := tab.INSERT("sessions", attributessession, valuesession)
			if errorsession != nil {
				fmt.Println("something wrong with insert session", errorsession)
				fmt.Println("error", error)
				auth.Snippets(w, http.StatusInternalServerError)
				return

			}
			// creation of the session
			auth.CreateSession(w, creds.id, tab)
			//redirecting the user to their home page
			http.Redirect(w, r, "/home", http.StatusSeeOther)
		} else {

			if !auth.FieldsLimited(name, 2, 15) || !auth.FieldsLimited(surname, 2, 15) || !auth.FieldsLimited(username, 2, 15) {
				messageE = "the name, surname and username must be between 2 to 15 characters"
				message := Message{Errormessage: messageE, CreateForm: formcreate}
				auth.DisplayFilewithexecute(w, "templates/createacount.html", message, http.StatusBadRequest)
				return
			} else if !auth.FieldsLimited(email, 10, 133) {
				messageE = "the Email must be between 10 to 132 characters"
				message := Message{Errormessage: messageE, CreateForm: formcreate}
				auth.DisplayFilewithexecute(w, "templates/createacount.html", message, http.StatusBadRequest)
				return
			} else {
				messageE = "the password and confirmpassword must be between 8 to 15 characters"
				message := Message{Errormessage: messageE, CreateForm: formcreate}
				auth.DisplayFilewithexecute(w, "templates/createacount.html", message, http.StatusBadRequest)
				return
			}
			// return

		}

	} else {
		// if neither method post nor method get is used
		auth.Snippets(w, http.StatusMethodNotAllowed)
		return

	}
}

func LoginPage(w http.ResponseWriter, r *http.Request, tab db.Db) {
	// check if the user is not already logged in to be able to access this page
	auth.CheckCookie(w, r, tab)
	// method verification
	if r.Method == "GET" {
		fmt.Println("affichage simple de la page du login")
		auth.DisplayFile(w, "templates/register.html")
		return
	} else if r.Method == "POST" {
		fmt.Println("requete post")
		// retrieving query data
		username := r.FormValue("username")
		password := r.FormValue("password")
		creds := Credentials{}
		//checks that the fields in the query are not null
		if username != "" && password != "" {
			//change
			if auth.NotAllow(username) {
				fmt.Println("here")
				message := "presence de \" ou ' ou injection XSS dans le username"
				formlogin := Register{Username: "", Password: password, Message: message}
				auth.DisplayFilewithexecute(w, "templates/register.html", formlogin, http.StatusBadRequest)
				return

			}
			if auth.NotAllow(password) {
				message := "presence de \" ou ' ou injection XSS dans le mdp"
				formlogin := Register{Username: username, Password: "", Message: message}
				auth.DisplayFilewithexecute(w, "templates/register.html", formlogin, http.StatusBadRequest)
				return
			}
			//give the user the possibility to enter an email or a nickname
			giveUsername := auth.GetDatafromBA(tab.Doc, username, "username", db.User)
			giveEmail := auth.GetDatafromBA(tab.Doc, username, "email", db.User)

			if giveEmail {
				// fmt.Println("email given")
				values := "WHERE email =" + "'" + username + "'"
				// samePassword, errpassword = tab.GetData("email", db.User, values)
				replaceEmailbyusername, err, _ := auth.HelpersBA("users", tab, "username", values, "")
				if err != nil {
					if err == sql.ErrNoRows {
						fmt.Println("erreur sql dans login page")
						auth.Snippets(w, http.StatusUnauthorized)
						return
					}
					fmt.Println("erreur interne dans login page")
					auth.Snippets(w, http.StatusInternalServerError)
					return
				}
				// fmt.Println("replace", replaceEmailbyusername)
				creds = Credentials{Username: replaceEmailbyusername, Password: password}

			} else if giveUsername {
				creds = Credentials{Username: username, Password: password}
				// samePassword, errpassword = tab.GetData("password", "users", values)
				// fmt.Println(2)

			}
			if !giveEmail && !giveUsername {
				message := "mauvais credentials"
				formlogin := Register{Username: username, Password: password, Message: message}
				fmt.Println("user,pass", username, password)
				// auth.Snippets(w, http.StatusUnauthorized)
				auth.DisplayFilewithexecute(w, "templates/register.html", formlogin, http.StatusUnauthorized)
				return
			}
			values := "WHERE username =" + "'" + creds.Username + "'"
			samePassword, errpassword, _ := auth.HelpersBA("users", tab, "password", values, "")
			// fmt.Println("same", samePassword)
			if errpassword != nil {
				if errpassword == sql.ErrNoRows {
					fmt.Println("pas de rows dans login page")
					auth.Snippets(w, http.StatusUnauthorized)
				}
				fmt.Println("autre erreur interne login page")
				auth.Snippets(w, http.StatusInternalServerError)
				return
			}
			// var store Credentials

			store := Credentials{Username: creds.Username, Password: samePassword}
			if !auth.CheckPasswordHash(password, store.Password) {
				fmt.Println("probleme hashage")
				// fmt.Printf("user,pass", username, password)

				message := "le mot de passe ne correspond pas"
				formlogin := Register{Username: username, Password: password, Message: message}
				// auth.Snippets(w, http.StatusUnauthorized)
				auth.DisplayFilewithexecute(w, "templates/register.html", formlogin, http.StatusUnauthorized)
				return
			}
			iduser, _, _ := auth.HelpersBA("users", tab, "id_user", "WHERE username='"+creds.Username+"'", "")
			auth.CreateSession(w, iduser, tab)

			http.Redirect(w, r, "/home", http.StatusSeeOther)

		} else {
			//champs vide
			// auth.Snippets(w, http.StatusBadRequest)
			fmt.Println("credentials vides")
			// fmt.Printf("user,pass", username, password)

			message := "credentials vides"
			formlogin := Register{Username: username, Password: password, Message: message}
			auth.DisplayFilewithexecute(w, "templates/register.html", formlogin, http.StatusBadRequest)
		}
		//type de methode
	} else {
		auth.Snippets(w, http.StatusMethodNotAllowed)
	}
}

func LogOutHandler(w http.ResponseWriter, r *http.Request, tab db.Db) {
	if r.Method != "POST" {
		auth.Snippets(w, http.StatusMethodNotAllowed)
		return
	}
	_, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			http.Redirect(w, r, "/", http.StatusSeeOther)
			auth.Snippets(w, http.StatusUnauthorized)
			return
		}
		// For any other type of error, return an internal server error
		auth.Snippets(w, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func HomeHandler(w http.ResponseWriter, r *http.Request, tab db.Db) {

	if r.Method != "GET" && r.Method != "POST" {
		auth.Snippets(w, http.StatusMethodNotAllowed)
		return
	}
	c, errc := r.Cookie("session_token")
	if errc != nil {
		fmt.Println("pas de cookie session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return

	} else {

		s, err, _ := auth.HelpersBA("sessions", tab, "user_id", "WHERE id_session='"+c.Value+"'", "")
		fmt.Println("here", s, "error", err)
		if err != nil {
			fmt.Println("erreur du serveur", err)
		}
		if s != "" {
			fmt.Println("cookie valide,affichage de /home", s)

			Communication(w, r, s, "/home")
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func Error404Handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/error.html"))
	tmpl.Execute(w, nil)
}
