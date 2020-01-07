package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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

	slackWebHook := ""
	if slackWebHook = os.Getenv("SLACK_WEBHOOK"); len(slackWebHook) == 0 {
		log.Fatal("Missing SLACK_WEBHOOK environment variable")
	}

	client := github.NewClient(token)

	auditEntries, err := client.FetchAllAuditLogEntries(organisation)
	if err != nil {
		log.Fatalf("Failed to fetch audit log entries: %v", err)
	}

	json, err := json.MarshalIndent(auditEntries, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", json)
}
