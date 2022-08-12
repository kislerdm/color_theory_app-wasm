//go:build !unittest
// +build !unittest

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var dir, port string

func init() {
	dir = os.Getenv("DIR_WEB")
	if dir == "" {
		log.Fatalln("'DIR_WEB' with web assets to be server must be set as envvar")
	}

	port = os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
}

func handler() http.Handler {
	log.Printf("open for dev server: http://localhost:%s\n", port)
	log.Println(dir)
	return http.FileServer(http.Dir(dir))
}

func main() {
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler()); err != nil {
		log.Fatalln("Failed to start server", err)
	}
}
