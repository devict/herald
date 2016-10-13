package main

import (
	"errors"
	"log"
	"math/rand"

	"github.com/nlopes/slack"
	"gopkg.in/mgo.v2"
)

func selectChannel(coll *mgo.Collection, api *slack.Client, ignore Channels) (slack.Channel, error) {
	ch, err := api.GetChannels(true)
	if err != nil {
		return slack.Channel{}, err
	}

	var old Channels
	if err := coll.Find(nil).All(&old); err != nil {
		return slack.Channel{}, err
	}

	all := ch
	for len(ch) > 0 {
		n := rand.Intn(len(ch))

		log.Print("Considering ", ch[n].Name)
		if old.Contains(ch[n].Name) || ignore.Contains(ch[n].Name) {
			log.Print("Skip it.")
			ch = append(ch[:n], ch[n+1:]...)
			continue
                } else if ch[n].Name = cnf.DestChan {
			log.Print("Skip it. Don't announce destination channel")
			ch = append(ch[:n], ch[n+1:]...)
			continue
                }

		return ch[n], nil
	}

	// Every channel has already been announced. Truncate the collection then
	// send a random one that isn't in the 'ignored' collection.
	log.Print("Couldn't find a new one")
	if err := coll.DropCollection(); err != nil {
		return slack.Channel{}, err
	}

	for len(all) > 0 {
		n := rand.Intn(len(all))

		log.Print("Considering ", all[n].Name)
		if ignore.Contains(all[n].Name) {
			log.Print("Skip it.")
			all = append(all[:n], all[n+1:]...)
			continue
		}
		return all[n], nil
	}

	return slack.Channel{}, errors.New("Unable to find a channel")
}
