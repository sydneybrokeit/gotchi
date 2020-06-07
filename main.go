package main

import (
	"golang.org/x/oauth2"
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

	}
}
