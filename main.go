package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nlopes/slack"
)

func run(api *slack.Client) int {
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Println("Event Received")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			log.Println("hello event")

		case *slack.MessageEvent:
			log.Println("Message: %v\n", ev)
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello World", ev.Channel))

		case *slack.InvalidAuthEvent:
			log.Println("Invalid credentials")
			return 1
		}
	}

	return 1
}

func main() {
	api := slack.New(os.Getenv("SLACK_TOKEN"))
	os.Exit(run(api))
}
