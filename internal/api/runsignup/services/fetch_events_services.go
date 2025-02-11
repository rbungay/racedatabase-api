package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/rbungay/racedatabase-api/config"
	"github.com/rbungay/racedatabase-api/internal/api/runsignup/constants"
	"github.com/rbungay/racedatabase-api/internal/api/runsignup/models"
)


func FetchEvents(state, city, eventType, startDate, endDate, minDistance, maxDistance, zipcode, radius string) ([]models.Event, error) {

	apiURL := config.GetEnv("RUNSIGNUP_API_URL", "")
	apiKey := config.GetEnv("RUNSIGNUP_API_KEY", "")
	apiSecret := config.GetEnv("RUNSIGNUP_API_SECRET", "")


	params := url.Values{}
	params.Set("api_key", apiKey)               
	params.Set("api_secret", apiSecret)         
	params.Set("format", "json")                
	params.Set("state", state)                  

	if eventType != "" {
		if _, isValid := constants.validEventTypes[eventType]; ! isValid {
			return nil, fmt.Errorf("invalid event_type: %s. Must be one of: %v", eventType, validEventTypes)
		} 
		params.Set("event_type", eventType)
	}

	if city != "" {
		params.Set("city", city) 
	}

	if startDate != "" {
		params.Set("start_date", startDate) 
	}
	if endDate != "" {
		params.Set("end_date", endDate) 
	}
	if minDistance != "" {
		params.Set("min_distance", minDistance) 
	}
	if maxDistance != "" {
		params.Set("max_distance", maxDistance) 
	}
	if zipcode != "" {
		params.Set("zipcode", zipcode) 
	}
	if radius != "" {
		params.Set("radius", radius) 
	}

	
	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())
	fmt.Println("Making API request to:", fullURL) 

	
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	
	fmt.Println("Raw API Response:", string(body)) 

	
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

	
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	
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
