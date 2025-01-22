package panicrecovery

import (
	"net/http"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
)

var logger, _ = zap.NewProduction()

// HTTPRecover is an HTTP middleware that recovers from panics and logs the error.
func HTTPRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Structual log with zap
				logger.Error("Recovered from panic",
					zap.Any("error", err),
					zap.String("stack_trace", string(debug.Stack())),
					zap.Time("timestamp", time.Now()))

				// Respond with 500 Internal Server Error.
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

type HTTPErrorHandler func(w http.ResponseWriter, r *http.Request, err interface{})

func HTTPRecoverWithHandler(next http.Handler, errorHandler HTTPErrorHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if errorHandler != nil {
					errorHandler(w, r, err)
				} else {
					logger.Error("Recovered from panic",
						zap.Any("error", err),
						zap.String("stack_trace", string(debug.Stack())),
						zap.Time("timestamp", time.Now()))

					// Default error response
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
