package event

import (
	"encoding/json"
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
func Process(events []github.Node, firestoreProject, slackAlertsChannel, slackWebHookURL string) {
	process(events, nil, firestoreProject, slackAlertsChannel, slackWebHookURL)
}

// ProcessWithCredentials processes the passed slice of GitHub audit events, creating Slack alerts in the passed Slack channel for events of interest.
func ProcessWithCredentials(events []github.Node, firestoreCredentials, firestoreProject, slackAlertsChannel, slackWebHookURL string) {
	process(events, &firestoreCredentials, firestoreProject, slackAlertsChannel, slackWebHookURL)
}

func process(events []github.Node, firestoreCredentials *string, firestoreProject, slackAlertsChannel, slackWebHookURL string) {
	for _, e := range events {
		timestamp := formatTime(e.CreatedAt)
		id := e.ID
		action := e.Action
		text := ""
		jsonData, err := json.Marshal(e)
		if err != nil {
			log.Fatalf("Failed marshalling event to JSON: %v", err)
		}

		switch e.Action {

		// OAuth events.
		case "oauth_application.create":
			text = fmt.Sprintf(github.MessageForEvent(action), e.OauthApplicationName, e.OrganizationName, formatActor(e.Actor, false))

		// Organisation events.
		case "org.add_billing_manager":
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.Actor, true), formatActor(e.User, false), e.OrganizationName)
		case "org.add_member":
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.User, true), e.OrganizationName)
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
			text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.Actor, true), formatActorOrEmail(e.User, e.Email, false), e.OrganizationName)
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
		case "org.restore_member":
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

			// A repo.change_merge_setting event is fired with a null merge setting when a new repo is created, so only log explicit merge setting changes.
			if len(e.MergeType) > 0 {
				text = fmt.Sprintf(github.MessageForEvent(action), formatActor(e.Actor, true), e.RepositoryName, strings.ToLower(e.MergeType))
			} else {
				text = ""
			}
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

			// Using fmt rather than log so the output goes to STDOUT rather than STDERR.
			fmt.Printf("Unknown GitHub event: %s\n", action)
		}

		var client *firestore.Client

		if firestoreCredentials != nil {
			client = firestore.NewClientWithCredentials(firestoreProject, *firestoreCredentials)
		} else {
			client = firestore.NewClient(firestoreProject)
		}

		if !client.DocExists(id, timestamp, action) && len(text) > 0 {
			logJSON(jsonData)
			postSlackMessage(timestamp, text, slackAlertsChannel, slackWebHookURL)
		}

		err = client.SaveDoc(id, timestamp, action)
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

func formatActorOrEmail(actor github.Actor, email string, capitalise bool) string {
	if len(email) > 0 {
		return fmt.Sprintf("*%s*", email)
	}

	return formatActor(actor, capitalise)
}

func logJSON(jsonData []byte) {

	// Using fmt rather than log so the output goes to STDOUT rather than STDERR.
	fmt.Println(string(jsonData))
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
