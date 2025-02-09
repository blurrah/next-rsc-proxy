package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type Server struct {
	proxy *httputil.ReverseProxy
}

func NewServer() *Server {
	url, err := url.Parse(os.Getenv("TARGET_URL"))
	if err != nil {
		log.Fatalf("Failed to parse target URL: %v", err)
	}
	return &Server{
		proxy: httputil.NewSingleHostReverseProxy(url),
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	originalURL := r.URL.String()

	// Skip modification for .rsc requests as they're already RSC payloads
	if strings.HasSuffix(r.URL.Path, ".rsc") {
		s.proxy.ServeHTTP(w, r)
		return
	}

	if r.Header.Get("RSC") == "1" {
		targetQuery := r.URL.Query()
		// Add `_rsc` query param if not present
		if !targetQuery.Has("_rsc") {
			// TODO: Generate a key based on the Next-Router-State-Tree instead of 1
			targetQuery.Add("_rsc", "1")
			r.URL.RawQuery = targetQuery.Encode()

			w.Header().Add("X-Forwarded-URL", r.URL.String())
		}
	}

	rw := &responseWriter{
		ResponseWriter: w,
		originalURL:    originalURL,
		modifiedURL:    r.URL.String(),
	}

	s.proxy.ServeHTTP(rw, r)
}

func (s *Server) Start() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return http.ListenAndServe(":"+port, s)
}

type responseWriter struct {
	http.ResponseWriter
	originalURL string
	modifiedURL string
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.Header().Add("X-Rsc-Proxy-Original-Url", rw.originalURL)
	rw.Header().Add("X-Rsc-Proxy-Modified-Url", rw.modifiedURL)
	rw.ResponseWriter.WriteHeader(statusCode)
}
