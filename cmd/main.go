package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/rbungay/racedatabase-api/internal/api/runsignup/services"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system env variables.")
	} else {
		log.Println(".env file loaded successfully.")
	}

	// Test fetch for New Jersey races
	events, err := services.FetchEvents("NJ", "", "", "", "", "", "", "", "")
	if err != nil {
		log.Fatalf("Error fetching events: %v", err)
	}
	fmt.Printf("Successfully fetched %d events from New Jersey\n", len(events))

	// Comment out existing server code
	/*
	http.HandleFunc("/runsignup/events", handlers.RunSignupEventsHandler)

	http.HandleFunc("/runsignup/race/", handlers.RunSignupRaceDetailsHandler(services.FetchRaceDetails))


	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}


	fmt.Println("Server is running on http://localhost:" + port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	*/
}
