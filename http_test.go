package panicrecovery

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPRecover(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.HandlerFunc
		expectedStatus int
	}{
		{
			name: "no panic",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "with panic",
			handler: func(w http.ResponseWriter, r *http.Request) {
				panic("test panic")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := HTTPRecover(tt.handler)
			server := httptest.NewServer(handler)
			defer server.Close()

			resp, err := http.Get(server.URL)
			if err != nil {
				t.Fatalf("Fialed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, but got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestHTTPRecoverWithHandler(t *testing.T) {
	customStatus := http.StatusServiceUnavailable
	customHandler := func(w http.ResponseWriter, r *http.Request, err any) {
		w.WriteHeader(customStatus)
	}

	tests := []struct {
		name           string
		handler        http.HandlerFunc
		errorHandler   HTTPErrorHandler
		expectedStatus int
	}{
		{
			name: "custom error handler",
			handler: func(w http.ResponseWriter, r *http.Request) {
				panic("test panic")
			},
			errorHandler:   customHandler,
			expectedStatus: customStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := HTTPRecoverWithHandler(tt.handler, tt.errorHandler)
			server := httptest.NewServer(handler)
			defer server.Close()

			resp, err := http.Get(server.URL)
			if err != nil {
				t.Fatalf("Fialed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, but got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}
