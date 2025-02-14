package services

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"os"

	"github.com/rbungay/racedatabase-api/config"
)

func mockRunSignupAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := `{
		"races": [{
			"race": {
				"race_id": 12345,
				"name": "Test Race",
				"next_date": "2025-01-01",
				"next_end_date": "2025-01-02",
				"url": "https://example.com",
				"external_race_url": "https://external-example.com",
				"logo_url": "https://example.com/logo.png",
				"event_type": "running_race",
				"address": {
					"city": "New York",
					"state": "NY",
					"zipcode": "10001"
				}
			}
		}]
	}`
	w.Write([]byte(response))
}

func mockRunSignupAPI_Fail(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "API Error", http.StatusInternalServerError)
}

func TestFetchEvents_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(mockRunSignupAPI))
	defer mockServer.Close()

	originalAPIURL := config.GetEnv("RUNSIGNUP_API_URL", "")
	os.Setenv("RUNSIGNUP_API_URL", mockServer.URL)
	defer os.Setenv("RUNSIGNUP_API_URL", originalAPIURL)

	events, err := FetchEvents("NY", "New York", "running_race", "2025-01-01", "2025-12-31", "", "", "", "")

	if err != nil {
		t.Fatalf("FetchEvents failed: %v", err)
	}
	if len(events) == 0 {
		t.Fatalf("Expected events, got none")
	}

	expectedName := "Test Race"
	expectedExternalURL := "https://external-example.com"
	expectedLogoURL := "https://example.com/logo.png"

	if events[0].Name != expectedName {
		t.Errorf("Unexpected event name: got %v, want %v", events[0].Name, expectedName)
	}
	if events[0].ExternalURL != expectedExternalURL {
		t.Errorf("Unexpected external URL: got %v, want %v", events[0].ExternalURL, expectedExternalURL)
	}
	if events[0].LogoURL != expectedLogoURL {
		t.Errorf("Unexpected logo URL: got %v, want %v", events[0].LogoURL, expectedLogoURL)
	}
}

func TestFetchEvents_Failure(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(mockRunSignupAPI_Fail))
	defer mockServer.Close()

	originalAPIURL := config.GetEnv("RUNSIGNUP_API_URL", "")
	os.Setenv("RUNSIGNUP_API_URL", mockServer.URL)
	defer os.Setenv("RUNSIGNUP_API_URL", originalAPIURL)

	events, err := FetchEvents("NY", "", "", "", "", "", "", "", "")

	if err == nil {
		t.Fatalf("Expected an error but got none")
	}
	if len(events) != 0 {
		t.Errorf("Expected no events on failure, but got %d", len(events))
	}
}

func TestFetchEvents_InvalidEventType(t *testing.T) {
	events, err := FetchEvents("NY", "", "invalid_event", "", "", "", "", "", "")

	if err == nil {
		t.Fatalf("Expected an error for invalid event type but got none")
	}
	if len(events) != 0 {
		t.Errorf("Expected no events with invalid event type, but got %d", len(events))
	}
}

func mockRunSignupAPI_PartialFail(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("event_type") == "running_race" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := `{
		"races": [{
			"race": {
				"race_id": 67890,
				"name": "Valid Triathlon",
				"next_date": "2025-06-15",
				"next_end_date": "2025-06-16",
				"url": "https://example.com/triathlon",
				"external_race_url": "https://external-triathlon.com",
				"logo_url": "https://example.com/triathlon-logo.png",
				"event_type": "triathlon",
				"address": {
					"city": "Los Angeles",
					"state": "CA",
					"zipcode": "90001"
				}
			}
		}]
	}`
	w.Write([]byte(response))
}

func TestFetchEvents_PartialFailure(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(mockRunSignupAPI_PartialFail))
	defer mockServer.Close()

	originalAPIURL := config.GetEnv("RUNSIGNUP_API_URL", "")
	os.Setenv("RUNSIGNUP_API_URL", mockServer.URL)
	defer os.Setenv("RUNSIGNUP_API_URL", originalAPIURL)

	events, err := FetchEvents("CA", "", "", "", "", "", "", "", "")

	if err == nil {
		t.Fatalf("Expected an error due to partial failure but got none")
	}
	if len(events) == 0 {
		t.Errorf("Expected some successful events, got none")
	}

	expectedName := "Valid Triathlon"
	expectedExternalURL := "https://external-triathlon.com"
	expectedLogoURL := "https://example.com/triathlon-logo.png"

	if events[0].Name != expectedName {
		t.Errorf("Unexpected event name: got %v, want %v", events[0].Name, expectedName)
	}
	if events[0].ExternalURL != expectedExternalURL {
		t.Errorf("Unexpected external URL: got %v, want %v", events[0].ExternalURL, expectedExternalURL)
	}
	if events[0].LogoURL != expectedLogoURL {
		t.Errorf("Unexpected logo URL: got %v, want %v", events[0].LogoURL, expectedLogoURL)
	}
}
