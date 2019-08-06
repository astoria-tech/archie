package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nlopes/slack"

	"github.com/astoria-arc/archie/config"
)

var rtm *slack.RTM

func main() {

	// Get slack token
	slackToken := os.Getenv("ARCHIE_SLACK_TOKEN")
	if slackToken == "" {
		fmt.Printf("Error: Requires ARCHIE_SLACK_TOKEN (should be begin with \"xoxb-\")")
		os.Exit(-1)
	}

	// Get config path
	configPath := os.Getenv("ARCHIE_CONFIG")
	if configPath == "" {
		configPath = "config.yaml"
	}

	// Check if we will be noisy
	_, slackDebug := os.LookupEnv("ARCHIE_DEBUG")

	//Load config
	config, err := config.Load(configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//Should start with xoxb
	rtm = slack.New(
		slackToken,
		slack.OptionDebug(slackDebug),
		slack.OptionLog(log.New(os.Stdout, "SLACK: ", log.Lshortfile|log.LstdFlags)),
	).NewRTM()

	go rtm.ManageConnection()
	slackLoop(config)
}
