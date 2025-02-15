package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rbungay/racedatabase-api/internal/api/runsignup/models"
 
)

func RunSignupRaceDetailsHandler(fetchRaceDetails func(int) (*models.RaceDetails, error)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
            return
        }

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

        raceDetails, err := fetchRaceDetails(raceID)
        if err != nil {
            http.Error(w, "Error fetching race details: "+err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(raceDetails)
    }
}
