package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rbungay/racedatabase-api/internal/api/runsignup/services"
	"github.com/rbungay/racedatabase-api/internal/api/runsignup/models"

)

var mockFetchEvents = func(state, city, eventType, startDate, endDate, minDistance, maxDistance, zipcode, radius string) ([]models.Event, error) {
	return []models.Event{
		{
			ID:         12345,
			Name:       "Test Race",
			URL:        "https://example.com",
			ExternalURL: "https://external.com",
			LogoURL:    "https://example.com/logo.png",
			Category:   "Runs",
		},
	}, nil
}

func TestRunSignupEventsHandler_ValidRequest(t *testing.T) {
	FetchEventsFunc = mockFetchEvents
	defer func() { FetchEventsFunc = services.FetchEvents}()

	req, err := http.NewRequest("GET", "/runsignup/events?state=NY&event_type=running_race", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RunSignupEventsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var events []models.Event
	if err := json.Unmarshal(rr.Body.Bytes(), &events); err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	if len(events) == 0 {
		t.Errorf("Expected events, got empty response")
	}
	if events[0].Name != "Test Race" {
		t.Errorf("Unexpected event name: got %v want %v", events[0].Name, "Test Race")
	}

	if events[0].ExternalURL != "https://external.com" {
		t.Errorf("Unexpected External URL: got %v want %v", events[0].ExternalURL, "https://external.com")
	}
	
	if events[0].LogoURL != "https://example.com/logo.png" {
		t.Errorf("Unexpected Logo URL: got %v want %v", events[0].LogoURL, "https://example.com/logo.png")
	}
	
}
