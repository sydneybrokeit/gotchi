package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"time"
)

const TwitchWSS = "wss://pubsub-edge.twitch.tv"

var wsConn *websocket.Conn
var refreshTicker time.Ticker
var scope = []string{}
var userId string

func Reconnect() {
	log.Debugf("reconnecting...")
	var err error
	wsConn, _, err = websocket.DefaultDialer.Dial(TwitchWSS, nil)
	if err != nil {
		log.Fatalf("could not connect to WSS %v", err)
	}
	err = Subscribe()
	if err != nil {
		log.Fatal("could not subscribe")
	}
	return
}

func Subscribe() (err error) {
	scopes := []string{
		fmt.Sprintf("channel-bits-events-v2.%s", userId),
		// fmt.Sprintf("channel-subscribe-events-v1.%s", channelId),
	}
	data := WSData{
		AuthToken: outerToken.AccessToken,
		Topics:    scopes,
	}
	subscribeMessage := WSMessage{
		Type: "LISTEN",
		Data: data,
	}
	log.Debugf("%v", subscribeMessage)
	err = wsConn.WriteJSON(subscribeMessage)
	if err != nil {
		log.Errorf("couldn't subscribe properly: %v", err)
	} else {
		log.Debugf("subscribe to channel")
	}
	return
}

func Ping() {
	log.Debugf("starting ping!")
	pingMsg := WSMessage{Type: "PING"}
	err := wsConn.WriteJSON(pingMsg)
	if err != nil {
		log.Errorf("could not PING %v", err)
	} else {
		log.Debug("pinged successfully!")
	}
}

type WSMessage struct {
	Type  string `json:"type"`
	Nonce string `json:"nonce,omitempty"`
	Data  WSData `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

type WSData struct {
	// used for the request for topics
	Topics []string `json:"topics,omitempty"`
	// authtoken used to subscribe
	AuthToken string `json:"auth_token,omitempty"`
	// used for incoming messages
	Topic string `json:"topic,omitempty"`
	// string representation of JSON data (wtaf is this, Twitch, _why_.)
	MessageString string `json:"message,omitempty"`
	Message       WSEventMessage
}

// Fields outside Data are used for subscription messages, apparently?
type WSEventMessage struct {
	Data             WSEventMessageData `json:"data,omitempty"`
	UserName         string             `json:"user_name,omitempty"`
	DisplayName      string             `json:"display_name,omitempty"`
	ChannelName      string             `json:"channel_name,omitempty"`
	UserId           string             `json:"user_id,omitempty"`
	ChannelId        string             `json:"channel_id,omitempty"`
	Time             time.Time          `json:"time,omitempty"`
	SubPlan          string             `json:"sub_plan,omitempty"`
	SubPlanName      string             `json:"sub_plan_name,omitempty"`
	CumulativeMonths int                `json:"cumulative-months,omitempty"`
	StreakMonths     int                `json:"streak-months,omitempty"`
	// used in gift subs and anonymous gift subs
	Months int `json:"months"`
	// context is used for gift subs
	// subgift for a regular gift
	// anonsubgift for an anonymous gift
	Context    string `json:"context,omitempty"`
	SubMessage struct {
		Message string `json:"message,omitempty"`
		Emotes  []struct {
			Start int `json:"start,omitempty"`
			End   int `json:"end,omitempty"`
			Id    int `json:"id,omitempty"`
		} `json:"emotes,omitempty"`
	} `json:"sub_message,omitempty"`
	// only used for gifts
	RecipientId          string `json:"recipient_id,omitempty"`
	RecipientUserName    string `json:"recipient_user_name,omitempty"`
	RecipientDisplayName string `json:"recipient_display_name,omitempty"`
}

type WSEventMessageData struct {
	BitsUsed      int       `json:"bits_used,omitempty"`
	ChannelId     string    `json:"channel_id,omitempty"`
	ChatMessage   string    `json:"chat_message,omitempty"`
	Context       string    `json:"context,omitempty"`
	IsAnonymous   bool      `json:"is_anonymous,omitempty"`
	MessageId     string    `json:"message_id,omitempty"`
	MessageType   string    `json:"message_type,omitempty"`
	Time          time.Time `json:"time,omitempty"`
	TotalBitsUsed int       `json:"total_bits_used,omitempty"`
	UserId        string    `json:"user_id,omitempty"`
	UserName      string    `json:"user_name,omitempty"`
	Version       string    `json:"version,omitempty"`
}

func Pinger() {
	refreshTicker = *time.NewTicker(180 * time.Second)
	go func() {
		for {
			select {
			case <-refreshTicker.C:
				Ping()
			}
		}
	}()
}
