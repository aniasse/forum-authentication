package Route

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	Github "forum/Authentication"
)

func HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	// redirecting user to githubAuth interface
	parameter := url.Values{}
	parameter.Set("client_id", Github.GitClientID)
	parameter.Set("redirect_uri",  Github.GitRedirectURI)
	parameter.Set("scope", "user:email") // ask permission to user's email
	parameter.Set("response_type", "code")
	redirectURL :=  Github.GitAuthURL + "?" + parameter.Encode()
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	// Retrieving permission code
	code := r.URL.Query().Get("code")
	// fmt.Println("code is here", code)

	// exchnaging code with token access
	tokenURL := "https://github.com/login/oauth/access_token"
	data := url.Values{}
	data.Set("client_id",  Github.GitClientID)
	data.Set("client_secret",  Github.GitClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri",  Github.GitRedirectURI)
	data.Set("grant_type", "authorization_code")

	tokenResp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tokenResp.Body.Close()

	//--reading and storing the response
	tokenData, err := ioutil.ReadAll(tokenResp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read token data: %s", err), http.StatusInternalServerError)
		return
	}
	fmt.Println("tokenData", string(tokenData))
	accessToken := strings.Split(strings.Split(string(tokenData), "=")[1], "&")[0]
	fmt.Println("token ->", accessToken)

	// Using the access token to fetch user's information
	client := &http.Client{}
	githubReq, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	githubReq.Header.Set("Authorization", "token "+accessToken)

	userInfoResp, errinf := client.Do(githubReq)
	if errinf != nil {
		http.Error(w, errinf.Error(), http.StatusInternalServerError)
		return
	}
	defer userInfoResp.Body.Close()

	var userResp map[string]interface{}
	err = json.NewDecoder(userInfoResp.Body).Decode(&userResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// debugging
	for i, v := range userResp {
		fmt.Printf("%v : %v\n", i, v)
	}

	var final = struct {
		Name  string
		Email string
		Id    string
	}{
		Name:  userResp["login"].(string),
		Email: userResp["email"].(string),
		Id:    userResp["id"].(string),
	}

	t, _ := template.ParseFiles("templates/success.html")
	t.Execute(w, final)
}
