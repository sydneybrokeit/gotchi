package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"strings"
)

var outerToken *oauth2.Token
var StartGotchiChan chan bool
var thisGotchi *Gotchi
var Started bool
var IsReadyToHatchChan chan bool
var HatchChan chan bool

func init() {
	// log.SetLevel(log.DebugLevel)
}

func main() {
	go StartWSS()
	thisGotchi = new(Gotchi)
	thisGotchi.Hatched = false
	StartGotchiChan = make(chan bool)
	IsReadyToHatchChan = make(chan bool)
	go OauthWebInterface()
	<-StartGotchiChan
	Started = true
	internalService = newService(*extensionId, userId, secret)
	<-IsReadyToHatchChan
	thisGotchi.ReadyToHatch = true
	<-HatchChan
	// go thisGotchi.Do()
	var err error
	userId, err = getUserIdFromName("sydneythedev")
	if err != nil {
		log.Fatal("Couldn't get user ID: ", err)
	}
	log.Infof("using user id %v", userId)
	log.Infof("start gotchi!")
	Reconnect()
	Pinger()
	for {
		var message WSMessage
		err := wsConn.ReadJSON(&message)
		if err != nil {
			log.Errorf("couldn't read message: %v", err)
		} else {
			err := ProcessMessage(message)
			if err != nil {
				log.Errorf("couldn't process message: %v", err)
				continue
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
	if message.Type == "MESSAGE" {
		log.Infof("received message for topic %v", message.Data.Topic)
		err = json.Unmarshal([]byte(message.Data.MessageString),
			&message.Data.Message)
		if err != nil {
			return
		}
		topic := strings.Split(message.Data.Topic, ".")[0]
		switch {
		case topic == "channel-bits-events-v2":
			err = ProcessBitsEvent(message)
		case topic == "channel-subscribe-events-v1":
			err = ProcessSubscribeEvent(message)
		default:
			return
		}
	}
	log.Debugf("%v", message)
	return
}

func ProcessBitsEvent(message WSMessage) (err error) {
	return
}

func ProcessSubscribeEvent(message WSMessage) (err error) {
	return
}
