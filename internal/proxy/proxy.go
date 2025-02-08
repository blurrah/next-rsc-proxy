package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
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
	if r.Header.Get("RSC") == "1" {
		targetQuery := r.URL.Query()
		// Add `_rsc` query param if not present
		if !targetQuery.Has("_rsc") {
			// TODO: Generate a key based on the Next-Router-State-Tree instead of 1
			targetQuery.Set("_rsc", "1")
			r.URL.RawQuery = targetQuery.Encode()
		}

	}
	s.proxy.ServeHTTP(w, r)
}

func (s *Server) Start() error {
	return http.ListenAndServe(":8080", s)
}
