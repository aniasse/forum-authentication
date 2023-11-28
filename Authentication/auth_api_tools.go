package auth

const (
	//-- google side
	GoClientID     = "744272664350-36saa891gr1j1sv19v1n3ug41mb1ujgn.apps.googleusercontent.com" // our application id
	GoClientSecret = "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"                                      // our secret client id
	GoRedirectURI  = "http://localhost:8080/auth/google/callback"                               // redirection after granted access
	GoAuthURL      = "https://accounts.google.com/o/oauth2/auth"                                // url to ask for access permission
	GoTokenURL     = "https://accounts.google.com/o/oauth2/token"                               // url to exchange permission with access token
	GoUserInfoURL  = "https://www.googleapis.com/oauth2/v2/userinfo"                            // url to exchange token with user info

	// -- github side
	GitClientID     = "9b8a06490595c8ed0010"                     // our application id
	GitClientSecret = "77aa4b05065248d269efe72d55ca2bbbdae7fa2b" // our secret client id //77aa4b05065248d269efe72d55ca2bbbdae7fa2b
	GitRedirectURI  = "http://localhost:8080/auth/github/callback"
	GitAuthURL      = "https://github.com/login/oauth/authorize"
)
