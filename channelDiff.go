package main

import (
	"fmt"
	"log"

	"github.com/google/go-cmp/cmp"
	"github.com/nlopes/slack"
	"gopkg.in/mgo.v2"
)

func announceChannelDifferences(coll *mgo.Collection, api *slack.Client, dest string) error {
	log.Print("Announce channel differences since last run:")

	// Get current channels (and their IDs)
	rawCurrent, err := api.GetChannels(true) // true ignores archived channels
	if err != nil {
		return err
	}
	current := make([]Channel, len(rawCurrent))
	for i, v := range rawCurrent {
		current[i] = Channel{Name: v.Name, ID: v.ID}
	}

	// Get the channel list from last run
	var old []Channel
	if err := coll.Find(nil).All(&old); err != nil {
		return err
	}

	if old == nil {
		// first run!, write current channels to collection and exit
		log.Print("No announcement: nothing to compare against. comparison will happen at next run.")
		return writeChannelsToColl(coll, current)
	}

	// Compare
	diff := cmp.Diff(old, current)
	if diff == "" {
		log.Print("No announcement: no differences.")
		return nil
	}

	// Print differences
	log.Print("Differences:")
	log.Print(diff)
	log.Print()
	log.Printf("Sending message to %s\n", dest)
	params := slack.NewPostMessageParameters()
	params.AsUser = true
	params.EscapeText = false
	if _, _, err := api.PostMessage(dest, fmt.Sprintf("Channel changes:\n```%s```", diff), params); err != nil {
		return err
	}
	log.Print("Sent")

	// Saving current channel list
	return writeChannelsToColl(coll, current)
}

func writeChannelsToColl(coll *mgo.Collection, chans []Channel) error {
	if err := coll.DropCollection(); err != nil {
		return err
	}
	for _, v := range chans {
		if err := coll.Insert(v); err != nil {
			return err
		}
	}
	return nil
}
