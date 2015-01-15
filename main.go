package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/api/generate", GenerateHandler)
	http.HandleFunc("/api/update", ProfileHandler)
	http.HandleFunc("/api/refresh", RefreshHandler)
	http.HandleFunc("/api/remove", RemoveHandler)
	http.HandleFunc("/signout", SignOutHandler)
	http.HandleFunc("/book", BookHandler)
	http.HandleFunc("/", DefaultHandler)

	log.Fatal(http.ListenAndServe(address(), nil))
}

// Retrieve the web server address from the environment variable ENIGMA_SERVER if possible.
// If the environment variable is not set then default to "localhost:8080".
func address() string {
	env := os.Getenv("ENIGMA_SERVER")
	if env == "" {
		return "localhost:8080"
	}
	return env
}
