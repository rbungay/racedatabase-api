package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rbungay/racedatabase-api/internal/api/runsignup/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system env variables.")
	}

	http.HandleFunc("/runsignup/events", handlers.RunSignupEventsHandler)

	port:= os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}


	fmt.Println("Server is running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}