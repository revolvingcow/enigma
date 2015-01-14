package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/generate", GenerateHandler)
	http.HandleFunc("/api/update", ProfileHandler)
	http.HandleFunc("/api/refresh", RefreshHandler)
	http.HandleFunc("/api/remove", RemoveHandler)
	http.HandleFunc("/book", BookHandler)
	http.HandleFunc("/", DefaultHandler)

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
