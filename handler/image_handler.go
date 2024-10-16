package handler

import (
	"API_Call_GVA_and_LLMAI/repositories"
	"API_Call_GVA_and_LLMAI/services"
	"bytes"
	"encoding/json"
	"image/png" // Add this package if you want to use PNG
	"log"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/labstack/echo/v4"
)

type ImageRequest struct {
	ImageURL string `json:"image_url"`
}

type ImageResponse struct {
	Message string `json:"message"`
}

func ProcessImage(c echo.Context) error {
	var req ImageRequest

	// Step 1: Decode the request
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil || req.ImageURL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Step 2: Check if the image result has already been saved in Firestore
	existingResult, exists := repositories.CheckFirestore(req.ImageURL)
	if exists {
		log.Printf("Document exists for %s, returning cached result", req.ImageURL)
		return c.JSON(http.StatusOK, existingResult)
	}

	// Step 3: Load pixel data from image_url
	pixels, err := loadImagePixels(req.ImageURL)
	if err != nil {
		log.Printf("Image could not be loaded, assuming it's safe: %s", req.ImageURL)
		repositories.SaveResult(req.ImageURL, "", "Image is whitelisted") // Change this
		return c.JSON(http.StatusOK, ImageResponse{Message: "Image is whitelisted"})
	}

	// Step 4: Send pixels to the AI model
	llmResult, err := services.AnalyzeImageWithLLM(pixels) // Change this
	if err != nil {
		log.Printf("Error analyzing image with LLM: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error analyzing image"})
	}

	// Check for trigger words in the response
	if hasTriggerWords(llmResult) {
		log.Printf("Trigger words found in response for %s, marking as bad", req.ImageURL)
		repositories.SaveResult(req.ImageURL, llmResult, "Image is bad") // Change this
		return c.JSON(http.StatusOK, ImageResponse{Message: "Image is bad"})
	}

	// If no trigger words, consider the image safe
	log.Printf("No trigger words found for %s, marking as safe", req.ImageURL)
	repositories.SaveResult(req.ImageURL, llmResult, "Image is whitelisted") // Change this
	return c.JSON(http.StatusOK, ImageResponse{Message: "Image is whitelisted"})
}

// Function to load pixel data from image_url and convert it to []byte
func loadImagePixels(imageURL string) ([]byte, error) {
	img, err := imaging.Open(imageURL)
	if err != nil {
		return nil, err
	}

	// Convert image.Image to []byte
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil { // Change this if you are using a different format
		return nil, err
	}

	return buf.Bytes(), nil
}

// Check if the response from LLM contains trigger words
func hasTriggerWords(response string) bool {
	triggerWords := []string{"bad", "blocked", "unsafe"} // Replace with actual trigger words
	for _, word := range triggerWords {
		if strings.Contains(response, word) {
			return true
		}
	}
	return false
}
