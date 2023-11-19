package Route

import (
	"encoding/json"
	"fmt"
	Google "forum/Authentication"
	auth "forum/Authentication"
	db "forum/Database"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofrs/uuid/v5"
)

// handleGoogleLogin redirects the user to the google auth interface
func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=profile email&response_type=code", Google.GoAuthURL, Google.GoClientID, Google.GoRedirectURI)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleGoogleCallback is called by Google after authentication and returns a token that can be used for API calls
func HandleCallback(w http.ResponseWriter, r *http.Request, tab db.Db) {
	auth.CheckCookie(w, r, tab)

	code := r.URL.Query().Get("code") //retireving the code for access permission
	if code == "" {
		http.Error(w, "Code missing", http.StatusBadRequest)
		return
	}
	// establishing the post request to exchange the permission with an access token
	data := url.Values{}   // we use url.values in order to ensure well url encoding and more security against injections
	data.Set("code", code) // setting the permission code
	data.Set("client_id", Google.GoClientID)
	data.Set("client_secret", Google.GoClientSecret)
	// telling to the google api that the request is based upon a permission code
	//meaning that we have the user's consentment
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", Google.GoRedirectURI) // setting url where goolgle api will send its response

	tokenResp, err := http.Post(Google.GoTokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to exchange token: %s", err), http.StatusInternalServerError)
		return
	}
	defer tokenResp.Body.Close()
	//--reading and storing the response
	tokenData, err := io.ReadAll(tokenResp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read token data: %s", err), http.StatusInternalServerError)
		return
	}

	var token map[string]interface{}
	json.Unmarshal(tokenData, &token)

	accessToken := token["access_token"].(string) //  retrievieng the access token in the token response body

	userInfoResp, err := http.Get(fmt.Sprintf("%s?access_token=%s", Google.GoUserInfoURL, accessToken))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch user info: %s", err), http.StatusInternalServerError)
		return
	}
	defer userInfoResp.Body.Close()

	userInfoData, err := io.ReadAll(userInfoResp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read user info data: %s", err), http.StatusInternalServerError)
		return
	}

	var userInfo map[string]interface{}
	json.Unmarshal(userInfoData, &userInfo)
	fmt.Println("user info here", userInfo)

	var final = struct {
		Name       string
		FamilyName string
		Email      string
		Id         string
	}{
		FamilyName: userInfo["family_name"].(string),
		Name:       userInfo["given_name"].(string),
		Email:      userInfo["email"].(string),
		Id:         userInfo["id"].(string),
	}
	foundEmail := auth.GetDatafromBA(tab.Doc, final.Email, "email", db.User)
	fmt.Println("find it           ", final.Name)

	// verifier si le user existe deja sinon lui creer un compte dans les deux cas redirections vers /home
	if foundEmail {
		iduser, _, _ := auth.HelpersBA("users", tab, "id_user", "WHERE email='"+final.Email+"'", "")
		auth.CreateSession(w, iduser, tab)
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		newid, err := uuid.NewV4()
		if err != nil {
			fmt.Println("erreur avec le uuid niveau create account")
			auth.Snippets(w, http.StatusInternalServerError)
			return
		}
		// password hash
		hashpassword, errorhash := auth.HashPassword(final.Id)
		if errorhash != nil {
			fmt.Println("error hash")
			auth.Snippets(w, http.StatusInternalServerError)
			return
		}
		//creation pseudo
		username := auth.GenerateUsername(final.Name, tab)

		values := "('" + newid.String() + "','" + final.Email + "','" + final.Name + "','" + username + "','" + final.FamilyName + "','" + hashpassword + "','../static/front-tools/images/profil.jpeg','../static/front-tools/images/mur.png')"
		attributes := "(id_user,email,name,username,surname, password,pp,pc)"
		error := tab.INSERT(db.User, attributes, values)
		if error != nil {
			fmt.Println("something wrong")
			fmt.Println("error", error)
			auth.Snippets(w, http.StatusInternalServerError)
			return

		}
		valuesession := "('" + newid.String() + "')"
		attributessession := "(user_id)"
		errorsession := tab.INSERT("sessions", attributessession, valuesession)
		if errorsession != nil {
			fmt.Println("something wrong with insert session", errorsession)
			fmt.Println("error", error)
			auth.Snippets(w, http.StatusInternalServerError)
			return

		}
		// fmt.Println("un w", w)
		// creation of the session
		auth.CreateSession(w, newid.String(), tab)
		// fmt.Println("deux w", w)

		//redirecting the user to their home page
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
	// t, _ := template.ParseFiles("templates/success.html")
	// t.Execute(w, final)
}
