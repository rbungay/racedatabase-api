package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rbungay/racedatabase-api/internal/api/runsignup/services"
)

var FetchEventsFunc = services.FetchEvents


func RunSignupEventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	
	state := r.URL.Query().Get("state")
	city := r.URL.Query().Get("city")
	eventType := r.URL.Query().Get("event_type")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	minDistance := r.URL.Query().Get("min_distance")
	maxDistance := r.URL.Query().Get("max_distance")
	zipcode := r.URL.Query().Get("zipcode")
	radius := r.URL.Query().Get("radius")

	
	if state == "" {
		http.Error(w, "State parameter is required", http.StatusBadRequest)
		return
	}

	
	events, err := FetchEventsFunc(state, city, eventType, startDate, endDate, minDistance, maxDistance, zipcode, radius)
	if err != nil {
		http.Error(w, "Error fetching events: "+err.Error(), http.StatusInternalServerError)
		return
	}

	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)

	fmt.Printf("Fetched events for state: %s\n", state)
}
