package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/astoria-arc/archie/config"
	"github.com/astoria-arc/archie/msgs"

	"github.com/nlopes/slack"
)

var (
	botTagString string
)

//SlackLoop processes events
func slackLoop(config *config.Config) {

	// Load messages
	msgs.LoadMessages(&config.Messages)

	//Ref: https://rsmitty.github.io/Slack-Bot/
	for {
		select {
		case msg := <-rtm.IncomingEvents:

			switch e := msg.Data.(type) {
			case *slack.MessageEvent:
				respond(e)

			case *slack.ConnectedEvent:
				fmt.Printf(
					"Connected, counter: %d\n",
					e.ConnectionCount,
				)
				botTagString = fmt.Sprintf("<@%s>", rtm.GetInfo().User.ID)

			case *slack.LatencyReport:
				fmt.Printf("Current latency: %v\n", e.Value)

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", e.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				os.Exit(1)

			// Ignored event types
			case
				*slack.PresenceChangeEvent,
				*slack.HelloEvent,
				*slack.AccountsChangedEvent,
				*slack.BotAddedEvent,
				*slack.BotChangedEvent,
				*slack.ChannelArchiveEvent,
				*slack.ChannelCreatedEvent,
				*slack.ChannelHistoryChangedEvent,
				*slack.ChannelJoinedEvent,
				*slack.ChannelLeftEvent,
				*slack.ChannelMarkedEvent,
				*slack.ChannelRenameEvent,
				*slack.ChannelUnarchiveEvent,
				*slack.CommandsChangedEvent,
				*slack.DNDUpdatedEvent,
				*slack.EmailDomainChangedEvent,
				*slack.EmojiChangedEvent,
				*slack.FileChangeEvent,
				*slack.FileCommentAddedEvent,
				*slack.FileCommentEditedEvent,
				*slack.FileCommentDeletedEvent,
				*slack.FileCreatedEvent,
				*slack.FileSharedEvent,
				*slack.FileUnsharedEvent,
				*slack.GroupCloseEvent,
				*slack.GroupHistoryChangedEvent,
				*slack.GroupJoinedEvent,
				*slack.GroupLeftEvent,
				*slack.GroupMarkedEvent,
				*slack.GroupOpenEvent,
				*slack.GroupRenameEvent,
				*slack.GroupUnarchiveEvent,
				*slack.IMCloseEvent,
				*slack.IMCreatedEvent,
				*slack.IMHistoryChangedEvent,
				*slack.IMMarkedEvent,
				*slack.ManualPresenceChangeEvent,
				*slack.MemberJoinedChannelEvent,
				*slack.MemberLeftChannelEvent,
				*slack.PinAddedEvent,
				*slack.PinRemovedEvent,
				*slack.PrefChangeEvent,
				*slack.ReactionAddedEvent,
				*slack.ReactionRemovedEvent,
				*slack.ReconnectUrlEvent,
				*slack.StarAddedEvent,
				*slack.StarRemovedEvent,
				*slack.SubteamCreatedEvent,
				*slack.SubteamMembersChangedEvent,
				*slack.SubteamSelfAddedEvent,
				*slack.SubteamSelfRemovedEvent,
				*slack.SubteamUpdatedEvent,
				*slack.TeamDomainChangeEvent,
				*slack.TeamJoinEvent,
				*slack.TeamMigrationStartedEvent,
				*slack.TeamPrefChangeEvent,
				*slack.UserChangeEvent,
				*slack.UserTypingEvent,
				*slack.ConnectingEvent,
				*slack.AckMessage:

			default:
				fmt.Printf("Unexpected [%[1]T]: %[1]v\n", msg.Data)
			}
		}
	}
}

func respond(e *slack.MessageEvent) {

	// Check if I was tagged
	botTagged := strings.Contains(e.Msg.Text, botTagString)

	// Get channel name
	c, err := rtm.GetConversationInfo(e.Channel, true)
	if err != nil {
		fmt.Printf("Msg get channel info error: %s\n", err)
		fmt.Println(e)
		return
	}

	// Only respond to channel messages if I'm tagged
	if c.Name != "" && !botTagged {
		return
	}

	// Do not reply to my own messages
	if e.User == rtm.GetInfo().User.ID {
		return
	}

	// Do not reply to Slackbot
	if e.User == "USLACKBOT" {
		return
	}

	// Remove tag from message
	msgText := strings.Replace(e.Text, botTagString, "", -1)

	// Get user info
	u, err := rtm.GetUserInfo(e.User)
	if err != nil {
		fmt.Printf("Msg get user info error: %s\n", err)
		return
	}

	// Get response
	response, err := msgs.Response(msgText)
	if err != nil {
		fmt.Printf("Error: Responding to %s/%s: %s", u.Name, c.Name, err)
		return
	}

	if response == "" {
		fmt.Printf("Error: Responding to %s/%s: Got empty message back", u.Name, c.Name)
	}

	// ASSERT: We have a message to send back

	// Massage message

	if strings.Contains(response, "$YOU") {

		// Replace $YOU
		response = strings.Replace(response, "$YOU", fmt.Sprintf("<@%s>", u.Name), -1)

	} else if c.Name != "" {

		// Tag user since we're in a room and haven't tagged him yet
		response = fmt.Sprintf("<@%s> %s", u.Name, response)
	}

	rtm.SendMessage(rtm.NewOutgoingMessage(response, e.Channel))
}
