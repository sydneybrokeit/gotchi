package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/hmschreck/envy"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
	"net/http"
)

const (
	stateCallbackKey = "oauth-state-callback"
	oauthSessionName = "oauth-session"
	oauthTokenKey    = "oauth-token"
)

// client id from Twitch Applications
var clientIDf = flag.String("clientid", "dummy", "client ID provided for the application by twitch")
var clientSecretf = flag.String("clientsecret", "dummy", "client secret provided for the application by twitch")
var redirectURLf = flag.String("redirecturl", "dummy", "redirect URL as authorized by twitch")
var port = flag.String("port", ":8080", "port to use, ex.: :8080")
var clientID string
var clientSecret string
var redirectURL string

var (
	scopes       = []string{"user:read:email", "chat:read", "bits:read"}
	oauth2Config *oauth2.Config
	cookieSecret = []byte("Please use a more sensible secret than this one")
	cookieStore  = sessions.NewCookieStore(cookieSecret)
)

func init() {
	flag.Parse()
	if *clientIDf == "dummy" {
		clientID = envy.GetEnv("CLIENT_ID", "dummy")
		if clientID == "dummy" {
			log.Fatalf("invalid client id given")
		}
	} else {
		clientID = *clientIDf
	}
	if *clientSecretf == "dummy" {
		clientSecret = envy.GetEnv("CLIENT_SECRET", "dummy")
		if clientSecret == "dummy" {
			log.Fatalf("invalid client secret given")
		}
	} else {
		clientSecret = *clientSecretf
	}
	if *redirectURLf == "dummy" {
		redirectURL = envy.GetEnv("REDIRECT_URL", "http://localhost:8080/redirect")
		if redirectURL == "dummy" {
			log.Fatalf("invalid redirect URL")
		}
	} else {
		redirectURL = *redirectURLf
	}
}

func OauthWebInterface() {
	// Gob encoding for gorilla/sessions
	gob.Register(&oauth2.Token{})
	oauth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		Endpoint:     twitch.Endpoint,
		RedirectURL:  redirectURL,
	}

	var middleware = func(h Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) (err error) {
			// parse POST body, limit request size
			if err = r.ParseForm(); err != nil {
				return err
			}

			return h(w, r)
		}
	}

	// errorHandling is a middleware that centralises error handling.
	// this prevents a lot of duplication and prevents issues where a missing
	// return causes an error to be printed, but functionality to otherwise continue
	// see https://blog.golang.org/error-handling-and-go
	var errorHandling = func(handler func(w http.ResponseWriter, r *http.Request) error) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := handler(w, r); err != nil {
				var errorString string = "Something went wrong! Please try again."
				var errorCode int = 500

				log.Println(err)
				w.Write([]byte(errorString))
				w.WriteHeader(errorCode)
				return
			}
		})
	}

	var handleFunc = func(path string, handler Handler) {
		http.Handle(path, errorHandling(middleware(handler)))
	}

	handleFunc("/", HandleRoot)
	handleFunc("/login", HandleLogin)
	handleFunc("/redirect", HandleOAuth2Callback)
	handleFunc("/hatch", HandleHatch)

	fmt.Println("Started running on http://localhost:8080")
	fmt.Println(http.ListenAndServe(":8080", nil))

}
