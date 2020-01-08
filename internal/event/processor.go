package event

import (
	"fmt"
	"log"
	"time"

	"github.com/ONSdigital/github-auditor/pkg/github"
	"github.com/ONSdigital/github-auditor/pkg/slack"
)

const slackRateLimitPause = 5 * time.Second

// Process processes the passed slice of GitHub audit events, creating Slack alerts in the passed Slack channel for events of interest.
func Process(events []github.Node, slackAlertsChannel string, slackWebHookURL string) {
	for _, e := range events {
		switch e.Action {
		case "oauth_application.create":
			text := fmt.Sprintf(github.MessageForEvent(e.Action), e.OauthApplicationName, e.OrganizationName, e.ActorLogin)
			postSlackMessage(e.CreatedAt, text, slackAlertsChannel, slackWebHookURL)
		case "repo.add_member":
			text := fmt.Sprintf(github.MessageForEvent(e.Action), e.ActorLogin, e.RepositoryName)
			postSlackMessage(e.CreatedAt, text, slackAlertsChannel, slackWebHookURL)
		default:
			log.Printf("Unknown GitHub event: %s", e.Action)
		}
	}
}

func postSlackMessage(timestamp, text, slackAlertsChannel, slackWebHookURL string) {
	payload := slack.Payload{
		Text:      fmt.Sprintf("_%s_\n%s\n\n", formatTime(timestamp), text),
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

func formatTime(s string) string {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		log.Fatalf("Unable to parse time '%s': %v", t, err)
	}

	return t.Format("Monday 02 Jan 2006 15:04:05 MST")
}
