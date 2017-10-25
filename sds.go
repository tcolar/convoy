package convoy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

type SDSResp struct {
	Hosts []Host `json:"hosts"`
}

type Host struct {
	IPAddress string `json:"ip_address"`
	Port      int    `json:"port"`
	Tags      Tags   `json:"tags"`
}

type Tags struct {
	AZ                  string `json:"az,omitempty"`
	Canary              string `json:"canary,omitempty"`
	LoadBalancingWeight int    `json:"load_balancing_weight,omitempty"`
}

func (s *Server) GetService(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	svcName := path[len(path)-1]

	catalog := s.ConsulAPI.Catalog()
	q := consulapi.QueryOptions{}
	service, _, err := catalog.Service(svcName, "", &q)
	if err != nil {
		s.error(w, r, fmt.Sprintf("Failed to retrieve Consul service meta for '%s': %s", svcName, err.Error()))
		return
	}

	hosts := ToHosts(service)

	hostsBytes, err := json.Marshal(hosts)
	if err != nil {
		s.error(w, r, fmt.Sprintf("Failed to marshal JSON : %s", err.Error()))
		return
	}

	w.Write(hostsBytes)
}

func ToHosts(services []*consulapi.CatalogService) SDSResp {
	sdsResp := SDSResp{}
	for _, service := range services {
		sdsResp.Hosts = append(sdsResp.Hosts, Host{
			IPAddress: service.Address,
			Port:      service.ServicePort,
			Tags:      ToTags(service.ServiceTags),
		})
	}
	return sdsResp
}

func ToTags(tags []string) Tags {
	serviceTags := Tags{}

	for _, tag := range tags {
		// TODO : Example
		if strings.HasPrefix(tag, "az:") {
			serviceTags.AZ = tag
		}
		if strings.HasPrefix(tag, "canary:") {
			serviceTags.Canary = tag
		}
		if strings.HasPrefix(tag, "loadbalancingweight:") {
			weight, err := strconv.ParseInt(tag, 10, 64)
			if err != nil {
				serviceTags.LoadBalancingWeight = int(weight)
			}
		}
	}
	return serviceTags
}
