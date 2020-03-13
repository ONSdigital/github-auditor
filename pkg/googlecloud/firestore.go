package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// NewClient instantiates a new Firestore client for the passed GCP project.
func NewClient(projectID string) *Client {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to instantiate Firestore client in project %s: %v", projectID, err)
	}

	return &Client{
		projectID: projectID,
		context:   &ctx,
		client:    client,
	}
}

// NewClientWithCredentials instantiates a new Firestore client for the passed GCP project using the passed path to a JSON service account key file.
func NewClientWithCredentials(projectID, credentialsFile string) *Client {
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
	snapshot, err := c.client.Collection(firestoreCollection).Doc(id).Get(*c.context)
	if status.Code(err) == codes.NotFound {
		return false
	}

	data := snapshot.Data()
	if data == nil {
		return false
	}

	ts, err := snapshot.DataAt("timestamp")
	if err != nil && status.Code(err) != codes.NotFound {
		log.Fatalf("Failed to retrieve timestamp value from Firestore document %s/%s: %v", firestoreCollection, id, err)
	}

	a, err := snapshot.DataAt("action")
	if err != nil && status.Code(err) != codes.NotFound {
		log.Fatalf("Failed to retrieve action value from Firestore document %s/%s: %v", firestoreCollection, id, err)
	}

	return timestamp == ts && action == a
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
