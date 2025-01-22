package panicrecovery

import (
	"net/http"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
)

type HTTPErrorHandler func(w http.ResponseWriter, r *http.Request, err interface{})

// HTTPRecover is an HTTP middleware that recovers from panics and logs the error.
func HTTPRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer handleHTTPPanic(w, r, nil)
		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// HTTPRecoverWithHandler adds customizable error responses
func HTTPRecoverWithHandler(next http.Handler, errorHandler HTTPErrorHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer handleHTTPPanic(w, r, errorHandler)
		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func handleHTTPPanic(w http.ResponseWriter, r *http.Request, errorHandler HTTPErrorHandler) {
	if err := recover(); err != nil {
		Logger.Error("Recovered from panic",
			zap.Any("error", err),
			zap.String("stack_trace", string(debug.Stack())),
			zap.Time("timestamp", time.Now()),
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method))

		if errorHandler != nil {
			errorHandler(w, r, err)
			return
		}
		// Default error response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
