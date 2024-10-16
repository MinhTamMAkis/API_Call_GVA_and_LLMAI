package cmd

import (
	"API_Call_GVA_and_LLMAI/handler"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
)

type ApiServer struct {
	e *echo.Echo
}

func NewApiServer() *ApiServer {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Check the FIRESTORE_PROJECT_ID environment variable
	projectID := os.Getenv("FIRESTORE_PROJECT_ID")
	log.Printf("FIRESTORE_PROJECT_ID: %s", projectID)
	if projectID == "" {
		log.Fatal("FIRESTORE_PROJECT_ID environment variable not set")
	}

	// Check the GOOGLE_APPLICATION_CREDENTIALS environment variable
	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credentialsFile == "" {
		log.Fatal("GOOGLE_APPLICATION_CREDENTIALS environment variable not set")
	}

	// Initialize Echo instance
	e := echo.New()

	// Declare ApiServer
	server := &ApiServer{
		e: e,
	}

	// Register routes
	server.registerRoutes()

	return server
}

func (s *ApiServer) Start(port string) {
	// Start the server
	log.Printf("Starting server at port %s", port)
	if err := s.e.Start(port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}
}

func (s *ApiServer) registerRoutes() {
	// Register API routes
	s.e.POST("/api/v1/process_image", handler.ProcessImage)
}

func (s *ApiServer) Stop() {
	// Function to stop the server when necessary
	if err := s.e.Shutdown(nil); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}
}
