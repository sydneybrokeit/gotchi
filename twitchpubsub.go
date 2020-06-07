package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"time"
)

const TwitchWSS = "wss://pubsub-edge.twitch.tv"
var wsConn *twitchWS

type twitchWS struct {
	Conn *websocket.Conn
	Ticker time.Ticker
	Scope []string
	UserID string
}

func (ws *twitchWS) Reconnect() {
	log.Debugf("reconnecting...")
	var err error
	ws.Conn, _, err = websocket.DefaultDialer.Dial(TwitchWSS, nil)
	if err != nil {
		log.Fatal("could not connect to WSS")
	}
	err = ws.Subscribe()
	if err != nil {
		log.Fatal("could not subscribe")
	}
	return
}

func (ws *twitchWS) Subscribe() (err error) {
	scopes := []string{
		fmt.Sprintf("channel-bits-events-v2.%s", ws.UserID),
		// fmt.Sprintf("channel-subscribe-events-v1.%s", channelId),
	}
	data := WSData{
		AuthToken: outerToken.AccessToken,
		Topics: scopes,
	}
	subscribeMessage := WSMessage{
		Type: "LISTEN",
		Data: data,
	}
	log.Debug("%v", subscribeMessage)
	err = ws.Conn.WriteJSON(subscribeMessage)
	if err != nil {
		log.Errorf("couldn't subscribe properly: %v", err)
	} else {
		log.Debugf("subscribe to channel")
	}
	return
}

func (ws *twitchWS) Ping() {
	log.Debugf("starting ping!")
	pingMsg := WSMessage{Type: "PING"}
	err := ws.Conn.WriteJSON(pingMsg)
	if err != nil {
		log.Errorf("could not PING %v", err)
	} else {
		log.Debug("pinged successfully!")
	}
}

type WSMessage struct {
	Type string `json:"type"`
	Nonce string `json:"nonce,omitempty"`
	Data WSData `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

type WSData struct {
	Topics []string `json:"topics,omitempty"`
	AuthToken string `json:"auth_token,omitempty"`
}

func (ws *twitchWS) Pinger() {
	ws.Ticker = *time.NewTicker(60 * time.Second)
	go func() {
		for {
			select {
			case <-ws.Ticker.C:
				ws.Ping()
			}
		}
	}()
}
