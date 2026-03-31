package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/learncodexx/stacklog"
)

// This example demonstrates the usage of the logging package with proper hover documentation.
// When you hover over these function calls in your IDE, you should see detailed descriptions
// of what each function does and how to use it.

func main() {
	// Hover over Init - shows full documentation about one-line setup
	stacklog.Init("ExampleService")

	// Hover over these shortcuts - shows when and how to use them
	stacklog.Startup("Starting example application")
	stacklog.Config("Configuration loaded")
	stacklog.Database("Connected to database")

	// Error logging examples - hover shows error handling patterns
	if err := loadConfig(); err != nil {
		stacklog.ConfigError("Failed to load configuration", err)
		return
	}

	// HTTP setup - hover shows middleware setup information
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			stacklog.SystemError("Unhandled error", err)
			return ctx.Status(500).JSON(fiber.Map{"error": "Internal server error"})
		},
	})

	// Hover over HTTP - shows automatic request grouping info
	app.Use(stacklog.HTTP())

	app.Post("/users", createUserHandler)
	app.Listen(":8080")
}

func createUserHandler(c *fiber.Ctx) error {
	ctx := c.UserContext()

	// Hover over API - shows context-aware logging info
	stacklog.API(ctx, "Processing user creation request")

	// Local API logger - hover shows CheckAndLogAPI pattern
	localAPI := stacklog.Local("UserHandler")
	var err error
	defer localAPI.CheckAndLogAPI(ctx, &err, "User creation failed")

	// Simulate some work
	if err = validateUserData(); err != nil {
		return err // Will be automatically logged with stack trace
	}

	// Hover over APIError - shows automatic grouping behavior 
	if err = createUser(); err != nil {
		stacklog.APIError(ctx, "Database operation failed", err)
		return err
	}

	stacklog.API(ctx, "User created successfully")
	return c.JSON(fiber.Map{"status": "success"})
}

// Example utility functions - hover shows error checking patterns
func loadConfig() error {
	// Hover over CheckError - shows simple error checking
	if stacklog.CheckError(someOperation(), "Configuration validation failed") {
		return fmt.Errorf("config error")
	}
	return nil
}

func someOperation() error {
	return nil // Placeholder
}

func validateUserData() error {
	return nil // Placeholder  
}

func createUser() error {
	return nil // Placeholder
}