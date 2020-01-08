package event

import (
	"fmt"
	"log"
	"time"

	"github.com/ONSdigital/github-auditor/pkg/github"
	"github.com/ONSdigital/github-auditor/pkg/slack"
)

const slackRateLimitPause = 3 * time.Second

// Process does something
func Process(events []github.Node, slackAlertsChannel string, slackWebHookURL string) {
	for _, e := range events {
		switch e.Action {
		case "oauth_application.create":
			text := fmt.Sprintf(github.MessageForEvent(e.Action), e.OauthApplicationName, e.OrganizationName, e.ActorLogin)
			postSlackMessage(text, slackAlertsChannel, slackWebHookURL)
		case "org.add_member":
			text := fmt.Sprintf(github.MessageForEvent(e.Action), e.OauthApplicationName, e.OrganizationName, e.ActorLogin)
			postSlackMessage(text, slackAlertsChannel, slackWebHookURL)
		default:
			log.Printf("Unknown GitHub event: %s", e.Action)
		}
	}
}

func postSlackMessage(text string, slackAlertsChannel string, slackWebHookURL string) {
	payload := slack.Payload{
		Text:      text,
		Username:  "GitHub Auditor Bot",
		Channel:   slackAlertsChannel,
		IconEmoji: ":github:",
	}

	time.Sleep(slackRateLimitPause)
	err := slack.Send(slackWebHookURL, payload)
	if err != nil {
		log.Fatalf("Failed to send Slack message: %v", err)
	}
}
