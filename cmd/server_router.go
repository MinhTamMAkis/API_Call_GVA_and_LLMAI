package cmd

import (
	"API_Call_GVA_and_LLMAI/handler"
	"github.com/labstack/echo/v4"
)

// InitRouter initializes the router for the API
func InitRouter() *echo.Echo {
	e := echo.New()

	// Define the route for the image processing API
	e.POST("/api/v1/process_image", handler.ProcessImage)

	return e
}
