package main

import (
	"golang.org/x/oauth2"
	log "github.com/sirupsen/logrus"
)

var outerToken *oauth2.Token
var StartGotchiChan chan bool

func init() {
}

func main() {
	go OauthWebInterface()
	<-StartGotchiChan
	wsConn.Pinger()
	for {
		var message WSMessage
		err := wsConn.Conn.ReadJSON(&message)
		if err != nil {
			log.Errorf("couldn't read message: %v", err)
		} else {
			err := ProcessMessage(message)
			if err != nil {
				log.Errorf("couldn't process message: %v", err)
			}
		}
	}
}

// triggers
// "!feedMeal" -- feed a meal, increasing hunger
// "!feedSnack" -- feed a snack (must be earned by random game)
//    - check to see if user has a snack
//    - if snack, give, increasing both hunger and affection
// "!guessLeft"/"!guessRight" -- play the guessing game
//    - select left or right at random
//    - if correct, grant user 1 snacc
// general donation (no trigger) -- "read message"
// "!pet"? -- increase affection by 1 bar
func ProcessMessage(message WSMessage) (err error) {
	return
}
