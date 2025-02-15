package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
	"net/http"

	"github.com/rbungay/racedatabase-api/config"
	"github.com/rbungay/racedatabase-api/internal/api/runsignup/constants"
	"github.com/rbungay/racedatabase-api/internal/api/runsignup/models"
)

func FetchRaceDetails(raceID int) (*models.RaceDetails, error) {
	apiURL := config.GetEnv("RUNSIGNUP_API_URL", "")
	apiKey := config.GetEnv("RUNSIGNUP_API_KEY", "")
	apiSecret := config.GetEnv("RUNSIGNUP_API_SECRET", "")

	fullURL := fmt.Sprintf("%s/race/%d?api_key=%s&api_secret=%s&format=json", apiURL, raceID, apiKey, apiSecret)
	fmt.Println("Fetching race details from:", fullURL)

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

	var data struct {
		Race struct {
			ID         int    `json:"race_id"`
			Name       string `json:"name"`
			URL        string `json:"url"`
			ExternalURL string `json:"external_race_url"`
			LogoURL    string `json:"logo_url"`
			Timezone   string `json:"timezone"`
			Events     []struct {
				EventID       int    `json:"event_id"`
				Name          string `json:"name"`
				StartTime     string `json:"start_time"`
				EndTime       string `json:"end_time"`
				EventType     string `json:"event_type"`
				Distance      string `json:"distance"`
				RegOpens      string `json:"registration_opens"`
				RegPeriods    []struct {
					Opens   string `json:"registration_opens"`
					Closes  string `json:"registration_closes"`
					Fee     string `json:"race_fee"`
					ProcFee string `json:"processing_fee"`
				} `json:"registration_periods"`
			} `json:"events"`
		} `json:"race"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	raceDetails := &models.RaceDetails{
		ID:         data.Race.ID,
		Name:       data.Race.Name,
		URL:        data.Race.URL,
		ExternalURL: data.Race.ExternalURL,
		LogoURL:    data.Race.LogoURL,
		Timezone:   data.Race.Timezone,
		Events:     []models.EventDetails{},
	}

	for _, event := range data.Race.Events {
		category, exists := constants.EventTypeToCategory[event.EventType]
		if !exists {
			category = constants.CategoryOther
		}

		eventDetails := models.EventDetails{
			EventID:    event.EventID,
			Name:       event.Name,
			StartTime:  event.StartTime,
			EndTime:    event.EndTime,
			EventType:  event.EventType,
			Distance:   event.Distance,
			RegOpens:   event.RegOpens,
			Category:   string(category), // ✅ FIXED: Convert EventCategory to string
			RegPeriods: []models.RegistrationPeriod{},
		}

		for _, regPeriod := range event.RegPeriods {
			eventDetails.RegPeriods = append(eventDetails.RegPeriods, models.RegistrationPeriod{
				Opens:    regPeriod.Opens,
				Closes:   regPeriod.Closes,
				Fee:      regPeriod.Fee,
				ProcFee:  regPeriod.ProcFee,
			})
		}

		raceDetails.Events = append(raceDetails.Events, eventDetails)
	}

	return raceDetails, nil
}
