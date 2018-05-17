package main

import (
	"log"
	"math/rand"
	"strings"
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

	ignore := Channels{
		Channel{
			Name: strings.TrimLeft(cnf.DestChan, "#"),
		},
	}

	c, err := selectChannel(coll, api, ignore)
	if err != nil {
		log.Fatal(err)
	}

	if err := announceChannel(coll, api, cnf.DestChan, c); err != nil {
		log.Fatal(err)
	}

	coll = db.DB(cnf.MongoDB).C("lastChannelList")
	if err := announceChannelDifferences(coll, api, cnf.DiffChan); err != nil {
		log.Fatal(err)
	}
}
