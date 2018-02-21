package main

import (
	"log"
	"os"
)

func main() {
	channel := os.Getenv("CHANNEL")
	if channel == "" {
		log.Panic("unset env CHANNEL")
	}

	slackToken := os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		log.Panic("unset env SLACK_TOKEN")
	}

	slack := MakeSlackRtm(slackToken)

	server := MakeServer(8080, channel, slack)
	server.Start()
}
