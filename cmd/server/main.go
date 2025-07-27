package main

import (
	"log"
	"time"

	"palu-wiki-backend/configs"          // Import configs package
	"palu-wiki-backend/internal/handler" // Import handler package
	"palu-wiki-backend/internal/repository"
	"palu-wiki-backend/internal/service" // Import service package
	"palu-wiki-backend/pkg/gemini"       // Import gemini package

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := configs.LoadConfig()                                                        // Load configuration
	log.Printf("Loaded Gemini API Key (first 5 chars): %s*****", cfg.GeminiAPIKey[:5]) // Log first 5 chars for security
	repository.InitDB()                                                                // Initialize database connection

	// Initialize Gemini client
	geminiClient, err := gemini.NewClient(cfg.GeminiAPIKey)
	if err != nil {
		log.Fatalf("Failed to initialize Gemini client: %v", err)
	}
	defer geminiClient.Close() // Ensure client is closed when main exits

	gin.SetMode(gin.DebugMode) // Set Gin to debug mode for more detailed logs
	r := gin.Default()

	// Initialize repositories and services
	guideRepo := repository.NewGuideRepository(repository.DB)
	updateService := service.NewUpdateService(guideRepo, geminiClient) // Pass geminiClient

	// Initialize handlers
	geminiHandler := handler.NewGeminiHandler(geminiClient)
	adminHandler := handler.NewAdminHandler(guideRepo, updateService) // Initialize AdminHandler

	// Register routes
	apiGroup := r.Group("/api/v1")
	{
		apiGroup.POST("/gemini/generate", geminiHandler.GenerateContent)
		apiGroup.POST("/admin/guides/topic", adminHandler.CreateGuideTopic) // New API for creating guide topics
		apiGroup.GET("/admin/guides", adminHandler.GetGuidesForAdmin)       // New API for admin to view guides

		// Public Guide APIs for frontend
		guideHandler := handler.NewGuideHandler(guideRepo)
		apiGroup.GET("/guides", guideHandler.GetGuides)
		apiGroup.GET("/guides/:id", guideHandler.GetGuideByID)
	}

	// Start periodic update task
	go func() {
		// Fetch updates immediately on startup
		log.Println("Starting initial Steam news fetch...")
		updates, err := updateService.FetchSteamNews("https://store.steampowered.com/news/app/1623730")
		if err != nil {
			log.Printf("Initial Steam news fetch failed: %v", err)
		} else {
			if err := updateService.ProcessUpdates(updates); err != nil {
				log.Printf("Initial Steam news processing failed: %v", err)
			}
		}

		// Then fetch updates periodically (e.g., every 1 hour)
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			log.Println("Fetching Steam news...")
			updates, err := updateService.FetchSteamNews("https://store.steampowered.com/news/app/1623730")
			if err != nil {
				log.Printf("Steam news fetch failed: %v", err)
				continue
			}
			if err := updateService.ProcessUpdates(updates); err != nil {
				log.Printf("Steam news processing failed: %v", err)
			}
		}
	}()

	log.Printf("Server started on :%s", cfg.Port)
	if err := r.Run("0.0.0.0:" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
