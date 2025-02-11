package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/rbungay/racedatabase-api/config"
)

// Event represents a race event from RunSignup
type Event struct {
	ID        int    `json:"race_id"`
	Name      string `json:"name"`
	StartDate string `json:"next_date"`
	EndDate   string `json:"next_end_date"`
	URL       string `json:"url"`
	City      string `json:"city"`
	State     string `json:"state"`
	Zipcode   string `json:"zipcode"`
}

// FetchEventsByState fetches events from RunSignup filtered by state
func FetchEventsByState(state string) ([]Event, error) {
	apiURL := config.GetEnv("RUNSIGNUP_API_URL", "")
	apiKey := config.GetEnv("RUNSIGNUP_API_KEY", "")
	apiSecret := config.GetEnv("RUNSIGNUP_API_SECRET", "")

	// Build query parameters
	params := url.Values{}
	params.Set("api_key", apiKey)
	params.Set("api_secret", apiSecret)
	params.Set("format", "json")
	params.Set("state", state)
	params.Set("events", "F")
	params.Set("race_headings", "F")
	params.Set("race_links", "F")
	params.Set("include_waiver", "F")
	params.Set("include_multiple_waivers", "F")
	params.Set("include_event_days", "F")
	params.Set("include_extra_date_info", "F")
	params.Set("page", "1")
	params.Set("results_per_page", "50")
	params.Set("start_date", "today")
	params.Set("only_partner_races", "F")
	params.Set("search_start_date_only", "F")
	params.Set("only_races_with_results", "F")
	params.Set("distance_units", "K")

	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())
	fmt.Println("Making API request to:", fullURL)

	// Make HTTP request
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Debug: Print raw API response
	fmt.Println("Raw API Response:", string(body))

	// Define a struct that matches the JSON response
	var data struct {
		Races []struct {
			Race struct {
				ID        int    `json:"race_id"`
				Name      string `json:"name"`
				StartDate string `json:"next_date"`
				EndDate   string `json:"next_end_date"`
				URL       string `json:"url"`
				Address   struct {
					City    string `json:"city"`
					State   string `json:"state"`
					Zipcode string `json:"zipcode"`
				} `json:"address"`
			} `json:"race"`
		} `json:"races"`
	}

	// Parse JSON response
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Extract relevant race details
	var events []Event
	for _, race := range data.Races {
		events = append(events, Event{
			ID:        race.Race.ID,
			Name:      race.Race.Name,
			StartDate: race.Race.StartDate,
			EndDate:   race.Race.EndDate,
			URL:       race.Race.URL,
			City:      race.Race.Address.City,
			State:     race.Race.Address.State,
			Zipcode:   race.Race.Address.Zipcode,
		})
	}

	return events, nil
}
