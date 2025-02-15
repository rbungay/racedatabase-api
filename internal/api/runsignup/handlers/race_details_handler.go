package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/rbungay/racedatabase-api/internal/api/runsignup/services"
)

// RunSignupRaceDetailsHandler handles requests for fetching race details
func RunSignupRaceDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure only GET requests are allowed
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract race ID from the URL
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/runsignup/race/"), "/")
	if len(pathParts[0]) == 0 {
		http.Error(w, "Race ID is required", http.StatusBadRequest)
		return
	}

	raceID, err := strconv.Atoi(pathParts[0])
	if err != nil {
		http.Error(w, "Invalid race ID format", http.StatusBadRequest)
		return
	}

	// Fetch race details using the service
	raceDetails, err := services.FetchRaceDetails(raceID)
	if err != nil {
		http.Error(w, "Error fetching race details: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(raceDetails)
}
