package auth

const (
	ClientID     = "744272664350-36saa891gr1j1sv19v1n3ug41mb1ujgn.apps.googleusercontent.com" // our application id
	ClientSecret = "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"                                      // our secret client id
	RedirectURI  = "http://localhost:3000/auth/google/callback"                               // redirection after granted access
	AuthURL      = "https://accounts.google.com/o/oauth2/auth"                                // url to ask for access permission
	TokenURL     = "https://accounts.google.com/o/oauth2/token"                               // url to exchange permission with access token
	UserInfoURL  = "https://www.googleapis.com/oauth2/v2/userinfo"                            // url to exchange token with user info
)


// func main() {
// 	http.HandleFunc("/", handleIndex)
// 	http.HandleFunc("/auth/google/login", handleGoogleLogin)
// 	http.HandleFunc("/auth/google/callback", handleCallback)
