package main

import (
	"fmt"
	"log"
	"net/http"
	"mattmajestic/twitch-go/internal/handler"
	"mattmajestic/twitch-go/internal/services"
)

func main() {
	// Load environment variables
	services.LoadEnv()

	// Define routes
	http.HandleFunc("/", handler.HomeHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/followers", handler.FollowersHandler)

	// Start server
	fmt.Println("Starting server at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
