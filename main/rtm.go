package main

import (
	"github.com/nlopes/slack"
	"log"
	"os"
)

func MakeSlackRtm(authToken string) *slack.RTM {

	api := slack.New(authToken)
	logger := log.New(os.Stdout, "slack: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(false)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	return rtm
}
