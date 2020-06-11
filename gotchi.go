package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

const DefaultFoodMax = 4
const DefaultLoveMax = 4
const DefaultDeplete = time.Duration(15 * time.Minute)

type Gotchi struct {
	Type            string         `json:"type"`
	Hatched         bool           `json:"hatched"`
	Food            int            `json:"food"`
	FoodMax         int            `json:"food_max"`
	FoodTicker      *time.Ticker   `json:"-"`
	Love            int            `json:"love"`
	LoveMax         int            `json:"love_max"`
	LoveTicker      *time.Ticker   `json:"-"`
	Level           int            `json:"level"`
	DepleteInt      int            `json:"deplete_int"`
	DepleteDuration time.Duration  `json:"deplete_duration"`
	Inventory       map[string]int `json:"inventory"`
}

func StartGotchi(species string, foodmax string, lovemax string, deplete string,
	food string, love string) (output *Gotchi, err error) {
	output = new(Gotchi)
	output.Type = species
	output.FoodMax, err = strconv.Atoi(foodmax)
	if err != nil {
		log.Warningf("invalid food max, setting default %d : %v", DefaultFoodMax, err)
		output.FoodMax = DefaultFoodMax
	}
	output.LoveMax, err = strconv.Atoi(lovemax)
	if err != nil {
		log.Warningf("invalid love max, setting default %d: %v", DefaultLoveMax, err)
		output.LoveMax = DefaultLoveMax
	}
	duration, err := strconv.Atoi(deplete)
	if err != nil {
		log.Warningf("invalid deplete duration, setting default %v : %v", DefaultDeplete, err)
		output.DepleteDuration = DefaultDeplete
	} else {
		output.DepleteDuration = time.Duration(duration) * time.Minute
	}
	output.Food, err = strconv.Atoi(food)
	if err != nil {
		output.Food = output.FoodMax / 2
	}
	output.Love, err = strconv.Atoi(love)
	if err != nil {
		output.Love = output.LoveMax / 2
	}
	output.FoodTicker = time.NewTicker(output.DepleteDuration)
	output.LoveTicker = time.NewTicker(output.DepleteDuration)
	output.Inventory = make(map[string]int)
	return
}

func (this *Gotchi) Hatch() {
	log.Warning("hatching!")
	this.Hatched = true
}

func (this *Gotchi) Do() {
	go this.PrintOutLoop()
}

func (this *Gotchi) Print() {
	output, err := json.MarshalIndent(this, "", "\t")
	if err != nil {
		log.Errorf("couldn't print - %v", err)
	}
	log.Debugf("%s", string(output))
}

type GotchiMessage struct {
	Type string `json:"type"`
	State Gotchi `json:"state,omitempty"`
}

func (this *Gotchi) UpdateAll() {
	update := GotchiMessage{State: *this, Type: "REFRESH"}
	jsonRepr, err := json.Marshal(update)
	if err != nil {
		log.Errorf("couldn't unmarshal: %v", err)
	}
	log.Print(jsonRepr)
	internalService.send(userId, string(jsonRepr))
}

func (this *Gotchi) PrintOutLoop() {
	printoutTicker := time.NewTicker(15 * time.Second)
	for {
		select {
		case <-printoutTicker.C:
			this.Print()
			this.UpdateAll()
		}
	}
}
