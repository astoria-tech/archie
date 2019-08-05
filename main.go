package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nlopes/slack"

	"github.com/belminf/archie/config"
)

var rtm *slack.RTM

func main() {

	//Load config
	config, err := config.Load("config.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//Should start with xoxb
	rtm = slack.New(
		config.SlackConfig.Token,
		slack.OptionDebug(config.SlackConfig.Debug),
		slack.OptionLog(log.New(os.Stdout, "SLACK: ", log.Lshortfile|log.LstdFlags)),
	).NewRTM()

	go rtm.ManageConnection()
	slackLoop(config)
}
