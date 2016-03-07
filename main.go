package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/nlopes/slack"
	"gopkg.in/mgo.v2"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	tkn := os.Getenv("SLACK_TOKEN")
	if tkn == "" {
		log.Fatal("Env var SLACK_TOKEN was not set")
	}
	api := slack.New(tkn)

	url := os.Getenv("MONGOLAB_URI")
	if url == "" {
		log.Fatal("Env var MONGOLAB_URI was not set")
	}

	log.Print("Connecting to ", url)
	db, err := mgo.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	coll := db.DB("herald").C("channels")

	c, err := selectChannel(coll, api)
	if err != nil {
		log.Fatal(err)
	}

	if err := announceChannel(coll, api, c); err != nil {
		log.Fatal(err)
	}
}
