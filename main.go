package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

// CLI Arguments
type Config struct {
	dev  *bool
	port *int
}

var config Config

func setup() {
	// Parse CLI Arguments
	config.dev = flag.Bool("dev", false, "run in development mode")
	config.port = flag.Int("p", 8080, "port to run server on")
	flag.Parse()
}

func main() {
	setup()

	http.Handle("/", http.FileServer(http.Dir("./dist")))

	portStr := fmt.Sprintf(":%d", *config.port)

	log.Printf("[INFO] Server running on %s\n", portStr)
	log.Fatal(http.ListenAndServe(portStr, nil))
}
