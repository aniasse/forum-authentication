package auth

import (
	"fmt"
	"regexp"

	// auth "forum/Authentication"
	db "forum/Database"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateSession allows you to create a session for the current user
func CreateSession(w http.ResponseWriter, iduser string, tab db.Db) {
	token, err := uuid.NewV4()
	if err != nil {
		fmt.Println("erreur uuid dans la creation de session")
		Snippets(w, http.StatusInternalServerError)
	}
	sessionToken := token.String()
	expiresAt := time.Now().Add(1800 * time.Second)
	fmt.Println("expire a", expiresAt.String())

	//update session dans la base de données
	error := tab.UPDATE("sessions", "id_session='"+sessionToken+"',expireat='"+expiresAt.String()+"'", "WHERE user_id="+"'"+iduser+"'")
	if error != nil {
		fmt.Println("erreur update", error)
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
		Path:    "/",
		// SameSite: http.SameSiteDefaultMode,
	})
}

// Snippets is a function which allows you to return an error page,
//
//	it receives an http.ResponseWriter as an argument from the handler function
//	and the the status of the error to be specified.
func Snippets(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	error_file := template.Must(template.ParseFiles("templates/error.html"))
	error_file.Execute(w, strconv.Itoa(statusCode))
	fmt.Println("⚠ ERROR ⚠: ", statusCode, "❌ ")
}

// the function DisplayFile allows you to display the page given as an argument
//
// while managing any errors that may occur.
func DisplayFile(w http.ResponseWriter, templatePath string) {
	file, errExecutionFile := template.ParseFiles(templatePath)
	if errExecutionFile != nil {
		fmt.Println("Probléme de parsing ou d'execution de fichier", templatePath)
		Snippets(w, http.StatusInternalServerError)
		return
	}

	errExecutionFile = file.Execute(w, nil)
	if errExecutionFile != nil {
		fmt.Println("Probléme de parsing ou d'execution de fichier", templatePath)
		Snippets(w, http.StatusInternalServerError)
		return
	}
	fmt.Println("✅ File displays!", templatePath)
}

func DisplayFilewithexecute(w http.ResponseWriter, templatePath string, execute interface{}, status int) {
	w.WriteHeader(status)
	file, errparsefile := template.ParseFiles(templatePath)
	if errparsefile != nil {
		fmt.Println("du mal a parser")
		Snippets(w, http.StatusInternalServerError)
		return
	}
	errExecutionFile := file.Execute(w, execute)
	if errExecutionFile != nil {
		fmt.Println("ne peut executer", errExecutionFile)
		Snippets(w, http.StatusInternalServerError)
		return
	}
	fmt.Println("✅ File displays!", templatePath)
}

// validation email user
func ValidMailAddress(address string) (string, bool) {

	regex := "^[A-Za-z0-9._%+-]{2,}@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$"

	// Test de la chaîne "peach" avec la regex
	match, err := regexp.MatchString(regex, address)
	// Vérification des erreurs
	if err != nil {
		fmt.Println("Erreur lors de la correspondance de la regex:", err)
	}

	// Affichage du résultat
	fmt.Println(match, "de l'email", address)
	return address, match
}

func CheckCookie(w http.ResponseWriter, r *http.Request, tab db.Db) {
	c, errc := r.Cookie("session_token")
	if errc != nil {
		fmt.Println("pas de cookie session")
	} else {

		idviasession, err, _ := HelpersBA("sessions", tab, "user_id", "WHERE id_session='"+c.Value+"'", "")
		// fmt.Println("here", s, "error", err)
		if err != nil {
			fmt.Println("erreur du serveur", err)
		}
		if idviasession != "" {
			fmt.Println("cookie valide,affichage de /home", idviasession)
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}
	}
}

func FieldsLimited(field string, min, max int) bool {
	return len(field) >= min && len(field) < max
}

func NotAllow(s string) bool {
	return strings.Contains(s, "'") || strings.Contains(s, "\"")
}
func GenerateUsername(name string, tab db.Db) string {
	username := name + strconv.Itoa(rand.Intn(101))
	_, _, confirmusername := HelpersBA("users", tab, "username", " WHERE username='"+username+"'", username)
	if confirmusername {
		return GenerateUsername(name, tab)
	}
	return username
}
func Familyname(name string) (string, string) {
	name = strings.Trim(name, " ")
	arrayname := strings.Split(name, " ")
	limit := len(arrayname)
	if limit == 1 {
		return name, name
	}
	lastfamilyname := arrayname[limit-1]
	familynames := arrayname[0 : limit-1]
	firstfamilyname := strings.Join(familynames, " ")
	return firstfamilyname, lastfamilyname
}
