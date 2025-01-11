package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

func main() {
	// Initialize database
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create repository and service instances using factory
	factory := NewMetadataFactory(db)
	handler := NewMetadataHandler(factory.CreateService())

	// Initialize Fiber app
	app := fiber.New()

	// Enable CORS for all origins
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Setup routes
	setupRoutes(app, handler)

	// Start server
	log.Fatal(app.Listen(":3000"))
}

func setupRoutes(app *fiber.App, handler *MetadataHandler) {
	api := app.Group("/api")

	api.Post("/metadata", handler.Create)
	api.Get("/metadata/:key", handler.Get)
	api.Put("/metadata/:key", handler.Update)
	api.Patch("/metadata/:key", handler.PatchUpdate)
	api.Delete("/metadata/:key", handler.Delete)
	api.Delete("/metadata", handler.DeleteAll)
	api.Put("/metadata", handler.Upsert)
}
