package main

import (
	"bytes"
	"log"
	"text/template"

	"github.com/nlopes/slack"
	"gopkg.in/mgo.v2"
)

const tpl = `:trumpet::trumpet:
Hear ye! Hear ye!

Hast thou visited the <#{{ .ID }}|{{ .Name }}> channel recently?

Purpose: {{ if .Purpose.Value }}{{ .Purpose.Value }}{{ else }}_(not set)_{{ end }}
:trumpet::trumpet:`

var t *template.Template

func init() {
	t = template.Must(template.New("msg").Parse(tpl))
}

func announceChannel(coll *mgo.Collection, api *slack.Client, dest string, c slack.Channel) error {
	var b bytes.Buffer
	if err := t.Execute(&b, c); err != nil {
		return err
	}

	log.Printf("Sending this message to %s\n%s", dest, b.String())

	params := slack.NewPostMessageParameters()
	params.AsUser = true
	params.EscapeText = false
	if _, _, err := api.PostMessage(dest, b.String(), params); err != nil {
		return err
	}

	if err := coll.Insert(Channel{Name: c.Name}); err != nil {
		return err
	}

	return nil
}
