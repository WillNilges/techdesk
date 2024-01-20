package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type Config struct {
	OrgName    string
	LogoURL    string
	FaviconURL string

	TechTeamID      string
	TechAccessToken string
	TechAppToken    string
	TechChannelID   string

	SupeTeamID      string
	SupeAccessToken string
	SupeAppToken    string
	SupeChannelID   string

	SlackBotID string
}

// Useful global variables
var config Config

// Slack Shit

type SlackBot struct {
	API    *slack.Client
	Socket *socketmode.Client
	BotID  string
	ChannelID string
	TeamID string
}

func NewSlackBot(accessToken string, appToken string) (bot SlackBot, err error) {
	// Login to Tech Slack 
	bot.API = slack.New(accessToken, slack.OptionAppLevelToken(appToken))
	bot.Socket = socketmode.New(bot.API,
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)
	techAuthTestResponse, err := bot.API.AuthTest()
	bot.BotID = techAuthTestResponse.UserID

	return bot, nil
}

// TechDesk object is the puppetmaster >:)
type TechDesk struct {
	Tech SlackBot
	TechHandler *socketmode.SocketmodeHandler

	Supervisor SlackBot
	SupervisorHandler *socketmode.SocketmodeHandler
}

func NewTechDesk() (app TechDesk, err error) {
	app.Tech, err = NewSlackBot(config.TechAccessToken, config.TechAppToken)
	if err != nil {
		return app, err
	}
	app.TechHandler = socketmode.NewSocketmodeHandler(app.Tech.Socket)
	app.TechHandler.Handle(socketmode.EventTypeConnecting, middlewareConnecting)
	app.TechHandler.Handle(socketmode.EventTypeConnectionError, middlewareConnectionError)
	app.TechHandler.Handle(socketmode.EventTypeConnected, middlewareConnected)

	app.TechHandler.HandleEvents(slackevents.AppMention, app.techMention)

	app.Supervisor, err = NewSlackBot(config.TechAccessToken, config.TechAppToken)
	if err != nil {
		return app, err
	}
	app.SupervisorHandler = socketmode.NewSocketmodeHandler(app.Supervisor.Socket)
	app.SupervisorHandler.Handle(socketmode.EventTypeConnecting, middlewareConnecting)
	app.SupervisorHandler.Handle(socketmode.EventTypeConnectionError, middlewareConnectionError)
	app.SupervisorHandler.Handle(socketmode.EventTypeConnected, middlewareConnected)

	return app, nil
}

func middlewareConnecting(evt *socketmode.Event, client *socketmode.Client) {
	fmt.Println("Connecting to Slack with Socket Mode...")
}

func middlewareConnectionError(evt *socketmode.Event, client *socketmode.Client) {
	fmt.Println("Connection failed. Retrying later...")
}

func middlewareConnected(evt *socketmode.Event, client *socketmode.Client) {
	fmt.Println("Connected to Slack with Socket Mode.")
}

func (app *TechDesk) techMention(evt *socketmode.Event, client *socketmode.Client) {
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		fmt.Printf("Ignored %+v\n", evt)
		return
	}

	client.Ack(*evt.Request)

	ev, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.AppMentionEvent)
	if !ok {
		fmt.Printf("Ignored %+v\n", ev)
		return
	}

	fmt.Printf("We have been mentionned in %v\n", ev.Channel)
	_, _, err := client.Client.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
	if err != nil {
		fmt.Printf("failed posting message: %v", err)
	}
}


// util functions

func (bot *SlackBot) getConversationHistory(oldTS string, newTS string, limit int) (c []slack.Message, err error) {
	log.Println("Fetching channel history...")
	params := slack.GetConversationHistoryParameters{
		ChannelID: bot.ChannelID,
		Oldest:    oldTS,
		Latest:    newTS,
		Inclusive: true,
		Limit:     limit,
	}

	var history *slack.GetConversationHistoryResponse
	history, err = bot.Socket.GetConversationHistory(&params)
	return history.Messages, nil
}

// Main Shit

func init() {
	// Load environment variables one way or another
	err := godotenv.Load()
	if err != nil {
		log.Println("Couldn't load .env file")
	}

	config.SlackRecvTeamID = os.Getenv("CSP_SLACK_TEAMID")
	config.SlackRecvAccessToken = os.Getenv("CSP_SLACK_ACCESS_TOKEN")
	config.SlackRecvAppToken = os.Getenv("CSP_SLACK_APP_TOKEN")
	config.SlackRecvStatusChannelID = os.Getenv("CSP_SLACK_STATUS_CHANNEL")
}

func main() {

}
