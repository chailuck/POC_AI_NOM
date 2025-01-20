// cmd/server/main.go
package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/your-username/tmf632-service/internal/config"
	"github.com/your-username/tmf632-service/internal/database"
	"github.com/your-username/tmf632-service/internal/handlers"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize handlers
	h := handlers.NewHandler(db, cfg)

	// Routes
	api := e.Group("/tmf-api/partyManagement/v4")
	api.POST("/individual", h.CreateIndividual)
	api.GET("/individual/:id", h.GetIndividual)
	api.PUT("/individual/:id", h.UpdateIndividual)
	api.DELETE("/individual/:id", h.DeleteIndividual)
	api.GET("/individual", h.ListIndividuals)

	// Start server
	e.Logger.Fatal(e.Start(":" + cfg.ServerPort))
}
