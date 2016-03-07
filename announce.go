package main

import (
	"bytes"
	"html/template"
	"log"

	"github.com/nlopes/slack"
)

const tpl = `Hey, everybody. Have you seen the <#{{ .ID }}|{{ .Name }}> channel recently?

{{ if .Purpose.Value -}}
Purpose: {{ .Purpose.Value }}
{{ else }}
Purpose: _(not set)_
{{ end }}`

var t *template.Template

func init() {
	t = template.Must(template.New("msg").Parse(tpl))
}

func announceChannel(api *slack.Client, c slack.Channel) error {
	var b bytes.Buffer
	if err := t.Execute(&b, c); err != nil {
		return err
	}

	log.Printf("Sending message\n%s", b.String())

	params := slack.NewPostMessageParameters()
	params.AsUser = true
	params.EscapeText = false
	_, _, err := api.PostMessage("herald-test", b.String(), params)

	return err
}