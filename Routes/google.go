package Route

import (
	"encoding/json"
	"fmt"
	auth "forum/Authentication"
	db "forum/Database"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// handleGoogleLogin redirects the user to the google auth interface
func HandleGoogleLogin(w http.ResponseWriter, r *http.Request, tab db.Db) {
	auth.CheckCookie(w, r, tab)
	if r.Method != "GET" {
		fmt.Printf("❌ cannot access to /auth/google/login page with %s method", r.Method)
		auth.Snippets(w, 405)
		return
	}
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=profile email&response_type=code", auth.GoAuthURL, auth.GoClientID, auth.GoRedirectURI)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleGoogleCallback is called by Google after authentication and returns a token that can be used for API calls
func HandleCallback(w http.ResponseWriter, r *http.Request, tab db.Db) {
	auth.CheckCookie(w, r, tab)
	if r.Method != "GET" {
		fmt.Printf("❌ cannot access to /auth/google/callback page with %s method", r.Method)
		auth.Snippets(w, 405)
		return
	}
	code := r.URL.Query().Get("code") //retireving the code for access permission
	if code == "" {
		http.Error(w, "Code missing", http.StatusBadRequest)
		return
	}
	// establishing the post request to exchange the permission with an access token
	data := url.Values{}   // we use url.values in order to ensure well url encoding and more security against injections
	data.Set("code", code) // setting the permission code
	data.Set("client_id", auth.GoClientID)
	data.Set("client_secret", auth.GoClientSecret)
	// telling to the google api that the request is based upon a permission code
	//meaning that we have the user's consentment
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", auth.GoRedirectURI) // setting url where goolgle api will send its response

	tokenResp, err := http.Post(auth.GoTokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
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

	userInfoResp, err := http.Get(fmt.Sprintf("%s?access_token=%s", auth.GoUserInfoURL, accessToken))
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
		Name       interface{}
		FamilyName interface{}
		Email      interface{}
		Id         interface{}
	}{
		Name:       userInfo["given_name"],
		FamilyName: userInfo["family_name"],
		Email:      userInfo["email"],
		Id:         userInfo["id"],
	}
	fmt.Println("     given         ", final.Email, final.Name, final.Id)
	if final.Email != nil && final.Id != nil && final.Name != nil {
		Email := (final.Email).(string)
		Name := (final.Name).(string)
		Id := (final.Id).(string)
		FamilyName := ""
		if final.FamilyName == nil {
			FamilyName = Name
		} else {
			FamilyName = (final.FamilyName).(string)
		}

		Connection0auth(tab, Email, Name, FamilyName, w, r, Id)
	} else {
		//pas d'email
		message := "missing personal information in Google"

		// message := "connecting to the forum requires a valid email address and personal details, please go to your google email settings to work the magic. See you soon!"
		formlogin := Register{Username: "", Password: "", Message: message}
		auth.DisplayFilewithexecute(w, "templates/register.html", formlogin, http.StatusBadRequest)
		return
	}
}
