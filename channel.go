package main

// Channel provides a light internal representation of a Slack Channel
// appropriate for storing in mongo.
type Channel struct {
	Name string `bson:"name"`
}

// Channels is a collection of Channels.
type Channels []Channel

// Contains returns true if the collection contains a channel with the same
// name.
func (c Channels) Contains(name string) bool {
	for _, x := range c {
		if name == x.Name {
			return true
		}
	}
	return false
}
