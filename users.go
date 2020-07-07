package main

import (
	"fmt"
	"github.com/imroc/req"
)

type TwitchUser struct {
	Data []struct {
		ID              string `json:"id"`
		Login           string `json:"login"`
		DisplayName     string `json:"display_name"`
		Type            string `json:"type"`
		BroadcasterType string `json:"broadcaster_type"`
		Description     string `json:"description"`
		ProfileImageURL string `json:"profile_image_url"`
		OfflineImageURL string `json:"offline_image_url"`
		ViewCount       int    `json:"view_count"`
		Email           string `json:"email"`
	} `json:"data"`
}

// pass in a username and get a userID
func getUserIdFromName(name string) (userId string, err error) {
	var user TwitchUser
	requestUrl := fmt.Sprintf("https://api.twitch.tv/helix/users?login=%s", name)
	AuthString := fmt.Sprintf("Bearer %s", outerToken.AccessToken)
	authHeader := req.Header{
		"Client-ID":     clientID,
		"Authorization": AuthString,
	}
	resp, err := req.Get(requestUrl, authHeader)
	if err != nil {
		return
	}
	err = resp.ToJSON(&user)
	if err != nil {
	}
	return user.Data[0].ID, nil
}
