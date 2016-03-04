package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/nlopes/slack"
)

const tpl = `Hey, everybody. Have you seen the #{{ .Name }} channel?

{{ if .Purpose.Value -}}
Purpose: {{ .Purpose.Value }}
{{- else -}}
Purpose: **(not set)**
{{- end }}
{{ if .Topic.Value -}}
Current Topic: {{ .Topic.Value }}
{{- else -}}
Current Topic: **(not set)**
{{ end }}`

func main() {
	t := template.Must(template.New("msg").Parse(tpl))

	tkn := os.Getenv("SLACK_TOKEN")
	if tkn == "" {
		log.Fatal("Env var SLACK_TOKEN was not set")
	}

	api := slack.New(tkn)
	ch, err := api.GetChannels(true)
	if err != nil {
		log.Fatal(err)
		return
	}
	var b bytes.Buffer
	for _, c := range ch {
		b.Reset()
		t.Execute(&b, c)
		fmt.Println("--------")
		fmt.Print(b.String())
		fmt.Println("--------")
	}
}
