package convoy

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	consulapi "github.com/hashicorp/consul/api"
)

// Server is the convoy server, a shim between Envoy and Consul
type Server struct {
	Port         int
	ConsulAPI    *consulapi.Client
	srv          *http.Server
	QueryOptions consulapi.QueryOptions
	ConsulKeys   ConsulKeys
}

type ConsulKeys struct {
	Keys map[string]consulapi.KVPairs
	sync.RWMutex
}

// NewServer creates a Convoy server
func NewServer(port int) *Server {

	consul, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to connect to Consul : %s", err.Error()))
		os.Exit(1)
	}

	return &Server{
		Port:         port,
		ConsulAPI:    consul,
		QueryOptions: consulapi.QueryOptions{},
		ConsulKeys: ConsulKeys{
			Keys: map[string]consulapi.KVPairs{},
		},
	}
}

// Start starts the server, in the background
func (s *Server) Start() {
	s.srv = &http.Server{Addr: fmt.Sprintf(":%d", s.Port)}

	// Need to fork GetEnvoyKeys
	go s.GetConsulKeys()

	http.HandleFunc("/v1/registration/", s.GetService)

	http.HandleFunc("/v1/clusters/", s.GetClusters)

	http.HandleFunc("/v1/routes/", s.GetRoutes)

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

func (s *Server) error(w http.ResponseWriter, r *http.Request, msg string) {
	log.Printf(msg)
	w.WriteHeader(500)
	w.Write([]byte(consulUnavailable))
}
