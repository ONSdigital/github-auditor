package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ONSdigital/github-auditor/internal/event"
	"github.com/ONSdigital/github-auditor/pkg/github"
)

func main() {
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

	event.Process(events, slackAlertsChannel, slackWebHookURL)

	json, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", json)
}
