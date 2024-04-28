package main

import (
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	_ "github.com/kanagavel07/cake" // Blank import so the init() of GCP function framework functions runs
)

func main() {
	// If PORT environment variable exists, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	// By default, listening on all interfaces. If testing locally, running with
	// LOCAL_ONLY=true to avoid triggering firewall warnings and
	// exposing the server outside of our own local machine.
	hostname := ""

	if localOnly := os.Getenv("LOCAL_ONLY"); localOnly == "true" {
		hostname = "127.0.0.1"
	}

	// StartHostPort serves an HTTP server with registered function(s) on the given host and port,
	// It also creates the new *http.ServeMux internally which allows for straightforward routing.
	if err := funcframework.StartHostPort(hostname, port); err != nil {
		log.Fatalf("funcframework.StartHostPort: %v\n", err)
	}
}
