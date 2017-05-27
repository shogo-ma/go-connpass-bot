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

func (bot *SlackBot) help(user, channel string) error {
	contents := "```" + `
help:
	help
events [word]:
	return recent events from connpass` + "```"

	bot.sendMessage(fmt.Sprintf("<@%s>", user), channel)
	bot.sendMessage(contents, channel)

	return nil
}

func (bot *SlackBot) events(user, text, channel string) error {
	// make search query
	qs := strings.Split(text, " ")
	fmt.Println(qs)

	var keyword string
	search_word := qs[len(qs)-1]
	if search_word == "events" {
		keyword = ""
	} else {
		keyword = search_word
	}

	cps, err := models.Request(
		&models.Params{
			Keyword: keyword,
			Count:   eventNum,
			Order:   3, // latest
		})

	if err != nil {
		return err
	}

	// mention
	bot.sendMessage(fmt.Sprintf("<@%s>", user), channel)
	for _, event := range cps.Events {
		contents := fmt.Sprintf("%s\n%s\n",
			event.Title,
			event.EventURL,
		)
		bot.sendMessage(contents, channel)
	}

	return nil
}

func (bot *SlackBot) notFound(user, channel string) {
	contents := fmt.Sprintf("<@%s> Command Not Found", user)
	bot.sendMessage(contents, channel)
}

func (bot *SlackBot) sendMessage(contents, channel string) {
	bot.rtm.SendMessage(
		bot.rtm.NewOutgoingMessage(
			contents,
			channel,
		),
	)
}

func (bot *SlackBot) handleResponse(text, channel, user string) {
	if strings.Contains(text, "help") {
		bot.help(user, channel)
	} else if strings.Contains(text, "events") {
		err := bot.events(user, text, channel)
		if err != nil {
			panic(err)
		}
	} else {
		bot.notFound(user, channel)
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
				slack_bot.handleResponse(ev.Text, ev.Channel, ev.User)
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
