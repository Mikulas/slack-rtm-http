package main

import (
	"github.com/nlopes/slack"
	"log"
	"os"
)

const (
	SlackToken = "xoxb-294382616742-StFn3kQijJMrk5Al6y5fPY5Q"
)

func MakeSlackRtm() *slack.RTM {

	api := slack.New(SlackToken)
	logger := log.New(os.Stdout, "slack: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(false)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	return rtm
}
