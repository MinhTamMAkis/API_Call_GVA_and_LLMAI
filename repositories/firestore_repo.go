package repositories

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
)

var client *firestore.Client

// Initialize Firestore client
func init() {
	var err error
	projectID := os.Getenv("FIRESTORE_PROJECT_ID") // Get project ID from environment variable
	if projectID == "" {
		log.Fatal("FIRESTORE_PROJECT_ID environment variable not set")
	}
	client, err = firestore.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
}

// Check if the document exists in Firestore
func CheckFirestore(imageURL string) (map[string]interface{}, bool) {
	doc, err := client.Collection("image_checks").Doc(imageURL).Get(context.Background())
	if err != nil {
		if status.Code(err) == codes.NotFound { // Use status.Code to check the status code
			return nil, false // Document does not exist
		}
		log.Printf("Failed to get document: %v", err)
		return nil, false
	}
	return doc.Data(), true
}

// Save the result to Firestore
func SaveResult(imageURL string, llmResult string, conclusion string) {
	_, err := client.Collection("image_filtering_results").Doc(imageURL).Set(context.Background(), map[string]interface{}{
		"image_url":                   imageURL,
		"gemini1.5Flash8bFlashResult": llmResult,
		"conclusion":                  conclusion,
	})
	if err != nil {
		log.Printf("Failed to save result: %v", err)
	}
}
