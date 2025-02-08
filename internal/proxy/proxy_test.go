package proxy

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestServer_ServeHTTP(t *testing.T) {
	// Create a mock target server
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Echo back the query parameters for testing
		w.Write([]byte(r.URL.RawQuery))
	}))
	defer targetServer.Close()

	// Set the target URL for the proxy
	os.Setenv("TARGET_URL", targetServer.URL)

	// Create our proxy server
	proxyServer := NewServer()

	tests := []struct {
		name          string
		rscHeader     bool
		initialQuery  string
		expectedQuery string
	}{
		{
			name:          "No RSC header",
			rscHeader:     false,
			initialQuery:  "",
			expectedQuery: "",
		},
		{
			name:          "RSC header with no query",
			rscHeader:     true,
			initialQuery:  "",
			expectedQuery: "_rsc=1",
		},
		{
			name:          "RSC header with existing query",
			rscHeader:     true,
			initialQuery:  "foo=bar",
			expectedQuery: "_rsc=1&foo=bar",
		},
		{
			name:          "RSC header with existing _rsc",
			rscHeader:     true,
			initialQuery:  "_rsc=2",
			expectedQuery: "_rsc=2",
		},
		{
			name:          "RSC header with multiple queries",
			rscHeader:     true,
			initialQuery:  "foo=bar&baz=qux",
			expectedQuery: "_rsc=1&baz=qux&foo=bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test request
			req := httptest.NewRequest("GET", "/?"+tt.initialQuery, nil)
			if tt.rscHeader {
				req.Header.Set("RSC", "1")
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			proxyServer.ServeHTTP(rr, req)

			// Check response
			if rr.Code != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, http.StatusOK)
			}

			// The mock server echoes back the query string, so we can verify it
			if got := rr.Body.String(); got != tt.expectedQuery {
				t.Errorf("handler returned unexpected query string: got %v want %v",
					got, tt.expectedQuery)
			}
		})
	}
}
