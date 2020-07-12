package main

import (
	"math/rand"
	"time"
)

func GuessingGame(message WSMessage, directionGuess string) (err error) {
	directionResult, _ := GetLeftOrRight()
	SendGuess(directionResult, directionGuess)
	// not sending a message to avoid weird state
	if directionGuess == directionResult {
		thisGotchi.Food += 1
		thisGotchi.Love += 1
	}
	return
}

func SendGuess(directionResult string, guess string) (err error) {
	message := GotchiMessage{
		Type: "guess",
		Guess: GuessStruct{
			Result: directionResult,
			Guess:  guess,
		},
	}
	SendMessage(message)
	return
}

func GuessLeft(message WSMessage) (err error) {
	err = GuessingGame(message, "left")
	return
}

func GuessRight(message WSMessage) (err error) {
	err = GuessingGame(message, "right")
	return
}

func GetLeftOrRight() (direction string, err error) {
	rand.Seed(time.Now().UnixNano())
	result := rand.Int()
	if result % 2 == 0 {
		direction = "left"
	} else {
		direction = "right"
	}
	return
}

type GuessStruct struct {
	Guess string `json:"guess"`
	Result string `json:"result"`
}