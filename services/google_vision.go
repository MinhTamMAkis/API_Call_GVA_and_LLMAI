package services

import (
	vision "cloud.google.com/go/vision/apiv1"
	"context"
)

func ProcessWithGoogleVision(imageURL string) (map[string]interface{}, error) {
	ctx := context.Background()
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	image := vision.NewImageFromURI(imageURL)
	annotations, err := client.DetectLabels(ctx, image, nil, 10)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for _, annotation := range annotations {
		result[annotation.Description] = annotation.Score
	}

	return result, nil
}

func HasTriggerWord(visionResult map[string]interface{}) bool {
	// Define a list of trigger words
	triggerWords := []string{"violence", "explicit"}

	for word := range visionResult {
		for _, trigger := range triggerWords {
			if word == trigger {
				return true
			}
		}
	}
	return false
}
