package convoy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	consulServiceDiscoveryPath = "v1/catalog/service/"
	consulUnavailable          = "Consul service unavailable"
)

func (s *Server) registration(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	svcName := path[len(path)-1]

	resp, err := s.consulClient.Get(
		fmt.Sprintf("%s/%s/%s", s.ConsulBaseUrl, consulServiceDiscoveryPath, svcName))

	if err != nil {
		log.Printf("Service discovery call failed at %s: %s", s.ConsulBaseUrl.String(), err.Error())
		w.WriteHeader(500)
		w.Write([]byte(consulUnavailable))
		return
	}

	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
