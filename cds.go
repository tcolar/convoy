package convoy

import (
	"encoding/json"
	"fmt"
	"net/http"

	consulapi "github.com/hashicorp/consul/api"
)

type CDSResp struct {
	Clusters []Cluster `json:"clusters"`
}

type Cluster struct {
	Name                          string            `json:"name"`                                        //required
	Type                          string            `json:"type"`                                        //required
	ConnectTimeoutMS              int               `json:"connect_timeout_ms"`                          //required
	PerConnectionBufferLimitBytes int               `json:"per_connection_buffer_limit_bytes,omitempty"` //optional
	LBType                        string            `json:"lb_type"`                                     //required
	ServiceName                   string            `json:"service_name"`                                //required
	HealthCheck                   *HealthCheck      `json:"health_check,omitempty"`                      //optional
	MaxRequestsPerConnection      int               `json:"max_requests_per_connection,omitempty"`       //optional
	CircuitBreakers               *CircuitBreakers  `json:"circuit_breakers,omitempty"`                  //optional
	SSLContext                    *SSLContext       `json:"ssl_context,omitempty"`                       //optional
	Features                      string            `json:"features,omitempty"`                          //optional
	HTTP2Settings                 *HTTP2Settings    `json:"http2_settings,omitempty"`                    //optional
	CleanupIntervalMS             int               `json:"cleanup_interval_ms,omitempty"`               //optional
	DNSRefreshRateMS              int               `json:"dns_refresh_rate_ms,omitempty"`               //optional
	DNSLookupFamily               string            `json:"dns_lookup_family,omitempty"`                 //optional
	DNSResolvers                  []string          `json:"dns_resolvers,omitempty"`                     //optional
	OutlierDetection              *OutlierDetection `json:"outlier_detection,omitempty"`                 //optional
	//Hosts                         []Host           `json:"hosts"`
}

type HealthCheck struct {
	Type string `json:"type,omitempty""`
}

type CircuitBreakers struct {
}

type SSLContext struct {
}

type HTTP2Settings struct {
}

type OutlierDetection struct {
}

func (s *Server) GetClusters(w http.ResponseWriter, r *http.Request) {
	// These fields are not necessary and will be ignored
	// path := strings.Split(r.URL.Path, "/")
	// serviceCluster := path[len(path)-2]
	// serviceNode := path[len(path)-1]

	catalog := s.ConsulAPI.Catalog()
	q := consulapi.QueryOptions{}
	services, _, err := catalog.Services(&q)
	if err != nil {
		s.error(w, r, fmt.Sprintf("Failed to list Consul registered services: %s", err.Error()))
		return
	}

	clusters := ToClusters(services)
	clusterBytes, err := json.Marshal(clusters)
	if err != nil {
		s.error(w, r, fmt.Sprintf("Failed to marshal JSON : %s", err.Error()))
		return
	}

	w.Write(clusterBytes)
}

func ToClusters(services map[string][]string) CDSResp {
	cdsResp := CDSResp{}

	for serviceName, _ := range services {
		cdsResp.Clusters = append(cdsResp.Clusters, Cluster{
			Name:             serviceName,
			Type:             "sds",
			ConnectTimeoutMS: 250,
			LBType:           "round_robin",
			ServiceName:      serviceName,
		})
	}
	return cdsResp
}
