package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ONSdigital/github-auditor/internal/event"
	"github.com/ONSdigital/github-auditor/pkg/github"
)

func main() {
	firestoreCredentials := os.Getenv("FIRESTORE_CREDENTIALS")

	firestoreProject := ""
	if firestoreProject = os.Getenv("FIRESTORE_PROJECT"); len(firestoreProject) == 0 {
		log.Fatal("Missing FIRESTORE_PROJECT environment variable")
	}

	token := ""
	if token = os.Getenv("GITHUB_TOKEN"); len(token) == 0 {
		log.Fatal("Missing GITHUB_TOKEN environmental variable")
	}

	organisation := ""
	if organisation = os.Getenv("GITHUB_ORG_NAME"); len(organisation) == 0 {
		log.Fatal("Missing GITHUB_ORG_NAME environmental variable")
	}

	slackAlertsChannel := ""
	if slackAlertsChannel = os.Getenv("SLACK_ALERTS_CHANNEL"); len(slackAlertsChannel) == 0 {
		log.Fatal("Missing SLACK_ALERTS_CHANNEL environment variable")
	}

	slackWebHookURL := ""
	if slackWebHookURL = os.Getenv("SLACK_WEBHOOK"); len(slackWebHookURL) == 0 {
		log.Fatal("Missing SLACK_WEBHOOK environment variable")
	}

	client := github.NewClient(token)
	events, err := client.FetchAllAuditEvents(organisation)
	if err != nil {
		log.Fatalf("Failed to fetch audit log entries: %v", err)
	}

	// Using fmt rather than log so the output goes to STDOUT rather than STDERR.
	fmt.Printf("Audit log API query returned %d results\n", len(events))

	if len(firestoreCredentials) > 0 {
		event.ProcessWithCredentials(events, firestoreCredentials, firestoreProject, slackAlertsChannel, slackWebHookURL)
	} else {
		event.Process(events, firestoreProject, slackAlertsChannel, slackWebHookURL)
	}
}
