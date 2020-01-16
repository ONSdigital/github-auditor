package event

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ONSdigital/github-auditor/pkg/github"
	firestore "github.com/ONSdigital/github-auditor/pkg/googlecloud"
	"github.com/ONSdigital/github-auditor/pkg/slack"
)

const slackRateLimitPause = 5 * time.Second

// Process processes the passed slice of GitHub audit events, creating Slack alerts in the passed Slack channel for events of interest.
func Process(events []github.Node, firestoreCredentials, firestoreProject, slackAlertsChannel, slackWebHookURL string) {
	for _, e := range events {
		timestamp := formatTime(e.CreatedAt)
		id := e.ID
		action := e.Action
		text := ""

		switch e.Action {

		// OAuth events.
		case "oauth_application.create":
			text = fmt.Sprintf(github.MessageForEvent(action), e.OauthApplicationName, e.OrganizationName, formatActor(e.Actor, false))

		// Organisation events.
		case "org.add_billing_manager":
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.Actor, true), formatActor(e.User, false), e.OrganizationName)
		case "org.add_member":
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.Actor, true), formatActor(e.User, false), e.RepositoryName)
		case "org.block_user":
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.BlockedUser, true), formatActor(e.Actor, false), e.OrganizationName)
		case "org.create":
			text = fmt.Sprintf(github.MessageForEvent(action), e.OrganizationName, formatActor(e.Actor, false))
		case "org.disable_saml":
			text = fmt.Sprintf(github.MessageForEvent(action), e.OrganizationName, formatActor(e.Actor, false))
		case "org.disable_two_factor_requirement":
			text = fmt.Sprintf(github.MessageForEvent(action), e.OrganizationName, formatActor(e.Actor, false))
		case "org.enable_oauth_app_restrictions":
			text = fmt.Sprintf(github.MessageForEvent(action), e.OrganizationName, formatActor(e.Actor, false))
		case "org.enable_saml":
			text = fmt.Sprintf(github.MessageForEvent(action), e.OrganizationName, formatActor(e.Actor, false))
		case "org.enable_two_factor_requirement":
			text = fmt.Sprintf(github.MessageForEvent(action), e.OrganizationName, formatActor(e.Actor, false))
		case "org.invite_member":
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.Actor, true), formatActor(e.User, false), e.OrganizationName)
		case "org.oauth_app_access_approved":
			text = fmt.Sprintf(github.MessageForEvent(action), e.OauthApplicationName, e.OrganizationName, formatActor(e.Actor, false))
		case "org.oauth_app_access_denied":
			text = fmt.Sprintf(github.MessageForEvent(action), e.OauthApplicationName, e.OrganizationName, formatActor(e.Actor, false))
		case "org.oauth_app_access_requested":
			text = fmt.Sprintf(github.MessageForEvent(action), e.OauthApplicationName, e.OrganizationName, formatActor(e.Actor, false))
		case "org.remove_billing_manager":
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.Actor, true), formatActor(e.User, false), e.OrganizationName)
		case "org.remove_member":
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.Actor, true), formatActor(e.User, false), e.OrganizationName)
		case "org.remove_outside_collaborator":
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.Actor, true), formatActor(e.User, false), e.OrganizationName)
		case "org.update_member":
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.Actor, true), formatActor(e.User, false), strings.ToLower(e.PermissionWas), strings.ToLower(e.Permission), e.OrganizationName)

		// Repo events.
		case "repo.access":
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.Actor, true), e.RepositoryName, strings.ToLower(e.Visibility))
		case "repo.add_member":
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.Actor, true), formatActor(e.User, false), e.RepositoryName)
		case "repo.archived":
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.Actor, true), e.RepositoryName)
		case "repo.change_merge_setting":
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.Actor, true), e.RepositoryName, strings.ToLower(e.MergeType))
		case "repo.create":
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.Actor, true), e.RepositoryName, strings.ToLower(e.Visibility))
		case "repo.destroy":
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.Actor, true), e.RepositoryName)
		case "repo.remove_member":
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.Actor, true), formatActor(e.User, false), e.RepositoryName)

		// Team events.
		case "team.add_member":
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.Actor, true), formatActor(e.User, false), e.TeamName)
		case "team.add_repository":
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.Actor, true), e.TeamName, e.RepositoryName)
		case "team.remove_member":
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.Actor, true), formatActor(e.User, false), e.TeamName)
		case "team.remove_repository":
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.Actor, true), e.TeamName, e.RepositoryName)
		default:
			log.Printf("Unknown GitHub event: %s", action)
		}

		client := firestore.NewClient(firestoreProject, firestoreCredentials)
		if !client.DocExists(id, timestamp, action) && len(text) > 0 {
			postSlackMessage(timestamp, text, slackAlertsChannel, slackWebHookURL)
		}

		err := client.SaveDoc(id, timestamp, action)
		if err != nil {
			log.Fatalf("Failed to save document to Firestore: %v", err)
		}
	}
}

func formatActor(actor github.Actor, capitalise bool) string {
	actorName := ""

	switch actor.Type {
	case "Bot":
		actorName = fmt.Sprintf("bot *%s*", actor.Login)
		if capitalise {
			actorName = fmt.Sprintf("Bot *%s*", actor.Login)
		}

	case "Organization":
		actorName = fmt.Sprintf("org *%s*", actor.Name)
		if capitalise {
			actorName = fmt.Sprintf("Org *%s*", actor.Name)
		}

	case "User":
		actorName = fmt.Sprintf("user *%s*", actor.Login)
		if capitalise {
			actorName = fmt.Sprintf("User *%s*", actor.Login)
		}

		if len(actor.Name) > 0 {
			actorName = fmt.Sprintf("%s (%s)", actorName, actor.Name)
		}
	}

	return actorName
}

func postSlackMessage(timestamp, text, slackAlertsChannel, slackWebHookURL string) {
	payload := slack.Payload{
		Text:      fmt.Sprintf("_%s_\n%s\n\n", timestamp, text),
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
		log.Fatalf("Failed to parse time '%s': %v", t, err)
	}

	return t.Format("Monday 02 Jan 2006 15:04:05 MST")
}
