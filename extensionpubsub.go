package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
	log "github.com/sirupsen/logrus"
)

const (
	authHeaderName      string = "Authorization"
	authHeaderPrefix    string = "Bearer "
	authHeaderPrefixLen int    = len(authHeaderPrefix)
	minLegalTokenLength int    = authHeaderPrefixLen + 5
)

var secret []byte

type contextKeyType string
var extensionId = flag.String("extclientid", "dummy", "extension ")
var secretKey = flag.String("extsecret", "dummy", "extension secret")
func init() {
	flag.Parse()
	if *secretKey == "" || *secretKey == "dummy" {
		log.Fatalf("couldn't get a real extension secret")
	}
	var err error
	secret, err = base64.StdEncoding.DecodeString(*secretKey)
	if err != nil {

	}
}
var internalService *service
type service struct {
	parser    jwt.Parser
	clientID  string
	ownerID   string
	secret    []byte
	nextPongs map[string]time.Time
	mutex     sync.Mutex
}

type pubSubMessage struct {
	ContentType string   `json:"content_type"`
	Targets     []string `json:"targets"`
	Message     string   `json:"message"`
}

func newService(clientID string, ownerID string, secret []byte) *service {
	return &service{
		parser:    jwt.Parser{ValidMethods: []string{"HS256"}},
		clientID:  clientID,
		ownerID:   ownerID,
		secret:    secret,
		nextPongs: make(map[string]time.Time),
	}
}

func (s *service) inCooldown(channelID string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if next, found := s.nextPongs[channelID]; found && next.After(time.Now()) {
		return true
	}

	s.nextPongs[channelID] = time.Now().Add(time.Second)
	return false
}

func (s *service) send(channelID, message string) {
	if s.inCooldown(channelID) { // don't spam PubSub or you'll be rate limited
		log.Println("Service is in cooldown")
		return
	}

	m := pubSubMessage{
		"application/json",
		[]string{"broadcast"},
		message,
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.twitch.tv/extensions/message/%v", channelID), b)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("Client-Id", s.clientID)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("%s%v", authHeaderPrefix, s.newJWT(channelID)))

	log.Printf("Message: %s via PubSub for channel %s\n", message, channelID)
	res, err := http.DefaultClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil { }
	log.Warningf("received status %v:%v", res.StatusCode, string(body))

	if err != nil {
		log.Println(err)
	}
}

func (s *service) newJWT(channelID string) string {
	var expiration = time.Now().Add(time.Minute * 3).Unix()

	claims := jwtAuthClaims{
		UserID:    s.ownerID,
		ChannelID: channelID,
		Role:      "external",
		Permissions: jwtPermissions{
			Send: []string{"broadcast"},
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		log.Println(err)
	}

	log.Printf("Generated JWT: %s\n", tokenString)

	return tokenString
}

type jwtAuthClaims struct {
	OpaqueUserID string         `json:"opaque_user_id,omitempty"`
	UserID       string         `json:"user_id"`
	ChannelID    string         `json:"channel_id,omitempty"`
	Role         string         `json:"role"`
	Permissions  jwtPermissions `json:"pubsub_perms"`
	jwt.StandardClaims
}

type jwtPermissions struct {
	Send   []string `json:"send,omitempty"`
	Listen []string `json:"listen,omitempty"`
}