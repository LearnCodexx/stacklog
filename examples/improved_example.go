package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/learncodexx/stacklog"
)

// Example demonstrating both global and instance-based logging approaches

func main() {
	// Option 1: Global logging (backward compatible)
	stacklog.Init("ExampleService")
	stacklog.Startup("Starting example application with improvements")

	// Option 2: Instance-based logging (preferred for new code)
	userLogger := stacklog.NewStacklog("UserService")
	orderLogger := stacklog.NewStacklog("OrderService")

	// Custom error patterns
	stacklog.AddErrorMapping("insufficient_funds", "You don't have enough balance for this transaction.", false)
	stacklog.AddErrorMapping("item_out_of_stock", "This item is currently out of stock.", false)

	// Setup HTTP with automatic cleanup
	app := fiber.New()
	app.Use(stacklog.HTTP())

	// User service endpoints
	app.Post("/users", func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		
		// Using instance logger
		userLogger.API(ctx, "Creating new user")
		
		// Simulate some work
		if err := createUser(ctx, userLogger); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		
		userLogger.API(ctx, "User created successfully")
		return c.JSON(fiber.Map{"status": "success"})
	})

	// Order service endpoints  
	app.Post("/orders", func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		
		// Using instance logger
		orderLogger.API(ctx, "Processing new order")
		
		if err := processOrder(ctx, orderLogger); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": stacklog.ErrorPattern(err)})
		}
		
		return c.JSON(fiber.Map{"status": "order processed"})
	})

	// Monitor memory usage
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		
		for range ticker.C {
			count := stacklog.GetRequestLogCount()
			fmt.Printf("Active request logs: %d\n", count)
		}
	}()

	// Clean shutdown
	defer func() {
		stacklog.StopCleanup()
		userLogger.Info("SYSTEM", "Application shutting down")
	}()

	fmt.Println("Server starting on :3000")
	app.Listen(":3000")
}

func createUser(ctx context.Context, logger *stacklog.Stacklog) error {
	// Simulate some validation logic
	logger.API(ctx, "Validating user data")
	
	// Simulate database check
	logger.API(ctx, "Checking if email exists")
	
	// Return success
	return nil
}

func processOrder(ctx context.Context, logger *stacklog.Stacklog) error {
	logger.API(ctx, "Validating order data")
	
	// Simulate a custom error that will be translated
	return fmt.Errorf("insufficient_funds")
}