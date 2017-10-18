package convoy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	consulServiceDiscoveryPath = "v1/catalog/service/"
	consulUnavailable          = "Consul service unavailable"
)

type ConsulServices []ConsulService

func (c ConsulServices) ToHosts() []Host {
	hosts := []Host{}
	for _, service := range c {
		hosts = append(hosts, Host{
			IPAddress: service.Address,
			Port:      service.ServicePort,
			Tags:      service.ToTags(),
		})
	}
	return hosts
}

type ConsulService struct {
	Address     string
	Datacenter  string
	ServiceName string
	ServicePort int
	ServiceTags []string
}

func (c *ConsulService) ToTags() Tags {
	tags := Tags{}
	for _, tag := range c.ServiceTags {
		// TODO : Example
		if strings.HasPrefix(tag, "az-") {
			tags.AZ = tag
		}
	}
	return tags
}

func (s *Server) registration(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	svcName := path[len(path)-1]

	resp, err := s.consulClient.Get(
		fmt.Sprintf("%s/%s/%s", s.ConsulBaseURL, consulServiceDiscoveryPath, svcName))

	if err != nil {
		s.error(w, r, fmt.Sprintf("Service discovery call failed at %s: %s", s.ConsulBaseURL.String(), err.Error()))
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.error(w, r, fmt.Sprintf("Service discovery call failed at %s: %s", s.ConsulBaseURL.String(), err.Error()))
		return
	}

	var consulService ConsulServices
	err = json.Unmarshal(data, &consulService)
	if err != nil {
		s.error(w, r, fmt.Sprintf("Failed to unmarshal JSON : %s", err.Error()))
		return
	}

	hosts := consulService.ToHosts()
	hostsBytes, err := json.Marshal(hosts)
	if err != nil {
		s.error(w, r, fmt.Sprintf("Failed to marshal JSON : %s", err.Error()))
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(hostsBytes)
}

func (s *Server) error(w http.ResponseWriter, r *http.Request, msg string) {
	log.Printf(msg)
	w.WriteHeader(500)
	w.Write([]byte(consulUnavailable))
}
