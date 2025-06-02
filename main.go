package main

import (
	"log"
	"net/http"

	"stargazer/api"
)

func main() {
	addr := ":8080"
	log.Printf("Starting server on %s", addr)
	err := http.ListenAndServe(addr, http.HandlerFunc(api.Handler))
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
