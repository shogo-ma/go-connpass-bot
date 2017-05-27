package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
	"github.com/shogo-ma/go-connpas-bot/models"
)

const eventNum = 10

type SlackBot struct {
	api  *slack.Client
	rtm  *slack.RTM
	ID   string
	Name string
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

func (bot *SlackBot) events(channel string) error {
	cps, err := models.Request(
		&models.Params{
			Count: eventNum,
			Order: 3, // latest
		})

	if err != nil {
		return err
	}

	contents := "```"
	for _, event := range cps.Events {
		contents += event.Title + "\n"
	}
	contents += "```"

	bot.rtm.SendMessage(
		bot.rtm.NewOutgoingMessage(
			contents,
			channel,
		),
	)

	return nil
}

func (bot *SlackBot) handleResponse(text, channel string) {
	if strings.Contains(text, "help") {
		bot.help(channel)
	} else if strings.Contains(text, "events") {
		err := bot.events(channel)
		if err != nil {
			panic(err)
		}
	}
}

func (bot *SlackBot) SetBotProfile(id, name string) {
	bot.ID = id
	bot.Name = name
}

func (bot *SlackBot) IsMention(text string) bool {
	bot_name := fmt.Sprintf("<@%s>", bot.ID)
	if strings.Contains(text, bot_name) {
		return true
	}

	return false
}

func run(api *slack.Client) int {
	api.SetDebug(true)

	slack_bot := NewSlackBot(api)
	go slack_bot.rtm.ManageConnection()

	for msg := range slack_bot.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			slack_bot.SetBotProfile(
				ev.Info.User.ID,
				ev.Info.User.Name,
			)

		case *slack.MessageEvent:
			if ev.Type == "message" && slack_bot.IsMention(ev.Text) {
				slack_bot.handleResponse(ev.Text, ev.Channel)
			}

		case *slack.InvalidAuthEvent:
			log.Println("Invalid Credentials")
			return 1
		}
	}

	return 1
}

func main() {
	api := slack.New(os.Getenv("SLACK_TOKEN"))
	os.Exit(run(api))
}
