package contextTimeoutMiddleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TimeoutMiddleware(timeoutValue int) fiber.Handler {
	// Convert timeout to a duration
	timeout := time.Duration(timeoutValue) * time.Second

	return func(c *fiber.Ctx) error {
		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(c.Context(), timeout)
		defer cancel()

		// Set the new context to the Fiber context
		c.SetUserContext(ctx)

		// Channel to handle when the handler completes
		done := make(chan error, 1)

		// Run the handler in a goroutine
		go func() {
			done <- c.Next()
		}()

		// Wait for handler to complete or timeout
		select {
		case err := <-done:
			// Handler completed
			return err
		case <-ctx.Done():
			// Context timeout triggered
			return c.Status(fiber.StatusGatewayTimeout).JSON(fiber.Map{
				"error": "Request timed out",
			})
		}
	}
}
