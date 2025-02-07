package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rbungay/racedatabase-api/internal/api/runsignup/services"
)


func RunSignupEventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}


	state := r.URL.Query().Get("state")
	if state == "" {
		http.Error(w, "State parameter is required", http.StatusBadRequest)
		return
	}


	events, err := services.FetchEventsByState(state)
	if err != nil {
		http.Error(w, "Error fetching events: "+err.Error(), http.StatusInternalServerError)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)

	fmt.Printf("Fetched events for state: %s\n", state)
}
