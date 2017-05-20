package main

import (
	"log"
	"os"

	"github.com/nlopes/slack"
)

type SlackBot struct {
	api *slack.Client
	rtm *slack.RTM
}

func NewSlackBot(api *slack.Client) *SlackBot {
	return &SlackBot{
		api: api,
		rtm: api.NewRTM(),
	}
}

func (bot *SlackBot) help(channel string) error {
	contents := "```" + `
help:
	help
show:
	return resent events from connpass` + "```"

	bot.rtm.SendMessage(
		bot.rtm.NewOutgoingMessage(
			contents,
			channel,
		),
	)

	return nil
}

func (bot *SlackBot) show() error {
	return nil
}

func (bot *SlackBot) handleResponse(text, channel string) {
	if text == "help" {
		bot.help(channel)
	} else if text == "show" {
		bot.rtm.SendMessage(bot.rtm.NewOutgoingMessage(text, channel))
	}
}

func run(api *slack.Client) int {
	api.SetDebug(true)

	slack_bot := NewSlackBot(api)
	go slack_bot.rtm.ManageConnection()

	for msg := range slack_bot.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if ev.Type == "message" {
				slack_bot.handleResponse(ev.Text, ev.Channel)
			}

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
