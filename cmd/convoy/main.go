package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/tcolar/convoy"
)

func main() {

	var port int
	var consulBaseURL string

	flag.IntVar(&port, "port", 3500, "Port number to bind to")
	flag.StringVar(&consulBaseURL, "consul", "http://127.0.0.1:8500", "Consul base url, ie: http://consul.acme.com")

	flag.Parse()

	consul, err := url.Parse(consulBaseURL)
	if err != nil {
		log.Fatalf("Failed to parse consul url %s : %s", consulBaseURL, err.Error())
		os.Exit(1)
	}

	server := convoy.NewServer(port, consul)
	server.Start()

	// Run until we get a syscall.SIGTERM (Ctrl+C)
	c := make(chan os.Signal, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt)
	<-c

	// Attempt a clean shutdown
	server.Stop()
}
