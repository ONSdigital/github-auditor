package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

type (

	// Client wraps a Google Firestore client.
	Client struct {
		projectID       string
		credentialsFile string
		context         *context.Context
		client          *firestore.Client
	}
)

const firestoreCollection = "github-auditor"

// NewClient instantiates a new Firestore client for the passed GCP project using the passed path to a JSON service account key file.
func NewClient(projectID, credentialsFile string) *Client {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.Fatalf("Failed to instantiate Firestore client in project %s using credentials file %s: %v", projectID, credentialsFile, err)
	}

	return &Client{
		projectID:       projectID,
		credentialsFile: credentialsFile,
		context:         &ctx,
		client:          client,
	}
}

// DocExists returns whether the Firestore document with the passed ID containing the passed timestamp and action exists.
func (c Client) DocExists(id, timestamp, action string) bool {
	// snapshot, err := c.client.Collection(firestoreCollection).Doc(id).Get(context.Background())
	// if err != nil {
	// 	log.Fatalf("Failed to retrieve document %s/%s from Firestore: %v", firestoreCollection, id, err)
	// }

	// data := snapshot.Data()
	// if data == nil {
	// 	return false
	// }

	// id, err := snapshot.DataAt("id")
	return false
}

// SaveDoc creates or updates the Firestore document with the passed ID, setting its contents to the passed timestamp and action.
func (c Client) SaveDoc(id, timestamp, action string) error {
	doc := c.client.Collection(firestoreCollection).Doc(id)
	_, err := doc.Set(*c.context, map[string]interface{}{
		"timestamp": timestamp,
		"action":    action,
	})

	return err
}