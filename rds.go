package convoy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

type RDSResp struct {
	ValidateClusters        bool                `json:"valide_clusters,omitempty"`
	VirtualHosts            []VirtualHost       `json:"virtual_hosts"`
	InternalOnlyHeaders     []string            `json:"internal_only_headers,omitempty"`
	ResponseHeadersToAdd    []map[string]string `json:"response_headers_to_add,omitempty"`
	ResponseHeadersToRemove []string            `json:"response_headers_to_remove,omitempty"`
	RequestHeadersToAdd     []map[string]string `json:"request_headers_to_add,omitempty"`
}

type VirtualHost struct {
	Name                string              `json:"name"`
	Domains             []string            `json:"domains"`
	Routes              []Route             `json:"routes"`
	CORs                *CORs               `json:"cors,omitempty"`
	RequireSSL          string              `json:"require_ssl,omitempty"`
	VirtualClusters     []VirtualCluster    `json:"virtual_clusters,omitempty"`
	RateLimits          []RateLimit         `json:"rate_limits,omitempty"`
	RequestHeadersToAdd []map[string]string `json:"request_headers_to_add,omitempty"`
}

type CORs struct {
}

type Route struct {
	Prefix    string `json:"prefix"`
	Cluster   string `json:"cluster"`
	TimeOutMS int    `json:"timeout_ms,omitempty"`
}

type VirtualCluster struct {
}

type RateLimit struct {
}

func (s *Server) GetRoutes(w http.ResponseWriter, r *http.Request) {
	catalog := s.ConsulAPI.Catalog()
	q := consulapi.QueryOptions{}
	services, _, err := catalog.Services(&q)
	if err != nil {
		s.error(w, r, fmt.Sprintf("Failed to list Consul registered services: %s", err.Error()))
		return
	}

	routes := s.ToRoutes(services)
	clusterBytes, err := json.Marshal(routes)
	if err != nil {
		s.error(w, r, fmt.Sprintf("Failed to marshal JSON : %s", err.Error()))
		return
	}

	w.Write(clusterBytes)
}

func (s *Server) ToRoutes(services map[string][]string) RDSResp {
	rdsResp := RDSResp{}

	log.Println("Fetching routes")

	for serviceName := range services {
		//Default
		vHost := VirtualHost{
			Name:    serviceName,
			Domains: []string{fmt.Sprintf("%s.*", serviceName)},
			Routes:  []Route{{Prefix: "/", Cluster: serviceName, TimeOutMS: 250}},
		}

		// Fetch envoy keys for this service
		s.ConsulKeys.RLock()
		defer s.ConsulKeys.RUnlock()

		serviceKeys := s.ConsulKeys.Keys[serviceName]
		for _, kv := range serviceKeys {
			names := strings.Split(kv.Key, "/")
			discoveryType := names[2]
			if !strings.EqualFold(discoveryType, "rds") {
				continue
			}
			config := strings.Split(kv.Key, fmt.Sprintf("%s/%s/%s/", names[0], names[1], discoveryType))[1]
			switch config {
			case "domains":
				vHost.Domains = strings.Split(string(kv.Value), ",")
			case "routes":
			case "require_ssl":
			case "virtual_clusters":
			case "rate_limits":
			case "request_headers_to_add":
			}
		}

		rdsResp.VirtualHosts = append(rdsResp.VirtualHosts, vHost)
	}
	return rdsResp
}
