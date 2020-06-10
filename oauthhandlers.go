package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
)


type Handler func(http.ResponseWriter, *http.Request) error

func HandleRoot(w http.ResponseWriter, r *http.Request) (err error) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Method == http.MethodPost && thisGotchi.Hatched != true {
		log.Warningf("starting the hatching!")
		thisGotchi, err = StartGotchi(r.FormValue("species"),
			r.FormValue("maxfood"),
			r.FormValue("maxlove"),
			r.FormValue("deplete"),
			r.FormValue("initial_food"),
			r.FormValue("initial_love"))
		if err != nil {
			log.Errorf("%v", err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		thisGotchi.Hatch()
		IsHatchedChan <- true
	}
	w.WriteHeader(http.StatusOK)
	if Started != true {
		w.Write([]byte(`<html><body><a href="/login">Login using Twitch</a></body></html>`))
	} else if thisGotchi.Hatched != true {
		w.Write([]byte(`
			<html><body><form method="POST">
			<label>Species</label><br />
			<select name="species">
				<option value="cat">Kitty!</option>
				<option value="dog">Pupper!</option>
			</select>
			<label>Max Food:</label><br />
			<input type="number" value="4" name="maxfood"><br />
			<label>Max Love:</label><br />
			<input type="number" value="4" name="maxlove"><br />
			<label>Countdown to Deplete (minutes)</label><br />
			<input type="number" value="15" name="deplete"><br />
			<label>Initial Food</label><br />
			<input type="number" value="2" name="initial_food"><br />
			<label>Initial Love</label><br />
			<input type="number" value="2" name="initial_love"><br />
			<input type="submit" value="Hatch!">
		`))
	} else {
		w.Write([]byte(`<html><body>I'm doing the thing!</body></html>`))
	}
	return
}

func HandleLogin(w http.ResponseWriter, r *http.Request) (err error) {
	session, err := cookieStore.Get(r, oauthSessionName)
	if err != nil {
		log.Printf("corrupted session %s -- generated new", err)
		err = nil
	}

	var tokenBytes [255]byte
	if _, err := rand.Read(tokenBytes[:]); err != nil {
		return err
	}

	state := hex.EncodeToString(tokenBytes[:])

	session.AddFlash(state, stateCallbackKey)

	if err = session.Save(r, w); err != nil {
		return
	}

	http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusTemporaryRedirect)

	return
}

func HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) (err error) {
	session, err := cookieStore.Get(r, oauthSessionName)
	if err != nil {
		log.Printf("corrupted session %s -- generated new", err)
		err = nil
	}

	// ensure we flush the csrf challenge even if the request is ultimately unsuccessful
	defer func() {
		if err := session.Save(r, w); err != nil {
			log.Printf("error saving session: %s", err)
		}
	}()

	switch stateChallenge, state := session.Flashes(stateCallbackKey), r.FormValue("state"); {
	case state == "", len(stateChallenge) < 1:
		err = errors.New("missing state challenge")
	case state != stateChallenge[0]:
		err = fmt.Errorf("invalid oauth state, expected '%s', got '%s'\n", state, stateChallenge[0])
	}

	if err != nil {
		return err
	}

	token, err := oauth2Config.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		return
	}

	// add the oauth token to session
	session.Values[oauthTokenKey] = token
	outerToken = token
	log.Debugf("Access token: %s\n", token.AccessToken)
	log.Debugf("%v", token)
	StartGotchiChan <- true
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return
}