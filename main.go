package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/nlopes/slack"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	tkn := os.Getenv("SLACK_TOKEN")
	if tkn == "" {
		log.Fatal("Env var SLACK_TOKEN was not set")
	}

	api := slack.New(tkn)

	c, err := getRandomChannel(api)
	if err != nil {
		log.Print(err)
	}

	if err := announceChannel(api, c); err != nil {
		log.Print(err)
	}
}

func getRandomChannel(api *slack.Client) (slack.Channel, error) {
	ch, err := api.GetChannels(true)
	if err != nil {
		return slack.Channel{}, err
	}

	return ch[rand.Intn(len(ch))], nil
}
