package main

import (
	"log"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/google/go-cmp/cmp"
	"gopkg.in/mgo.v2"
)

type ChannelAndID struct {
	Channel Channel
	ID string
}

func announceChannelDifferences(coll *mgo.Collection, api *slack.Client, dest string) error {
	log.Print("Announce channel differences since last run:")


	// Get current channels (and their IDs)
	rawCurrent, err := api.GetChannels(true) // true ignores archived channels
	if err != nil {
		return err
	}
	current := make([]ChannelAndID, 0, 50)
	for _, v := range rawCurrent {
		current = append(current, ChannelAndID{Channel: Channel{v.Name}, ID: v.ID})
	}


	// Get the channel list from last run
	var old []ChannelAndID
	if err := coll.Find(nil).All(&old); err != nil {
		return err
	}
	if old == nil {
		// first run!, write current channels to collection and exit
		if err := writeChannelsToColl(coll, current); err != nil {
			return err
		}
		log.Print("No announcement: nothing to compare against. comparison will happen at next run.")
		return nil
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
	msg := fmt.Sprintf("Channel changes:\n%s", diff)
	if _, _, err := api.PostMessage(dest, msg, params); err != nil {
		return err
	}
	log.Print("Sent")


	// Saving current channel list
	if err := writeChannelsToColl(coll, current); err != nil {
		return err
	}


	return nil
}



func writeChannelsToColl(coll *mgo.Collection, chans []ChannelAndID) error {
	coll.DropCollection()
	for _, v := range chans {
		if err := coll.Insert(v); err != nil {
			return err
		}
	}
	return nil
}

