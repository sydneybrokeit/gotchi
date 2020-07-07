package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

const DefaultFoodMax = 4
const DefaultLoveMax = 4
const DefaultDeplete = time.Duration(15 * time.Minute)

type Gotchi struct {
	// species
	Type            string         `json:"type"`
	// is configured
	ReadyToHatch	bool			`json:"ready_to_hatch"`
	// has hatched
	Hatched         bool           `json:"hatched"`
	// current food
	Food            int            `json:"food"`
	// maximum food
	FoodMax         int            `json:"food_max"`
	// timer for food depletion
	FoodTicker      *time.Ticker   `json:"-"`
	// current love
	Love            int            `json:"love"`
	// maximum love
	LoveMax         int            `json:"love_max"`
	// timer for love depletion
	LoveTicker      *time.Ticker   `json:"-"`
	// level (1,2,3)
	Level           int            `json:"level"`
	DepleteInt      int            `json:"deplete_int"`
	// basic depletion
	DepleteDuration time.Duration  `json:"deplete_duration"`
	// a map of user ids and their count of snacks
	Inventory       map[string]int `json:"inventory"`
}

// from data passed in, provide a ready-to-hatch Gotchi
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

// hatch the Gotchi
func (this *Gotchi) Hatch() {
	log.Warning("hatching!")
	this.Hatched = true
	message := GotchiMessage{
		Type: "HATCH",
		State: this,
	}
	SendMessage(message)
	go FoodTimers()
}


// handle the timers for food and love depletion
func FoodTimers() {
	for {
		select {
		case <- thisGotchi.LoveTicker.C:
			ChangeLove(-1)
		case <- thisGotchi.FoodTicker.C:
			ChangeFood(-1)
		}
	}
}

// add the passed in amount to food count and update
func ChangeLove(amount int) {
	log.Debug("decrementing love")
	thisGotchi.LoveTicker = time.NewTicker(thisGotchi.DepleteDuration)
	thisGotchi.Love += amount
	// bounds checking
	if thisGotchi.Love < 0 {
		thisGotchi.Love = 0
	} else if thisGotchi.Love > thisGotchi.LoveMax {
		thisGotchi.Love = thisGotchi.LoveMax
	} else {
		message := GotchiMessage{
			Type: "DELTA",
			Delta: Delta{
				Type: "love",
				Amount: amount,
			},
		}
		SendMessage(message)
	}
}

// add the passed in amount to love count and update
func ChangeFood(amount int) {
	log.Debug("decrementing food")
	thisGotchi.FoodTicker = time.NewTicker(thisGotchi.DepleteDuration)
	thisGotchi.Food += amount
	// bounds checking
	if thisGotchi.Food < 0 {
		thisGotchi.Food = 0
	} else if thisGotchi.Food > thisGotchi.FoodMax {
		thisGotchi.Food = thisGotchi.FoodMax
	} else {
		message := GotchiMessage{
			Type: "DELTA",
			Delta: Delta{
				Type: "food",
				Amount: amount,
			},
		}
		SendMessage(message)
	}
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

type Delta struct {
	Type string `json:"type"`
	Amount int `json:"amount"`
}

type GotchiMessage struct {
	Type string `json:"type"`
	State *Gotchi `json:"state,omitempty"`
	Delta Delta `json:"delta,omitempty"`
}

func (this *Gotchi) UpdateAll() {
	update := GotchiMessage{State: this, Type: "REFRESH"}
	jsonRepr, err := json.Marshal(update)
	if err != nil {
		log.Errorf("couldn't unmarshal: %v", err)
	}
	log.Print(jsonRepr)
	// hub.broadcast <- jsonRepr
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



// /hatch endpoint
func HandleHatch(w http.ResponseWriter, r *http.Request) (err error){
	if thisGotchi.ReadyToHatch != true || thisGotchi.Hatched == true {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		thisGotchi.Hatch()
		HatchChan <- true
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	return
}

// debugging endpoint for increasing "food"
func HandleFeed(w http.ResponseWriter, r *http.Request) (err error) {
	ChangeFood(1)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return
}

// debugging endpoint for increasing "love"
func HandleLove(w http.ResponseWriter, r *http.Request) (err error) {
	ChangeLove(1)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect )
	return
}
