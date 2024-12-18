package contextTimeoutMiddleware

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func TimeoutMiddleware(timeoutValue int) mux.MiddlewareFunc {
	// Convert timeout to a duration
	timeout := time.Duration(timeoutValue) * time.Second

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a context with a timeout
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			// Replace the request context with the new timeout context
			r = r.WithContext(ctx)

			// Channel to handle when the handler completes
			done := make(chan struct{})

			// Run the handler in a goroutine
			go func() {
				next.ServeHTTP(w, r)
				close(done)
			}()

			// Wait for handler to complete or timeout
			select {
			case <-done:
				// Handler completed successfully
			case <-ctx.Done():
				// Context timeout triggered
				http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			}
		})
	}
}
