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

	slack := MakeSlackRtm()

	server := MakeServer(8080, channel, slack)
	server.Start()
}
