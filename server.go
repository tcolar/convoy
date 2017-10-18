package convoy

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// Server is the convoy server, a shim between Envoy and Consul
type Server struct {
	Port          int
	ConsulBaseURL *url.URL
	consulClient  *http.Client
	srv           *http.Server
}

// NewServer creates a Convoy server
func NewServer(port int, consulBaseURL *url.URL) *Server {
	return &Server{
		Port:          port,
		ConsulBaseURL: consulBaseURL,
		consulClient:  &http.Client{},
	}
}

// Start starts the server, in the background
func (s *Server) Start() {
	s.srv = &http.Server{Addr: fmt.Sprintf(":%d", s.Port)}

	http.HandleFunc("/v1/registration/", s.registration)

	log.Printf("Starting convoy on :%d\n", s.Port)

	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			log.Fatalf("Convoy: ListenAndServe() error: %s", err.Error())
		}
	}()
}

// Stop gracefully stops the server
func (s *Server) Stop() {
	log.Println("Stopping convoy")
	if s.srv != nil {
		s.srv.Shutdown(nil)
	}
}
