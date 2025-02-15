package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rbungay/racedatabase-api/internal/api/runsignup/services"
)

// RunSignupRaceDetailsHandler handles requests for fetching race details
func RunSignupRaceDetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get race_id from query params
	raceIDStr := r.URL.Query().Get("race_id")
	if raceIDStr == "" {
		http.Error(w, "race_id parameter is required", http.StatusBadRequest)
		return
	}

	raceID, err := strconv.Atoi(raceIDStr)
	if err != nil {
		http.Error(w, "Invalid race_id format", http.StatusBadRequest)
		return
	}

	// Fetch race details
	raceDetails, err := services.FetchRaceDetails(raceID)
	if err != nil {
		http.Error(w, "Error fetching race details: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(raceDetails)
}
