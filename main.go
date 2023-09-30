package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
	}
}

type OauthTokens struct {
	SlackBotToken string
	SlackAppToken string
}

func loadEnv() {
	err := godotenv.Load("app.env")
	if err != nil {
		return
	}
}

func getSlackTokens() OauthTokens {

	var slackOauthTokens OauthTokens
	slackOauthTokens.SlackBotToken = os.Getenv("SLACK_BOT_TOKEN")
	slackOauthTokens.SlackAppToken = os.Getenv("SLACK_APP_TOKEN")
	return slackOauthTokens
}

func main() {
	loadEnv()

	errBot := os.Setenv("SLACK_BOT_TOKEN", getSlackTokens().SlackBotToken)
	if errBot != nil {
		return
	}
	errApp := os.Setenv("SLACK_APP_TOKEN", getSlackTokens().SlackAppToken)
	if errApp != nil {
		return
	}

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	go printCommandEvents(bot.CommandEvents())

	definition := &slacker.CommandDefinition{
		Description: "Yob Calculator",
		Handler: func(context slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			year := request.Param("year")
			yob, err := strconv.Atoi(year)
			if err != nil {
				println("error converting year to integer")
			}
			age := 2023 - yob
			r := fmt.Sprintf("age is %d", age)
			response.Reply(r)
		},
	}

	bot.Command("my yob is <year>", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
