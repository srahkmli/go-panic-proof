package panicrecovery

import (
	"log"
	"net/http"
	"runtime/debug"
)

// HTTPRecover is an HTTP middleware that recovers from panics and logs the error.
func HTTPRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic error with stack trace.
				log.Printf("Recovered from panic: %v\nStack Trace: %s", err, string(debug.Stack()))

				// Respond with 500 Internal Server Error.
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
