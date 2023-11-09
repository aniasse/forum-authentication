package Route

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	Google "forum/Authentication"
)

// handleGoogleLogin redirects the user to the google auth interface
func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=profile email&response_type=code", Google.AuthURL, Google.ClientID, Google.RedirectURI)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleGoogleCallback is called by Google after authentication and returns a token that can be used for API calls
func HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code") //retireving the code for access permission
	if code == "" {
		http.Error(w, "Code missing", http.StatusBadRequest)
		return
	}
	// establishing the post request to exchange the permission with an access token
	data := url.Values{}   // we use url.values in order to ensure well url encoding and more security against injections
	data.Set("code", code) // setting the permission code
	data.Set("client_id", Google.ClientID)
	data.Set("client_secret", Google.ClientSecret)
	// telling to the google api that the request is based upon a permission code
	//meaning that we have the user's consentment
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", Google.RedirectURI) // setting url where goolgle api will send its response

	tokenResp, err := http.Post(Google.TokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to exchange token: %s", err), http.StatusInternalServerError)
		return
	}
	defer tokenResp.Body.Close()
	//--reading and storing the response
	tokenData, err := ioutil.ReadAll(tokenResp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read token data: %s", err), http.StatusInternalServerError)
		return
	}

	var token map[string]interface{}
	json.Unmarshal(tokenData, &token)

	accessToken := token["access_token"].(string) //  retrievieng the access token in the token response body

	userInfoResp, err := http.Get(fmt.Sprintf("%s?access_token=%s", Google.UserInfoURL, accessToken))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch user info: %s", err), http.StatusInternalServerError)
		return
	}
	defer userInfoResp.Body.Close()

	userInfoData, err := ioutil.ReadAll(userInfoResp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read user info data: %s", err), http.StatusInternalServerError)
		return
	}

	var userInfo map[string]interface{}
	json.Unmarshal(userInfoData, &userInfo)
	fmt.Println("user info here", userInfo)

	var final = struct {
		Name  string
		Email string
		Id    string
	}{
		Name:  userInfo["name"].(string),
		Email: userInfo["email"].(string),
		Id:    userInfo["id"].(string),
	}

	t, _ := template.ParseFiles("templates/success.html")
	t.Execute(w, final)
}
