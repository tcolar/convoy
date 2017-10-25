package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/tcolar/convoy"
)

func main() {

	var port int

	flag.IntVar(&port, "port", 3500, "Port number to bind to")
	flag.Parse()

	server := convoy.NewServer(port)
	server.Start()

	// Run until we get a syscall.SIGTERM (Ctrl+C)
	c := make(chan os.Signal, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt)
	<-c

	// Attempt a clean shutdown
	server.Stop()
}
