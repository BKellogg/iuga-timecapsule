package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	defaultHost = ""
	defaultPort = "80"
)

func main() {

	host := os.Getenv("HOST")
	if len(host) == 0 {
		fmt.Printf("no HOST environment variable set, defaulting to %s\n", defaultHost)
		host = defaultHost
	}
	port := os.Getenv("PORT")
	if len(port) == 0 {
		fmt.Printf("no PORT environment variable set, defaulting to %s\n", defaultPort)
		port = defaultPort
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	log.Fatal(http.ListenAndServe(addr, nil))
}
