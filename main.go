package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"text/template"
	"time"

	"github.com/nlopes/slack"
)

const tpl = `Hey, everybody. Have you seen the #{{ .Name }} channel recently?

{{ if .Purpose.Value -}}
Purpose: {{ .Purpose.Value }}
{{ else }}
Purpose: _(not set)_
{{ end }}`

func main() {
	t := template.Must(template.New("msg").Parse(tpl))

	tkn := os.Getenv("SLACK_TOKEN")
	if tkn == "" {
		log.Fatal("Env var SLACK_TOKEN was not set")
	}

	api := slack.New(tkn)

	c, err := getRandomChannel(api)
	if err != nil {
		log.Fatal(err)
	}

	var b bytes.Buffer
	t.Execute(&b, c)
	fmt.Println("--------")
	fmt.Print(b.String())
	fmt.Println("--------")
}

func getRandomChannel(api *slack.Client) (slack.Channel, error) {
	ch, err := api.GetChannels(true)
	if err != nil {
		return slack.Channel{}, err
	}

	rand.Seed(time.Now().UnixNano())

	return ch[rand.Intn(len(ch))], nil
}
