package services

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"os"

	"github.com/rbungay/racedatabase-api/config"
)

func mockRaceDetailsAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := `{
		"race": {
			"race_id": 12345,
			"name": "Test Race",
			"url": "https://example.com",
			"external_race_url": "https://external-example.com",
			"logo_url": "https://example.com/logo.png",
			"timezone": "America/New_York",
			"events": [{
				"event_id": 98765,
				"name": "5K Run",
				"start_time": "2025-06-01 08:00",
				"end_time": "2025-06-01 10:00",
				"event_type": "running_race",
				"distance": "5K",
				"registration_opens": "2025-01-01 00:00",
				"registration_periods": [{
					"registration_opens": "2025-01-01 00:00",
					"registration_closes": "2025-05-30 23:59",
					"race_fee": "$30.00",
					"processing_fee": "$2.50"
				}]
			}]
		}
	}`
	w.Write([]byte(response))
}

func TestFetchRaceDetails_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(mockRaceDetailsAPI))
	defer mockServer.Close()

	originalAPIURL := config.GetEnv("RUNSIGNUP_API_URL", "")
	os.Setenv("RUNSIGNUP_API_URL", mockServer.URL)
	defer os.Setenv("RUNSIGNUP_API_URL", originalAPIURL)

	raceDetails, err := FetchRaceDetails(12345)

	if err != nil {
		t.Fatalf("FetchRaceDetails failed: %v", err)
	}
	if raceDetails == nil {
		t.Fatalf("Expected race details, got nil")
	}

	expectedName := "Test Race"
	expectedURL := "https://example.com"
	expectedExternalURL := "https://external-example.com"
	expectedLogoURL := "https://example.com/logo.png"
	expectedTimezone := "America/New_York"

	if raceDetails.Name != expectedName {
		t.Errorf("Unexpected race name: got %v, want %v", raceDetails.Name, expectedName)
	}
	if raceDetails.URL != expectedURL {
		t.Errorf("Unexpected URL: got %v, want %v", raceDetails.URL, expectedURL)
	}
	if raceDetails.ExternalURL != expectedExternalURL {
		t.Errorf("Unexpected external URL: got %v, want %v", raceDetails.ExternalURL, expectedExternalURL)
	}
	if raceDetails.LogoURL != expectedLogoURL {
		t.Errorf("Unexpected logo URL: got %v, want %v", raceDetails.LogoURL, expectedLogoURL)
	}
	if raceDetails.Timezone != expectedTimezone {
		t.Errorf("Unexpected timezone: got %v, want %v", raceDetails.Timezone, expectedTimezone)
	}
}
