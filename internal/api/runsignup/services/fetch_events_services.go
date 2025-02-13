package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"sync"

	"github.com/rbungay/racedatabase-api/config"
	"github.com/rbungay/racedatabase-api/internal/api/runsignup/constants"
	"github.com/rbungay/racedatabase-api/internal/api/runsignup/models"
)

func FetchEvents(state, city, eventType, startDate, endDate, minDistance, maxDistance, zipcode, radius string) ([]models.Event,error){
	var wg sync.WaitGroup
	var mu sync.Mutex
	var allEvents []models.Event
	var errorList []error

	if eventType != "" && !constants.ValidEventTypes[eventType] {
		return nil, fmt.Errorf("invalid event_type: %s. Must be one of: %v", eventType, constants.ValidEventTypes)
	}

	for eventType := range constants.ValidEventTypes{
		wg.Add(1)
		go func(eventType string){
			defer wg.Done()
			fmt.Println("Fetching event type:", eventType)

			events, err := fetchEventsFromAPI(state,city,eventType, startDate, endDate, minDistance, maxDistance, zipcode, radius)
			if err != nil{
				mu.Lock()
				errorList = append(errorList, fmt.Errorf("%s: %v", eventType, err))
				mu.Unlock()
				return
			}

			mu.Lock()
			allEvents = append(allEvents,events...)
			mu.Unlock()
		}(eventType)
	}

	wg.Wait()

	if len(errorList)>0 {
		return allEvents, fmt.Errorf("some event types failed to fetch: %v", errorList)
	}

	return allEvents, nil
}

func fetchEventsFromAPI(state, city, eventType, startDate, endDate, minDistance, maxDistance, zipcode, radius string) ([]models.Event, error) {
	apiURL := config.GetEnv("RUNSIGNUP_API_URL", "")
	apiKey := config.GetEnv("RUNSIGNUP_API_KEY", "")
	apiSecret := config.GetEnv("RUNSIGNUP_API_SECRET", "")

	params := url.Values{}
	params.Set("api_key", apiKey)               
	params.Set("api_secret", apiSecret)         
	params.Set("format", "json")                
      
	if state != "" {
		params.Set("state", state)
	} else {
		return nil, fmt.Errorf("state paramater is required")
	}

	if eventType != "" {
		if _, isValid := constants.ValidEventTypes[eventType]; ! isValid {
			return nil, fmt.Errorf("invalid event_type: %s. Must be one of: %v", eventType, constants.ValidEventTypes)
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status %d - response: %s", resp.StatusCode, string(body))
	}
	
	if config.GetEnv("ENV", "development") == "development" {
		// fmt.Println("Raw API Response:", string(body))
	}
	
	var data struct {
		Races []struct {
			Race struct {
				ID        int    `json:"race_id"`
				Name      string `json:"name"`
				StartDate string `json:"next_date"`
				EndDate   string `json:"next_end_date"`
				URL       string `json:"url"`
				EventType string `json:"event_type"`
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
	
	var events []models.Event
	for _, race := range data.Races {
	
		finalEventType := race.Race.EventType

		if finalEventType == "" {
			finalEventType = eventType
		}

		category, exists := constants.EventTypeToCategory[finalEventType]
		if !exists {
			category = "Other"
		}

		events = append(events, models.Event{
			ID:        race.Race.ID,
			Name:      race.Race.Name,
			StartDate: race.Race.StartDate,
			EndDate:   race.Race.EndDate,
			URL:       race.Race.URL,
			City:      race.Race.Address.City,
			State:     race.Race.Address.State,
			Zipcode:   race.Race.Address.Zipcode,
			EventType: finalEventType,
			Category: category,
		})
	}

	return events, nil
}
