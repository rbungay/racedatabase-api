package services

import (
	"net/http"
	"net/http/httptest"
	"testing"
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


func TestFetchEvents_Success(t *testing.T) {
	// Step 1: Start a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(mockRunSignupAPI))
	defer mockServer.Close()

	// Step 2: Temporarily override API URL
	originalAPIURL := config.GetEnv("RUNSIGNUP_API_URL", "")
	config.SetEnv("RUNSIGNUP_API_URL", mockServer.URL) // Assume SetEnv exists
	defer config.SetEnv("RUNSIGNUP_API_URL", originalAPIURL)

	// Step 3: Call FetchEvents with test parameters
	events, err := FetchEvents("NY", "New York", "running_race", "2025-01-01", "2025-12-31", "", "", "", "")

	// Step 4: Validate results
	if err != nil {
		t.Fatalf("FetchEvents failed: %v", err)
	}
	if len(events) == 0 {
		t.Fatalf("Expected events, got none")
	}

	// Step 5: Verify data
	expectedName := "Test Race"
	if events[0].Name != expectedName {
		t.Errorf("Unexpected event name: got %v, want %v", events[0].Name, expectedName)
	}
}


func mockRunSignupAPI_Fail(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "API Error", http.StatusInternalServerError)
}

func TestFetchEvents_Failure(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(mockRunSignupAPI_Fail))
	defer mockServer.Close()

	// Override API URL temporarily
	originalAPIURL := config.GetEnv("RUNSIGNUP_API_URL", "")
	config.SetEnv("RUNSIGNUP_API_URL", mockServer.URL)
	defer config.SetEnv("RUNSIGNUP_API_URL", originalAPIURL)

	events, err := FetchEvents("NY", "", "", "", "", "", "", "", "")

	if err == nil {
		t.Fatalf("Expected an error but got none")
	}
	if len(events) != 0 {
		t.Errorf("Expected no events on failure, but got %d", len(events))
	}
}
