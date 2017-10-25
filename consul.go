package convoy

import (
	"fmt"
	"log"
	"strings"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

const (
	consulUnavailable = "Consul service unavailable"
)

func (s *Server) GetConsulKeys() {

	ticker := time.NewTicker(time.Millisecond * 500)

	kv := s.ConsulAPI.KV()

	for range ticker.C {

		log.Println("Fetching consul keys")

		kvPairs, queryMeta, err := kv.List("convoy", &s.QueryOptions)
		if err != nil {

		}
		s.QueryOptions.WaitIndex = queryMeta.LastIndex

		log.Println(fmt.Sprintf("Found %d keys in Consul", len(kvPairs)))

		s.ConsulKeys.Lock()

		s.ConsulKeys.Keys = map[string]consulapi.KVPairs{}

		for _, kv := range kvPairs {
			// 'convoy/<service_name>/<cds,rds,sds>/*'
			names := strings.Split(kv.Key, "/")
			serviceName := names[1]
			log.Println(fmt.Sprintf("Processing service %s", serviceName))
			s.ConsulKeys.Keys[serviceName] = append(s.ConsulKeys.Keys[serviceName], kv)
		}

		s.ConsulKeys.Unlock()
	}
}
