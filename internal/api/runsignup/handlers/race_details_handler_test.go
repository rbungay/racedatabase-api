package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rbungay/racedatabase-api/internal/api/runsignup/models"
)

// Mock FetchRaceDetails function
var mockFetchRaceDetails = func(raceID int) (*models.RaceDetails, error) {
	return &models.RaceDetails{
		ID:         raceID,
		Name:       "Test Race",
		URL:        "https://example.com",
		ExternalURL: "https://external-example.com",
		LogoURL:    "https://example.com/logo.png",
		Timezone:   "America/New_York",
		Events: []models.EventDetails{
			{
				EventID:   98765,
				Name:      "5K Run",
				StartTime: "2025-06-01 08:00",
				EndTime:   "2025-06-01 10:00",
				EventType: "running_race",
				Distance:  "5K",
				RegOpens:  "2025-01-01 00:00",
				Category:  "Running",
				RegPeriods: []models.RegistrationPeriod{
					{
						Opens:   "2025-01-01 00:00",
						Closes:  "2025-05-30 23:59",
						Fee:     "$30.00",
						ProcFee: "$2.50",
					},
				},
			},
		},
	}, nil
}

// Mock function for error case
var mockFetchRaceDetailsError = func(raceID int) (*models.RaceDetails, error) {
	return nil, errors.New("failed to fetch race details")
}

func TestRunSignupRaceDetailsHandler_Success(t *testing.T) {
	req, err := http.NewRequest("GET", "/runsignup/race?race_id=12345", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := RunSignupRaceDetailsHandler(mockFetchRaceDetails)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var raceDetails models.RaceDetails
	if err := json.Unmarshal(rr.Body.Bytes(), &raceDetails); err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	if raceDetails.ID != 12345 {
		t.Errorf("Unexpected race ID: got %v, want %v", raceDetails.ID, 12345)
	}
	if raceDetails.Name != "Test Race" {
		t.Errorf("Unexpected race name: got %v, want %v", raceDetails.Name, "Test Race")
	}
}

func TestRunSignupRaceDetailsHandler_MissingRaceID(t *testing.T) {
	req, err := http.NewRequest("GET", "/runsignup/race", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := RunSignupRaceDetailsHandler(mockFetchRaceDetails)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := "race_id parameter is required\n"
	if rr.Body.String() != expected {
		t.Errorf("Unexpected response body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestRunSignupRaceDetailsHandler_InvalidRaceID(t *testing.T) {
	req, err := http.NewRequest("GET", "/runsignup/race?race_id=invalid", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := RunSignupRaceDetailsHandler(mockFetchRaceDetails)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := "Invalid race_id format\n"
	if rr.Body.String() != expected {
		t.Errorf("Unexpected response body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestRunSignupRaceDetailsHandler_FetchError(t *testing.T) {
	req, err := http.NewRequest("GET", "/runsignup/race?race_id=12345", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := RunSignupRaceDetailsHandler(mockFetchRaceDetailsError)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := "Error fetching race details: failed to fetch race details\n"
	if rr.Body.String() != expected {
		t.Errorf("Unexpected response body: got %v want %v", rr.Body.String(), expected)
	}
}
