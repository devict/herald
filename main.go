package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/nlopes/slack"
	"gopkg.in/mgo.v2"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	rand.Seed(time.Now().UnixNano())

	cnf, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := mgo.Dial(cnf.MongoURI)
	if err != nil {
		log.Fatal(err)
	}
	coll := db.DB(cnf.MongoDB).C("channels")

	api := slack.New(cnf.SlackToken)

	c, err := selectChannel(coll, api)
	if err != nil {
		log.Fatal(err)
	}

	if err := announceChannel(coll, api, cnf.DestChan, c); err != nil {
		log.Fatal(err)
	}
}
